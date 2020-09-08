package couchdb

import (
	"context"
	"log"
	"strconv"

	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabase,
		Update: UpdateDatabase,
		Read:   ReadDatabase,
		Delete: DeleteDatabase,
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

func CreateDatabase(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("name").(string)
	log.Println("Executing CreateDB:", dbName)

	err = client.CreateDB(context.Background(), dbName, extractClusterOptions(d.Get("clustering")))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("security"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			db := client.DB(context.Background(), dbName)
			err := db.SetSecurity(context.Background(), extractDatabaseSecurity(vs[0]))
			if err != nil {
				return err
			}
		}
	}

	d.SetId(dbName)

	return ReadDatabase(d, meta)
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("name").(string)

	if d.HasChange("security") {
		log.Println("Executing SetSecurity on:", dbName)
		db := client.DB(context.Background(), dbName)
		if v, ok := d.GetOk("security"); ok {
			vs := v.([]interface{})
			if len(vs) == 1 {
				err := db.SetSecurity(context.Background(), extractDatabaseSecurity(vs[0]))
				if err != nil {
					return err
				}
			}
		} else {
			err := db.SetSecurity(context.Background(), extractDatabaseSecurity(nil))
			if err != nil {
				return err
			}
		}
	}


	return ReadDatabase(d, meta)
}

func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Id()
	log.Println("Executing DBsStats on:", dbName)
	dbStates, err := client.DBsStats(context.Background(), []string{dbName} )
	if err != nil {
		return err
	}

	if len(dbStates) > 0 {
		state := dbStates[0]
		d.Set("document_count",  strconv.FormatInt(state.DocCount, 16))
		d.Set("document_deletion_count", strconv.FormatInt(state.DeletedCount, 16))
		d.Set("disk_size", strconv.FormatInt(state.DiskSize, 16))
		d.Set("data_size", strconv.FormatInt(state.ActiveSize, 16))
	}

	log.Println("Executing Security on:", dbName)
	db := client.DB(context.Background(), dbName)
	if err != nil {
		return err
	}

	sec, err := db.Security(context.Background())
	if err != nil {
		return err
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

	return nil
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Id()
	log.Println("Executing DestroyDB:", dbName)
	err = client.DestroyDB(context.Background(), dbName)
	if err == nil {
		d.SetId("")
	}

	return err
}

func extractClusterOptions(v interface{}) (ret kivik.Options) {
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