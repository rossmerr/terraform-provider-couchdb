package couchdb

import (
	"context"
	"fmt"
	"github.com/rossmerr/couchdb_go/client/document"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchDBDatabaseReplication_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBDatabaseReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDatabaseReplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseReplicationExists("couchdb_database_replication.test"),
				),
			},
			{
				Config: testAccCouchDBDatabaseReplication_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseReplicationExists("couchdb_database_replication.test"),
				),
			},
		},
	})
}

func testAccCouchDBDatabaseReplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("replication ID is not set")
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		params := document.NewDocInfoParams().WithDb(rs.Primary.Attributes["database"]).WithDocid(rs.Primary.ID)
		_, err := client.Document.DocInfo(params)


		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCouchDBDatabaseReplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}




		params := document.NewDocInfoParams().WithDb(rs.Primary.Attributes["database"]).WithDocid(rs.Primary.ID)
		_, err := client.Document.DocInfo(params)
		if err == nil {
			return fmt.Errorf("replication still exists")
		}


	}

	return nil
}

var testAccCouchDBDatabaseReplication = `
resource "couchdb_database" "test" {
	name = "test"
}
resource "couchdb_database_replication" "test" {
	name = "test"
	source = "${couchdb_database.test.name}"
	target = "bar"
	create_target = true
	continuous = true
} 
`

var testAccCouchDBDatabaseReplication_update = `
resource "couchdb_database" "test" {
	name = "test"
}
resource "couchdb_database_replication" "test" {
	name = "test"
	
	source = "${couchdb_database.test.name}"
	target = "bar"
	create_target = true
	continuous = true
	context {
		user = "admin"
	}
	filter = "documents/by_author"
	query_params {
		author = "alex"
	}
} 
`
