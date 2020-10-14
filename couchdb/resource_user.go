package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RossMerr/couchdb_go/client/document"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const usersDB = "_users"
const userPrefix = "org.couchdb.user:"

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


	user := &tuser{
		ID:       userPrefix + uuid.New().String(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
	}

	params := document.NewPostParams().WithDb(usersDB).WithBody(user)
	created, accepted, err := client.Document.Post(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to create to Document")
	}

	if created != nil && created.Payload.Ok {
		d.Set("revision", created.Payload.Rev)
	}

	if accepted != nil && accepted.Payload.Ok {
		d.Set("revision", accepted.Payload.Rev)
	}

	d.SetId(user.ID)
	return userRead(ctx, d, meta)
}

func userRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}


	params := document.NewDocGetParams().WithDb(usersDB).WithDocid(d.Id())
	ok, err := client.Document.DocGet(params)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read Document")
	}

	doc := ok.Payload.(map[string]interface{})

	d.Set("revision", ok.ETag)

	raw, err := json.Marshal(doc)

	var user tuser
	err = json.Unmarshal(raw, &user)
	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read User")
	}

	d.Set("revision", ok.ETag)
	d.Set("roles", user.Roles)

	return diags
}

func userUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	user := &tuser{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Type:     "user",
		Roles:    stringsFromSet(d.Get("roles")),
		Password: d.Get("password").(string),
		Revision: d.Get("revision").(string),
	}

	rev := d.Get("revision").(string)
	params := document.NewDocPutParams().WithDb(usersDB).WithDocid(d.Id()).WithRev(&rev).WithBody(user)
	created, accepted, err := client.Document.DocPut(params)
	if err != nil {
		AppendDiagnostic(diags, err, "Unable to update Document")
	}

	if created != nil {
		d.Set("revision", created.ETag)
	}

	if accepted != nil {
		d.Set("revision", accepted.ETag)
	}

	return userRead(ctx, d, meta)
}

func userDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	rev := d.Get("revision").(string)

	params := document.NewDocDeleteParams().WithDb(usersDB).WithRev(&rev).WithDocid(d.Id())
	ok, accepted, err := client.Document.DocDelete(params)
	if err != nil {
		return AppendDiagnostic(diags, fmt.Errorf("docID: %s \nrev: %s \n%s", d.Id(), d.Get("revision").(string), err.Error()), "Unable to delete User")
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

type tuser struct {
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
	Password string   `json:"password"`
	Revision string   `json:"_rev,omitempty"`
}
