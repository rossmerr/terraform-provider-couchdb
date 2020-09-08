package couchdb

import (
	"context"

	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const usersDB = "_users"

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: UserCreate,
		Read:   UserRead,
		Update: UserUpdate,
		Delete: UserDelete,

		Schema: map[string]*schema.Schema{
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password",
			},
			"roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User roles",
			},
		},
	}
}

func UserCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	db := client.DB(context.Background(), usersDB)

	user := &tuser{
		ID:       kivik.UserPrefix + d.Id(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
	}

	_, err = db.Put(context.Background(), user.ID, user)
	if err != nil {
		return err
	}

	d.SetId(user.ID)
	return UserRead(d, meta)
}

func UserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	db := client.DB(context.Background(), usersDB)

	row := db.Get(context.Background(), d.Id())

	var user tuser
	if err = row.ScanDoc(&user); err != nil {
		return err
	}

	d.Set("revision", row.Rev)
	d.Set("roles", user.Roles)

	return nil
}

func UserUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	db := client.DB(context.Background(), usersDB)

	user := &tuser{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
		Revision: d.Get("revision").(string),
	}

	_, err = db.Put(context.Background(), user.ID, user)
	if err != nil {
		return err
	}

	return UserRead(d, meta)
}

func UserDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := connectToCouchDB(meta.(*CouchDBConfiguration))
	if err != nil {
		return err
	}

	db := client.DB(context.Background(), usersDB)

	_, err = db.Delete(context.Background(), d.Id(), d.Get("revision").(string))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

type tuser struct {
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
	Password string   `json:"password"`
	Revision string   `json:"_rev,omitempty"`
}
