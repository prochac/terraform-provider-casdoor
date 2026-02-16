// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// enableBuiltInUserCreation enables user creation in the built-in org by
// setting HasPrivilegeConsent to true.
func enableBuiltInUserCreation(t *testing.T, config CasdoorTestConfig) {
	t.Helper()

	client := casdoorsdk.NewClient(
		config.Endpoint,
		config.ClientID,
		config.ClientSecret,
		config.Certificate,
		config.OrganizationName,
		config.ApplicationName,
	)

	org, err := client.GetOrganization(config.OrganizationName)
	if err != nil {
		t.Fatalf("Failed to get organization: %v", err)
	}
	if org == nil {
		t.Fatalf("Organization %q not found", config.OrganizationName)
	}

	org.HasPrivilegeConsent = true
	_, err = client.UpdateOrganization(org)
	if err != nil {
		t.Fatalf("Failed to update organization: %v", err)
	}
}

func TestAccUserResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	enableBuiltInUserCreation(t, config)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_user.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccUserResourceConfig(config.OrganizationName, rName, "Test User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "owner", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "type", "normal-user"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccUserResourceConfig(config.OrganizationName, rName, "Updated User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated User"),
				),
			},
		},
	})
}

func TestAccUserResource_import(t *testing.T) {
	config := setupTestConfig(t)
	enableBuiltInUserCreation(t, config)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_user.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccUserResourceConfig(config.OrganizationName, rName, "Test User"),
			},
			// ImportState testing
			{
				Config:                  testAccProviderConfig(config) + testAccUserResourceConfig(config.OrganizationName, rName, "Test User"),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           config.OrganizationName + "/" + rName,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "password_salt", "access_key", "access_secret", "totp_secret", "recovery_codes", "id_card"},
			},
		},
	})
}

func testAccUserResourceConfig(owner, name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_user" "test" {
  owner        = %q
  name         = %q
  display_name = %q
}
`, owner, name, displayName)
}
