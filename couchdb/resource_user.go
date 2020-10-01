package couchdb

import (
	"context"
	"fmt"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const usersDB = "_users"

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: userCreate,
		ReadContext:   userRead,
		UpdateContext: userUpdate,
		DeleteContext: userDelete,

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

func userCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	user := &tuser{
		ID:       kivik.UserPrefix + uuid.New().String(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
	}

	_, err := db.Put(ctx, user.ID, user)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create User")
	}

	d.SetId(user.ID)
	return userRead(ctx, d, meta)
}

func userRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	row := db.Get(ctx, d.Id())

	var user tuser
	if err := row.ScanDoc(&user); err != nil {
		diags = AppendDiagnostic(diags, err, "Unable to read User")
		return diags
	}

	d.Set("revision", row.Rev)
	d.Set("roles", user.Roles)

	return diags
}

func userUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	user := &tuser{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
		Revision: d.Get("revision").(string),
	}

	_, err := db.Put(ctx, user.ID, user)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to update User")
	}

	return userRead(ctx, d, meta)
}

func userDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		return append(diags, *dd)
	}
	defer db.Close(ctx)

	_, err := db.Delete(ctx, d.Id(), d.Get("revision").(string))
	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("docID: %s \nrev: %s \n%s", d.Id(), d.Get("revision").(string), err.Error()), "Unable to delete User")
	}
	d.SetId("")

	return diags
}

type tuser struct {
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
	Password string   `json:"password"`
	Revision string   `json:"_rev,omitempty"`
}
