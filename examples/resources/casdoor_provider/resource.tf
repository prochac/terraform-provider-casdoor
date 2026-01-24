# =============================================================================
# OAuth Providers
# =============================================================================

# Google OAuth provider
resource "casdoor_provider" "google" {
  owner        = "my-organization"
  name         = "provider-google"
  display_name = "Google"
  category     = "OAuth"
  type         = "Google"

  client_id     = var.google_client_id
  client_secret = var.google_client_secret
  scopes        = "openid profile email"
}

# GitHub OAuth provider
resource "casdoor_provider" "github" {
  owner        = "my-organization"
  name         = "provider-github"
  display_name = "GitHub"
  category     = "OAuth"
  type         = "GitHub"

  client_id     = var.github_client_id
  client_secret = var.github_client_secret
  scopes        = "user:email read:user"
}

# Custom OAuth provider (for self-hosted or non-standard OAuth2 servers)
resource "casdoor_provider" "custom_oauth" {
  owner        = "my-organization"
  name         = "provider-custom-oauth"
  display_name = "My Custom IdP"
  category     = "OAuth"
  type         = "Custom"

  client_id     = var.custom_oauth_client_id
  client_secret = var.custom_oauth_client_secret

  # Custom endpoints (required for Custom OAuth)
  custom_auth_url      = "https://idp.example.com/oauth2/authorize"
  custom_token_url     = "https://idp.example.com/oauth2/token"
  custom_user_info_url = "https://idp.example.com/oauth2/userinfo"
  scopes               = "openid profile email"

  # Custom logo for login page
  custom_logo = "https://example.com/logo.png"

  # Map IdP user attributes to Casdoor user fields
  user_mapping = {
    "id"          = "sub"
    "username"    = "preferred_username"
    "displayName" = "name"
    "email"       = "email"
    "avatarUrl"   = "picture"
  }
}

# =============================================================================
# Email Providers
# =============================================================================

# Default (SMTP) email provider
resource "casdoor_provider" "email_smtp" {
  owner        = "my-organization"
  name         = "provider-email-smtp"
  display_name = "SMTP Email"
  category     = "Email"
  type         = "Default"

  # SMTP server settings
  host        = "smtp.example.com"
  port        = 587
  disable_ssl = false

  # SMTP credentials
  client_id     = "smtp-username@example.com"
  client_secret = var.smtp_password

  # Email template
  title   = "Your Verification Code"
  content = "Your verification code is: %s"

  # Test email recipient (optional, for testing)
  receiver = ""
}

# SendGrid email provider
resource "casdoor_provider" "email_sendgrid" {
  owner        = "my-organization"
  name         = "provider-email-sendgrid"
  display_name = "SendGrid Email"
  category     = "Email"
  type         = "SendGrid"

  # SendGrid API credentials
  # client_id is the sender email address
  client_id     = "noreply@example.com"
  client_secret = var.sendgrid_api_key

  # Email template
  title   = "Your Verification Code"
  content = "Your verification code is: %s"

  # Test email recipient (optional, for testing)
  receiver = ""
}

# =============================================================================
# Storage Providers
# =============================================================================

# AWS S3 storage provider
resource "casdoor_provider" "storage_s3" {
  owner        = "my-organization"
  name         = "provider-storage-s3"
  display_name = "AWS S3 Storage"
  category     = "Storage"
  type         = "AWS S3"

  # AWS credentials
  client_id     = var.aws_access_key_id
  client_secret = var.aws_secret_access_key

  # S3 bucket configuration
  endpoint    = "https://s3.us-east-1.amazonaws.com"
  bucket      = "my-casdoor-bucket"
  path_prefix = "uploads/"
  region_id   = "us-east-1"

  # Public URL domain for uploaded files (optional, for CDN)
  domain = "https://cdn.example.com"
}

# MinIO storage provider (S3-compatible)
resource "casdoor_provider" "storage_minio" {
  owner        = "my-organization"
  name         = "provider-storage-minio"
  display_name = "MinIO Storage"
  category     = "Storage"
  type         = "MinIO"

  # MinIO credentials
  client_id     = var.minio_access_key
  client_secret = var.minio_secret_key

  # MinIO server configuration
  endpoint          = "https://minio.example.com"
  intranet_endpoint = "http://minio.internal:9000"
  bucket            = "casdoor"
  path_prefix       = "uploads/"

  # Public URL domain for uploaded files
  domain = "https://minio.example.com/casdoor"
}

# Local File System storage provider
resource "casdoor_provider" "storage_local" {
  owner        = "my-organization"
  name         = "provider-storage-local"
  display_name = "Local Storage"
  category     = "Storage"
  type         = "Local File System"

  # Public URL domain where files are served
  domain = "https://casdoor.example.com/files"
}

# =============================================================================
# Payment Providers
# =============================================================================

# Stripe payment provider
resource "casdoor_provider" "payment_stripe" {
  owner        = "my-organization"
  name         = "provider-payment-stripe"
  display_name = "Stripe"
  category     = "Payment"
  type         = "Stripe"

  # Stripe API credentials
  # client_id is the Publishable Key
  client_id = var.stripe_publishable_key
  # client_secret is the Secret Key
  client_secret = var.stripe_secret_key
}

# PayPal payment provider
resource "casdoor_provider" "payment_paypal" {
  owner        = "my-organization"
  name         = "provider-payment-paypal"
  display_name = "PayPal"
  category     = "Payment"
  type         = "PayPal"

  # PayPal API credentials
  # client_id is the PayPal Client ID
  client_id = var.paypal_client_id
  # client_secret is the PayPal Secret
  client_secret = var.paypal_secret
}

# =============================================================================
# Notification Providers
# =============================================================================

# Custom HTTP notification provider (webhook)
resource "casdoor_provider" "notification_http" {
  owner        = "my-organization"
  name         = "provider-notification-http"
  display_name = "Custom Webhook"
  category     = "Notification"
  type         = "Custom HTTP"

  # HTTP method: GET or POST
  method = "POST"

  # Endpoint URL to send notifications to
  receiver = "https://webhook.example.com/notify"

  # Parameter name (sent as query param for GET, or in body for POST)
  title = "message"

  # Notification content template
  content = "Casdoor notification: %s"
}

# Slack notification provider
resource "casdoor_provider" "notification_slack" {
  owner        = "my-organization"
  name         = "provider-notification-slack"
  display_name = "Slack Notifications"
  category     = "Notification"
  type         = "Slack"

  # Slack webhook URL or bot token
  client_secret = var.slack_webhook_url

  # Slack channel ID or user ID
  receiver = "C01234567"

  # Notification content template
  content = "Casdoor notification: %s"
}
