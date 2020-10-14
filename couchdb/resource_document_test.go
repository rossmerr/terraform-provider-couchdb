package couchdb

import (
	"context"
	"fmt"
	"github.com/RossMerr/couchdb_go/client/document"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccCouchDBDocument(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBDocumentDestroy,
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

func testAccCouchDBDocumentDestroy(s *terraform.State) error {
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
			return fmt.Errorf("document still exists")
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
