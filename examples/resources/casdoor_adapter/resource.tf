# Example: Adapter using the same database as Casdoor
resource "casdoor_adapter" "api-adapter" {
  owner       = "built-in"
  name        = "api-adapter-built-in"
  table       = "casbin_api_rule"
  use_same_db = true
}

# Example: Adapter with custom database configuration
resource "casdoor_adapter" "custom_db" {
  owner         = "my-organization"
  name          = "adapter-custom-db"
  type          = "Database"
  database_type = "mysql"
  host          = "db.example.com"
  port          = 3306
  user          = "casbin_user"
  password      = var.db_password
  database      = "casbin"
  table         = "casbin_rule"
  is_enabled    = true
}

# Example: Adapter with PostgreSQL
resource "casdoor_adapter" "postgres" {
  owner         = "my-organization"
  name          = "adapter-postgres"
  type          = "Database"
  database_type = "postgres"
  host          = "postgres.example.com"
  port          = 5432
  user          = "casbin_user"
  password      = var.db_password
  database      = "casbin"
  table         = "casbin_rule"
  is_enabled    = true
}
