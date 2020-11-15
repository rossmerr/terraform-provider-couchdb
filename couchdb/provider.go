package couchdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/rossmerr/couchdb_go/client/server"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/rossmerr/couchdb_go/client"
)

type CouchDBConfiguration struct {
	Endpoint string
	Scheme   string
	Username string
	Password string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHDB_ENDPOINT", nil),
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if value == "" {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Validate endpoint",
							Detail:   "Endpoint must not be an empty string",
						})
					}

					return diags
				},
			},
			"scheme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
				ForceNew: true,
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHDB_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHDB_PASSWORD", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"couchdb_database":                 resourceDatabase(),
			"couchdb_database_replication":     resourceDatabaseReplication(),
			"couchdb_user":                     resourceUser(),
			"couchdb_document":                 resourceDocument(),
			"couchdb_bulk_documents":           resourceBulkDocuments(),
			"couchdb_database_design_document": resourceDesignDocument(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	return &CouchDBConfiguration{
		Endpoint: d.Get("endpoint").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Scheme:   d.Get("scheme").(string),
	}, diags
}

func connectToCouchDB(ctx context.Context, conf *CouchDBConfiguration) (*apiclient.CouchdbGo, *diag.Diagnostic) {
	var client *apiclient.CouchdbGo
	var err error
	// When provisioning a database server there can often be a lag between
	// when Terraform thinks it's available and when it is actually available.
	// This is particularly acute when provisioning a server and then immediately
	// trying to provision a database on it.
	retryError := resource.RetryContext(ctx, time.Minute, func() *resource.RetryError {
		transport := httptransport.New(conf.Endpoint, "", []string{conf.Scheme})
		transport.DefaultAuthentication = httptransport.BasicAuth(conf.Username, conf.Password)
		client = apiclient.New(transport, strfmt.Default)

		_, err = client.Server.Up(server.NewUpParams())
		if err != nil {
			return resource.RetryableError(err)
		}

		return nil
	})

	if retryError != nil {
		return client, &diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to connect to Server",
			Detail:   fmt.Sprintf("Could not connect to server: %s", retryError),
		}
	}

	return client, nil
}
