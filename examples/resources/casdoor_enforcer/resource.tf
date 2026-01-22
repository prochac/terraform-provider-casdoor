# Basic enforcer
resource "casdoor_enforcer" "main" {
  owner        = "my-organization"
  name         = "enforcer-main"
  display_name = "Main Enforcer"
  description  = "Primary authorization enforcer"
  model        = "my-organization/model-rbac"
  is_enabled   = true
}

# Enforcer with adapter
resource "casdoor_enforcer" "with_adapter" {
  owner        = "my-organization"
  name         = "enforcer-db"
  display_name = "Database Enforcer"
  description  = "Enforcer with database adapter for policy storage"
  model        = "my-organization/model-acl"
  adapter      = "my-organization/adapter-db"
  is_enabled   = true
}

# Enforcer referencing the model resource
resource "casdoor_enforcer" "api" {
  owner        = "my-organization"
  name         = "enforcer-api"
  display_name = "API Enforcer"
  description  = "Enforcer for API authorization"
  model        = "${casdoor_model.acl.owner}/${casdoor_model.acl.name}"
  is_enabled   = true

  depends_on = [casdoor_model.acl]
}
