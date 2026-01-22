# Basic token
resource "casdoor_token" "example" {
  owner        = "my-organization"
  name         = "token-example"
  application  = "my-app"
  organization = "my-organization"
  user         = "john.doe"
  expires_in   = 7200
  token_type   = "Bearer"
  scope        = "read,write"
}
