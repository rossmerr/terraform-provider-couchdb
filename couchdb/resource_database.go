package couchdb

import (
	"context"
	"strconv"

	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateDatabase,
		UpdateContext: UpdateDatabase,
		ReadContext:   ReadDatabase,
		DeleteContext: DeleteDatabase,
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
				ForceNew: 	 true,
				Description: "Whether to create a partitioned database",
			},
			"security": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Security configuration of the database",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admins": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database administrators",
						},
						"admin_roles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database administration roles",
						},
						"members": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database members",
						},
						"member_roles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database member roles",
						},
					},
				},
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
			"disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of storage disk",
			},
			"data_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of database data",
			},
		},
	}
}

func CreateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	dbName := d.Get("name").(string)
	options := extractClusterOptions(d.Get("clustering"))
	options["partitioned"] = d.Get("partitioned").(bool)

	err = client.CreateDB(ctx, dbName, options)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create DB",
			Detail:   err.Error(),
		})
		return diags
	}

	if v, ok := d.GetOk("security"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			db := client.DB(ctx, dbName)
			err := db.SetSecurity(ctx, extractDatabaseSecurity(vs[0]))
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to set security on DB",
					Detail:   err.Error(),
				})
				return diags
			}
		}
	}

	d.SetId(dbName)

	return ReadDatabase(ctx, d, meta)
}

func UpdateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{})  diag.Diagnostics {
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
	dbName := d.Get("name").(string)

	if d.HasChange("security") {
		db, dd := connectToDB(ctx, client, dbName)
		if dd != nil {
			diags = append(diags, *dd)
			return diags
		}

		if v, ok := d.GetOk("security"); ok {
			vs := v.([]interface{})
			if len(vs) == 1 {
				err := db.SetSecurity(ctx, extractDatabaseSecurity(vs[0]))
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Unable to set security on DB",
						Detail:   err.Error(),
					})
					return diags
				}
			}
		} else {
			err := db.SetSecurity(ctx, extractDatabaseSecurity(nil))
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to set security on DB",
					Detail:   err.Error(),
				})
				return diags
			}
		}
	}


	return ReadDatabase(ctx, d, meta)
}

func ReadDatabase(ctx context.Context, d *schema.ResourceData, meta interface{})  diag.Diagnostics {
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

	dbName := d.Id()
	dbStates, err := client.DBsStats(ctx, []string{dbName} )
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read DB states",
			Detail:   err.Error(),
		})
		return diags
	}

	if len(dbStates) > 0 {
		state := dbStates[0]
		d.Set("document_count",  strconv.FormatInt(state.DocCount, 16))
		d.Set("document_deletion_count", strconv.FormatInt(state.DeletedCount, 16))
		d.Set("disk_size", strconv.FormatInt(state.DiskSize, 16))
		d.Set("data_size", strconv.FormatInt(state.ActiveSize, 16))
	}

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	sec, err := db.Security(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read security on DB",
			Detail:   err.Error(),
		})
		return diags
	}

	security := []map[string][]string{
		{
			"admins":       sec.Admins.Names,
			"admin_roles":  sec.Admins.Roles,
			"members":      sec.Members.Names,
			"member_roles": sec.Members.Roles,
		},
	}
	d.Set("security", security)

	return diags
}

func DeleteDatabase(ctx context.Context, d *schema.ResourceData, meta interface{})  diag.Diagnostics {
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

	dbName := d.Id()
	err = client.DestroyDB(ctx, dbName)
	if err == nil {
		d.SetId("")
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to delete DB",
		Detail:   err.Error(),
	})

	return diags
}

func extractClusterOptions(v interface{}) (kivik.Options) {
	ret := kivik.Options{}
	vs := v.([]interface{})
	if len(vs) != 1 {
		return ret
	}
	vi := vs[0].(map[string]interface{})
	ret["replicas"]= vi["replicas"].(int)
	ret["shards"] = vi["shards"].(int)
	return ret
}

func extractDatabaseSecurity(d interface{}) *kivik.Security {
	security, ok := d.(map[string]interface{})
	if !ok {
		return &kivik.Security{}
	}

	return &kivik.Security{
		Admins: kivik.Members{
			Names: stringsFromSet(security["admins"]),
			Roles: stringsFromSet(security["admin_roles"]),
		},
		Members: kivik.Members{
			Names: stringsFromSet(security["members"]),
			Roles: stringsFromSet(security["member_roles"]),
		},
	}
}