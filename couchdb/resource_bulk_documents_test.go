package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RossMerr/couchdb_go/client/database"
	"github.com/RossMerr/couchdb_go/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccCouchDBBulkDocuments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBBulkDocumentsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBBulkDocuments,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBBulkDocumentsExists("couchdb_bulk_documents.test"),
				),
			},
		},
	})
}

func testAccCouchDBBulkDocumentsExists(n string) resource.TestCheckFunc {
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

		var revisions map[string]string
		err := json.Unmarshal([]byte(rs.Primary.ID), &revisions)
		if err != nil {
			return err
		}

		docs := []*models.BasicDoc{}
		for id, rev := range revisions {
			docs = append(docs, &models.BasicDoc{
				ID: id,
				Rev: rev,
			})
		}

		body := database.BulkGetBody{
			Docs: docs,
		}

		params := database.NewBulkGetParams().WithDb(rs.Primary.Attributes["database"]).WithBody(body)
		_, err = client.Database.BulkGet(params)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCouchDBBulkDocumentsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client, dd := connectToCouchDB(context.Background(), testAccProvider.Meta().(*CouchDBConfiguration))
		if dd != nil {
			return fmt.Errorf(dd.Detail)
		}


		var revisions map[string]string
		err := json.Unmarshal([]byte(rs.Primary.ID), &revisions)
		if err != nil {
			return err
		}

		docs := []*models.BasicDoc{}
		for id, rev := range revisions {
			docs = append(docs, &models.BasicDoc{
				ID: id,
				Rev: rev,
			})
		}

		body := database.BulkGetBody{
			Docs: docs,
		}

		params := database.NewBulkGetParams().WithDb(rs.Primary.Attributes["database"]).WithBody(body)
		_, err = client.Database.BulkGet(params)

		if err == nil {
			return fmt.Errorf("documents still exists")
		}
	}

	return nil
}

var testAccCouchDBBulkDocuments = `
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_bulk_documents" "test" {
  database = "test"
  docs = <<EOF
  [
    {
          "_id": "9391913b56c655881fa57d60830008ac",
          "description": "An Italian-American dish that usually consists of spaghetti, tomato sauce and meatballs.",
          "ingredients": [
              "spaghetti",
              "tomato sauce",
              "meatballs"
          ],
          "name": "Spaghetti"
      }
  ]
EOF
}`
