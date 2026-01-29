// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCertResource_basic(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_cert.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(config) + testAccCertResourceConfig(rName, "Test Cert"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Cert"),
					resource.TestCheckResourceAttr(resourceName, "owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scope", "JWT"),
					resource.TestCheckResourceAttr(resourceName, "type", "x509"),
					resource.TestCheckResourceAttr(resourceName, "crypto_algorithm", "RS256"),
					resource.TestCheckResourceAttr(resourceName, "bit_size", "4096"),
					resource.TestCheckResourceAttr(resourceName, "expire_in_years", "20"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(config) + testAccCertResourceConfig(rName, "Updated Cert"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated Cert"),
				),
			},
		},
	})
}

func TestAccCertResource_import(t *testing.T) {
	config := setupTestConfig(t)
	rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resourceName := "casdoor_cert.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(config),
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccProviderConfig(config) + testAccCertResourceConfig(rName, "Test Cert"),
			},
			// ImportState testing
			{
				Config:                               testAccProviderConfig(config) + testAccCertResourceConfig(rName, "Test Cert"),
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        rName,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"private_key"},
			},
		},
	})
}

const testCertPEM = `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAK3MN0KQGsQiMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RDQTAEFW0yNTAxMDEwMDAwMDBaFw0yNjAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RDQTBcMA0GCSqGSIb3DQEBAQUAAwsAMEgCQQC7o96Gv3MR0Iqx3MD0QNKL
HnGGPJkGRvEjBJSFQM5QngVeDB0JsMB/tITFocOFIiD5RVFmGa6s1YODpazVimEd
AgMBAAGjUDBOMB0GA1UdDgQWBBQ5sFpVHs3gNUaXMSBPbJNEX7ZTGB8GA1UdIwQY
MBaAFDmwWlUezeA1RpcxIE9sk0RftlMYMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcN
AQELBQADQQBxXV7VEOEBi6dAeL0E5MFv3NvNfvg6Y+VlJzb4dLUKYMFkLAGO78Q
Z3MPm0e6T3RBMcfB6DGwKkbq5WHFupoT
-----END CERTIFICATE-----`

const testPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBogIBAAJBALuj3oa/cxHQirHcwPRA0osecYY8mQZG8SMElIVAzlCeBV4MHQmw
wH+0hMWhw4UiIPlFUWYZrqzVg4OlrNWKYR0CAwEAAQJAGz5MnQb3BD4ydHDJcdR0
rEMCMQC7o96Gv3MR0Iqx3MD0QNKLHZ3LzTDIDAK7o96Gv3MR0Iqx3MD0QNKLHQ
IhAO+eMIAEsncP1rGVUEXiFiIxLb3K2K0oF7a0w8x8R25AiEA2HfWC7gIKJ0sNF
QJBAfj3o96Gv3MR0Iqx3MD0QNKLHnGGPJk4yB0dJbEMC7xLzL2D5SZ1A0oRqx3
IhAO9dPKGv3MR0Iqx3MD0QNKLHZ3LzCIDAK7o96SZ1A0oRqx3MD0QNKLHZ3LzL
IgEAl87TqD0nH9DFz3BRsF3KxJJvQBXW5z0MF3VqB8JGRp0=
-----END RSA PRIVATE KEY-----`

func testAccCertResourceConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "casdoor_cert" "test" {
  owner        = "admin"
  name         = %q
  display_name = %q
  certificate  = %q
  private_key  = %q
}
`, name, displayName, testCertPEM, testPrivateKeyPEM)
}
