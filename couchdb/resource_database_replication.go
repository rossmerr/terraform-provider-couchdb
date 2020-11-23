package couchdb

import (
	"context"
	"encoding/json"
	"github.com/rossmerr/couchdb_go/client/document"
	"github.com/rossmerr/couchdb_go/client/server"
	"github.com/rossmerr/couchdb_go/models"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const replicatorDB = "_replicator"


func resourceDatabaseReplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: databaseReplicationCreate,
		ReadContext:   databaseReplicationRead,
		DeleteContext: databaseReplicationDelete,
		UpdateContext: databaseReplicationUpdate,
		Schema: map[string]*schema.Schema{
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"create_target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Create target if it does not exist?",
			},
			"continuous": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep the replication permanently running?",
			},
			"source": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Source of the replication",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization": {
										Type:     schema.TypeString,
										Optional: true,
										Default: "",
									},
								},
							},
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"target": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Target of the replication",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization": {
										Type:     schema.TypeString,
										Optional: true,
										Default: "",
									},
								},
							},
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}



func databaseReplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	replicate := &models.Replicate{
		CreateTarget: d.Get("create_target").(bool),
		Continuous: d.Get("continuous").(bool),
		Filter:  d.Get("filter").(string),
		Source: extractRequest(d.Get("source")),
		Target:  extractRequest(d.Get("target")),
	}

	params := server.NewReplicationParams().WithBody(replicate)
	ok, accepted, err := client.Server.Replication(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if ok != nil && ok.Payload != nil && ok.Payload.Ok {
		d.SetId(ok.Payload.ID)
		d.Set("revision", ok.Payload.Rev)
	}

	if accepted != nil && accepted.Payload != nil && accepted.Payload.Ok {
		d.SetId(accepted.Payload.ID)
		d.Set("revision", accepted.Payload.Rev)
	}

	return databaseReplicationRead(ctx, d, meta)
}

func databaseReplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	rev := d.Get("revision").(string)
	params := document.NewDocGetParams().WithDb(replicatorDB).WithDocid(d.Id()).WithRev(&rev)
	reps, err := client.Document.DocGet(params)
	if err != nil {
		return diag.FromErr(err)
	}

	raw, err := json.Marshal(reps.Payload)

	var replicate models.Replicate
	err = json.Unmarshal(raw, &replicate)

	d.Set("continuous", replicate.Continuous)
	d.Set("create_target", replicate.CreateTarget)
	d.Set("revision", reps.ETag)

	return diags
}

func databaseReplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}
	rev := d.Get("revision").(string)
	params := document.NewDocDeleteParams().WithDb(replicatorDB).WithRev(&rev).WithDocid(d.Id())
	ok, accepted, err := client.Document.DocDelete(params)
	if err != nil {
		AppendDiagnostic(diags, err, "Unable to delete Document")
	}

	d.SetId("")

	if ok != nil {
		d.Set("revision", ok.ETag)
	}

	if accepted != nil {
		d.Set("revision", accepted.ETag)
	}

	return diags
}

func databaseReplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := databaseReplicationDelete(ctx, d, meta)
	if err != nil {
		return err
	}
	return databaseReplicationCreate(ctx, d, meta)
}


func extractRequest(v interface{}) *models.Request {
	vs := v.([]interface{})
	if len(vs) != 1 {
		return nil
	}
	vi := vs[0].(map[string]interface{})

	request := &models.Request{
		URL: strfmt.URI(vi["url"].(string)),
		Headers: extractHeaders(vi["headers"]),
	}

	return request
}

func extractHeaders(v interface{}) *models.RequestHeaders {
	vs := v.([]interface{})
	if len(vs) != 1 {
		return nil
	}
	vi := vs[0].(map[string]interface{})
	headers := &models.RequestHeaders{
		Authorization: vi["authorization"].(string),
	}

	return headers
}