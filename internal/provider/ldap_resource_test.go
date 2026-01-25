// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLdapResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rID := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_ldap.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccLdapResourceConfig(rID, "Test LDAP Server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", rID),
					resource.TestCheckResourceAttr(resourceName, "owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "server_name", "Test LDAP Server"),
					resource.TestCheckResourceAttr(resourceName, "host", "ldap.example.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccLdapResourceConfig(rID, "Updated LDAP Server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", rID),
					resource.TestCheckResourceAttr(resourceName, "server_name", "Updated LDAP Server"),
				),
			},
		},
	})
}

func TestAccLdapResource_withSsl(t *testing.T) {
	config := setupTestConfig(t)
	rID := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_ldap.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing with SSL configuration
			{
				Config: testAccProviderConfig(config) + testAccLdapResourceConfigWithSsl(rID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", rID),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_self_signed_cert", "true"),
					resource.TestCheckResourceAttr(resourceName, "port", "636"),
				),
			},
		},
	})
}

func TestAccLdapResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rID := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_ldap.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccLdapResourceConfig(rID, "Test LDAP Server"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccLdapResourceConfig(rID, "Test LDAP Server"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rID,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"password"},
				ImportStateVerifyIdentifierAttribute: "id",
			},
		},
	})
}

func testAccLdapResourceConfig(id, serverName string) string {
	return fmt.Sprintf(`
resource "casdoor_ldap" "test" {
  id          = %q
  owner       = "admin"
  server_name = %q
  host        = "ldap.example.com"
  port        = 389
  username    = "cn=admin,dc=example,dc=com"
  password    = "admin_password"
  base_dn     = "dc=example,dc=com"
}
`, id, serverName)
}

func testAccLdapResourceConfigWithSsl(id string) string {
	return fmt.Sprintf(`
resource "casdoor_ldap" "test" {
  id                     = %q
  owner                  = "admin"
  server_name            = "Secure LDAP Server"
  host                   = "ldaps.example.com"
  port                   = 636
  enable_ssl             = true
  allow_self_signed_cert = true
  username               = "cn=admin,dc=example,dc=com"
  password               = "admin_password"
  base_dn                = "dc=example,dc=com"
  filter                 = "(objectClass=posixAccount)"
}
`, id)
}
