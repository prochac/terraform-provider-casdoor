# Basic organization
resource "casdoor_organization" "example" {
  name         = "my-organization"
  display_name = "My Organization"
}

# Organization with additional settings
resource "casdoor_organization" "advanced" {
  name         = "advanced-org"
  display_name = "Advanced Organization"

  website_url    = "https://example.com"
  logo           = "https://example.com/logo.png"
  favicon        = "https://example.com/favicon.ico"
  password_type  = "bcrypt"
  default_avatar = "https://example.com/default-avatar.png"

  tags      = ["production", "main"]
  languages = ["en", "de", "fr"]

  enable_soft_deletion  = true
  is_profile_public     = false
  use_email_as_username = true
}
