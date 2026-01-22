# Basic permission
resource "casdoor_permission" "read_users" {
  owner        = "my-organization"
  name         = "read-users"
  display_name = "Read Users"
  description  = "Permission to read user data"
  effect       = "Allow"
  is_enabled   = true

  resources = ["users/*"]
  actions   = ["Read"]
}

# Permission for a specific role
resource "casdoor_permission" "admin_access" {
  owner        = "my-organization"
  name         = "admin-access"
  display_name = "Admin Access"
  description  = "Full administrative access"
  effect       = "Allow"
  is_enabled   = true

  roles = ["my-organization/admins"]

  resources = ["*"]
  actions   = ["Read", "Write", "Admin"]
}

# Permission with Casbin model
resource "casdoor_permission" "api_access" {
  owner         = "my-organization"
  name          = "api-access"
  display_name  = "API Access"
  description   = "Access to API endpoints"
  effect        = "Allow"
  is_enabled    = true
  model         = "my-organization/api-model"
  resource_type = "API"

  users = ["my-organization/api-user"]

  resources = ["/api/v1/*"]
  actions   = ["GET", "POST"]
}
