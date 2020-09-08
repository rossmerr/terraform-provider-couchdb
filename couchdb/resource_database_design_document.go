package couchdb

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


func resourceDesignDocument() *schema.Resource {
	return &schema.Resource{
		CreateContext: DesignDocumentCreate,
		ReadContext:   DesignDocumentRead,
		UpdateContext: DesignDocumentUpdate,
		DeleteContext: DesignDocumentDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database to associate design with",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the design document",
			},
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "javascript",
				Description: "Language of map/ reduce functions",
			},
			"view": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A view inside the design document",
				Set: func(v interface{}) int {
					view := v.(map[string]interface{})
					name := view["name"].(string)
					id := 0
					for _, b := range md5.Sum([]byte(name)) {
						id += int(b)
					}
					return id
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the view",
						},
						"map": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Map function",
						},
						"reduce": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reduce functionn",
						},
					},
				},
			},
		},
	}
}


func DesignDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("database").(string)
	db := client.DB(ctx, dbName)
	if db.Err() != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to DB",
			Detail:   db.Err().Error(),
		})
	}

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	if vs, ok := d.GetOk("view"); ok {

		i := map[string]interface{}{}
		ddoc := map[string]interface{}{
			"_id": docId,
			"views": i,
			"language": d.Get("language").(string),
		}

		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})

			i[view["name"].(string)] = map[string]interface{}{
				 "map": view["map"].(string),
				 "reduce": view["reduce"].(string),
			}
		}

		rev, err := db.Put(ctx, docId, ddoc)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create design doc",
				Detail:   err.Error(),
			})
			return diags
		}

		d.Set("revision", rev)
	}

	d.SetId(docId)

	return DesignDocumentRead(ctx, d, meta)
}

func DesignDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("database").(string)
	db := client.DB(ctx, dbName)
	if db.Err() != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to DB",
			Detail:   db.Err().Error(),
		})
	}
	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	row := db.Get(ctx, docId)

	var ddoc map[string]map[string]interface{}
	if err = row.ScanDoc(&ddoc); err != nil {
		return diag.FromErr(err)
	}

	d.Set("language", ddoc["language"])
	d.Set("view", ddoc["views"])
	d.Set("revision", row.Rev)

	return diags
}

func DesignDocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("database").(string)
	db := client.DB(ctx, dbName)
	if db.Err() != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to DB",
			Detail:   db.Err().Error(),
		})
	}

	if vs, ok := d.GetOk("view"); ok {

		i := map[string]interface{}{}
		ddoc := map[string]interface{}{
			"_id":  d.Id(),
			"_rev":  d.Get("revision").(string),
			"views": i,
			"language": d.Get("language").(string),
		}

		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})

			i[view["name"].(string)] = map[string]interface{}{
				"map": view["map"].(string),
				"reduce": view["reduce"].(string),
			}
		}

		rev, err := db.Put(ctx, d.Id(), ddoc)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update design doc",
				Detail:   err.Error(),
			})
			return diags
		}

		d.Set("revision", rev)
	}

	return diags
}

func DesignDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, err := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("database").(string)
	db := client.DB(ctx, dbName)
	if db.Err() != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to DB",
			Detail:   db.Err().Error(),
		})
	}
	_, err = db.Delete(ctx, d.Id(), d.Get("revision").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete design doc",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")
	return diags
}

