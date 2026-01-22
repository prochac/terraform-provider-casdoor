# GitHub OAuth provider
resource "casdoor_provider" "github" {
  owner         = "my-organization"
  name          = "provider-github"
  display_name  = "GitHub"
  category      = "OAuth"
  type          = "GitHub"
  client_id     = var.github_client_id
  client_secret = var.github_client_secret
  scopes        = "user:email,read:user"
}

# Google OAuth provider
resource "casdoor_provider" "google" {
  owner         = "my-organization"
  name          = "provider-google"
  display_name  = "Google"
  category      = "OAuth"
  type          = "Google"
  client_id     = var.google_client_id
  client_secret = var.google_client_secret
  scopes        = "openid,profile,email"
}

# SAML provider
resource "casdoor_provider" "saml_idp" {
  owner        = "my-organization"
  name         = "provider-saml"
  display_name = "Corporate SSO"
  category     = "SAML"
  type         = "SAML"
  issuer_url   = "https://idp.example.com/saml"
  metadata     = file("path/to/saml-metadata.xml")

  enable_sign_authn_request = true
  cert                      = "my-organization/saml-cert"
}

# Email provider (SMTP)
resource "casdoor_provider" "email" {
  owner        = "my-organization"
  name         = "provider-email"
  display_name = "Email Service"
  category     = "Email"
  type         = "SMTP"
  host         = "smtp.example.com"
  port         = 587
  disable_ssl  = false

  client_id     = "smtp-username"
  client_secret = var.smtp_password

  title   = "Verification Code"
  content = "Your verification code is: %s"
}

# AWS S3 storage provider
resource "casdoor_provider" "s3" {
  owner        = "my-organization"
  name         = "provider-s3"
  display_name = "AWS S3 Storage"
  category     = "Storage"
  type         = "AWS S3"

  client_id     = var.aws_access_key_id
  client_secret = var.aws_secret_access_key
  region_id     = "us-east-1"
  bucket        = "my-casdoor-bucket"
  endpoint      = "https://s3.us-east-1.amazonaws.com"
  path_prefix   = "avatars/"
}
