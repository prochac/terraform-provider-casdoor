# Basic user
resource "casdoor_user" "example" {
  owner        = "my-organization"
  name         = "john.doe"
  display_name = "John Doe"
  email        = "john.doe@example.com"
  phone        = "+1234567890"
  password     = "secure-password-123"
}

# Admin user with additional attributes
resource "casdoor_user" "admin" {
  owner        = "built-in"
  name         = "admin-user"
  display_name = "Admin User"
  email        = "admin@example.com"
  is_admin     = true
  type         = "normal-user"

  affiliation = "Example Corp"
  title       = "System Administrator"
  language    = "en"
  country_code = "US"

  groups = ["admins", "developers"]
}
