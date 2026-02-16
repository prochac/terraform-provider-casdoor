// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPricingResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_pricing.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPricingResourceConfig(rName, "Test Pricing"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Pricing"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "trial_duration", "0"),
					resource.TestCheckResourceAttr(resourceName, "application", ""),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPricingResourceConfig(rName, "Updated Pricing"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Pricing"),
				),
			},
		},
	})
}

func TestAccPricingResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_pricing.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccPricingResourceConfig(rName, "Test Pricing"),
			},
			// ImportState testing
			{
				Config:            testAccProviderConfig(config) + testAccPricingResourceConfig(rName, "Test Pricing"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "built-in/" + rName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPricingResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_pricing" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
}
`, name, displayName)
}
