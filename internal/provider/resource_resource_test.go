// Copyright (c) HashiCorp, Inc.

package provider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ensureStorageProvider creates a "Local File System" storage provider in Casdoor
// and links it to the test application so that resource uploads work.
func ensureStorageProvider(t *testing.T, config CasdoorTestConfig) {
	t.Helper()

	client := casdoorsdk.NewClient(
		config.Endpoint,
		config.ClientID,
		config.ClientSecret,
		config.Certificate,
		config.OrganizationName,
		config.ApplicationName,
	)

	providerName := "provider-storage-local"

	// Check if provider already exists.
	existing, err := client.GetProvider(providerName)
	if err != nil {
		t.Fatalf("Failed to check for storage provider: %v", err)
	}

	if existing == nil {
		// Create a Local File System storage provider.
		_, err = client.AddProvider(&casdoorsdk.Provider{
			Owner:       "admin",
			Name:        providerName,
			DisplayName: "Local File System",
			Category:    "Storage",
			Type:        "Local File System",
		})
		if err != nil {
			t.Fatalf("Failed to create storage provider: %v", err)
		}
	}

	// Get the application and add the storage provider if not already present.
	app, err := client.GetApplication("app-built-in")
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}

	found := false
	for _, p := range app.Providers {
		if p.Name == providerName {
			found = true
			break
		}
	}

	if !found {
		app.Providers = append(app.Providers, &casdoorsdk.ProviderItem{
			Owner: "admin",
			Name:  providerName,
			Provider: &casdoorsdk.Provider{
				Owner:       "admin",
				Name:        providerName,
				DisplayName: "Local File System",
				Category:    "Storage",
				Type:        "Local File System",
			},
		})

		_, err = client.UpdateApplication(app)
		if err != nil {
			t.Fatalf("Failed to update application with storage provider: %v", err)
		}
	}
}

func TestAccResourceResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	ensureStorageProvider(t, config)

	rSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	fileName := fmt.Sprintf("/casdoor/tf-test-%s.txt", rSuffix)
	resourceName := "casdoor_resource.test"
	content := "hello world"
	contentB64 := base64.StdEncoding.EncodeToString([]byte(content))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccResourceResourceConfig(fileName, contentB64),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "owner", config.OrganizationName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "user", "admin"),
					resource.TestCheckResourceAttr(resourceName, "tag", "txt"),
					resource.TestCheckResourceAttr(resourceName, "file_name", fileName),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccResourceResource_import(t *testing.T) {
	config := setupTestConfig(t)
	ensureStorageProvider(t, config)

	rSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	fileName := fmt.Sprintf("/casdoor/tf-test-%s.txt", rSuffix)
	resourceName := "casdoor_resource.test"
	content := "hello world"
	contentB64 := base64.StdEncoding.EncodeToString([]byte(content))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccResourceResourceConfig(fileName, contentB64),
			},
			// ImportState testing
			{
				Config:                  testAccProviderConfig(config) + testAccResourceResourceConfig(fileName, contentB64),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content_base64", "file_name"},
			},
		},
	})
}

func testAccResourceResourceConfig(fileName, contentB64 string) string {
	return fmt.Sprintf(`
resource "casdoor_resource" "test" {
  user           = "admin"
  tag            = "txt"
  file_name      = %q
  content_base64 = %q
}
`, fileName, contentB64)
}
