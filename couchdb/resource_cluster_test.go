package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/rossmerr/couchdb_go/client/database"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchDBCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccCouchDBDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCouchDBCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBClusterExists("couchdb_cluster.test"),
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

func testAccCouchDBClusterExists(n string) resource.TestCheckFunc {
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

var testAccCouchDBCluster = `
resource "couchdb_cluster" "test" {
	action = "finish_cluster"
}`
