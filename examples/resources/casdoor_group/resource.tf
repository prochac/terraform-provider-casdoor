# Basic group
resource "casdoor_group" "engineering" {
  owner        = "my-organization"
  name         = "engineering"
  display_name = "Engineering"
  type         = "Virtual"
  is_enabled   = true
}

# Group with users
resource "casdoor_group" "backend" {
  owner        = "my-organization"
  name         = "backend"
  display_name = "Backend Team"
  type         = "Virtual"
  is_top_group = false
  is_enabled   = true

  users = [
    "my-organization/alice",
    "my-organization/bob",
  ]
}
