# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OpenTofu/Terraform Provider for [Casdoor](https://casdoor.org), an Identity and Access Management (IAM) platform. Maps Terraform resources to the [Casdoor Go SDK](https://github.com/casdoor/casdoor-go-sdk) APIs.

## Tech Stack

- **Language:** Go 1.25+
- **Framework:** [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) (NOT SDKv2)
- **Client:** [Casdoor Go SDK](https://github.com/casdoor/casdoor-go-sdk) (uses
  [my fork](https://github.com/prochac/casdoor-go-sdk) via `replace` in `go.mod`)
- **Testing:** [Testcontainers for Go](https://github.com/testcontainers/testcontainers-go) (integration tests via Docker)

## Architectural Guidelines

### 1. Provider Design

- **Directory Structure:** Follow standard `internal/provider`, `examples/`,
  `docs/` layout.
- **Naming:** All resources must be prefixed with `casdoor_` (e.g.,
  `casdoor_organization`).
- **Client Handling:** The `provider.Configure` method must initialize the
  Casdoor SDK client and store it in `resp.DataSourceData` and
  `resp.ResourceData`.
- **Data Handling:**
    - **Strictly** use Terraform Framework types (`types.String`, `types.Bool`,
      etc.) in struct models.
    - **Never** use `d.Set()` or `schema.Resource`. Use `resp.State.Set`,
      `resp.Diagnostics`, and struct-based plans/states.
    - Map Terraform models to Casdoor SDK structs inside `Create`, `Read`,
      `Update` methods.

### 2. Resource Implementation Patterns

- **Types:** Strictly use Framework types (`types.String`, `types.Int64`) in
  model structs. Do not use Go primitives in the model.
- **CRUD Operations:**
    - **Create:** Convert Plan model -> SDK struct -> Call SDK `AddX`.
    - **Read:** Call SDK `GetX` -> Convert SDK struct -> State model.
    - **Update:** Convert Plan model -> SDK struct -> Call SDK `UpdateX`.
    - **Delete:** Call SDK `DeleteX`.
- **Import:** Implement `ImportState` interface for all resources.

### 3. Testing Strategy

- **Acceptance Tests Only:** Focus on `resource.Test` (Acc tests).
    - **Unit Tests:** minimal, only for complex logic.
- **Testcontainers:** The test harness (`testhelper_test.go`) uses
  `testcontainers-go` to spin up a `casdoor/casdoor:latest` container,
  dynamically maps the container port, and configures the provider factory
  against the ephemeral endpoint.
- **Data Isolation:** Use random suffixes for resource names in tests (e.g.,
  `acctest.RandStringFromCharSet`) to prevent collisions.

## Build & Run Commands

```bash
go build .                                    # Build provider
golangci-lint run                             # Lint
go generate ./...                             # Generate docs (terraform-plugin-docs)
CASDOOR_TEST_LOCAL=1 TF_ACC=1 go test -v -timeout 30m ./...                 # Run all tests (requires Docker)
CASDOOR_TEST_LOCAL=1 TF_ACC=1 go test -v -timeout 30m ./internal/provider -run TestAccOrganization  # Run single test
```

## Code Style

- Always handle errors (check `diag.Diagnostics`).
- Use `t.Helper()` in test helpers.
- Run `gofumpt` before committing.
