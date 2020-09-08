package couchdb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"couchdb": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"COUCHDB_ENDPOINT", "COUCHDB_USERNAME", "COUCHDB_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("COUCHDB_ENDPOINT, COUCHDB_USERNAME and COUCHDB_PASSWORD must be set for acceptance tests")
		}
	}

	err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}