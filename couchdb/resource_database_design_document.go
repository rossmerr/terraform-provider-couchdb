package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDesignDocument() *schema.Resource {
	return &schema.Resource{
		CreateContext: designDocumentCreate,
		ReadContext:   designDocumentRead,
		UpdateContext: designDocumentUpdate,
		DeleteContext: designDocumentDelete,

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
			"views": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The views inside the design document (wrap in <<EOF { } EOF)",
			},
			"indexes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The indexes inside the design document (wrap in <<EOF { } EOF)",
			},
		},
	}
}

func designDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	viewsDoc := map[string]interface{}{}

	if interfaceViews, ok := d.GetOk("views"); ok {
		if err := json.Unmarshal([]byte(interfaceViews.(string)), &viewsDoc); err != nil {
			return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
		}

	}

	indexesDoc := map[string]interface{}{}
	if interfaceIndexes, ok := d.GetOk("indexes"); ok {
		if err := json.Unmarshal([]byte(interfaceIndexes.(string)), &indexesDoc); err != nil {
			return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
		}
	}


	doc := map[string]interface{}{}
	doc["views"] = viewsDoc
	doc["indexes"] = indexesDoc
	doc["_id"] = docId
	doc["language"] = d.Get("language").(string)

	rev, err := db.Put(ctx, docId, doc)
	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("%s \nDesign Doc:- \n%s", err.Error(), d.Get("view").(string)), "Unable to create design doc")
	}

	d.Set("revision", rev)

	d.SetId(docId)

	return designDocumentRead(ctx, d, meta)
}

func designDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	row := db.Get(ctx, docId)

	var designDoc tdesignDoc
	if err := row.ScanDoc(&designDoc); err != nil {
		return diag.FromErr(err)
	}

	d.Set("language", designDoc.Language)

	if designDoc.Views != nil {
		b, err := json.Marshal(designDoc.Views)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("views", string(b))
	}

	if designDoc.Indexes != nil {
		b, err := json.Marshal(designDoc.Indexes)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("indexes", string(b))
	}

	d.Set("revision", row.Rev)



	return diags
}

func designDocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	viewsDoc := map[string]interface{}{}
	err := json.Unmarshal([]byte(d.Get("views").(string)), &viewsDoc)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	indexesDoc := map[string]interface{}{}
	err = json.Unmarshal([]byte(d.Get("indexes").(string)), &indexesDoc)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to unmarshal JSON")
	}

	doc := map[string]interface{}{}
	doc["_rev"] = d.Get("revision").(string)
	doc["views"] = viewsDoc
	doc["indexes"] = indexesDoc
	doc["_id"] = d.Id()
	doc["language"] = d.Get("language").(string)

	rev, err := db.Put(ctx, d.Id(), doc)
	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("%s \nDesign Doc:- \n%s", err.Error(), d.Get("view").(string)), "Unable to update design doc")
	}

	d.Set("revision", rev)

	return diags
}

func designDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
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

	rev, err := db.Delete(ctx, d.Id(), d.Get("revision").(string))

	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("docID: %s \nrev: %s \n%s", d.Id(), d.Get("revision").(string), err.Error()), "Unable to delete design doc")
	}

	d.SetId("")
	d.Set("revision", rev)

	return diags
}

type tdesignDoc struct {
	ID       string            `json:"_id"`
	Rev      string            `json:"_rev,omitempty"`
	Views    map[string]interface{}  `json:"views"`
	Indexes  map[string]interface{} `json:"indexes"`
	Language string            `json:"language"`
}
