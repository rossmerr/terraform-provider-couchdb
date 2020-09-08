package couchdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type CouchDBConfiguration struct {
	Endpoint        string
	Username        string
	Password        string
	MaxConnLifetime time.Duration
	MaxOpenConns    int
}


func Provider() terraform.ResourceProvider {
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
			"couchdb_user":                     resourceUser(),
			"couchdb_database_design_document": resourceDesignDocument(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &CouchDBConfiguration{
		Endpoint:        d.Get("endpoint").(string),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
		MaxConnLifetime: time.Duration(d.Get("max_conn_lifetime_sec").(int)) * time.Second,
		MaxOpenConns:    d.Get("max_open_conns").(int),
	}, nil
}

func connectToCouchDB(conf *CouchDBConfiguration) (*kivik.Client, error) {

	var client *kivik.Client
	var err error

	// When provisioning a database server there can often be a lag between
	// when Terraform thinks it's available and when it is actually available.
	// This is particularly acute when provisioning a server and then immediately
	// trying to provision a database on it.
	retryError := resource.Retry(5*time.Minute, func() *resource.RetryError {

		client, err = kivik.New("couch", conf.Endpoint)
		if err != nil {
			return resource.RetryableError(err)
		}

		_, err = client.Ping(context.Background())
		if err != nil {
			return resource.RetryableError(err)
		}

		err := client.Authenticate(context.Background(), couchdb.BasicAuth(conf.Username, conf.Password))
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
