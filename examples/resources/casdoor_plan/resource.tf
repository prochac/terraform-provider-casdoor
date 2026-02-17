# Basic subscription plan
resource "casdoor_plan" "basic" {
  owner        = "my-organization"
  name         = "plan-basic"
  display_name = "Basic Plan"
  description  = "Basic tier with limited features"
  price        = 9.99
  currency     = "USD"
  period       = "Monthly"
  is_enabled   = true
}

# Premium plan with a role
resource "casdoor_plan" "premium" {
  owner        = "my-organization"
  name         = "plan-premium"
  display_name = "Premium Plan"
  description  = "Premium tier with all features"
  price        = 29.99
  currency     = "USD"
  period       = "Monthly"
  is_enabled   = true
  role         = "my-organization/premium-users"
}
