# Authentication Method 1: OAuth Application Credentials (recommended for production)
# To obtain these credentials, use Casdoor's Data Initialization (init_data.json)
# to inject known values from day one. See:
# https://casdoor.org/docs/deployment/data-initialization/
provider "casdoor" {
  endpoint          = "https://casdoor.example.com"
  client_id         = "your-client-id"
  client_secret     = "your-client-secret"
  certificate       = file("path/to/certificate.pem")
  organization_name = "built-in"
  application_name  = "app-built-in"
}

# Authentication Method 2: Admin Username/Password (convenient for development)
# provider "casdoor" {
#   endpoint          = "https://casdoor.example.com"
#   organization_name = "built-in"
#   application_name  = "app-built-in"
#   username          = "admin"
#   password          = var.casdoor_admin_password
# }
