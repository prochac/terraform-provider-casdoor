# Basic application
resource "casdoor_application" "example" {
  name         = "my-app"
  display_name = "My Application"
  organization = "my-organization"
}

# Application with OAuth settings
resource "casdoor_application" "oauth_app" {
  name         = "oauth-app"
  display_name = "OAuth Application"
  organization = casdoor_organization.example.name

  logo         = "https://example.com/app-logo.png"
  homepage_url = "https://myapp.example.com"
  description  = "My OAuth-enabled application"

  enable_password = true
  enable_sign_up  = true

  redirect_uris = [
    "https://myapp.example.com/callback",
    "https://myapp.example.com/oauth/callback",
  ]

  token_format            = "JWT"
  expire_in_hours         = 24
  refresh_expire_in_hours = 168
}

# Reference the generated OAuth credentials
output "client_id" {
  value = casdoor_application.oauth_app.client_id
}

output "client_secret" {
  value     = casdoor_application.oauth_app.client_secret
  sensitive = true
}
