// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTokenResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccTokenResourceConfig(config, rName, "read", 7200),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "application", config.ApplicationName),
					resource.TestCheckResourceAttr(resourceName, "organization", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "user", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scope", "read"),
					resource.TestCheckResourceAttr(resourceName, "token_type", "Bearer"),
					resource.TestCheckResourceAttr(resourceName, "expires_in", "7200"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccTokenResourceConfig(config, rName, "read write", 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "read write"),
					resource.TestCheckResourceAttr(resourceName, "expires_in", "3600"),
				),
			},
		},
	})
}

func TestAccTokenResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccTokenResourceConfig(config, rName, "read", 7200),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccTokenResourceConfig(config, rName, "read", 7200),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"code", "access_token", "refresh_token"},
			},
		},
	})
}

func testAccTokenResourceConfig(config CasdoorTestConfig, name, scope string, expiresIn int) string {
	return fmt.Sprintf(`
resource "casdoor_token" "test" {
  owner        = %q
  name         = %q
  application  = %q
  organization = %q
  user         = "admin"
  scope        = %q
  expires_in   = %d
}
`, "admin", name, config.ApplicationName, config.OrganizationName, scope, expiresIn)
}
