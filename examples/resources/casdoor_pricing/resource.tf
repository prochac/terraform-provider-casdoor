# Basic pricing configuration
resource "casdoor_pricing" "standard" {
  owner        = "my-organization"
  name         = "pricing-standard"
  display_name = "Standard Pricing"
  description  = "Standard pricing for the SaaS product"
  application  = "my-app"
  is_enabled   = true

  plans = [
    "plan-basic",
    "plan-premium",
  ]
}

# Pricing with trial period
resource "casdoor_pricing" "with_trial" {
  owner          = "my-organization"
  name           = "pricing-trial"
  display_name   = "Trial Pricing"
  application    = "my-app"
  trial_duration = 14
  is_enabled     = true

  plans = [
    "plan-basic",
  ]
}
