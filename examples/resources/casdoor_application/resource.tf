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

  logo         = "https://cdn.casbin.org/img/casdoor-logo_1185x256.png"
  homepage_url = "https://myapp.example.com"
  description  = "My OAuth-enabled application"

  enable_password = true
  enable_sign_up  = true

  redirect_uris = [
    "https://myapp.example.com/callback",
    "https://myapp.example.com/oauth/callback",
  ]

  grant_types = [
    "authorization_code",
    "password",
    "client_credentials",
    "token",
    "id_token",
    "refresh_token",
  ]

  cert                    = "cert-built-in"
  token_format            = "JWT"
  expire_in_hours         = 24
  refresh_expire_in_hours = 168
}

# Full application with signin methods, signup items, and providers
resource "casdoor_application" "full_app" {
  name         = "full-app"
  display_name = "Full Application"
  organization = casdoor_organization.example.name

  logo        = "https://cdn.casbin.org/img/casdoor-logo_1185x256.png"
  description = "Application with all configuration options"

  # Authentication settings
  enable_password = true
  enable_sign_up  = true
  disable_signin  = false

  # OAuth settings
  redirect_uris = ["http://localhost:9000/callback"]
  grant_types = [
    "authorization_code",
    "password",
    "client_credentials",
    "token",
    "id_token",
    "refresh_token",
  ]

  cert                    = "cert-built-in"
  token_format            = "JWT"
  expire_in_hours         = 168
  refresh_expire_in_hours = 168

  # UI settings
  form_offset = 2

  # Providers (identity providers for the application)
  providers {
    name        = "provider_captcha_default"
    can_sign_up = false
    can_sign_in = false
    can_unlink  = false
    prompted    = false
  }

  # Signin methods
  signin_methods {
    name         = "Password"
    display_name = "Password"
    rule         = "All"
  }
  signin_methods {
    name         = "Verification code"
    display_name = "Verification code"
    rule         = "All"
  }
  signin_methods {
    name         = "WebAuthn"
    display_name = "WebAuthn"
    rule         = "None"
  }

  # Signup form items
  signup_items {
    name     = "ID"
    visible  = false
    required = true
    rule     = "Random"
  }
  signup_items {
    name     = "Username"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Display name"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Password"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Confirm password"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Email"
    visible  = true
    required = true
    rule     = "Normal"
  }
  signup_items {
    name     = "Phone"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Agreement"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Signup button"
    visible  = true
    required = true
    rule     = "None"
  }
  signup_items {
    name     = "Providers"
    visible  = true
    required = true
    rule     = "None"
    custom_css = <<-EOT
      .provider-img {
        width: 30px;
        margin: 5px;
      }
      .provider-big-img {
        margin-bottom: 10px;
      }
    EOT
  }
}

# Reference the generated OAuth credentials
output "client_id" {
  value = casdoor_application.oauth_app.client_id
}

output "client_secret" {
  value     = casdoor_application.oauth_app.client_secret
  sensitive = true
}
