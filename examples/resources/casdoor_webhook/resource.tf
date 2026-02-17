# Basic webhook for user events
resource "casdoor_webhook" "user_events" {
  owner        = "my-organization"
  name         = "webhook-user-events"
  organization = "my-organization"
  url          = "https://api.example.com/webhooks/casdoor"
  method       = "POST"
  content_type = "application/json"
  is_enabled   = true

  events = [
    "signup",
    "login",
    "logout",
  ]
}

# Webhook with custom headers and extended user info
resource "casdoor_webhook" "audit" {
  owner        = "my-organization"
  name         = "webhook-audit"
  organization = "my-organization"
  url          = "https://audit.example.com/events"
  method       = "POST"
  content_type = "application/json"
  is_enabled   = true

  is_user_extended = true
  single_org_only  = true

  events = [
    "signup",
    "login",
    "logout",
    "update-user",
  ]

  headers {
    name  = "Authorization"
    value = "Bearer ${var.webhook_token}"
  }
}
