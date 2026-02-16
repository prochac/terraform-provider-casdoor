// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdapterResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_adapter.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccAdapterResourceConfig(config.OrganizationName, rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "use_same_db", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccAdapterResourceConfig(config.OrganizationName, rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "use_same_db", "false"),
				),
			},
		},
	})
}

func TestAccAdapterResource_withDatabase(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_adapter.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing with database configuration
			{
				Config: testAccProviderConfig(config) + testAccAdapterResourceConfigWithDatabase(config.OrganizationName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "type", "Database"),
					resource.TestCheckResourceAttr(resourceName, "database_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "host", "localhost"),
					resource.TestCheckResourceAttr(resourceName, "port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "database", "casbin"),
					resource.TestCheckResourceAttr(resourceName, "table", "casbin_rule"),
				),
			},
		},
	})
}

func TestAccAdapterResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_adapter.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccAdapterResourceConfig(config.OrganizationName, rName, true),
			},
			// ImportState testing
			{
				Config:            testAccProviderConfig(config) + testAccAdapterResourceConfig(config.OrganizationName, rName, true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     config.OrganizationName + "/" + rName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAdapterResourceConfig(owner, name string, useSameDb bool) string {
	return fmt.Sprintf(`
resource "casdoor_adapter" "test" {
  owner       = %q
  name        = %q
  use_same_db = %t
}
`, owner, name, useSameDb)
}

func testAccAdapterResourceConfigWithDatabase(owner, name string) string {
	return fmt.Sprintf(`
resource "casdoor_adapter" "test" {
  owner         = %q
  name          = %q
  table         = "casbin_rule"
  type          = "Database"
  database_type = "mysql"
  host          = "localhost"
  port          = 3306
  user          = "root"
  password      = "password"
  database      = "casbin"
}
`, owner, name)
}
