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
		CreateContext: DocumentCreate,
		ReadContext:   DocumentRead,
		UpdateContext: DocumentUpdate,
		DeleteContext: DocumentDelete,

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

func DocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to Server",
			Detail:   err.Error(),
		})
		return diags
	}

	dbName := d.Get("database").(string)

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	docId := d.Id()
	doc := map[string]interface{}{}
	err = json.Unmarshal([]byte(d.Get("doc").(string)), &doc)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal JSON",
			Detail:   err.Error(),
		})
		return diags
	}

	options := kivik.Options{}
	if d.Get("batch").(bool) {
		options["batch"] = true
	}

	rev, err := db.Put(ctx, docId, doc, options)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Document",
			Detail:   err.Error(),
		})
		return diags
	}
	d.Set("revision", rev)

	return DocumentRead(ctx, d, meta)
}


func DocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to Server",
			Detail:   err.Error(),
		})
		return diags
	}

	dbName := d.Get("database").(string)

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}


	rev := d.Get("revision").(string)

	rev, err = db.Delete(ctx, d.Id(), rev)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Document",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("revision", rev)

	return diags
}

func DocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to Server",
			Detail:   err.Error(),
		})
		return diags
	}

	dbName := d.Get("database").(string)

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	row := db.Get(ctx, d.Id())

	var doc  map[string]interface{}

	if err := row.ScanDoc(&doc); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Document",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("revision", row.Rev)

	raw, err := json.Marshal(doc)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal JSON",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("doc", raw)

	return diags
}

func DocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to Server",
			Detail:   err.Error(),
		})
		return diags
	}

	dbName := d.Get("database").(string)

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	docId := d.Get("docid").(string)
	if docId == "" {
		docId = uuid.New().String()
	}

	doc := map[string]interface{}{}
	err = json.Unmarshal([]byte(d.Get("doc").(string)), &doc)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal JSON",
			Detail:   err.Error(),
		})
		return diags
	}

	options := kivik.Options{}
	if d.Get("batch").(bool) {
		options["batch"] = true
	}

	rev, err := db.Put(ctx, docId, doc, options)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create to Document",
			Detail:   err.Error(),
		})
		return diags
	}
	d.Set("revision", rev)
	d.SetId(docId)

	return DocumentRead(ctx, d, meta)
}