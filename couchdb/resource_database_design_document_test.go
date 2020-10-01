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

func TestAccCouchDBDesignDocument_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDesignDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDesignDocument,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
				),
			},
			{
				Config: testAccCouchDBDesignDocument_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
				),
			},
		},
	})
}

func testAccCouchDBDesignDocumentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("database design document ID is not set")
		}


		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		db := client.DB(context.Background(), rs.Primary.Attributes["database"])
		row := db.Get(context.Background(), rs.Primary.ID)

		if row.Err != nil {
			return row.Err
		}

		return nil
	}
}

func testAccCouchDBDesignDocumentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}

		db := client.DB(context.Background(), rs.Primary.Attributes["database"])
		row := db.Get(context.Background(), rs.Primary.ID)

		var ddoc map[string]interface{}
		if err := row.ScanDoc(&ddoc); err != nil {
			switch kivik.StatusCode(err) {
			case http.StatusNotFound:
				return nil
			}
			return err
		}

	}

	return nil
}

var testAccCouchDBDesignDocument = `
resource "couchdb_database" "test" {
	name = "test"
}
resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name = "test"
	view {
		name = "test"
		map = "function(doc) { emit(doc._id, doc); }"
	}
}
`
var testAccCouchDBDesignDocument_update = `
resource "couchdb_database" "test" {
	name = "test"
}
resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name = "test"
	view {
		name = "cat"
		map = "function(doc) { emit(doc._id, doc); }"
	}
	view {
		name = "test"
		map = "function(doc) { emit(doc._id, doc); }"
	}
}
`
