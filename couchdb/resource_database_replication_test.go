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

func TestAccCouchDBDatabaseReplication_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDatabaseReplicationDestroy,
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

		client, err := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if err != nil {
			return err
		}

		db := client.DB(context.Background(), rs.Primary.Attributes["database"])
		row := db.Get(context.Background(), rs.Primary.ID)

		if row.Err != nil {
			return row.Err
		}

		return nil
	}
}

func testAccCouchDBDatabaseReplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client, err := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if err != nil {
			return err
		}

		db := client.DB(context.Background(), rs.Primary.Attributes["database"])
		row := db.Get(context.Background(), rs.Primary.ID)

		var rep map[string]interface{}
		if err = row.ScanDoc(&rep); err != nil {
			switch kivik.StatusCode(err) {
			case http.StatusNotFound:
				return nil
			}
			return err
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