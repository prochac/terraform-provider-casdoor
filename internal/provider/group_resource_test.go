// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccGroupResourceConfig(rName, "Test Group"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "is_top_group", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccGroupResourceConfig(rName, "Updated Group"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Group"),
				),
			},
		},
	})
}

func TestAccGroupResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccGroupResourceConfig(rName, "Test Group"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccGroupResourceConfig(rName, "Test Group"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"updated_time", "parent_name", "title", "key", "have_children"},
			},
		},
	})
}

func testAccGroupResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_group" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
}
`, name, displayName)
}
