// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccRoleResourceConfig(config.OrganizationName, rName, "Test Role"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Role"),
					resource.TestCheckResourceAttr(resourceName, "owner", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccRoleResourceConfig(config.OrganizationName, rName, "Updated Role"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Role"),
				),
			},
		},
	})
}

func TestAccRoleResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccRoleResourceConfig(config.OrganizationName, rName, "Test Role"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccRoleResourceConfig(config.OrganizationName, rName, "Test Role"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccRoleResourceConfig(owner, name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_role" "test" {
  owner        = %q
  name         = %q
  display_name = %q
}
`, owner, name, displayName)
}
