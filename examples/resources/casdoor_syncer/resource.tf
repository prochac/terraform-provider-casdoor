# Database syncer for user synchronization
resource "casdoor_syncer" "user_sync" {
  owner         = "my-organization"
  name          = "syncer-users"
  organization  = "my-organization"
  type          = "Database"
  host          = "db.example.com"
  port          = 3306
  user          = "sync_user"
  password      = var.db_password
  database_type = "mysql"
  database      = "users_db"
  table         = "users"
  sync_interval = 5
  is_read_only  = true
  is_enabled    = true
}

# Syncer with SSH tunnel
resource "casdoor_syncer" "remote_sync" {
  owner         = "my-organization"
  name          = "syncer-remote"
  organization  = "my-organization"
  type          = "Database"
  host          = "internal-db.local"
  port          = 5432
  user          = "sync_user"
  password      = var.db_password
  database_type = "postgres"
  database      = "users_db"
  table         = "users"
  ssh_type      = "password"
  ssh_host      = "bastion.example.com"
  ssh_port      = 22
  ssh_user      = "tunnel"
  ssh_password  = var.ssh_password
  sync_interval = 10
  is_enabled    = true
}
