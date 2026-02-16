// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModelResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_model.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccModelResourceConfig(rName, "Test Model", "A test model"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Model"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test model"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttrSet(resourceName, "model_text"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccModelResourceConfig(rName, "Updated Model", "An updated model"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Model"),
					resource.TestCheckResourceAttr(resourceName, "description", "An updated model"),
				),
			},
		},
	})
}

func TestAccModelResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_model.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccModelResourceConfig(rName, "Test Model", "A test model"),
			},
			// ImportState testing
			{
				Config:            testAccProviderConfig(config) + testAccModelResourceConfig(rName, "Test Model", "A test model"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "built-in/" + rName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccModelResourceConfig(name, displayName, description string) string {
	return fmt.Sprintf(`
resource "casdoor_model" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
  description  = %q
  model_text   = <<-EOT
  [request_definition]
  r = sub, obj, act

  [policy_definition]
  p = sub, obj, act

  [role_definition]
  g = _, _

  [policy_effect]
  e = some(where (p.eft == allow))

  [matchers]
  m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
  EOT
}
`, name, displayName, description)
}
