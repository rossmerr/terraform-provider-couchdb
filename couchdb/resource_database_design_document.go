package couchdb

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


func resourceDesignDocument() *schema.Resource {
	return &schema.Resource{
		Create: DesignDocumentCreate,
		Read:   DesignDocumentRead,
		Update: DesignDocumentUpdate,
		Delete: DesignDocumentDelete,

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


func DesignDocumentCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("database").(string)
	db := client.DB(context.Background(), dbName)

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

		log.Println("Executing DesignDocumentCreate:", docId, ddoc)
		rev, err := db.Put(context.Background(), docId, ddoc)
		if err != nil {
			return err
		}

		d.Set("revision", rev)
	}

	d.SetId(docId)

	return DesignDocumentRead(d, meta)
}

func DesignDocumentRead(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("database").(string)
	db := client.DB(context.Background(), dbName)

	docId := fmt.Sprintf("_design/%s", d.Get("name").(string))

	log.Println("Executing DesignDocumentRead:", docId)
	row := db.Get(context.Background(), docId)

	var ddoc map[string]map[string]interface{}
	if err = row.ScanDoc(&ddoc); err != nil {
		return err
	}

	d.Set("language", ddoc["language"])
	d.Set("view", ddoc["views"])
	d.Set("revision", row.Rev)

	return nil
}

func DesignDocumentUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("database").(string)
	db := client.DB(context.Background(), dbName)

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

		log.Println("Executing DesignDocumentUpdate:", d.Id(), ddoc)
		rev, err := db.Put(context.Background(), d.Id(), ddoc)
		if err != nil {
			return err
		}

		d.Set("revision", rev)
	}

	return nil
}

func DesignDocumentDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	dbName := d.Get("database").(string)
	db := client.DB(context.Background(), dbName)

	log.Println("Executing DesignDocumentDelete:", d.Id(), d.Get("revision").(string))
	_, err = db.Delete(context.Background(), d.Id(), d.Get("revision").(string))

	d.SetId("")
	return nil
}

