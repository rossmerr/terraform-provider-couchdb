package couchdb

import (
	"context"
	"fmt"
	"github.com/RossMerr/couchdb_go/client/design_documents"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchDBDesignDocument_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBDesignDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDesignDocument,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
				),
			},
			//{
			//	Config: testAccCouchDBDesignDocument_update,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
			//	),
			//},
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

		params := design_documents.NewDesignDocExistsParams().WithDb(rs.Primary.Attributes["database"]).WithDdoc(rs.Primary.ID)
		_, err := client.DesignDocuments.DesignDocExists(params)

		if err != nil {
			return err
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

		params := design_documents.NewDesignDocExistsParams().WithDb(rs.Primary.Attributes["database"]).WithDdoc(rs.Primary.ID)
		_, err := client.DesignDocuments.DesignDocExists(params)
		if err == nil {
			return fmt.Errorf("design document still exists")
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
	views = <<EOF
	{
		"test": {
			"map": "function(doc) { emit(doc._id, doc); }"
		}
	}
EOF
}
`
var testAccCouchDBDesignDocument_update = `
resource "couchdb_database" "test2" {
	name = "test"
}
resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name = "test"
	views = <<EOF
	{
		"cat" : {
			"map" : "function(doc) { emit(doc._id, doc); }"
		}
	}
EOF
}
`
