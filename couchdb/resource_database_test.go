package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/rossmerr/couchdb_go/client/database"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchDBDatabase_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCouchDBDatabase,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseExists("couchdb_database.test"),
				),
			},
			//resource.TestStep{
			//	Config: testAccCouchDBDatabase_security,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCouchDBDatabaseSecurity("couchdb_database.test"),
			//	),
			//},
			//resource.TestStep{
			//	Config: testAccCouchDBDatabase,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCouchDBDatabaseSecurity("couchdb_database.test"),
			//	),
			//},
		},
	})
}

func testAccCouchDBDatabaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("database ID is set")
		}
		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		params := database.NewExistsParams().WithDb(rs.Primary.ID)
		_, err := client.Database.Exists(params)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCouchDBDatabaseDestroy(s *terraform.State) error {
	client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
	if dd != nil {
		return fmt.Errorf(dd.Detail)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		params := database.NewExistsParams().WithDb(rs.Primary.ID)
		_, err := client.Database.Exists(params)
		if err == nil {
			return fmt.Errorf("DB still exists")
		}
	}

	return nil
}

//func testAccCouchDBDatabaseSecurity(n string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//
//		if rs.Primary.ID == "" {
//			return fmt.Errorf("database ID is not set")
//		}
//
//		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
//		if dd != nil {
//			return fmt.Errorf(dd.Detail)
//		}
//
//		db := client.DB(context.Background(), rs.Primary.ID)
//
//		sec, err := db.Security(context.Background())
//		if err != nil {
//			return err
//		}
//
//		if rs.Primary.Attributes["security.0.members.#"] != strconv.Itoa(len(sec.Members.Names)) {
//			return fmt.Errorf("expected %d members, got %s", len(sec.Members.Names), rs.Primary.Attributes["security.0.members.#"])
//		}
//		if rs.Primary.Attributes["security.0.member_roles.#"] != strconv.Itoa(len(sec.Members.Roles)) {
//			return fmt.Errorf("expected %d member roles, got %s", len(sec.Members.Roles), rs.Primary.Attributes["security.0.member_roles.#"])
//		}
//		if rs.Primary.Attributes["security.0.admins.#"] != strconv.Itoa(len(sec.Admins.Names)) {
//			return fmt.Errorf("expected %d admins, got %s", len(sec.Admins.Names), rs.Primary.Attributes["security.0.admins.#"])
//		}
//		if rs.Primary.Attributes["security.0.admin_roles.#"] != strconv.Itoa(len(sec.Admins.Roles)) {
//			return fmt.Errorf("expected %d admin roles, got %s", len(sec.Admins.Roles), rs.Primary.Attributes["security.0.admin_roles.#"])
//		}
//		return nil
//	}
//}

var testAccCouchDBDatabase_security = `
resource "couchdb_database" "test" {
	name = "test"
}
`

var testAccCouchDBDatabase = `
resource "couchdb_database" "test" {
	name = "test"
	clustering {
		shards   = 6
		replicas = 2
	}
}`
