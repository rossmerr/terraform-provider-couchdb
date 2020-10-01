package couchdb

import (
	"context"
	"encoding/json"
	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	docId := d.Id()
	doc := map[string]interface{}{}
	err := json.Unmarshal([]byte(d.Get("doc").(string)), &doc)

	if err != nil {
		AppendDiagnostic(diags, err,  "Unable to unmarshal JSON")
	}

	options := kivik.Options{}
	if d.Get("batch").(bool) {
		options["batch"] = true
	}

	doc["_rev"] = d.Get("revision").(string)

	rev, err := db.Put(ctx, docId, doc, options)
	if err != nil {
		AppendDiagnostic(diags, err,  "Unable to update Document")
	}
	d.Set("revision", rev)

	return documentRead(ctx, d, meta)
}


func documentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	rev := d.Get("revision").(string)

	rev, err := db.Delete(ctx, d.Id(), rev)

	if err != nil {
		AppendDiagnostic(diags, err,  "Unable to delete Document")
	}

	d.SetId("")
	d.Set("revision", rev)

	return diags
}

func documentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	row := db.Get(ctx, d.Id())

	var doc  map[string]interface{}

	if err := row.ScanDoc(&doc); err != nil {
		AppendDiagnostic(diags, err,  "Unable to read Document")
	}

	d.Set("revision", row.Rev)

	raw, err := json.Marshal(doc)

	if err != nil {
		AppendDiagnostic(diags, err,  "Unable to marshal JSON")
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

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	docId := d.Get("docid").(string)
	if docId == "" {
		docId = uuid.New().String()
	}

	doc := map[string]interface{}{}
	err := json.Unmarshal([]byte(d.Get("doc").(string)), &doc)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	options := kivik.Options{}
	if d.Get("batch").(bool) {
		options["batch"] = true
	}

	rev, err := db.Put(ctx, docId, doc, options)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create to Document")
	}

	d.Set("revision", rev)
	d.SetId(docId)

	return documentRead(ctx, d, meta)
}