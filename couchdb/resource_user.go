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
		CreateContext: UserCreate,
		ReadContext:   UserRead,
		UpdateContext: UserUpdate,
		DeleteContext: UserDelete,

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

func UserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	user := &tuser{
		ID:       kivik.UserPrefix + uuid.New().String(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
	}

	_, err = db.Put(ctx, user.ID, user)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(user.ID)
	return UserRead(ctx, d, meta)
}

func UserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}


	row := db.Get(ctx, d.Id())

	var user tuser
	if err := row.ScanDoc(&user); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("revision", row.Rev)
	d.Set("roles", user.Roles)

	return diags
}

func UserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	user := &tuser{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
		Revision: d.Get("revision").(string),
	}

	_, err = db.Put(ctx, user.ID, user)
	if err != nil {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update User",
			Detail:   err.Error(),
		})
		return diags
	}

	return UserRead(ctx, d, meta)
}

func UserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	db, dd := connectToDB(ctx, client, usersDB)
	if dd != nil {
		diags = append(diags, *dd)
		return diags
	}

	_, err = db.Delete(ctx, d.Id(), d.Get("revision").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unable to delete User",
			Detail:  fmt.Sprintf("docID: %s \nrev: %s \n%s", d.Id(), d.Get("revision").(string),  err.Error()),
		})
		return diags
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
