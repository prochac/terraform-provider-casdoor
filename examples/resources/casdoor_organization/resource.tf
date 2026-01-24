# Basic organization
resource "casdoor_organization" "example" {
  name         = "my-organization"
  display_name = "My Organization"
}

# Organization with all settings (based on Casdoor UI defaults)
resource "casdoor_organization" "full" {
  name         = "full-organization"
  display_name = "Full Organization"

  website_url = "https://door.casdoor.com"
  favicon     = "https://cdn.casbin.org/img/favicon.png"

  # Password settings
  password_type            = "bcrypt"
  password_options         = ["AtLeast6"]
  password_obfuscator_type = "Plain"
  password_expire_days     = 0

  # Regional settings
  country_codes = ["US"]
  languages     = ["en", "es", "fr", "de", "ja", "zh", "vi", "pt", "tr", "pl", "uk"]

  # User defaults
  default_avatar = "https://cdn.casbin.org/img/casbin.svg"

  # Feature flags
  enable_soft_deletion = false
  is_profile_public    = true
  enable_tour          = true
  disable_signin       = false

  # MFA settings
  mfa_remember_in_hours = 12

  # Account items - controls user profile fields visibility and editability
  account_items = [
    {
      name        = "Organization"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "ID"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Immutable"
    },
    {
      name        = "Name"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Display name"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Avatar"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "User type"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Password"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "Email"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Phone"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Country code"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Country/Region"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Location"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Address"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Affiliation"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Title"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "ID card type"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "ID card"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "ID card info"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Real name"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "ID verification"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "Homepage"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Bio"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Self"
    },
    {
      name        = "Tag"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Language"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Gender"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Birthday"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Education"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Score"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Karma"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Ranking"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Balance"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Balance credit"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Balance currency"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Signup application"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Register type"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Register source"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "API key"
      visible     = false
      modify_rule = "Self"
    },
    {
      name        = "Groups"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Admin"
    },
    {
      name        = "Roles"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Immutable"
    },
    {
      name        = "Permissions"
      visible     = true
      view_rule   = "Public"
      modify_rule = "Immutable"
    },
    {
      name        = "3rd-party logins"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "Properties"
      visible     = false
      view_rule   = "Admin"
      modify_rule = "Admin"
    },
    {
      name        = "Is online"
      visible     = true
      view_rule   = "Admin"
      modify_rule = "Admin"
    },
    {
      name        = "Is admin"
      visible     = true
      view_rule   = "Admin"
      modify_rule = "Admin"
    },
    {
      name        = "Is forbidden"
      visible     = true
      view_rule   = "Admin"
      modify_rule = "Admin"
    },
    {
      name        = "Is deleted"
      visible     = true
      view_rule   = "Admin"
      modify_rule = "Admin"
    },
    {
      name        = "Multi-factor authentication"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "WebAuthn credentials"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "Managed accounts"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
    {
      name        = "MFA accounts"
      visible     = true
      view_rule   = "Self"
      modify_rule = "Self"
    },
  ]
}

# Organization with theme customization
resource "casdoor_organization" "themed" {
  name         = "themed-org"
  display_name = "Themed Organization"

  theme_data = {
    theme_type    = "default"
    color_primary = "#1890ff"
    border_radius = 6
    is_compact    = false
    is_enabled    = true
  }
}

# Organization with MFA configuration
resource "casdoor_organization" "secure" {
  name         = "secure-org"
  display_name = "Secure Organization"

  mfa_remember_in_hours = 24

  mfa_items = [
    {
      name = "app"
      rule = "Optional"
    },
    {
      name = "email"
      rule = "Optional"
    },
    {
      name = "sms"
      rule = "Optional"
    },
  ]
}
