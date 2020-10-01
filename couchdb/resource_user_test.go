package couchdb

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchDBUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBUserExists("couchdb_user.test"),
				),
			},
		},
	})
}

func testAccCouchDBUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("user ID is not set")
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		db := client.DB(context.Background(), usersDB)

		row := db.Get(context.Background(), rs.Primary.ID)
		var user tuser
		if err := row.ScanDoc(&user); err != nil {
			return err
		}

		return testAccCouchDBUserWorks(client.DSN(), rs.Primary.Attributes["name"], rs.Primary.Attributes["password"], "developer")
	}
}

func testAccCouchDBUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		db := client.DB(context.Background(), usersDB)

		row := db.Get(context.Background(), rs.Primary.ID)

		var user tuser
		if err := row.ScanDoc(&user); err != nil {
			switch kivik.StatusCode(err) {
			case http.StatusNotFound:
				return nil
			}
			return err
		}
	}

	return nil
}

var testAccCouchDBUser = `
resource "couchdb_user" "test" {
	name = "test"
	password = "test"
	roles = ["developer"]
}`
