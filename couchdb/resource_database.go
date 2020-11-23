package couchdb

import (
	"context"
	"fmt"
	"github.com/rossmerr/couchdb_go/client/database"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDatabase,
		UpdateContext: updateDatabase,
		ReadContext:   readDatabase,
		DeleteContext: deleteDatabase,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the database",
			},
			"partitioned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Whether to create a partitioned database",
			},
			"clustering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "database clustering configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "Number of replicas",
						},
						"shards": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     8,
							Description: "Number of shards",
						},
					},
				},
			},
			"document_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of documents in database",
			},
			"document_deletion_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tombstones in database",
			},
			"sizes_active": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of live data inside the database, in bytes.",
			},
			"sizes_external": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The uncompressed size of database contents in bytes.",
			},
			"size_file": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the database file on disk in bytes. Views indexes are not included in the calculation",
			},
		},
	}
}

func createDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("name").(string)
	shards, replicas := extractClusterOptions(d.Get("clustering"))
	partitioned := d.Get("partitioned").(bool)

	params := database.NewPutParams().WithDb(dbName).WithN(replicas).WithQ(shards).WithPartitioned(&partitioned)
	_, _, err := client.Database.Put(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create DB")
	}

	d.SetId(dbName)

	return readDatabase(ctx, d, meta)
}

func updateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	return readDatabase(ctx, d, meta)
}

func readDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Id()

	params := database.NewGetParams().WithDb(dbName)
	response, err := client.Database.Get(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read DB states")
	}

	if response != nil && response.Payload != nil {
		d.Set("document_count", int(response.Payload.DocCount))
		d.Set("document_deletion_count", int(response.Payload.DocDelCount))
		d.Set("sizes_active", int(response.Payload.Sizes.Active))
		d.Set("sizes_external", int(response.Payload.Sizes.External))
		d.Set("size_file", int(response.Payload.Sizes.File))
	}

	return diags
}

func deleteDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Id()

	params := database.NewDeleteParams().WithDb(dbName)
	_ ,_, err := client.Database.Delete(params)
	if err == nil {
		d.SetId("")
		return diags
	}

	return AppendDiagnostic(diags, fmt.Errorf("dbName: %s \n%s", dbName, err.Error()), "Unable to delete DB")
}

func extractClusterOptions(v interface{}) (*int32, *int32) {
	vs := v.([]interface{})
	if len(vs) != 1 {
		return nil, nil
	}
	vi := vs[0].(map[string]interface{})
	replicas := int32(vi["replicas"].(int))
	shards := int32(vi["shards"].(int))
	return &shards, &replicas
}

