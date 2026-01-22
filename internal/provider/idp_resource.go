// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &IdpResource{}
	_ resource.ResourceWithConfigure   = &IdpResource{}
	_ resource.ResourceWithImportState = &IdpResource{}
)

type IdpResource struct {
	client *casdoorsdk.Client
}

type IdpResourceModel struct {
	Owner                  types.String `tfsdk:"owner"`
	Name                   types.String `tfsdk:"name"`
	DisplayName            types.String `tfsdk:"display_name"`
	Category               types.String `tfsdk:"category"`
	Type                   types.String `tfsdk:"type"`
	SubType                types.String `tfsdk:"sub_type"`
	Method                 types.String `tfsdk:"method"`
	ClientID               types.String `tfsdk:"client_id"`
	ClientSecret           types.String `tfsdk:"client_secret"`
	ClientID2              types.String `tfsdk:"client_id_2"`
	ClientSecret2          types.String `tfsdk:"client_secret_2"`
	Cert                   types.String `tfsdk:"cert"`
	CustomAuthURL          types.String `tfsdk:"custom_auth_url"`
	CustomTokenURL         types.String `tfsdk:"custom_token_url"`
	CustomUserInfoURL      types.String `tfsdk:"custom_user_info_url"`
	CustomLogo             types.String `tfsdk:"custom_logo"`
	Scopes                 types.String `tfsdk:"scopes"`
	UserMapping            types.Map    `tfsdk:"user_mapping"`
	Host                   types.String `tfsdk:"host"`
	Port                   types.Int64  `tfsdk:"port"`
	DisableSSL             types.Bool   `tfsdk:"disable_ssl"`
	Title                  types.String `tfsdk:"title"`
	Content                types.String `tfsdk:"content"`
	Receiver               types.String `tfsdk:"receiver"`
	RegionID               types.String `tfsdk:"region_id"`
	SignName               types.String `tfsdk:"sign_name"`
	TemplateCode           types.String `tfsdk:"template_code"`
	AppID                  types.String `tfsdk:"app_id"`
	Endpoint               types.String `tfsdk:"endpoint"`
	IntranetEndpoint       types.String `tfsdk:"intranet_endpoint"`
	Domain                 types.String `tfsdk:"domain"`
	Bucket                 types.String `tfsdk:"bucket"`
	PathPrefix             types.String `tfsdk:"path_prefix"`
	Metadata               types.String `tfsdk:"metadata"`
	IdP                    types.String `tfsdk:"idp"`
	IssuerURL              types.String `tfsdk:"issuer_url"`
	EnableSignAuthnRequest types.Bool   `tfsdk:"enable_sign_authn_request"`
	ProviderURL            types.String `tfsdk:"provider_url"`
}

func NewIdpResource() resource.Resource {
	return &IdpResource{}
}

func (r *IdpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

func (r *IdpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor identity provider (OAuth, SAML, etc.).",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this provider.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the provider.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"category": schema.StringAttribute{
				Description: "The category of the provider (e.g., 'OAuth', 'SAML', 'Email', 'SMS', 'Storage').",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the provider (e.g., 'Google', 'GitHub', 'SAML', 'AWS S3').",
				Required:    true,
			},
			"sub_type": schema.StringAttribute{
				Description: "The sub-type of the provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"method": schema.StringAttribute{
				Description: "The authentication method.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"client_id": schema.StringAttribute{
				Description: "The OAuth client ID.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"client_secret": schema.StringAttribute{
				Description: "The OAuth client secret.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"client_id_2": schema.StringAttribute{
				Description: "Secondary client ID (for some providers).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"client_secret_2": schema.StringAttribute{
				Description: "Secondary client secret (for some providers).",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"cert": schema.StringAttribute{
				Description: "The certificate name for this provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"custom_auth_url": schema.StringAttribute{
				Description: "Custom authorization URL for OAuth.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"custom_token_url": schema.StringAttribute{
				Description: "Custom token URL for OAuth.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"custom_user_info_url": schema.StringAttribute{
				Description: "Custom user info URL for OAuth.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"custom_logo": schema.StringAttribute{
				Description: "Custom logo URL for the provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"scopes": schema.StringAttribute{
				Description: "OAuth scopes (comma-separated).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"user_mapping": schema.MapAttribute{
				Description: "Mapping of provider user attributes to Casdoor user fields.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"host": schema.StringAttribute{
				Description: "Host for email/SMS providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"port": schema.Int64Attribute{
				Description: "Port for email/SMS providers.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"disable_ssl": schema.BoolAttribute{
				Description: "Whether to disable SSL.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"title": schema.StringAttribute{
				Description: "Title for email templates.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"content": schema.StringAttribute{
				Description: "Content for email/SMS templates.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"receiver": schema.StringAttribute{
				Description: "Receiver for notifications.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"region_id": schema.StringAttribute{
				Description: "Region ID for cloud providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"sign_name": schema.StringAttribute{
				Description: "Sign name for SMS providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"template_code": schema.StringAttribute{
				Description: "Template code for SMS providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"app_id": schema.StringAttribute{
				Description: "App ID for certain providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint for storage/cloud providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"intranet_endpoint": schema.StringAttribute{
				Description: "Intranet endpoint for storage providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"domain": schema.StringAttribute{
				Description: "Domain for the provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"bucket": schema.StringAttribute{
				Description: "Bucket name for storage providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"path_prefix": schema.StringAttribute{
				Description: "Path prefix for storage providers.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"metadata": schema.StringAttribute{
				Description: "Provider metadata (e.g., SAML metadata XML).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"idp": schema.StringAttribute{
				Description: "Identity provider identifier.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"issuer_url": schema.StringAttribute{
				Description: "SAML/OIDC issuer URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"enable_sign_authn_request": schema.BoolAttribute{
				Description: "Whether to sign SAML authentication requests.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"provider_url": schema.StringAttribute{
				Description: "The provider URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *IdpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*casdoorsdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *casdoorsdk.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *IdpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IdpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userMapping map[string]string
	if !plan.UserMapping.IsNull() {
		userMapping = make(map[string]string)
		resp.Diagnostics.Append(plan.UserMapping.ElementsAs(ctx, &userMapping, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	provider := &casdoorsdk.Provider{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		Category:               plan.Category.ValueString(),
		Type:                   plan.Type.ValueString(),
		SubType:                plan.SubType.ValueString(),
		Method:                 plan.Method.ValueString(),
		ClientId:               plan.ClientID.ValueString(),
		ClientSecret:           plan.ClientSecret.ValueString(),
		ClientId2:              plan.ClientID2.ValueString(),
		ClientSecret2:          plan.ClientSecret2.ValueString(),
		Cert:                   plan.Cert.ValueString(),
		CustomAuthUrl:          plan.CustomAuthURL.ValueString(),
		CustomTokenUrl:         plan.CustomTokenURL.ValueString(),
		CustomUserInfoUrl:      plan.CustomUserInfoURL.ValueString(),
		CustomLogo:             plan.CustomLogo.ValueString(),
		Scopes:                 plan.Scopes.ValueString(),
		UserMapping:            userMapping,
		Host:                   plan.Host.ValueString(),
		Port:                   int(plan.Port.ValueInt64()),
		DisableSsl:             plan.DisableSSL.ValueBool(),
		Title:                  plan.Title.ValueString(),
		Content:                plan.Content.ValueString(),
		Receiver:               plan.Receiver.ValueString(),
		RegionId:               plan.RegionID.ValueString(),
		SignName:               plan.SignName.ValueString(),
		TemplateCode:           plan.TemplateCode.ValueString(),
		AppId:                  plan.AppID.ValueString(),
		Endpoint:               plan.Endpoint.ValueString(),
		IntranetEndpoint:       plan.IntranetEndpoint.ValueString(),
		Domain:                 plan.Domain.ValueString(),
		Bucket:                 plan.Bucket.ValueString(),
		PathPrefix:             plan.PathPrefix.ValueString(),
		Metadata:               plan.Metadata.ValueString(),
		IdP:                    plan.IdP.ValueString(),
		IssuerUrl:              plan.IssuerURL.ValueString(),
		EnableSignAuthnRequest: plan.EnableSignAuthnRequest.ValueBool(),
		ProviderUrl:            plan.ProviderURL.ValueString(),
	}

	success, err := r.client.AddProvider(provider)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Provider",
			fmt.Sprintf("Could not create provider %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Provider",
			fmt.Sprintf("Casdoor returned failure when creating provider %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IdpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IdpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := r.client.GetProvider(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Provider",
			fmt.Sprintf("Could not read provider %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if provider == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(provider.Owner)
	state.Name = types.StringValue(provider.Name)
	state.DisplayName = types.StringValue(provider.DisplayName)
	state.Category = types.StringValue(provider.Category)
	state.Type = types.StringValue(provider.Type)
	state.SubType = types.StringValue(provider.SubType)
	state.Method = types.StringValue(provider.Method)
	state.ClientID = types.StringValue(provider.ClientId)
	state.ClientSecret = types.StringValue(provider.ClientSecret)
	state.ClientID2 = types.StringValue(provider.ClientId2)
	state.ClientSecret2 = types.StringValue(provider.ClientSecret2)
	state.Cert = types.StringValue(provider.Cert)
	state.CustomAuthURL = types.StringValue(provider.CustomAuthUrl)
	state.CustomTokenURL = types.StringValue(provider.CustomTokenUrl)
	state.CustomUserInfoURL = types.StringValue(provider.CustomUserInfoUrl)
	state.CustomLogo = types.StringValue(provider.CustomLogo)
	state.Scopes = types.StringValue(provider.Scopes)
	state.Host = types.StringValue(provider.Host)
	state.Port = types.Int64Value(int64(provider.Port))
	state.DisableSSL = types.BoolValue(provider.DisableSsl)
	state.Title = types.StringValue(provider.Title)
	state.Content = types.StringValue(provider.Content)
	state.Receiver = types.StringValue(provider.Receiver)
	state.RegionID = types.StringValue(provider.RegionId)
	state.SignName = types.StringValue(provider.SignName)
	state.TemplateCode = types.StringValue(provider.TemplateCode)
	state.AppID = types.StringValue(provider.AppId)
	state.Endpoint = types.StringValue(provider.Endpoint)
	state.IntranetEndpoint = types.StringValue(provider.IntranetEndpoint)
	state.Domain = types.StringValue(provider.Domain)
	state.Bucket = types.StringValue(provider.Bucket)
	state.PathPrefix = types.StringValue(provider.PathPrefix)
	state.Metadata = types.StringValue(provider.Metadata)
	state.IdP = types.StringValue(provider.IdP)
	state.IssuerURL = types.StringValue(provider.IssuerUrl)
	state.EnableSignAuthnRequest = types.BoolValue(provider.EnableSignAuthnRequest)
	state.ProviderURL = types.StringValue(provider.ProviderUrl)

	if len(provider.UserMapping) > 0 {
		userMapping, diags := types.MapValueFrom(ctx, types.StringType, provider.UserMapping)
		resp.Diagnostics.Append(diags...)
		state.UserMapping = userMapping
	} else {
		state.UserMapping = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IdpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IdpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userMapping map[string]string
	if !plan.UserMapping.IsNull() {
		userMapping = make(map[string]string)
		resp.Diagnostics.Append(plan.UserMapping.ElementsAs(ctx, &userMapping, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	provider := &casdoorsdk.Provider{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		Category:               plan.Category.ValueString(),
		Type:                   plan.Type.ValueString(),
		SubType:                plan.SubType.ValueString(),
		Method:                 plan.Method.ValueString(),
		ClientId:               plan.ClientID.ValueString(),
		ClientSecret:           plan.ClientSecret.ValueString(),
		ClientId2:              plan.ClientID2.ValueString(),
		ClientSecret2:          plan.ClientSecret2.ValueString(),
		Cert:                   plan.Cert.ValueString(),
		CustomAuthUrl:          plan.CustomAuthURL.ValueString(),
		CustomTokenUrl:         plan.CustomTokenURL.ValueString(),
		CustomUserInfoUrl:      plan.CustomUserInfoURL.ValueString(),
		CustomLogo:             plan.CustomLogo.ValueString(),
		Scopes:                 plan.Scopes.ValueString(),
		UserMapping:            userMapping,
		Host:                   plan.Host.ValueString(),
		Port:                   int(plan.Port.ValueInt64()),
		DisableSsl:             plan.DisableSSL.ValueBool(),
		Title:                  plan.Title.ValueString(),
		Content:                plan.Content.ValueString(),
		Receiver:               plan.Receiver.ValueString(),
		RegionId:               plan.RegionID.ValueString(),
		SignName:               plan.SignName.ValueString(),
		TemplateCode:           plan.TemplateCode.ValueString(),
		AppId:                  plan.AppID.ValueString(),
		Endpoint:               plan.Endpoint.ValueString(),
		IntranetEndpoint:       plan.IntranetEndpoint.ValueString(),
		Domain:                 plan.Domain.ValueString(),
		Bucket:                 plan.Bucket.ValueString(),
		PathPrefix:             plan.PathPrefix.ValueString(),
		Metadata:               plan.Metadata.ValueString(),
		IdP:                    plan.IdP.ValueString(),
		IssuerUrl:              plan.IssuerURL.ValueString(),
		EnableSignAuthnRequest: plan.EnableSignAuthnRequest.ValueBool(),
		ProviderUrl:            plan.ProviderURL.ValueString(),
	}

	success, err := r.client.UpdateProvider(provider)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Provider",
			fmt.Sprintf("Could not update provider %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Provider",
			fmt.Sprintf("Casdoor returned failure when updating provider %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IdpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IdpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider := &casdoorsdk.Provider{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteProvider(provider)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Provider",
			fmt.Sprintf("Could not delete provider %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Provider",
			fmt.Sprintf("Casdoor returned failure when deleting provider %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *IdpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
