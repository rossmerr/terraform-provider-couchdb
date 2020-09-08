package couchdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CouchDBConfiguration struct {
	Endpoint        string
	Username        string
	Password        string
	MaxConnLifetime time.Duration
	MaxOpenConns    int
}


func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHDB_ENDPOINT", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Endpoint must not be an empty string"))
					}

					return
				},
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
			"couchdb_database_design_document": resourceDesignDocument(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return &CouchDBConfiguration{
		Endpoint:        d.Get("endpoint").(string),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
	}, diags
}

func connectToCouchDB(ctx context.Context, conf *CouchDBConfiguration) (*kivik.Client, error) {

	var client *kivik.Client
	var err error

	// When provisioning a database server there can often be a lag between
	// when Terraform thinks it's available and when it is actually available.
	// This is particularly acute when provisioning a server and then immediately
	// trying to provision a database on it.
	retryError := resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		client, err = kivik.New("couch", conf.Endpoint)
		if err != nil {
			return resource.RetryableError(err)
		}

		_, err = client.Ping(context.Background())
		if err != nil {
			return resource.RetryableError(err)
		}

		err := client.Authenticate(ctx, couchdb.BasicAuth(conf.Username, conf.Password))
		if err != nil {
			return resource.RetryableError(err)
		}

		return nil
	})

	if retryError != nil {
		return nil, fmt.Errorf("Could not connect to server: %s", retryError)
	}

	return client, nil
}
