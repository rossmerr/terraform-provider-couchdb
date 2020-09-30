package couchdb

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"net/http"
	"testing"
)

func TestAccCouchDBDocument(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDocument,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDocumentExists("couchdb_document.test"),
				),
			},
		},
	})
}

func testAccCouchDBDocumentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("document ID is not set")
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

func testAccCouchDBDocumentDestroy(s *terraform.State) error {
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

		var ddoc map[string]interface{}
		if err = row.ScanDoc(&ddoc); err != nil {
			switch kivik.StatusCode(err) {
			case http.StatusNotFound:
				return nil
			}
			return err
		}

	}

	return nil
}

var testAccCouchDBDocument = `
resource "couchdb_database" "test" {
	name = "test"
}
resource "couchdb_document" "test" {
	database = "${couchdb_database.test.name}"
	doc = <<EOF
	{
		"description": "An Italian-American dish that usually consists of spaghetti, tomato sauce and meatballs.",
		"ingredients": [
			"spaghetti",
			"tomato sauce",
			"meatballs"
		],
		"name": "Spaghetti with meatballs"
	}
EOF
}`
