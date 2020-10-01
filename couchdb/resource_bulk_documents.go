package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
}

func bulkDocumentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	var docs []interface{}
	err = json.Unmarshal([]byte(d.Get("docs").(string)), &docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal docs JSON")
	}

	for i, raw := range docs {
		doc := raw.(map[string]interface{})
		id := doc["_id"].(string)
		doc["_rev"] = revisions[id]
		docs[i] = doc
	}

	row, err := db.BulkDocs(ctx, docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk update documents")
	}
	defer row.Close()

	if ok := row.Next(); ok {
		if row.UpdateErr() != nil {
			AppendDiagnostic(diags, row.UpdateErr(), "Unable to update/read Document")
		} else {
			revisions[row.ID()] = row.Rev()
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

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	var docs []interface{}
	for id, rev := range revisions {
		doc := map[string]interface{}{}
		doc["_id"] = id
		doc["_rev"] = rev
		doc["_deleted"] = true
		docs = append(docs, doc)
	}

	row, err := db.BulkDocs(ctx, docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk delete documents")
	}
	defer row.Close()

	if ok := row.Next(); ok {
		if row.UpdateErr() != nil {
			AppendDiagnosticWarning(diags, row.UpdateErr(), fmt.Sprintf("Unable to delete document: %s", row.ID()))
		}

		delete(revisions, row.ID())
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

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	var bulkRev []kivik.BulkGetReference

	//revisions := strings.Split(d.Get("revisions").(string), ",")
	//ids := strings.Split(d.Id(), ",")

	var revisions map[string]string
	err := json.Unmarshal([]byte(d.Id()), &revisions)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	for id, _ := range revisions {
		ref := kivik.BulkGetReference{
			ID: id,
		}
		bulkRev = append(bulkRev, ref)
	}

	rows, err := db.BulkGet(ctx, bulkRev)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read documents")
	}
	defer rows.Close()

	if rows.Err() != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  rows.Err().Error(),
		})
	}

	revisions = map[string]string{}
	if ok := rows.Next(); ok {

		var raw map[string]interface{}
		if err := rows.ScanDoc(&raw); err == nil {
			revisions[rows.ID()] = raw["_rev"].(string)
		} else {
			return AppendDiagnostic(diags, err, fmt.Sprintf("Unable to read document: %s", rows.ID()))
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

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	var docs []interface{}
	err := json.Unmarshal([]byte(d.Get("docs").(string)), &docs)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	for i, raw := range docs {
		doc := raw.(map[string]interface{})
		if _, ok := doc["_id"]; !ok {
			return AppendDiagnostic(diags, fmt.Errorf("doc no %d missing _id field for change tracking", i), "_id field required on each document")
		}
	}

	row, err := db.BulkDocs(ctx, docs)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to bulk create documents")
	}
	defer row.Close()

	revisions := map[string]string{}

	if ok := row.Next(); ok {
		if row.UpdateErr() != nil {
			AppendDiagnosticWarning(diags, row.UpdateErr(), fmt.Sprintf("Unable to read document: %s", row.ID()))
		} else {
			revisions[row.ID()] = row.Rev()
		}
	}

	byt, err := json.Marshal(revisions)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.SetId(string(byt))

	return
}
