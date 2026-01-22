# OpenTofu/Terraform Provider for Casdoor

This is the official OpenTofu/Terraform provider for
managing [Casdoor](https://casdoor.org/) resources. It allows you to define
Organizations, Applications, Users, and Roles as Infrastructure as Code.

## Requirements

- [OpenTofu](https://opentofu.org/) 1.11+
  or [Terraform](https://www.terraform.io/) 1.14+
- [Go](https://golang.org/doc/install) 1.25+ (for development)

## Installation

Add the following to your `main.tf`:

```hcl
terraform {
  required_providers {
    casdoor = {
      source = "registry.terraform.io/prochac/casdoor"
      # version = "0.1.0"
    }
  }
}

provider "casdoor" {
  endpoint = "http://localhost:8000" # Casdoor URL
  client_id     = "YOUR_CLIENT_ID"
  client_secret = "YOUR_CLIENT_SECRET"
  certificate   = "" # Optional: X.509 certificate for validation
}
```

## Usage Example

```hcl
resource "casdoor_organization" "example_org" {
  name         = "acme-corp"
  display_name = "Acme Corporation"
  website      = "[https://example.com](https://example.com)"
}

resource "casdoor_application" "app" {
  organization  = casdoor_organization.example_org.name
  name          = "crm-app"
  display_name  = "CRM System"
  client_id     = "generated-id"
  client_secret = "generated-secret"
  redirect_uris = [
    "[https://callback.example.com](https://callback.example.com)"
  ]
}
```

## Resource Implementation Status

| Resource     | Status | Terraform Name       | Notes                                  |
|--------------|--------|----------------------|----------------------------------------|
| Organization | 🚧     | casdoor_organization |                                        |
| Application  | 🚧     | casdoor_application  |                                        |
| User         | ❌      | casdoor_user         |                                        |
| Role         | ❌      | casdoor_role         |                                        |
| Permission   | ❌      | casdoor_permission   |                                        |
| Provider     | ❌      | casdoor_provider     | OAuth providers (Google, GitHub, etc.) |
| Token        | ❌      | casdoor_token        |                                        |
| Cert         | ❌      | casdoor_cert         |                                        |
| Model        | ❌      | casdoor_model        | Casbin model configuration             |
| Enforcer     | ❌      | casdoor_enforcer     |                                        |

Legend: ✔️ - Implemented & Tested 🚧 - Work In Progress ❌ - Not Implemented

## Debugging

1. Run the `terraform-provider-casdoor` binary with `--debug` flag.
2. Copy `TF_REATTACH_PROVIDERS` variable the binary prints to stdout.
3. Set with `export TF_...` in the shell.
4. Run `terraform init` or `tofu init`.

You can start a demo Casdoor container

```shell
docker run -d -p 8000:8000 --name casdoor casbin/casdoor-all-in-one:latest
```

And setup provider using default admin user.

```hcl
terraform {
  required_providers {
    casdoor = {
      source = "registry.terraform.io/prochac/casdoor"
    }
  }
}

provider "casdoor" {
  endpoint          = "http://localhost:8000"
  organization_name = "built-in"
  application_name  = "app-built-in"
  username          = "admin"
  password          = "123"
}
```
