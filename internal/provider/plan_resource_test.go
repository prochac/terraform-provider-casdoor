// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlanResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_plan.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPlanResourceConfig(rName, "Test Plan"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Plan"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "price", "0"),
					resource.TestCheckResourceAttr(resourceName, "currency", "USD"),
					resource.TestCheckResourceAttr(resourceName, "period", ""),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "product"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccPlanResourceConfig(rName, "Updated Plan"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Plan"),
				),
			},
		},
	})
}

func TestAccPlanResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_plan.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccPlanResourceConfig(rName, "Test Plan"),
			},
			// ImportState testing
			{
				Config:            testAccProviderConfig(config) + testAccPlanResourceConfig(rName, "Test Plan"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "built-in/" + rName,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPlanResource_productAdoption(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	productResourceName := "casdoor_product.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create a plan (auto-creates a backing product), then adopt the product.
			{
				Config: testAccProviderConfig(config) + testAccPlanProductAdoptionConfig(rName, "https://example.com/success"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(productResourceName, "managed_by_plan", "true"),
					resource.TestCheckResourceAttr(productResourceName, "success_url", "https://example.com/success"),
					resource.TestCheckResourceAttr(productResourceName, "owner", "built-in"),
				),
			},
			// Update the adopted product's success_url.
			{
				Config: testAccProviderConfig(config) + testAccPlanProductAdoptionConfig(rName, "https://example.com/updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(productResourceName, "managed_by_plan", "true"),
					resource.TestCheckResourceAttr(productResourceName, "success_url", "https://example.com/updated"),
				),
			},
		},
	})
}

func testAccPlanProductAdoptionConfig(name, successUrl string) string {
	return fmt.Sprintf(`
resource "casdoor_plan" "test" {
  owner        = "built-in"
  name         = %[1]q
  display_name = "Test Plan"
  currency     = "USD"
}

resource "casdoor_product" "test" {
  owner       = "built-in"
  name        = casdoor_plan.test.product
  success_url = %[2]q
}
`, name, successUrl)
}

func testAccPlanResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_plan" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
  currency     = "USD"
}
`, name, displayName)
}
