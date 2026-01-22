// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithConfigure   = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

type OrganizationResource struct {
	client *casdoorsdk.Client
}

type OrganizationResourceModel struct {
	Owner                  types.String `tfsdk:"owner"`
	Name                   types.String `tfsdk:"name"`
	DisplayName            types.String `tfsdk:"display_name"`
	WebsiteURL             types.String `tfsdk:"website_url"`
	Logo                   types.String `tfsdk:"logo"`
	LogoDark               types.String `tfsdk:"logo_dark"`
	Favicon                types.String `tfsdk:"favicon"`
	PasswordType           types.String `tfsdk:"password_type"`
	PasswordSalt           types.String `tfsdk:"password_salt"`
	PasswordOptions        types.List   `tfsdk:"password_options"`
	PasswordObfuscatorType types.String `tfsdk:"password_obfuscator_type"`
	PasswordObfuscatorKey  types.String `tfsdk:"password_obfuscator_key"`
	CountryCodes           types.List   `tfsdk:"country_codes"`
	DefaultAvatar          types.String `tfsdk:"default_avatar"`
	DefaultApplication     types.String `tfsdk:"default_application"`
	Tags                   types.List   `tfsdk:"tags"`
	Languages              types.List   `tfsdk:"languages"`
	MasterPassword         types.String `tfsdk:"master_password"`
	DefaultPassword        types.String `tfsdk:"default_password"`
	MasterVerificationCode types.String `tfsdk:"master_verification_code"`
	IPWhitelist            types.String `tfsdk:"ip_whitelist"`
	InitScore              types.Int64  `tfsdk:"init_score"`
	EnableSoftDeletion     types.Bool   `tfsdk:"enable_soft_deletion"`
	IsProfilePublic        types.Bool   `tfsdk:"is_profile_public"`
	UseEmailAsUsername     types.Bool   `tfsdk:"use_email_as_username"`
	EnableTour             types.Bool   `tfsdk:"enable_tour"`
	IPRestriction          types.String `tfsdk:"ip_restriction"`
}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

func (r *OrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor organization.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The owner of the organization. Defaults to 'admin'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("admin"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the organization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the organization.",
				Required:    true,
			},
			"website_url": schema.StringAttribute{
				Description: "The website URL of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"logo": schema.StringAttribute{
				Description: "The logo URL of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"logo_dark": schema.StringAttribute{
				Description: "The dark mode logo URL of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"favicon": schema.StringAttribute{
				Description: "The favicon URL of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_type": schema.StringAttribute{
				Description: "The password hashing algorithm. Valid values: plain, bcrypt, sha256-salt, md5-salt, etc.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("bcrypt"),
			},
			"password_salt": schema.StringAttribute{
				Description: "The salt used for password hashing.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_options": schema.ListAttribute{
				Description: "Password complexity options.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"password_obfuscator_type": schema.StringAttribute{
				Description: "The password obfuscator type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_obfuscator_key": schema.StringAttribute{
				Description: "The password obfuscator key.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"country_codes": schema.ListAttribute{
				Description: "List of allowed country codes.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"default_avatar": schema.StringAttribute{
				Description: "The default avatar URL for users.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"default_application": schema.StringAttribute{
				Description: "The default application name for this organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tags": schema.ListAttribute{
				Description: "Tags for the organization.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"languages": schema.ListAttribute{
				Description: "Supported languages for the organization.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"master_password": schema.StringAttribute{
				Description: "The master password for the organization.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"default_password": schema.StringAttribute{
				Description: "The default password for new users.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"master_verification_code": schema.StringAttribute{
				Description: "The master verification code.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"ip_whitelist": schema.StringAttribute{
				Description: "IP whitelist for the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"init_score": schema.Int64Attribute{
				Description: "Initial score for new users.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"enable_soft_deletion": schema.BoolAttribute{
				Description: "Whether soft deletion is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_profile_public": schema.BoolAttribute{
				Description: "Whether user profiles are public by default.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"use_email_as_username": schema.BoolAttribute{
				Description: "Whether to use email as username.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_tour": schema.BoolAttribute{
				Description: "Whether the tour guide is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ip_restriction": schema.StringAttribute{
				Description: "IP restriction rules.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *OrganizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert list types to Go slices.
	var passwordOptions, countryCodes, tags, languages []string

	if !plan.PasswordOptions.IsNull() {
		resp.Diagnostics.Append(plan.PasswordOptions.ElementsAs(ctx, &passwordOptions, false)...)
	}
	if !plan.CountryCodes.IsNull() {
		resp.Diagnostics.Append(plan.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
	}
	if !plan.Tags.IsNull() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if !plan.Languages.IsNull() {
		resp.Diagnostics.Append(plan.Languages.ElementsAs(ctx, &languages, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	org := &casdoorsdk.Organization{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		WebsiteUrl:             plan.WebsiteURL.ValueString(),
		Logo:                   plan.Logo.ValueString(),
		LogoDark:               plan.LogoDark.ValueString(),
		Favicon:                plan.Favicon.ValueString(),
		PasswordType:           plan.PasswordType.ValueString(),
		PasswordSalt:           plan.PasswordSalt.ValueString(),
		PasswordOptions:        passwordOptions,
		PasswordObfuscatorType: plan.PasswordObfuscatorType.ValueString(),
		PasswordObfuscatorKey:  plan.PasswordObfuscatorKey.ValueString(),
		CountryCodes:           countryCodes,
		DefaultAvatar:          plan.DefaultAvatar.ValueString(),
		DefaultApplication:     plan.DefaultApplication.ValueString(),
		Tags:                   tags,
		Languages:              languages,
		MasterPassword:         plan.MasterPassword.ValueString(),
		DefaultPassword:        plan.DefaultPassword.ValueString(),
		MasterVerificationCode: plan.MasterVerificationCode.ValueString(),
		IpWhitelist:            plan.IPWhitelist.ValueString(),
		InitScore:              int(plan.InitScore.ValueInt64()),
		EnableSoftDeletion:     plan.EnableSoftDeletion.ValueBool(),
		IsProfilePublic:        plan.IsProfilePublic.ValueBool(),
		UseEmailAsUsername:     plan.UseEmailAsUsername.ValueBool(),
		EnableTour:             plan.EnableTour.ValueBool(),
		IpRestriction:          plan.IPRestriction.ValueString(),
	}

	success, err := r.client.AddOrganization(org)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Organization",
			fmt.Sprintf("Could not create organization %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Organization",
			fmt.Sprintf("Casdoor returned failure when creating organization %q", plan.Name.ValueString()),
		)
		return
	}

	// Set list values to null if empty to match plan.
	if len(passwordOptions) == 0 {
		plan.PasswordOptions = types.ListNull(types.StringType)
	}
	if len(countryCodes) == 0 {
		plan.CountryCodes = types.ListNull(types.StringType)
	}
	if len(tags) == 0 {
		plan.Tags = types.ListNull(types.StringType)
	}
	if len(languages) == 0 {
		plan.Languages = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.GetOrganization(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organization",
			fmt.Sprintf("Could not read organization %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if org == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(org.Owner)
	state.Name = types.StringValue(org.Name)
	state.DisplayName = types.StringValue(org.DisplayName)
	state.WebsiteURL = types.StringValue(org.WebsiteUrl)
	state.Logo = types.StringValue(org.Logo)
	state.LogoDark = types.StringValue(org.LogoDark)
	state.Favicon = types.StringValue(org.Favicon)
	state.PasswordType = types.StringValue(org.PasswordType)
	state.PasswordSalt = types.StringValue(org.PasswordSalt)
	state.PasswordObfuscatorType = types.StringValue(org.PasswordObfuscatorType)
	state.PasswordObfuscatorKey = types.StringValue(org.PasswordObfuscatorKey)
	state.DefaultAvatar = types.StringValue(org.DefaultAvatar)
	state.DefaultApplication = types.StringValue(org.DefaultApplication)
	state.MasterPassword = types.StringValue(org.MasterPassword)
	state.DefaultPassword = types.StringValue(org.DefaultPassword)
	state.MasterVerificationCode = types.StringValue(org.MasterVerificationCode)
	state.IPWhitelist = types.StringValue(org.IpWhitelist)
	state.InitScore = types.Int64Value(int64(org.InitScore))
	state.EnableSoftDeletion = types.BoolValue(org.EnableSoftDeletion)
	state.IsProfilePublic = types.BoolValue(org.IsProfilePublic)
	state.UseEmailAsUsername = types.BoolValue(org.UseEmailAsUsername)
	state.EnableTour = types.BoolValue(org.EnableTour)
	state.IPRestriction = types.StringValue(org.IpRestriction)

	// Convert slices to list types.
	if len(org.PasswordOptions) > 0 {
		passwordOptions, diags := types.ListValueFrom(ctx, types.StringType, org.PasswordOptions)
		resp.Diagnostics.Append(diags...)
		state.PasswordOptions = passwordOptions
	} else {
		state.PasswordOptions = types.ListNull(types.StringType)
	}

	if len(org.CountryCodes) > 0 {
		countryCodes, diags := types.ListValueFrom(ctx, types.StringType, org.CountryCodes)
		resp.Diagnostics.Append(diags...)
		state.CountryCodes = countryCodes
	} else {
		state.CountryCodes = types.ListNull(types.StringType)
	}

	if len(org.Tags) > 0 {
		tags, diags := types.ListValueFrom(ctx, types.StringType, org.Tags)
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	} else {
		state.Tags = types.ListNull(types.StringType)
	}

	if len(org.Languages) > 0 {
		languages, diags := types.ListValueFrom(ctx, types.StringType, org.Languages)
		resp.Diagnostics.Append(diags...)
		state.Languages = languages
	} else {
		state.Languages = types.ListNull(types.StringType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert list types to Go slices.
	var passwordOptions, countryCodes, tags, languages []string

	if !plan.PasswordOptions.IsNull() {
		resp.Diagnostics.Append(plan.PasswordOptions.ElementsAs(ctx, &passwordOptions, false)...)
	}
	if !plan.CountryCodes.IsNull() {
		resp.Diagnostics.Append(plan.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
	}
	if !plan.Tags.IsNull() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if !plan.Languages.IsNull() {
		resp.Diagnostics.Append(plan.Languages.ElementsAs(ctx, &languages, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	org := &casdoorsdk.Organization{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		WebsiteUrl:             plan.WebsiteURL.ValueString(),
		Logo:                   plan.Logo.ValueString(),
		LogoDark:               plan.LogoDark.ValueString(),
		Favicon:                plan.Favicon.ValueString(),
		PasswordType:           plan.PasswordType.ValueString(),
		PasswordSalt:           plan.PasswordSalt.ValueString(),
		PasswordOptions:        passwordOptions,
		PasswordObfuscatorType: plan.PasswordObfuscatorType.ValueString(),
		PasswordObfuscatorKey:  plan.PasswordObfuscatorKey.ValueString(),
		CountryCodes:           countryCodes,
		DefaultAvatar:          plan.DefaultAvatar.ValueString(),
		DefaultApplication:     plan.DefaultApplication.ValueString(),
		Tags:                   tags,
		Languages:              languages,
		MasterPassword:         plan.MasterPassword.ValueString(),
		DefaultPassword:        plan.DefaultPassword.ValueString(),
		MasterVerificationCode: plan.MasterVerificationCode.ValueString(),
		IpWhitelist:            plan.IPWhitelist.ValueString(),
		InitScore:              int(plan.InitScore.ValueInt64()),
		EnableSoftDeletion:     plan.EnableSoftDeletion.ValueBool(),
		IsProfilePublic:        plan.IsProfilePublic.ValueBool(),
		UseEmailAsUsername:     plan.UseEmailAsUsername.ValueBool(),
		EnableTour:             plan.EnableTour.ValueBool(),
		IpRestriction:          plan.IPRestriction.ValueString(),
	}

	success, err := r.client.UpdateOrganization(org)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Organization",
			fmt.Sprintf("Could not update organization %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Organization",
			fmt.Sprintf("Casdoor returned failure when updating organization %q", plan.Name.ValueString()),
		)
		return
	}

	// Set list values to null if empty to match plan.
	if len(passwordOptions) == 0 {
		plan.PasswordOptions = types.ListNull(types.StringType)
	}
	if len(countryCodes) == 0 {
		plan.CountryCodes = types.ListNull(types.StringType)
	}
	if len(tags) == 0 {
		plan.Tags = types.ListNull(types.StringType)
	}
	if len(languages) == 0 {
		plan.Languages = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := &casdoorsdk.Organization{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteOrganization(org)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Organization",
			fmt.Sprintf("Could not delete organization %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Organization",
			fmt.Sprintf("Casdoor returned failure when deleting organization %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
