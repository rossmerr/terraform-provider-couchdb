package couchdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-kivik/kivik/v3"
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

func createDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("name").(string)
	options := extractClusterOptions(d.Get("clustering"))
	options["partitioned"] = d.Get("partitioned").(bool)

	err := client.CreateDB(ctx, dbName, options)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create DB")
	}

	d.SetId(dbName)

	// You can't edit the security object of the user database.
	if dbName == usersDB {
		return readDatabase(ctx, d, meta)
	}

	if v, ok := d.GetOk("security"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			db := client.DB(ctx, dbName)
			err := db.SetSecurity(ctx, extractDatabaseSecurity(vs[0]))
			if err != nil {
				return AppendDiagnostic(diags, err, "Unable to set security on DB")
			}
		}
	}

	return readDatabase(ctx, d, meta)
}

func updateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Get("name").(string)

	// You can't edit the security object of the user database.
	if dbName == usersDB {
		return readDatabase(ctx, d, meta)
	}

	if d.HasChange("security") {
		db, dd := connectToDB(ctx, client, dbName)
		if dd != nil {
			return append(diags, *dd)
		}
		defer db.Close(ctx)

		if v, ok := d.GetOk("security"); ok {
			vs := v.([]interface{})
			if len(vs) == 1 {
				err := db.SetSecurity(ctx, extractDatabaseSecurity(vs[0]))
				if err != nil {
					return AppendDiagnostic(diags, err, "Unable to update security on DB")
				}
			}
		} else {
			err := db.SetSecurity(ctx, &kivik.Security{})
			if err != nil {
				return AppendDiagnostic(diags, err, "Unable to clear security on DB")
			}
		}
	}

	return readDatabase(ctx, d, meta)
}

func readDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Id()
	dbStates, err := client.DBsStats(ctx, []string{dbName})
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read DB states")
	}

	if len(dbStates) > 0 {
		state := dbStates[0]
		d.Set("document_count", strconv.FormatInt(state.DocCount, 16))
		d.Set("document_deletion_count", strconv.FormatInt(state.DeletedCount, 16))
		d.Set("disk_size", strconv.FormatInt(state.DiskSize, 16))
		d.Set("data_size", strconv.FormatInt(state.ActiveSize, 16))
	}

	db, dd := connectToDB(ctx, client, dbName)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	// You can't edit the security object of the user database.
	if dbName == usersDB {
		return diags
	}

	sec, err := db.Security(ctx)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read security on DB")
	}

	if len(sec.Admins.Roles) > 0 ||
		len(sec.Admins.Names) > 0 ||
		len(sec.Members.Roles) > 0 ||
		len(sec.Members.Names) > 0 {
		security := []map[string][]string{
			{
				"admins":       sec.Admins.Names,
				"admin_roles":  sec.Admins.Roles,
				"members":      sec.Members.Names,
				"member_roles": sec.Members.Roles,
			},
		}
		d.Set("security", security)
	} else {
		d.Set("security", nil)
	}

	return diags
}

func deleteDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	dbName := d.Id()
	err := client.DestroyDB(ctx, dbName)
	if err == nil {
		d.SetId("")
		return diags
	}

	return AppendDiagnostic(diags, fmt.Errorf("dbName: %s \n%s", dbName, err.Error()), "Unable to delete DB")
}

func extractClusterOptions(v interface{}) kivik.Options {
	ret := kivik.Options{}
	vs := v.([]interface{})
	if len(vs) != 1 {
		return ret
	}
	vi := vs[0].(map[string]interface{})
	ret["replicas"] = vi["replicas"].(int)
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
