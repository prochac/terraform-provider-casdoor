// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdpResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_provider.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccIdpResourceConfig(rName, "Test Provider"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Provider"),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "category", "OAuth"),
					resource.TestCheckResourceAttr(resourceName, "type", "GitHub"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "test-client-id"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "test-client-secret"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccIdpResourceConfig(rName, "Updated Provider"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Provider"),
				),
			},
		},
	})
}

func TestAccIdpResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_provider.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccIdpResourceConfig(rName, "Test Provider"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccIdpResourceConfig(rName, "Test Provider"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"client_secret", "client_secret_2"},
			},
		},
	})
}

func testAccIdpResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_provider" "test" {
  owner         = "built-in"
  name          = %q
  display_name  = %q
  category      = "OAuth"
  type          = "GitHub"
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
}
`, name, displayName)
}
