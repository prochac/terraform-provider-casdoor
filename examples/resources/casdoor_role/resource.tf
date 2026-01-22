# Basic role
resource "casdoor_role" "developers" {
  owner        = "my-organization"
  name         = "developers"
  display_name = "Developers"
  description  = "Role for development team members"
  is_enabled   = true
}

# Role with users assigned
resource "casdoor_role" "admins" {
  owner        = "my-organization"
  name         = "admins"
  display_name = "Administrators"
  description  = "Role for system administrators"
  is_enabled   = true

  users = [
    "my-organization/admin-user",
    "my-organization/super-admin",
  ]
}

# Role hierarchy (role containing other roles)
resource "casdoor_role" "super_admins" {
  owner        = "my-organization"
  name         = "super-admins"
  display_name = "Super Administrators"
  description  = "Role with all admin privileges"
  is_enabled   = true

  roles = [
    "my-organization/admins",
    "my-organization/developers",
  ]
}
