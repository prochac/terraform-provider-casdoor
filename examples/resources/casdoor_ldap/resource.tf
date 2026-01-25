# Example: Basic LDAP configuration
resource "casdoor_ldap" "basic" {
  id          = "ldap-basic"
  owner       = "admin"
  server_name = "Corporate LDAP"
  host        = "ldap.example.com"
  port        = 389
  username    = "cn=admin,dc=example,dc=com"
  password    = var.ldap_password
  base_dn     = "dc=example,dc=com"
}

# Example: LDAP with SSL/TLS
resource "casdoor_ldap" "secure" {
  id                     = "ldap-secure"
  owner                  = "admin"
  server_name            = "Secure Corporate LDAP"
  host                   = "ldaps.example.com"
  port                   = 636
  enable_ssl             = true
  allow_self_signed_cert = false
  username               = "cn=admin,dc=example,dc=com"
  password               = var.ldap_password
  base_dn                = "ou=users,dc=example,dc=com"
  filter                 = "(objectClass=inetOrgPerson)"
}

# Example: LDAP with auto-sync and custom attributes
resource "casdoor_ldap" "full" {
  id                     = "ldap-full"
  owner                  = "admin"
  server_name            = "Full LDAP Configuration"
  host                   = "ldap.example.com"
  port                   = 389
  enable_ssl             = false
  username               = "cn=readonly,dc=example,dc=com"
  password               = var.ldap_password
  base_dn                = "ou=people,dc=example,dc=com"
  filter                 = "(&(objectClass=posixAccount)(uid=*))"
  filter_fields          = ["uid", "cn", "mail"]
  default_group          = "users"
  password_type          = "plain"
  auto_sync              = 60 # Sync every 60 minutes

  custom_attributes = {
    "displayName" = "name"
    "mail"        = "email"
    "mobile"      = "phone"
  }
}
