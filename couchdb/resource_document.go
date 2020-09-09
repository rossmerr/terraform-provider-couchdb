package couchdb

import (
	"context"
	"strconv"
	"time"

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
		},
	}
}

func DocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}


func DocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
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

	docId := strconv.FormatInt(time.Now().UnixNano(),16)
	doc := map[string]interface{}{ "_id": docId,  "name": "test"}

	rev, err := db.Put(ctx, docId, doc)
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