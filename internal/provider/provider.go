// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &CasdoorProvider{}

type CasdoorProvider struct {
	version string
}

type CasdoorProviderModel struct {
	Endpoint         types.String `tfsdk:"endpoint"`
	ClientID         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
	Certificate      types.String `tfsdk:"certificate"`
	OrganizationName types.String `tfsdk:"organization_name"`
	ApplicationName  types.String `tfsdk:"application_name"`
	// Alternative auth: admin login.
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CasdoorProvider{
			version: version,
		}
	}
}

func (p *CasdoorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "casdoor"
	resp.Version = p.version
}

func (p *CasdoorProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Casdoor provider allows you to manage Casdoor resources.

## Authentication

The provider supports two authWhen using username/passwordentication methods:

### 1. OAuth Application Credentials (recommended for production)

` + "```hcl" + `
provider "casdoor" {
  endpoint          = "https://casdoor.example.com"
  client_id         = "your-client-id"
  client_secret     = "your-client-secret"
  certificate       = file("path/to/cert.pem")
  organization_name = "built-in"
  application_name  = "app-built-in"
}
` + "```" + `

### 2. Admin Username/Password (convenient for development)

` + "```hcl" + `
provider "casdoor" {
  endpoint          = "https://casdoor.example.com"
  organization_name = "built-in"
  application_name  = "app-built-in"
  username          = "admin"
  password          = var.casdoor_admin_password
}
` + "```" + `

When using username/password, the provider will login and automatically fetch the application's OAuth credentials.
`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The Casdoor server endpoint URL (e.g., https://casdoor.example.com).",
				Required:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The OAuth2 client ID for the Casdoor application. Required if username is not set.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The OAuth2 client secret for the Casdoor application. Required if username is not set.",
				Optional:    true,
				Sensitive:   true,
			},
			"certificate": schema.StringAttribute{
				Description: "The X.509 certificate (public key) for JWT verification. Required if username is not set.",
				Optional:    true,
			},
			"organization_name": schema.StringAttribute{
				Description: "The organization name in Casdoor.",
				Required:    true,
			},
			"application_name": schema.StringAttribute{
				Description: "The application name in Casdoor.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Admin username for authentication. If set, the provider will login and fetch OAuth credentials automatically.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Admin password for authentication. Required if username is set.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *CasdoorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CasdoorProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientID, clientSecret, certificate string

	// Determine authentication method.
	useLoginAuth := !config.Username.IsNull() && config.Username.ValueString() != ""

	if useLoginAuth {
		// Validate password is provided.
		if config.Password.IsNull() || config.Password.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Password",
				"Password is required when using username authentication.",
			)
			return
		}

		// Login and fetch credentials.
		creds, err := fetchCredentialsViaLogin(
			config.Endpoint.ValueString(),
			config.OrganizationName.ValueString(),
			config.ApplicationName.ValueString(),
			config.Username.ValueString(),
			config.Password.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Authentication Failed",
				fmt.Sprintf("Failed to authenticate with Casdoor: %s", err),
			)
			return
		}

		clientID = creds.ClientID
		clientSecret = creds.ClientSecret
		certificate = creds.Certificate
	} else {
		// Use OAuth credentials directly.
		if config.ClientID.IsNull() || config.ClientID.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Client ID",
				"Either client_id or username must be provided for authentication.",
			)
			return
		}
		if config.ClientSecret.IsNull() || config.ClientSecret.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Client Secret",
				"client_secret is required when using OAuth authentication.",
			)
			return
		}
		if config.Certificate.IsNull() || config.Certificate.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Certificate",
				"certificate is required when using OAuth authentication.",
			)
			return
		}

		clientID = config.ClientID.ValueString()
		clientSecret = config.ClientSecret.ValueString()
		certificate = config.Certificate.ValueString()
	}

	client := casdoorsdk.NewClient(
		config.Endpoint.ValueString(),
		clientID,
		clientSecret,
		certificate,
		config.OrganizationName.ValueString(),
		config.ApplicationName.ValueString(),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

// appCredentials holds the OAuth credentials fetched via login.
type appCredentials struct {
	ClientID     string
	ClientSecret string
	Certificate  string
}

// fetchCredentialsViaLogin authenticates with username/password and fetches application credentials.
func fetchCredentialsViaLogin(endpoint, organization, application, username, password string) (*appCredentials, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Step 1: Login.
	loginPayload := map[string]string{
		"application":  application,
		"organization": organization,
		"username":     username,
		"password":     password,
		"type":         "login",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginResp, err := client.Post(
		endpoint+"/api/login",
		"application/json",
		strings.NewReader(string(loginBody)),
	)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer func() { _ = loginResp.Body.Close() }()

	var loginResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginResult); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}
	if loginResult.Status != "ok" {
		return nil, fmt.Errorf("login failed: %s", loginResult.Msg)
	}

	// Step 2: Get application details.
	appResp, err := client.Get(fmt.Sprintf("%s/api/get-application?id=admin/%s", endpoint, application))
	if err != nil {
		return nil, fmt.Errorf("get application request failed: %w", err)
	}
	defer func() { _ = appResp.Body.Close() }()

	var appResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			ClientID     string `json:"clientId"`
			ClientSecret string `json:"clientSecret"`
			Cert         string `json:"cert"`
		} `json:"data"`
	}
	if err := json.NewDecoder(appResp.Body).Decode(&appResult); err != nil {
		return nil, fmt.Errorf("failed to decode application response: %w", err)
	}

	// Step 3: Get certificate.
	certName := appResult.Data.Cert
	if certName == "" {
		certName = "cert-built-in"
	}

	certResp, err := client.Get(fmt.Sprintf("%s/api/get-cert?id=admin/%s", endpoint, certName))
	if err != nil {
		return nil, fmt.Errorf("get cert request failed: %w", err)
	}
	defer func() { _ = certResp.Body.Close() }()

	var certResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			Certificate string `json:"certificate"`
		} `json:"data"`
	}
	if err := json.NewDecoder(certResp.Body).Decode(&certResult); err != nil {
		return nil, fmt.Errorf("failed to decode cert response: %w", err)
	}

	certificate := certResult.Data.Certificate
	if certificate == "" {
		return nil, fmt.Errorf("certificate not found for application %s", application)
	}

	return &appCredentials{
		ClientID:     appResult.Data.ClientID,
		ClientSecret: appResult.Data.ClientSecret,
		Certificate:  certificate,
	}, nil
}

func (p *CasdoorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationResource,
		NewCertResource,
		NewEnforcerResource,
		NewIdpResource,
		NewModelResource,
		NewOrganizationResource,
		NewPermissionResource,
		NewRoleResource,
		NewTokenResource,
		NewUserResource,
	}
}

func (p *CasdoorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
