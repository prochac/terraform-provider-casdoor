// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_product.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccProductResourceConfig(rName, "Test Product", "9.99"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "price", "9.99"),
					resource.TestCheckResourceAttr(resourceName, "quantity", "100"),
					resource.TestCheckResourceAttr(resourceName, "currency", "USD"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccProductResourceConfig(rName, "Updated Product", "19.99"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Product"),
					resource.TestCheckResourceAttr(resourceName, "price", "19.99"),
				),
			},
		},
	})
}

func TestAccProductResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_product.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccProductResourceConfig(rName, "Test Product", "9.99"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccProductResourceConfig(rName, "Test Product", "9.99"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccProductResourceConfig(name, displayName, price string) string {
	return fmt.Sprintf(`
resource "casdoor_product" "test" {
  owner        = "built-in"
  name         = %q
  display_name = %q
  price        = %s
  quantity     = 100
  currency     = "USD"
}
`, name, displayName, price)
}
