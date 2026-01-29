// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnforcerResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_enforcer.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccEnforcerResourceConfig(rName, "Test Enforcer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-enforcer"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Enforcer"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttrSet(resourceName, "model"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccEnforcerResourceConfig(rName, "Updated Enforcer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-enforcer"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Enforcer"),
				),
			},
		},
	})
}

func TestAccEnforcerResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_enforcer.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccEnforcerResourceConfig(rName, "Test Enforcer"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccEnforcerResourceConfig(rName, "Test Enforcer"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName + "-enforcer",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccEnforcerResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_model" "test" {
  owner        = "built-in"
  name         = %[1]q
  display_name = "Test Model"
  model_text   = <<-EOT
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
  EOT
}

resource "casdoor_enforcer" "test" {
  owner        = "built-in"
  name         = "%[1]s-enforcer"
  display_name = %[2]q
  model        = "built-in/${casdoor_model.test.name}"
}
`, name, displayName)
}
