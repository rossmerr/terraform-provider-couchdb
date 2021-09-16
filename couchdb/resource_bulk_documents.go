package couchdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rossmerr/couchdb_go/client/database"
	"github.com/rossmerr/couchdb_go/models"
)

func resourceBulkDocuments() *schema.Resource {
	return &schema.Resource{
		CreateContext: bulkDocumentsCreate,
		ReadContext:   bulkDocumentsRead,
		UpdateContext: bulkDocumentsUpdate,
		DeleteContext: bulkDocumentsDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database to associate design with",
			},
			"docs": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON Document, the _id field is required (wrap in <<EOF { } EOF)",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "partition key to add to the document id",
			},
		},
	}
}

func bulkDocumentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	var docs []models.Document
	err = json.Unmarshal([]byte(d.Get("docs").(string)), &docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal docs JSON")
	}

	for i, raw := range docs {
		doc := raw.(map[string]interface{})

		if _, ok := doc["_id"]; !ok {
			if id, ok := doc["id"]; !ok {
				return AppendDiagnostic(diags, fmt.Errorf("doc no %d missing _id or id field for change tracking", i), "_id field required on each document")
			} else {
				doc["_id"] = id
			}
		}
		if d.Get("partition").(string) != "" {
			doc["_id"] = fmt.Sprintf("%s:%s", d.Get("partition").(string), doc["_id"])
		}

		id := doc["_id"].(string)
		doc["_rev"] = revisions[id]
		docs[i] = doc
	}

	body := &models.Body3{
		Docs: docs,
	}

	params := database.NewBulkDocsParams().WithDb(dbName).WithBody(body)
	created, err := client.Database.BulkDocs(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk update documents")
	}

	for _, item := range created.Payload {
		if item.Ok {
			revisions[item.ID] = item.Rev
		} else {
			AppendDiagnostic(diags, fmt.Errorf(item.Error), item.Reason)
		}
	}

	byt, err := json.Marshal(revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.SetId(string(byt))

	return
}

func bulkDocumentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	var docs []models.Document
	for id, rev := range revisions {
		doc := map[string]interface{}{}
		doc["_id"] = id
		doc["_rev"] = rev
		doc["_deleted"] = true
		docs = append(docs, doc)
	}

	body := &models.Body3{
		Docs: docs,
	}

	params := database.NewBulkDocsParams().WithDb(dbName).WithBody(body)
	created, err := client.Database.BulkDocs(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk delete documents")
	}

	for _, item := range created.Payload {
		if item.Ok {
			delete(revisions, item.ID)
		}
	}

	byt, err := json.Marshal(revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.SetId(string(byt))

	return
}

func bulkDocumentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	docs := []*models.BasicDoc{}
	for id, rev := range revisions {
		ref := &models.BasicDoc{
			ID:  id,
			Rev: rev,
		}
		docs = append(docs, ref)
	}

	body := &models.Body2{
		Docs: docs,
	}

	params := database.NewBulkGetParams().WithDb(dbName).WithBody(body)
	created, err := client.Database.BulkGet(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk read documents")
	}

	revisions = map[string]string{}
	for _, item := range created.Payload.Results {
		// docs contains a single-item array of objects,
		// each of which has either an error key and value describing the error,
		//or ok key and associated value of the requested document
		if len(item.Docs) > 0 {
			if item.Docs[0].Ok != nil {
				doc := item.Docs[0].Ok.(map[string]interface{})
				revisions[item.ID] = doc["_rev"].(string)
			} else {
				AppendDiagnostic(diags, fmt.Errorf(item.Docs[0].Error.Error), fmt.Sprintf("%s id: %s, rev: %s", item.Docs[0].Error.Reason, item.Docs[0].Error.ID, item.Docs[0].Error.Rev))
			}
		}
	}

	byt, err := json.Marshal(revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.SetId(string(byt))

	return diags
}

func bulkDocumentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	var docs []models.Document
	err := json.Unmarshal([]byte(d.Get("docs").(string)), &docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	for i, raw := range docs {
		doc := raw.(map[string]interface{})
		if _, ok := doc["_id"]; !ok {
			if id, ok := doc["id"]; !ok {
				return AppendDiagnostic(diags, fmt.Errorf("doc no %d missing _id or id field for change tracking", i), "_id field required on each document")
			} else {
				doc["_id"] = id
			}
		}
		if d.Get("partition").(string) != "" {
			doc["_id"] = fmt.Sprintf("%s:%s", d.Get("partition").(string), doc["_id"])
		}
	}

	body := &models.Body3{
		Docs: docs,
	}

	params := database.NewBulkDocsParams().WithDb(dbName).WithBody(body)
	created, err := client.Database.BulkDocs(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk create documents")
	}

	revisions := map[string]string{}
	for _, item := range created.Payload {
		if item.Ok {
			revisions[item.ID] = item.Rev
		} else {

		}
	}

	byt, err := json.Marshal(revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.SetId(string(byt))

	return
}
