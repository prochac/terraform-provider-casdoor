// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestCasdoorSDKApplicationClient tests the Casdoor SDK client directly for applications.
func TestCasdoorSDKApplicationClient(t *testing.T) {
	config := setupTestConfig(t)

	client := casdoorsdk.NewClient(
		config.Endpoint,
		config.ClientID,
		config.ClientSecret,
		config.Certificate,
		config.OrganizationName,
		config.ApplicationName,
	)

	// Test getting existing applications.
	apps, err := client.GetApplications()
	if err != nil {
		t.Fatalf("Failed to get applications: %v", err)
	}

	t.Logf("Found %d applications", len(apps))
	for _, app := range apps {
		t.Logf("  - %s/%s (%s)", app.Owner, app.Name, app.DisplayName)
	}

	// Test CRUD operations with a test application.
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	// Create.
	newApp := &casdoorsdk.Application{
		Owner:          "admin",
		Name:           rName,
		DisplayName:    "Test Application",
		Organization:   config.OrganizationName,
		Logo:           "https://cdn.casbin.org/img/casdoor-logo_1185x256.png",
		HomepageUrl:    "https://example.com",
		Description:    "Test application for terraform provider",
		EnablePassword: true,
		EnableSignUp:   true,
	}

	success, err := client.AddApplication(newApp)
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}
	if !success {
		t.Fatal("AddApplication returned false")
	}
	t.Logf("Created application: %s", rName)

	// Read.
	app, err := client.GetApplication(rName)
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}
	if app == nil {
		t.Fatal("GetApplication returned nil")
	}
	t.Logf("Read application: %s/%s (ClientID: %s)", app.Owner, app.Name, app.ClientId)

	// Update.
	app.DisplayName = "Updated Test Application"
	success, err = client.UpdateApplication(app)
	if err != nil {
		t.Fatalf("Failed to update application: %v", err)
	}
	if !success {
		t.Fatal("UpdateApplication returned false")
	}
	t.Logf("Updated application: %s", rName)

	// Verify update.
	app, err = client.GetApplication(rName)
	if err != nil {
		t.Fatalf("Failed to get updated application: %v", err)
	}
	if app.DisplayName != "Updated Test Application" {
		t.Fatalf("Expected DisplayName 'Updated Test Application', got '%s'", app.DisplayName)
	}

	// Delete - include Organization as Casdoor API requires it.
	success, err = client.DeleteApplication(&casdoorsdk.Application{
		Owner:        "admin",
		Name:         rName,
		Organization: config.OrganizationName,
	})
	if err != nil {
		t.Fatalf("Failed to delete application: %v", err)
	}
	if !success {
		t.Fatal("DeleteApplication returned false")
	}
	t.Logf("Deleted application: %s", rName)

	// Verify deletion.
	app, err = client.GetApplication(rName)
	if err != nil {
		t.Logf("GetApplication after delete returned error (expected): %v", err)
	}
	if app != nil {
		t.Fatal("GetApplication returned non-nil after deletion")
	}
}

func TestAccApplicationResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_application.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccProviderConfig(config) + testAccApplicationResourceConfig(rName, config.OrganizationName, "Test Application"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Application"),
					resource.TestCheckResourceAttr(resourceName, "owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "organization", config.OrganizationName),
					resource.TestCheckResourceAttr(resourceName, "enable_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_sign_up", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
				),
			},
			// Update and Read testing.
			{
				Config: testAccProviderConfig(config) + testAccApplicationResourceConfig(rName, config.OrganizationName, "Updated Application"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Application"),
				),
			},
		},
	})
}

func TestAccApplicationResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_application.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first.
			{
				Config: testAccProviderConfig(config) + testAccApplicationResourceConfig(rName, config.OrganizationName, "Test Application"),
			},
			// ImportState testing.
			{
				Config:                               testAccProviderConfig(config) + testAccApplicationResourceConfig(rName, config.OrganizationName, "Test Application"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccApplicationResource_withRedirectURIs(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_application.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig(config) + testAccApplicationResourceConfigWithRedirectURIs(rName, config.OrganizationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.0", "https://example.com/callback"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.1", "https://example.com/oauth"),
				),
			},
		},
	})
}

func testAccApplicationResourceConfig(name, organization, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_application" "test" {
  name         = %q
  display_name = %q
  organization = %q
}
`, name, displayName, organization)
}

func testAccApplicationResourceConfigWithRedirectURIs(name, organization string) string {
	return fmt.Sprintf(`
resource "casdoor_application" "test" {
  name         = %q
  display_name = "Test App with Redirect URIs"
  organization = %q
  redirect_uris = [
    "https://example.com/callback",
    "https://example.com/oauth"
  ]
}
`, name, organization)
}
