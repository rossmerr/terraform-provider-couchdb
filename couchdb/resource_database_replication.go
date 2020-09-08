package couchdb

import (
	"context"

	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabaseReplication() *schema.Resource {
	return &schema.Resource{
		Create: DatabaseReplicationCreate,
		Read:   DatabaseReplicationRead,
		Delete: DatabaseReplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the replication document",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source of the replication",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target of the replication",
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
			"replication_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal replication ID",
			},
			"replication_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "currennt replication state",
			},
			"replication_state_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current replication state transition reason",
			},
		},
	}
}

func DatabaseReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	options := kivik.Options{
		"id": d.Get("name").(string),
		"create_target" : d.Get("create_target").(bool),
		"continuous": d.Get("continuous").(bool),
		"filter": d.Get("filter").(string),
	}

	rep, err := client.Replicate(context.Background(), d.Get("target").(string), d.Get("source").(string), options)
	if err != nil {
		return err
	}
	d.SetId(rep.ReplicationID())

	return DatabaseReplicationRead(d, meta)
}

func DatabaseReplicationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	reps, err := client.GetReplications(context.Background())
	if err != nil {
		return err
	}

	for _, rep := range reps {
		if rep.ReplicationID() == d.Id() {
			d.Set("source", rep.Source)
			d.Set("target", rep.Target)
			d.Set("replication_id", rep.ReplicationID)
			d.Set("replication_state", rep.State())
			d.Set("replication_state_reason", rep.Source)
			break
		}
	}

	return nil
}


func DatabaseReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	reps, err := client.GetReplications(context.Background())
	if err != nil {
		return err
	}

	for _, rep := range reps {
		if rep.ReplicationID() == d.Id() {
			err = rep.Delete(context.Background())
			if err != nil {
				return err
			}
		}
	}

	d.SetId("")
	return nil
}