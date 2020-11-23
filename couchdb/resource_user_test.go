package couchdb

import (
	"context"
	"fmt"
	"github.com/rossmerr/couchdb_go/client/document"
	"testing"

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

		params := document.NewDocInfoParams().WithDb(usersDB).WithDocid(rs.Primary.ID)
		_, err := client.Document.DocInfo(params)

		if  err != nil {
			return err
		}

		return testAccCouchDBUserWorks(rs.Primary.Attributes["endpoint"], rs.Primary.Attributes["name"], rs.Primary.Attributes["password"])
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

		params := document.NewDocInfoParams().WithDb(rs.Primary.Attributes["database"]).WithDocid(rs.Primary.ID)
		_, err := client.Document.DocInfo(params)

		if err == nil {
			return fmt.Errorf("user still exists")
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
