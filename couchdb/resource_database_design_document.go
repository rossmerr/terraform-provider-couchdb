package couchdb

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

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
					_map := view["map"].(string)
					reduce := view["reduce"].(string)
					id := 0
					for _, b := range md5.Sum([]byte(name + _map + reduce)) {
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

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	if vs, ok := d.GetOk("view"); ok {

		designDoc := tdesignDoc{
			ID:       docId,
			Language: d.Get("language").(string),
			View:     Tview{},
		}

		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})
			designDoc.View.Name = view["name"].(string)
			designDoc.View.Map = strings.ReplaceAll(view["map"].(string), "\n", "")
			designDoc.View.Reduce = strings.ReplaceAll(view["reduce"].(string), "\n", "")
		}

		rev, err := db.Put(ctx, docId, designDoc)
		if err != nil {

			body := ""
			if b, err := json.Marshal(designDoc); err == nil {
				body = string(b)
			} else {
				body = err.Error()
			}

			url := fmt.Sprintf("%s/%s", dbName, docId)

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create design doc",
				Detail:   fmt.Sprintf("%s \nUrl: %s \nDesign Doc:- \n%s", err.Error(), url, body),
			})
			return diags
		}

		d.Set("revision", rev)
	}

	d.SetId(docId)

	return DesignDocumentRead(ctx, d, meta)
}

func DesignDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	row := db.Get(ctx, docId)

	var designDoc tdesignDoc
	if err := row.ScanDoc(&designDoc); err != nil {
		return diag.FromErr(err)
	}

	d.Set("language", designDoc.Language)

	view := []map[string]string{}
	v := map[string]string{
		"name":   designDoc.View.Name,
		"map":    designDoc.View.Map,
		"reduce": designDoc.View.Reduce,
	}
	view = append(view, v)

	d.Set("view", view)
	d.Set("revision", row.Rev)

	return diags
}

func DesignDocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if vs, ok := d.GetOk("view"); ok {
		designDoc := tdesignDoc{
			ID:       d.Id(),
			Language: d.Get("language").(string),
			View:     Tview{},
			Rev:      d.Get("revision").(string),
		}

		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})
			designDoc.View.Name = view["name"].(string)
			designDoc.View.Map = strings.ReplaceAll(view["map"].(string), "\n", "")
			designDoc.View.Reduce = strings.ReplaceAll(view["reduce"].(string), "\n", "")
		}

		rev, err := db.Put(ctx, d.Id(), designDoc)
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

type tdesignDoc struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev,omitempty"`
	View     Tview  `json:"view"`
	Language string `json:"language"`
}

type Tview struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
	Name   string `json:"name"`
}
