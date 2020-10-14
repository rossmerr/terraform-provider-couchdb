package couchdb

import (
	"context"
	apiclient "github.com/RossMerr/couchdb_go/client"
	"github.com/RossMerr/couchdb_go/client/server"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testProviderFactories map[string]func() (*schema.Provider, error)

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()

	testProviderFactories = map[string]func() (*schema.Provider, error){
		"couchdb": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {

	for _, name := range []string{"COUCHDB_ENDPOINT", "COUCHDB_USERNAME", "COUCHDB_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("COUCHDB_ENDPOINT, COUCHDB_USERNAME and COUCHDB_PASSWORD must be set for acceptance tests")
		}
	}

	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func testAccCouchDBUserWorks(endpoint, username, password string) error {
	transport := httptransport.New(endpoint, "", []string{"http"})
	transport.DefaultAuthentication = httptransport.BasicAuth(username, password)
	client := apiclient.New(transport, strfmt.Default)

	_, err := client.Server.Up(server.NewUpParams())
	if err != nil {
		return err
	}

	return nil
}
