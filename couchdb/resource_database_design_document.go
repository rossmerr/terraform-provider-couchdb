package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RossMerr/couchdb_go/client/design_documents"
	"github.com/RossMerr/couchdb_go/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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

	docId := fmt.Sprintf("%s", d.Get("name").(string))

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

	designDoc := &models.DesignDoc{
		Language:  d.Get("language").(string),
		Views: viewsDoc,
		Indexes: indexesDoc,
	}

	params := design_documents.NewDesignDocPutParams().WithDb(dbName).WithBody(designDoc).WithDdoc(docId)
	created, accepted, err := client.DesignDocuments.DesignDocPut(params)
	if err != nil {
		diags = AppendDiagnostic(diags, fmt.Errorf("designDoc"), fmt.Sprintf("%+v", designDoc))
		return AppendDiagnostic(diags, fmt.Errorf("%s \nDesign Doc:- \n%s", err.Error(), d.Get("view").(string)), "Unable to create design doc")
	}

	if created != nil {
		d.Set("revision", strings.Trim(created.ETag, "\""))
	}

	if accepted != nil {
		d.Set("revision", strings.Trim(accepted.ETag, "\""))
	}

	d.SetId(docId)

	return designDocumentRead(ctx, d, meta)
}

func designDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}
	dbName := d.Get("database").(string)
	docId := fmt.Sprintf("%s", d.Get("name").(string))
	rev := d.Get("revision").(string)

	params := design_documents.NewDesignDocGetParams().WithDb(dbName).WithDdoc(docId).WithRev(&rev)
	ok, err := client.DesignDocuments.DesignDocGet(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read Design Document")
	}

	if ok.Payload != nil {
		d.Set("language",  ok.Payload.Language)

		if ok.Payload.Views != nil {
			if data, err := json.Marshal(ok.Payload.Views); err == nil {
				d.Set("views", string(data))
			}
		}

		if ok.Payload.Indexes != nil {
			if data, err := json.Marshal(ok.Payload.Indexes); err == nil {
				d.Set("indexes", string(data))
			}
		}
	}

	d.Set("revision", strings.Trim(ok.ETag, "\""))
	return diags
}

func designDocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)

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


	designDoc := &models.DesignDoc{
		Language:  d.Get("language").(string),
		Views: viewsDoc,
		Indexes: indexesDoc,
	}

	params := design_documents.NewDesignDocPutParams().WithDb(dbName).WithBody(designDoc).WithDdoc(d.Id())
	created, accepted, err := client.DesignDocuments.DesignDocPut(params)

	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("%s \nDesign Doc:- \n%s", err.Error(), d.Get("view").(string)), "Unable to update design doc")
	}

	if created != nil {
		d.Set("revision", strings.Trim(created.ETag, "\""))
	}

	if accepted != nil {
		d.Set("revision", strings.Trim(accepted.ETag, "\""))
	}

	return diags
}

func designDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("database").(string)
	rev := d.Get("revision").(string)

	params := design_documents.NewDesignDocDeleteParams().WithDb(dbName).WithDdoc(d.Id()).WithRev(&rev)
	ok, accepted, err := client.DesignDocuments.DesignDocDelete(params)
	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("docID: %s \nrev: %s \n%s", d.Id(), d.Get("revision").(string), err.Error()), "Unable to delete design doc")
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
