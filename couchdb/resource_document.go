package couchdb

import (
	"context"
	"encoding/json"
	"github.com/RossMerr/couchdb_go/client/document"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceDocument() *schema.Resource {
	return &schema.Resource{
		CreateContext: documentCreate,
		ReadContext:   documentRead,
		UpdateContext: documentUpdate,
		DeleteContext: documentDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database to associate design with",
			},
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"docid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Document ID",
			},
			"doc": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON Document (wrap in <<EOF { } EOF)",
			},
			"batch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: " Stores document in batch mode, Batch mode is not suitable for critical data",
			},
		},
	}
}

func documentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	docId := d.Id()
	doc := map[string]interface{}{}
	err := json.Unmarshal([]byte(d.Get("doc").(string)), &doc)

	if err != nil {
		AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	rev := d.Get("revision").(string)

	batch := ""
	if d.Get("batch").(bool) {
		batch = "ok"
	}

	params := document.NewDocPutParams().WithDb(dbName).WithDocid(docId).WithRev(&rev).WithBatch(&batch).WithBody(doc)
	created, accepted, err := client.Document.DocPut(params)
	if err != nil {
		AppendDiagnostic(diags, err, "Unable to update Document")
	}

	if created != nil {
		d.Set("revision", strings.Trim(created.ETag, "\""))
	}

	if accepted != nil {
		d.Set("revision", strings.Trim(accepted.ETag, "\""))
	}

	return documentRead(ctx, d, meta)
}

func documentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)
	rev := d.Get("revision").(string)

	params := document.NewDocDeleteParams().WithDb(dbName).WithRev(&rev).WithDocid(d.Id())
	ok, accepted, err := client.Document.DocDelete(params)
	if err != nil {
		AppendDiagnostic(diags, err, "Unable to delete Document")
	}

	d.SetId("")

	if ok != nil {
		d.Set("revision", strings.Trim(ok.ETag, "\""))
	}

	if accepted != nil {
		d.Set("revision", strings.Trim(accepted.ETag, "\""))
	}

	return diags
}

func documentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

	params := document.NewDocGetParams().WithDb(dbName).WithDocid(d.Id())
	ok, err := client.Document.DocGet(params)
	if err != nil {
		AppendDiagnostic(diags, err, "Unable to read Document")
	}

	doc := ok.Payload.(map[string]interface{})

	d.Set("revision", strings.Trim(ok.ETag, "\""))

	raw, err := json.Marshal(doc)

	if err != nil {
		AppendDiagnostic(diags, err, "Unable to marshal JSON")
	}

	d.Set("doc", raw)

	return diags
}

func documentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)
	docId := d.Get("docid").(string)
	if docId == "" {
		docId = uuid.New().String()
	}

	batch := ""
	if d.Get("batch").(bool) {
		batch = "ok"
	}

	doc := map[string]interface{}{}
	err := json.Unmarshal([]byte(d.Get("doc").(string)), &doc)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	params := document.NewDocPutParams().WithDb(dbName).WithBody(doc).WithBatch(&batch).WithDocid(docId)
	created, accepted, err := client.Document.DocPut(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create to Document")
	}

	d.SetId(docId)

	if created != nil && created.Payload.Ok {
		d.Set("revision", strings.Trim(created.ETag, "\""))
	}

	if accepted != nil && accepted.Payload.Ok {
		d.Set("revision", strings.Trim(accepted.ETag, "\""))
	}

	return documentRead(ctx, d, meta)
}
