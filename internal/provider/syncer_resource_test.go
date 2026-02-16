// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyncerResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_syncer.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccSyncerResourceConfig(rName, "Database", "localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "organization", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "type", "Database"),
					resource.TestCheckResourceAttr(resourceName, "host", "localhost"),
					resource.TestCheckResourceAttr(resourceName, "port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "database_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "database", "testdb"),
					resource.TestCheckResourceAttr(resourceName, "table", "users"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccSyncerResourceConfig(rName, "Database", "127.0.0.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "host", "127.0.0.1"),
				),
			},
		},
	})
}

func TestAccSyncerResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_syncer.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccSyncerResourceConfig(rName, "Database", "localhost"),
			},
			// ImportState testing
			{
				Config:                  testAccProviderConfig(config) + testAccSyncerResourceConfig(rName, "Database", "localhost"),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "built-in/" + rName,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "ssh_password", "ssl_mode", "ssh_type", "ssh_host", "ssh_port", "ssh_user", "cert"},
			},
		},
	})
}

func testAccSyncerResourceConfig(name, syncerType, host string) string { // nolint: unparam
	return fmt.Sprintf(`
resource "casdoor_syncer" "test" {
  owner         = "built-in"
  name          = %q
  organization  = "built-in"
  type          = %q
  host          = %q
  port          = 3306
  user          = "root"
  password      = "password"
  database_type = "mysql"
  database      = "testdb"
  table         = "users"
  is_enabled    = false
}
`, name, syncerType, host)
}
