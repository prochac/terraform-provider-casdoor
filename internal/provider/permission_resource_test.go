// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_permission.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPermissionResourceConfig(rName, "Test Permission"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Permission"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "effect", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPermissionResourceConfig(rName, "Updated Permission"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Permission"),
				),
			},
		},
	})
}

func TestAccPermissionResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_permission.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccPermissionResourceConfig(rName, "Test Permission"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccPermissionResourceConfig(rName, "Test Permission"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccPermissionResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_permission" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
}
`, name, displayName)
}
