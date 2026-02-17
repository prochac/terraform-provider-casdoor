# Basic product
resource "casdoor_product" "saas_app" {
  owner        = "my-organization"
  name         = "product-saas"
  display_name = "SaaS Application"
  description  = "Cloud-based SaaS application"
  tag          = "software"
  currency     = "USD"
  price        = 49.99
  quantity     = 999
  state        = "Published"

  providers = [
    "provider-stripe",
  ]
}

# Rechargeable product
resource "casdoor_product" "credits" {
  owner        = "my-organization"
  name         = "product-credits"
  display_name = "Account Credits"
  description  = "Rechargeable account credits"
  currency     = "USD"
  price        = 10.0
  quantity     = 999
  is_recharge  = true
  state        = "Published"
}
