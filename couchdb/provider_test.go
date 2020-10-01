package couchdb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
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

func testAccCouchDBUserWorks(endpoint, username, password, expectedRole string) error {
	client, err := kivik.New("couch", endpoint)
	if err != nil {
		return err
	}

	err = client.Authenticate(context.Background(), couchdb.BasicAuth(username, password))
	if err != nil {
		return err
	}

	sess, err := client.Session(context.Background())
	if err != nil {
		return err
	}

	if sess.Name != username {
		return fmt.Errorf("expected user %s, but got %s", username, sess.Name)
	}
	if sess.Roles[0] != expectedRole {
		return fmt.Errorf("expected user role %s, but got %s", expectedRole, sess.Roles)
	}
	return nil
}
