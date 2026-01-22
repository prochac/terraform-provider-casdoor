# Auto-generated certificate (Casdoor generates the keys)
resource "casdoor_cert" "jwt_signing" {
  owner            = "my-organization"
  name             = "jwt-signing-cert"
  display_name     = "JWT Signing Certificate"
  scope            = "JWT"
  type             = "x509"
  crypto_algorithm = "RS256"
  bit_size         = 4096
  expire_in_years  = 20
}

# Certificate with provided keys
resource "casdoor_cert" "custom_cert" {
  owner            = "my-organization"
  name             = "custom-cert"
  display_name     = "Custom Certificate"
  scope            = "JWT"
  type             = "x509"
  crypto_algorithm = "RS256"
  bit_size         = 2048
  expire_in_years  = 10

  certificate = file("path/to/certificate.pem")
  private_key = file("path/to/private-key.pem")
}
