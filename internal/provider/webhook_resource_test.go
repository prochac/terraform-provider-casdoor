// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_webhook.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccWebhookResourceConfig(rName, "https://example.com/hook"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", "built-in"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/hook"),
					resource.TestCheckResourceAttr(resourceName, "method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "is_user_extended", "false"),
					resource.TestCheckResourceAttr(resourceName, "single_org_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccWebhookResourceConfig(rName, "https://example.com/hook-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/hook-updated"),
				),
			},
		},
	})
}

func TestAccWebhookResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_webhook.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccWebhookResourceConfig(rName, "https://example.com/hook"),
			},
			// ImportState testing
			{
				Config:            testAccProviderConfig(config) + testAccWebhookResourceConfig(rName, "https://example.com/hook"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "built-in/" + rName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWebhookResourceConfig(name, url string) string {
	return fmt.Sprintf(`
resource "casdoor_webhook" "test" {
  owner = "built-in"
  name  = %q
  url   = %q
}
`, name, url)
}
