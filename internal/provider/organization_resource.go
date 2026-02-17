// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithConfigure   = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

type OrganizationResource struct {
	client *casdoorsdk.Client
}

// AccountItemModel represents an account item configuration.
type AccountItemModel struct {
	Name       types.String `tfsdk:"name"`
	Visible    types.Bool   `tfsdk:"visible"`
	ViewRule   types.String `tfsdk:"view_rule"`
	ModifyRule types.String `tfsdk:"modify_rule"`
	Regex      types.String `tfsdk:"regex"`
}

// MfaItemModel represents an MFA item configuration.
type MfaItemModel struct {
	Name types.String `tfsdk:"name"`
	Rule types.String `tfsdk:"rule"`
}

// ThemeDataModel represents the theme configuration.
type ThemeDataModel struct {
	ThemeType    types.String `tfsdk:"theme_type"`
	ColorPrimary types.String `tfsdk:"color_primary"`
	BorderRadius types.Int64  `tfsdk:"border_radius"`
	IsCompact    types.Bool   `tfsdk:"is_compact"`
	IsEnabled    types.Bool   `tfsdk:"is_enabled"`
}

// AccountItemAttrTypes returns the attribute types for AccountItemModel.
func AccountItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"visible":     types.BoolType,
		"view_rule":   types.StringType,
		"modify_rule": types.StringType,
		"regex":       types.StringType,
	}
}

// MfaItemAttrTypes returns the attribute types for MfaItemModel.
func MfaItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"rule": types.StringType,
	}
}

// ThemeDataAttrTypes returns the attribute types for ThemeDataModel.
func ThemeDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"theme_type":    types.StringType,
		"color_primary": types.StringType,
		"border_radius": types.Int64Type,
		"is_compact":    types.BoolType,
		"is_enabled":    types.BoolType,
	}
}

type OrganizationResourceModel struct {
	ID                     types.String  `tfsdk:"id"`
	Owner                  types.String  `tfsdk:"owner"`
	Name                   types.String  `tfsdk:"name"`
	CreatedTime            types.String  `tfsdk:"created_time"`
	DisplayName            types.String  `tfsdk:"display_name"`
	WebsiteURL             types.String  `tfsdk:"website_url"`
	Logo                   types.String  `tfsdk:"logo"`
	LogoDark               types.String  `tfsdk:"logo_dark"`
	Favicon                types.String  `tfsdk:"favicon"`
	HasPrivilegeConsent    types.Bool    `tfsdk:"has_privilege_consent"`
	PasswordType           types.String  `tfsdk:"password_type"`
	PasswordSalt           types.String  `tfsdk:"password_salt"`
	PasswordOptions        types.List    `tfsdk:"password_options"`
	PasswordObfuscatorType types.String  `tfsdk:"password_obfuscator_type"`
	PasswordObfuscatorKey  types.String  `tfsdk:"password_obfuscator_key"`
	PasswordExpireDays     types.Int64   `tfsdk:"password_expire_days"`
	CountryCodes           types.List    `tfsdk:"country_codes"`
	DefaultAvatar          types.String  `tfsdk:"default_avatar"`
	DefaultApplication     types.String  `tfsdk:"default_application"`
	UserTypes              types.List    `tfsdk:"user_types"`
	Tags                   types.List    `tfsdk:"tags"`
	Languages              types.List    `tfsdk:"languages"`
	ThemeData              types.Object  `tfsdk:"theme_data"`
	MasterPassword         types.String  `tfsdk:"master_password"`
	DefaultPassword        types.String  `tfsdk:"default_password"`
	MasterVerificationCode types.String  `tfsdk:"master_verification_code"`
	IPWhitelist            types.String  `tfsdk:"ip_whitelist"`
	InitScore              types.Int64   `tfsdk:"init_score"`
	EnableSoftDeletion     types.Bool    `tfsdk:"enable_soft_deletion"`
	IsProfilePublic        types.Bool    `tfsdk:"is_profile_public"`
	UseEmailAsUsername     types.Bool    `tfsdk:"use_email_as_username"`
	EnableTour             types.Bool    `tfsdk:"enable_tour"`
	DisableSignin          types.Bool    `tfsdk:"disable_signin"`
	IPRestriction          types.String  `tfsdk:"ip_restriction"`
	NavItems               types.List    `tfsdk:"nav_items"`
	UserNavItems           types.List    `tfsdk:"user_nav_items"`
	WidgetItems            types.List    `tfsdk:"widget_items"`
	MfaItems               types.List    `tfsdk:"mfa_items"`
	MfaRememberInHours     types.Int64   `tfsdk:"mfa_remember_in_hours"`
	AccountItems           types.List    `tfsdk:"account_items"`
	OrgBalance             types.Float64 `tfsdk:"org_balance"`
	UserBalance            types.Float64 `tfsdk:"user_balance"`
	BalanceCredit          types.Float64 `tfsdk:"balance_credit"`
	BalanceCurrency        types.String  `tfsdk:"balance_currency"`
	AccountMenu            types.String  `tfsdk:"account_menu"`
	DcrPolicy              types.String  `tfsdk:"dcr_policy"`
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
			"id": schema.StringAttribute{
				Description: "The ID of the organization in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
			"created_time": schema.StringAttribute{
				Description: "The time when the organization was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"has_privilege_consent": schema.BoolAttribute{
				Description: "Whether the organization has privilege consent enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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
			"password_expire_days": schema.Int64Attribute{
				Description: "Number of days before password expires. 0 means no expiration.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
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
			"user_types": schema.ListAttribute{
				Description: "List of user types allowed in the organization.",
				Optional:    true,
				ElementType: types.StringType,
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
			"theme_data": schema.SingleNestedAttribute{
				Description: "Theme configuration for the organization.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"theme_type": schema.StringAttribute{
						Description: "The theme type (e.g., 'default', 'dark').",
						Optional:    true,
					},
					"color_primary": schema.StringAttribute{
						Description: "The primary color in hex format.",
						Optional:    true,
					},
					"border_radius": schema.Int64Attribute{
						Description: "The border radius in pixels.",
						Optional:    true,
					},
					"is_compact": schema.BoolAttribute{
						Description: "Whether to use compact mode.",
						Optional:    true,
					},
					"is_enabled": schema.BoolAttribute{
						Description: "Whether the theme is enabled.",
						Optional:    true,
					},
				},
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
			"disable_signin": schema.BoolAttribute{
				Description: "Whether sign-in is disabled for the organization.",
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
			"nav_items": schema.ListAttribute{
				Description: "List of navigation items.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_nav_items": schema.ListAttribute{
				Description: "List of user navigation items.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"widget_items": schema.ListAttribute{
				Description: "List of widget items.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"mfa_items": schema.ListNestedAttribute{
				Description: "List of MFA configurations.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the MFA method.",
							Required:    true,
						},
						"rule": schema.StringAttribute{
							Description: "The rule for the MFA method.",
							Optional:    true,
						},
					},
				},
			},
			"mfa_remember_in_hours": schema.Int64Attribute{
				Description: "Number of hours to remember MFA authentication.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(12),
			},
			"account_items": schema.ListNestedAttribute{
				Description: "List of account item configurations that control user profile fields visibility and editability.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the account item (field name).",
							Required:    true,
						},
						"visible": schema.BoolAttribute{
							Description: "Whether this field is visible.",
							Optional:    true,
						},
						"view_rule": schema.StringAttribute{
							Description: "Rule for viewing this field.",
							Optional:    true,
						},
						"modify_rule": schema.StringAttribute{
							Description: "Rule for modifying this field.",
							Optional:    true,
						},
						"regex": schema.StringAttribute{
							Description: "Regex pattern for field validation.",
							Optional:    true,
						},
					},
				},
			},
			"org_balance": schema.Float64Attribute{
				Description: "The organization balance.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"user_balance": schema.Float64Attribute{
				Description: "The user balance.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"balance_credit": schema.Float64Attribute{
				Description: "The balance credit.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"balance_currency": schema.StringAttribute{
				Description: "The balance currency.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_menu": schema.StringAttribute{
				Description: "The account menu configuration.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"dcr_policy": schema.StringAttribute{
				Description: "The dynamic client registration policy.",
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
	// Use empty slices (not nil) so JSON serializes as [] instead of null.
	passwordOptions := make([]string, 0)
	countryCodes := make([]string, 0)
	userTypes := make([]string, 0)
	tags := make([]string, 0)
	languages := make([]string, 0)
	navItems := make([]string, 0)
	userNavItems := make([]string, 0)
	widgetItems := make([]string, 0)

	if !plan.PasswordOptions.IsNull() {
		resp.Diagnostics.Append(plan.PasswordOptions.ElementsAs(ctx, &passwordOptions, false)...)
	}
	if !plan.CountryCodes.IsNull() {
		resp.Diagnostics.Append(plan.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
	}
	if !plan.UserTypes.IsNull() {
		resp.Diagnostics.Append(plan.UserTypes.ElementsAs(ctx, &userTypes, false)...)
	}
	if !plan.Tags.IsNull() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if !plan.Languages.IsNull() {
		resp.Diagnostics.Append(plan.Languages.ElementsAs(ctx, &languages, false)...)
	}
	if !plan.NavItems.IsNull() {
		resp.Diagnostics.Append(plan.NavItems.ElementsAs(ctx, &navItems, false)...)
	}
	if !plan.UserNavItems.IsNull() {
		resp.Diagnostics.Append(plan.UserNavItems.ElementsAs(ctx, &userNavItems, false)...)
	}
	if !plan.WidgetItems.IsNull() {
		resp.Diagnostics.Append(plan.WidgetItems.ElementsAs(ctx, &widgetItems, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nested objects.
	var themeData *casdoorsdk.ThemeData
	if !plan.ThemeData.IsNull() {
		var themeModel ThemeDataModel
		resp.Diagnostics.Append(plan.ThemeData.As(ctx, &themeModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		themeData = &casdoorsdk.ThemeData{
			ThemeType:    themeModel.ThemeType.ValueString(),
			ColorPrimary: themeModel.ColorPrimary.ValueString(),
			BorderRadius: int(themeModel.BorderRadius.ValueInt64()),
			IsCompact:    themeModel.IsCompact.ValueBool(),
			IsEnabled:    themeModel.IsEnabled.ValueBool(),
		}
	}

	mfaItems := make([]*casdoorsdk.MfaItem, 0)
	if !plan.MfaItems.IsNull() {
		var mfaModels []MfaItemModel
		resp.Diagnostics.Append(plan.MfaItems.ElementsAs(ctx, &mfaModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, m := range mfaModels {
			mfaItems = append(mfaItems, &casdoorsdk.MfaItem{
				Name: m.Name.ValueString(),
				Rule: m.Rule.ValueString(),
			})
		}
	}

	accountItems := make([]*casdoorsdk.AccountItem, 0)
	if !plan.AccountItems.IsNull() {
		var accountModels []AccountItemModel
		resp.Diagnostics.Append(plan.AccountItems.ElementsAs(ctx, &accountModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, a := range accountModels {
			accountItems = append(accountItems, &casdoorsdk.AccountItem{
				Name:       a.Name.ValueString(),
				Visible:    a.Visible.ValueBool(),
				ViewRule:   a.ViewRule.ValueString(),
				ModifyRule: a.ModifyRule.ValueString(),
				Regex:      a.Regex.ValueString(),
			})
		}
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	org := &casdoorsdk.Organization{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		CreatedTime:            createdTime,
		DisplayName:            plan.DisplayName.ValueString(),
		WebsiteUrl:             plan.WebsiteURL.ValueString(),
		Logo:                   plan.Logo.ValueString(),
		LogoDark:               plan.LogoDark.ValueString(),
		Favicon:                plan.Favicon.ValueString(),
		HasPrivilegeConsent:    plan.HasPrivilegeConsent.ValueBool(),
		PasswordType:           plan.PasswordType.ValueString(),
		PasswordSalt:           plan.PasswordSalt.ValueString(),
		PasswordOptions:        passwordOptions,
		PasswordObfuscatorType: plan.PasswordObfuscatorType.ValueString(),
		PasswordObfuscatorKey:  plan.PasswordObfuscatorKey.ValueString(),
		PasswordExpireDays:     int(plan.PasswordExpireDays.ValueInt64()),
		CountryCodes:           countryCodes,
		DefaultAvatar:          plan.DefaultAvatar.ValueString(),
		DefaultApplication:     plan.DefaultApplication.ValueString(),
		UserTypes:              userTypes,
		Tags:                   tags,
		Languages:              languages,
		ThemeData:              themeData,
		MasterPassword:         plan.MasterPassword.ValueString(),
		DefaultPassword:        plan.DefaultPassword.ValueString(),
		MasterVerificationCode: plan.MasterVerificationCode.ValueString(),
		IpWhitelist:            plan.IPWhitelist.ValueString(),
		InitScore:              int(plan.InitScore.ValueInt64()),
		EnableSoftDeletion:     plan.EnableSoftDeletion.ValueBool(),
		IsProfilePublic:        plan.IsProfilePublic.ValueBool(),
		UseEmailAsUsername:     plan.UseEmailAsUsername.ValueBool(),
		EnableTour:             plan.EnableTour.ValueBool(),
		DisableSignin:          plan.DisableSignin.ValueBool(),
		IpRestriction:          plan.IPRestriction.ValueString(),
		NavItems:               navItems,
		UserNavItems:           userNavItems,
		WidgetItems:            widgetItems,
		MfaItems:               mfaItems,
		MfaRememberInHours:     int(plan.MfaRememberInHours.ValueInt64()),
		AccountItems:           accountItems,
		OrgBalance:             plan.OrgBalance.ValueFloat64(),
		UserBalance:            plan.UserBalance.ValueFloat64(),
		BalanceCredit:          plan.BalanceCredit.ValueFloat64(),
		BalanceCurrency:        plan.BalanceCurrency.ValueString(),
		AccountMenu:            plan.AccountMenu.ValueString(),
		DcrPolicy:              plan.DcrPolicy.ValueString(),
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

	// Read back the organization to get server-generated values like CreatedTime.
	createdOrg, err := r.client.GetOrganization(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organization",
			fmt.Sprintf("Could not read organization %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdOrg != nil {
		plan.CreatedTime = types.StringValue(createdOrg.CreatedTime)
		plan.BalanceCurrency = types.StringValue(createdOrg.BalanceCurrency)
	}

	// Set list values to null if empty to match plan.
	if len(passwordOptions) == 0 {
		plan.PasswordOptions = types.ListNull(types.StringType)
	}
	if len(countryCodes) == 0 {
		plan.CountryCodes = types.ListNull(types.StringType)
	}
	if len(userTypes) == 0 {
		plan.UserTypes = types.ListNull(types.StringType)
	}
	if len(tags) == 0 {
		plan.Tags = types.ListNull(types.StringType)
	}
	if len(languages) == 0 {
		plan.Languages = types.ListNull(types.StringType)
	}
	if len(navItems) == 0 {
		plan.NavItems = types.ListNull(types.StringType)
	}
	if len(userNavItems) == 0 {
		plan.UserNavItems = types.ListNull(types.StringType)
	}
	if len(widgetItems) == 0 {
		plan.WidgetItems = types.ListNull(types.StringType)
	}
	if len(mfaItems) == 0 {
		plan.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}
	if len(accountItems) == 0 {
		plan.AccountItems = types.ListNull(types.ObjectType{AttrTypes: AccountItemAttrTypes()})
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())

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

	// Set scalar fields.
	state.ID = types.StringValue(org.Owner + "/" + org.Name)
	state.Owner = types.StringValue(org.Owner)
	state.Name = types.StringValue(org.Name)
	state.CreatedTime = types.StringValue(org.CreatedTime)
	state.DisplayName = types.StringValue(org.DisplayName)
	state.WebsiteURL = types.StringValue(org.WebsiteUrl)
	state.Logo = types.StringValue(org.Logo)
	state.LogoDark = types.StringValue(org.LogoDark)
	state.Favicon = types.StringValue(org.Favicon)
	state.HasPrivilegeConsent = types.BoolValue(org.HasPrivilegeConsent)
	state.PasswordType = types.StringValue(org.PasswordType)
	state.PasswordSalt = types.StringValue(org.PasswordSalt)
	state.PasswordObfuscatorType = types.StringValue(org.PasswordObfuscatorType)
	state.PasswordObfuscatorKey = types.StringValue(org.PasswordObfuscatorKey)
	state.PasswordExpireDays = types.Int64Value(int64(org.PasswordExpireDays))
	state.DefaultAvatar = types.StringValue(org.DefaultAvatar)
	state.DefaultApplication = types.StringValue(org.DefaultApplication)
	// MasterPassword, DefaultPassword, MasterVerificationCode are always masked
	// by Casdoor API ("***"). Preserve real values from state; on import (when
	// state is null) fall back to empty string.
	if org.MasterPassword == "***" {
		if state.MasterPassword.IsNull() {
			state.MasterPassword = types.StringValue("")
		}
	} else {
		state.MasterPassword = types.StringValue(org.MasterPassword)
	}
	if org.DefaultPassword == "***" {
		if state.DefaultPassword.IsNull() {
			state.DefaultPassword = types.StringValue("")
		}
	} else {
		state.DefaultPassword = types.StringValue(org.DefaultPassword)
	}
	if org.MasterVerificationCode == "***" {
		if state.MasterVerificationCode.IsNull() {
			state.MasterVerificationCode = types.StringValue("")
		}
	} else {
		state.MasterVerificationCode = types.StringValue(org.MasterVerificationCode)
	}

	state.IPWhitelist = types.StringValue(org.IpWhitelist)
	state.InitScore = types.Int64Value(int64(org.InitScore))
	state.EnableSoftDeletion = types.BoolValue(org.EnableSoftDeletion)
	state.IsProfilePublic = types.BoolValue(org.IsProfilePublic)
	state.UseEmailAsUsername = types.BoolValue(org.UseEmailAsUsername)
	state.EnableTour = types.BoolValue(org.EnableTour)
	state.DisableSignin = types.BoolValue(org.DisableSignin)
	state.IPRestriction = types.StringValue(org.IpRestriction)
	state.MfaRememberInHours = types.Int64Value(int64(org.MfaRememberInHours))
	state.OrgBalance = types.Float64Value(org.OrgBalance)
	state.UserBalance = types.Float64Value(org.UserBalance)
	state.BalanceCredit = types.Float64Value(org.BalanceCredit)
	state.BalanceCurrency = types.StringValue(org.BalanceCurrency)
	state.AccountMenu = types.StringValue(org.AccountMenu)
	state.DcrPolicy = types.StringValue(org.DcrPolicy)

	// Convert string slices to list types.
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

	if len(org.UserTypes) > 0 {
		userTypes, diags := types.ListValueFrom(ctx, types.StringType, org.UserTypes)
		resp.Diagnostics.Append(diags...)
		state.UserTypes = userTypes
	} else {
		state.UserTypes = types.ListNull(types.StringType)
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

	if len(org.NavItems) > 0 {
		navItems, diags := types.ListValueFrom(ctx, types.StringType, org.NavItems)
		resp.Diagnostics.Append(diags...)
		state.NavItems = navItems
	} else {
		state.NavItems = types.ListNull(types.StringType)
	}

	if len(org.UserNavItems) > 0 {
		userNavItems, diags := types.ListValueFrom(ctx, types.StringType, org.UserNavItems)
		resp.Diagnostics.Append(diags...)
		state.UserNavItems = userNavItems
	} else {
		state.UserNavItems = types.ListNull(types.StringType)
	}

	if len(org.WidgetItems) > 0 {
		widgetItems, diags := types.ListValueFrom(ctx, types.StringType, org.WidgetItems)
		resp.Diagnostics.Append(diags...)
		state.WidgetItems = widgetItems
	} else {
		state.WidgetItems = types.ListNull(types.StringType)
	}

	// Convert ThemeData to object type.
	if org.ThemeData != nil {
		themeObj, diags := types.ObjectValue(ThemeDataAttrTypes(), map[string]attr.Value{
			"theme_type":    types.StringValue(org.ThemeData.ThemeType),
			"color_primary": types.StringValue(org.ThemeData.ColorPrimary),
			"border_radius": types.Int64Value(int64(org.ThemeData.BorderRadius)),
			"is_compact":    types.BoolValue(org.ThemeData.IsCompact),
			"is_enabled":    types.BoolValue(org.ThemeData.IsEnabled),
		})
		resp.Diagnostics.Append(diags...)
		state.ThemeData = themeObj
	} else {
		state.ThemeData = types.ObjectNull(ThemeDataAttrTypes())
	}

	// Convert MfaItems to list of objects.
	if len(org.MfaItems) > 0 {
		mfaObjList := make([]attr.Value, 0, len(org.MfaItems))
		for _, m := range org.MfaItems {
			mfaObj, diags := types.ObjectValue(MfaItemAttrTypes(), map[string]attr.Value{
				"name": types.StringValue(m.Name),
				"rule": types.StringValue(m.Rule),
			})
			resp.Diagnostics.Append(diags...)
			mfaObjList = append(mfaObjList, mfaObj)
		}
		mfaList, diags := types.ListValue(types.ObjectType{AttrTypes: MfaItemAttrTypes()}, mfaObjList)
		resp.Diagnostics.Append(diags...)
		state.MfaItems = mfaList
	} else {
		state.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}

	// Convert AccountItems to list of objects.
	if len(org.AccountItems) > 0 {
		accountObjList := make([]attr.Value, 0, len(org.AccountItems))
		for _, a := range org.AccountItems {
			accountObj, diags := types.ObjectValue(AccountItemAttrTypes(), map[string]attr.Value{
				"name":        types.StringValue(a.Name),
				"visible":     types.BoolValue(a.Visible),
				"view_rule":   types.StringValue(a.ViewRule),
				"modify_rule": types.StringValue(a.ModifyRule),
				"regex":       types.StringValue(a.Regex),
			})
			resp.Diagnostics.Append(diags...)
			accountObjList = append(accountObjList, accountObj)
		}
		accountList, diags := types.ListValue(types.ObjectType{AttrTypes: AccountItemAttrTypes()}, accountObjList)
		resp.Diagnostics.Append(diags...)
		state.AccountItems = accountList
	} else {
		state.AccountItems = types.ListNull(types.ObjectType{AttrTypes: AccountItemAttrTypes()})
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
	// Use empty slices (not nil) so JSON serializes as [] instead of null.
	passwordOptions := make([]string, 0)
	countryCodes := make([]string, 0)
	userTypes := make([]string, 0)
	tags := make([]string, 0)
	languages := make([]string, 0)
	navItems := make([]string, 0)
	userNavItems := make([]string, 0)
	widgetItems := make([]string, 0)

	if !plan.PasswordOptions.IsNull() {
		resp.Diagnostics.Append(plan.PasswordOptions.ElementsAs(ctx, &passwordOptions, false)...)
	}
	if !plan.CountryCodes.IsNull() {
		resp.Diagnostics.Append(plan.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
	}
	if !plan.UserTypes.IsNull() {
		resp.Diagnostics.Append(plan.UserTypes.ElementsAs(ctx, &userTypes, false)...)
	}
	if !plan.Tags.IsNull() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if !plan.Languages.IsNull() {
		resp.Diagnostics.Append(plan.Languages.ElementsAs(ctx, &languages, false)...)
	}
	if !plan.NavItems.IsNull() {
		resp.Diagnostics.Append(plan.NavItems.ElementsAs(ctx, &navItems, false)...)
	}
	if !plan.UserNavItems.IsNull() {
		resp.Diagnostics.Append(plan.UserNavItems.ElementsAs(ctx, &userNavItems, false)...)
	}
	if !plan.WidgetItems.IsNull() {
		resp.Diagnostics.Append(plan.WidgetItems.ElementsAs(ctx, &widgetItems, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nested objects.
	var themeData *casdoorsdk.ThemeData
	if !plan.ThemeData.IsNull() {
		var themeModel ThemeDataModel
		resp.Diagnostics.Append(plan.ThemeData.As(ctx, &themeModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		themeData = &casdoorsdk.ThemeData{
			ThemeType:    themeModel.ThemeType.ValueString(),
			ColorPrimary: themeModel.ColorPrimary.ValueString(),
			BorderRadius: int(themeModel.BorderRadius.ValueInt64()),
			IsCompact:    themeModel.IsCompact.ValueBool(),
			IsEnabled:    themeModel.IsEnabled.ValueBool(),
		}
	}

	mfaItems := make([]*casdoorsdk.MfaItem, 0)
	if !plan.MfaItems.IsNull() {
		var mfaModels []MfaItemModel
		resp.Diagnostics.Append(plan.MfaItems.ElementsAs(ctx, &mfaModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, m := range mfaModels {
			mfaItems = append(mfaItems, &casdoorsdk.MfaItem{
				Name: m.Name.ValueString(),
				Rule: m.Rule.ValueString(),
			})
		}
	}

	accountItems := make([]*casdoorsdk.AccountItem, 0)
	if !plan.AccountItems.IsNull() {
		var accountModels []AccountItemModel
		resp.Diagnostics.Append(plan.AccountItems.ElementsAs(ctx, &accountModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, a := range accountModels {
			accountItems = append(accountItems, &casdoorsdk.AccountItem{
				Name:       a.Name.ValueString(),
				Visible:    a.Visible.ValueBool(),
				ViewRule:   a.ViewRule.ValueString(),
				ModifyRule: a.ModifyRule.ValueString(),
				Regex:      a.Regex.ValueString(),
			})
		}
	}

	org := &casdoorsdk.Organization{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		CreatedTime:            plan.CreatedTime.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		WebsiteUrl:             plan.WebsiteURL.ValueString(),
		Logo:                   plan.Logo.ValueString(),
		LogoDark:               plan.LogoDark.ValueString(),
		Favicon:                plan.Favicon.ValueString(),
		HasPrivilegeConsent:    plan.HasPrivilegeConsent.ValueBool(),
		PasswordType:           plan.PasswordType.ValueString(),
		PasswordSalt:           plan.PasswordSalt.ValueString(),
		PasswordOptions:        passwordOptions,
		PasswordObfuscatorType: plan.PasswordObfuscatorType.ValueString(),
		PasswordObfuscatorKey:  plan.PasswordObfuscatorKey.ValueString(),
		PasswordExpireDays:     int(plan.PasswordExpireDays.ValueInt64()),
		CountryCodes:           countryCodes,
		DefaultAvatar:          plan.DefaultAvatar.ValueString(),
		DefaultApplication:     plan.DefaultApplication.ValueString(),
		UserTypes:              userTypes,
		Tags:                   tags,
		Languages:              languages,
		ThemeData:              themeData,
		MasterPassword:         plan.MasterPassword.ValueString(),
		DefaultPassword:        plan.DefaultPassword.ValueString(),
		MasterVerificationCode: plan.MasterVerificationCode.ValueString(),
		IpWhitelist:            plan.IPWhitelist.ValueString(),
		InitScore:              int(plan.InitScore.ValueInt64()),
		EnableSoftDeletion:     plan.EnableSoftDeletion.ValueBool(),
		IsProfilePublic:        plan.IsProfilePublic.ValueBool(),
		UseEmailAsUsername:     plan.UseEmailAsUsername.ValueBool(),
		EnableTour:             plan.EnableTour.ValueBool(),
		DisableSignin:          plan.DisableSignin.ValueBool(),
		IpRestriction:          plan.IPRestriction.ValueString(),
		NavItems:               navItems,
		UserNavItems:           userNavItems,
		WidgetItems:            widgetItems,
		MfaItems:               mfaItems,
		MfaRememberInHours:     int(plan.MfaRememberInHours.ValueInt64()),
		AccountItems:           accountItems,
		OrgBalance:             plan.OrgBalance.ValueFloat64(),
		UserBalance:            plan.UserBalance.ValueFloat64(),
		BalanceCredit:          plan.BalanceCredit.ValueFloat64(),
		BalanceCurrency:        plan.BalanceCurrency.ValueString(),
		AccountMenu:            plan.AccountMenu.ValueString(),
		DcrPolicy:              plan.DcrPolicy.ValueString(),
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
	if len(userTypes) == 0 {
		plan.UserTypes = types.ListNull(types.StringType)
	}
	if len(tags) == 0 {
		plan.Tags = types.ListNull(types.StringType)
	}
	if len(languages) == 0 {
		plan.Languages = types.ListNull(types.StringType)
	}
	if len(navItems) == 0 {
		plan.NavItems = types.ListNull(types.StringType)
	}
	if len(userNavItems) == 0 {
		plan.UserNavItems = types.ListNull(types.StringType)
	}
	if len(widgetItems) == 0 {
		plan.WidgetItems = types.ListNull(types.StringType)
	}
	if len(mfaItems) == 0 {
		plan.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}
	if len(accountItems) == 0 {
		plan.AccountItems = types.ListNull(types.ObjectType{AttrTypes: AccountItemAttrTypes()})
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())

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
	importStateOwnerName(ctx, req, resp)
}
