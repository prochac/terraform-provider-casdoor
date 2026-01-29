// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &ApplicationResource{}
	_ resource.ResourceWithConfigure   = &ApplicationResource{}
	_ resource.ResourceWithImportState = &ApplicationResource{}
)

// ApplicationResource defines the resource implementation.
type ApplicationResource struct {
	client *casdoorsdk.Client
}

// ProviderItemModel represents a provider item configuration.
type ProviderItemModel struct {
	Owner        types.String `tfsdk:"owner"`
	Name         types.String `tfsdk:"name"`
	CanSignUp    types.Bool   `tfsdk:"can_sign_up"`
	CanSignIn    types.Bool   `tfsdk:"can_sign_in"`
	CanUnlink    types.Bool   `tfsdk:"can_unlink"`
	Prompted     types.Bool   `tfsdk:"prompted"`
	Rule         types.String `tfsdk:"rule"`
	SignupGroup  types.String `tfsdk:"signup_group"`
	CountryCodes types.List   `tfsdk:"country_codes"`
}

// JwtItemModel represents a JWT item configuration.
type JwtItemModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
	Type  types.String `tfsdk:"type"`
}

// SigninMethodModel represents a signin method configuration.
type SigninMethodModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Rule        types.String `tfsdk:"rule"`
}

// SignupItemModel represents a signup item configuration.
type SignupItemModel struct {
	Name        types.String `tfsdk:"name"`
	Visible     types.Bool   `tfsdk:"visible"`
	Required    types.Bool   `tfsdk:"required"`
	Prompted    types.Bool   `tfsdk:"prompted"`
	Type        types.String `tfsdk:"type"`
	Rule        types.String `tfsdk:"rule"`
	Label       types.String `tfsdk:"label"`
	Placeholder types.String `tfsdk:"placeholder"`
	Regex       types.String `tfsdk:"regex"`
	CustomCSS   types.String `tfsdk:"custom_css"`
	Options     types.List   `tfsdk:"options"`
}

// SigninItemModel represents a signin item configuration.
type SigninItemModel struct {
	Name        types.String `tfsdk:"name"`
	Visible     types.Bool   `tfsdk:"visible"`
	IsCustom    types.Bool   `tfsdk:"is_custom"`
	Label       types.String `tfsdk:"label"`
	Placeholder types.String `tfsdk:"placeholder"`
	Rule        types.String `tfsdk:"rule"`
	CustomCSS   types.String `tfsdk:"custom_css"`
}

// SamlItemModel represents a SAML attribute configuration.
type SamlItemModel struct {
	Name       types.String `tfsdk:"name"`
	NameFormat types.String `tfsdk:"name_format"`
	Value      types.String `tfsdk:"value"`
}

// ProviderItemAttrTypes returns the attribute types for ProviderItemModel.
func ProviderItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"owner":         types.StringType,
		"name":          types.StringType,
		"can_sign_up":   types.BoolType,
		"can_sign_in":   types.BoolType,
		"can_unlink":    types.BoolType,
		"prompted":      types.BoolType,
		"rule":          types.StringType,
		"signup_group":  types.StringType,
		"country_codes": types.ListType{ElemType: types.StringType},
	}
}

// JwtItemAttrTypes returns the attribute types for JwtItemModel.
func JwtItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
		"type":  types.StringType,
	}
}

// SigninMethodAttrTypes returns the attribute types for SigninMethodModel.
func SigninMethodAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"display_name": types.StringType,
		"rule":         types.StringType,
	}
}

// SignupItemAttrTypes returns the attribute types for SignupItemModel.
func SignupItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"visible":     types.BoolType,
		"required":    types.BoolType,
		"prompted":    types.BoolType,
		"type":        types.StringType,
		"rule":        types.StringType,
		"label":       types.StringType,
		"placeholder": types.StringType,
		"regex":       types.StringType,
		"custom_css":  types.StringType,
		"options":     types.ListType{ElemType: types.StringType},
	}
}

// SigninItemAttrTypes returns the attribute types for SigninItemModel.
func SigninItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"visible":     types.BoolType,
		"is_custom":   types.BoolType,
		"label":       types.StringType,
		"placeholder": types.StringType,
		"rule":        types.StringType,
		"custom_css":  types.StringType,
	}
}

// SamlItemAttrTypes returns the attribute types for SamlItemModel.
func SamlItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"name_format": types.StringType,
		"value":       types.StringType,
	}
}

// ApplicationResourceModel describes the resource data model.
type ApplicationResourceModel struct {
	// Core fields
	Owner        types.String `tfsdk:"owner"`
	Name         types.String `tfsdk:"name"`
	CreatedTime  types.String `tfsdk:"created_time"`
	DisplayName  types.String `tfsdk:"display_name"`
	Title        types.String `tfsdk:"title"`
	Favicon      types.String `tfsdk:"favicon"`
	Logo         types.String `tfsdk:"logo"`
	HomepageURL  types.String `tfsdk:"homepage_url"`
	Description  types.String `tfsdk:"description"`
	Organization types.String `tfsdk:"organization"`
	Cert         types.String `tfsdk:"cert"`
	DefaultGroup types.String `tfsdk:"default_group"`

	// Enable/disable flags
	EnablePassword               types.Bool `tfsdk:"enable_password"`
	EnableSignUp                 types.Bool `tfsdk:"enable_sign_up"`
	EnableSigninSession          types.Bool `tfsdk:"enable_signin_session"`
	EnableAutoSignin             types.Bool `tfsdk:"enable_auto_signin"`
	EnableCodeSignin             types.Bool `tfsdk:"enable_code_signin"`
	EnableExclusiveSignin        types.Bool `tfsdk:"enable_exclusive_signin"`
	EnableSamlCompress           types.Bool `tfsdk:"enable_saml_compress"`
	EnableSamlC14n10             types.Bool `tfsdk:"enable_saml_c14n10"`
	EnableSamlPostBinding        types.Bool `tfsdk:"enable_saml_post_binding"`
	EnableSamlAssertionSignature types.Bool `tfsdk:"enable_saml_assertion_signature"`
	DisableSamlAttributes        types.Bool `tfsdk:"disable_saml_attributes"`
	UseEmailAsSamlNameId         types.Bool `tfsdk:"use_email_as_saml_name_id"`
	EnableWebAuthn               types.Bool `tfsdk:"enable_web_authn"`
	EnableLinkWithEmail          types.Bool `tfsdk:"enable_link_with_email"`
	DisableSignin                types.Bool `tfsdk:"disable_signin"`
	IsShared                     types.Bool `tfsdk:"is_shared"`

	// OAuth/Token
	ClientID             types.String  `tfsdk:"client_id"`
	ClientSecret         types.String  `tfsdk:"client_secret"`
	RedirectURIs         types.List    `tfsdk:"redirect_uris"`
	TokenFormat          types.String  `tfsdk:"token_format"`
	TokenSigningMethod   types.String  `tfsdk:"token_signing_method"`
	TokenFields          types.List    `tfsdk:"token_fields"`
	TokenAttributes      types.List    `tfsdk:"token_attributes"`
	ExpireInHours        types.Float64 `tfsdk:"expire_in_hours"`
	RefreshExpireInHours types.Float64 `tfsdk:"refresh_expire_in_hours"`
	CookieExpireInHours  types.Int64   `tfsdk:"cookie_expire_in_hours"`
	GrantTypes           types.List    `tfsdk:"grant_types"`

	// SAML
	SamlReplyUrl      types.String `tfsdk:"saml_reply_url"`
	SamlHashAlgorithm types.String `tfsdk:"saml_hash_algorithm"`
	SamlAttributes    types.List   `tfsdk:"saml_attributes"`

	// URLs
	SignupUrl      types.String `tfsdk:"signup_url"`
	SigninUrl      types.String `tfsdk:"signin_url"`
	ForgetUrl      types.String `tfsdk:"forget_url"`
	AffiliationUrl types.String `tfsdk:"affiliation_url"`

	// UI/Appearance
	HeaderHtml              types.String `tfsdk:"header_html"`
	FooterHtml              types.String `tfsdk:"footer_html"`
	SignupHtml              types.String `tfsdk:"signup_html"`
	SigninHtml              types.String `tfsdk:"signin_html"`
	FormCss                 types.String `tfsdk:"form_css"`
	FormCssMobile           types.String `tfsdk:"form_css_mobile"`
	FormOffset              types.Int64  `tfsdk:"form_offset"`
	FormSideHtml            types.String `tfsdk:"form_side_html"`
	FormBackgroundUrl       types.String `tfsdk:"form_background_url"`
	FormBackgroundUrlMobile types.String `tfsdk:"form_background_url_mobile"`
	ThemeData               types.Object `tfsdk:"theme_data"`

	// Security
	IpRestriction          types.String `tfsdk:"ip_restriction"`
	IpWhitelist            types.String `tfsdk:"ip_whitelist"`
	FailedSigninLimit      types.Int64  `tfsdk:"failed_signin_limit"`
	FailedSigninFrozenTime types.Int64  `tfsdk:"failed_signin_frozen_time"`

	// Misc
	CodeResendTimeout    types.Int64  `tfsdk:"code_resend_timeout"`
	OrgChoiceMode        types.String `tfsdk:"org_choice_mode"`
	TermsOfUse           types.String `tfsdk:"terms_of_use"`
	Tags                 types.List   `tfsdk:"tags"`
	CertPublicKey        types.String `tfsdk:"cert_public_key"`
	ForcedRedirectOrigin types.String `tfsdk:"forced_redirect_origin"`
	Order                types.Int64  `tfsdk:"order"`

	// Complex nested types
	Providers     types.List `tfsdk:"providers"`
	SigninMethods types.List `tfsdk:"signin_methods"`
	SignupItems   types.List `tfsdk:"signup_items"`
	SigninItems   types.List `tfsdk:"signin_items"`
}

// NewApplicationResource creates a new Application resource.
func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

func (r *ApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor application.",
		Attributes: map[string]schema.Attribute{
			// Core fields
			"owner": schema.StringAttribute{
				Description: "The owner of the application. Defaults to 'admin'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("admin"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the application.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the application was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the application.",
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: "The title of the application.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"favicon": schema.StringAttribute{
				Description: "The favicon URL of the application.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"logo": schema.StringAttribute{
				Description: "The logo URL of the application.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"homepage_url": schema.StringAttribute{
				Description: "The homepage URL of the application.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "The description of the application.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"organization": schema.StringAttribute{
				Description: "The organization that owns this application.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cert": schema.StringAttribute{
				Description: "The certificate name used for signing tokens.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"default_group": schema.StringAttribute{
				Description: "The default group for new users.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},

			// Enable/disable flags
			"enable_password": schema.BoolAttribute{
				Description: "Whether password login is enabled. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"enable_sign_up": schema.BoolAttribute{
				Description: "Whether sign up is enabled. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"enable_signin_session": schema.BoolAttribute{
				Description: "Whether signin session is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_auto_signin": schema.BoolAttribute{
				Description: "Whether auto signin is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_code_signin": schema.BoolAttribute{
				Description: "Whether code signin is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_exclusive_signin": schema.BoolAttribute{
				Description: "Whether exclusive signin is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_saml_compress": schema.BoolAttribute{
				Description: "Whether SAML response compression is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_saml_c14n10": schema.BoolAttribute{
				Description: "Whether SAML C14N 1.0 is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_saml_post_binding": schema.BoolAttribute{
				Description: "Whether SAML POST binding is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_saml_assertion_signature": schema.BoolAttribute{
				Description: "Whether SAML assertion signature is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_saml_attributes": schema.BoolAttribute{
				Description: "Whether SAML attributes are disabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"use_email_as_saml_name_id": schema.BoolAttribute{
				Description: "Whether to use email as SAML NameID.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_web_authn": schema.BoolAttribute{
				Description: "Whether WebAuthn is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_link_with_email": schema.BoolAttribute{
				Description: "Whether linking with email is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_signin": schema.BoolAttribute{
				Description: "Whether signin is disabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_shared": schema.BoolAttribute{
				Description: "Whether the application is shared across organizations.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},

			// OAuth/Token
			"client_id": schema.StringAttribute{
				Description: "The OAuth client ID. Generated by Casdoor.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				Description: "The OAuth client secret. Generated by Casdoor.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"redirect_uris": schema.ListAttribute{
				Description: "The allowed redirect URIs for OAuth.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"token_format": schema.StringAttribute{
				Description: "The token format. Valid values: JWT, JWT-Empty.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("JWT"),
			},
			"token_signing_method": schema.StringAttribute{
				Description: "The token signing method (e.g., RS256).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"token_fields": schema.ListAttribute{
				Description: "Additional fields to include in the token.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"token_attributes": schema.ListNestedAttribute{
				Description: "Token attribute mappings.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the token attribute.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the token attribute.",
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the token attribute.",
							Optional:    true,
						},
					},
				},
			},
			"expire_in_hours": schema.Float64Attribute{
				Description: "The access token expiration time in hours. Defaults to 168 (7 days).",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(168),
			},
			"refresh_expire_in_hours": schema.Float64Attribute{
				Description: "The refresh token expiration time in hours. Defaults to 168 (7 days).",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(168),
			},
			"cookie_expire_in_hours": schema.Int64Attribute{
				Description: "The cookie expiration time in hours.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"grant_types": schema.ListAttribute{
				Description: "The allowed OAuth grant types.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},

			// SAML
			"saml_reply_url": schema.StringAttribute{
				Description: "The SAML reply URL (Assertion Consumer Service URL).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"saml_hash_algorithm": schema.StringAttribute{
				Description: "The SAML hash algorithm.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"saml_attributes": schema.ListNestedAttribute{
				Description: "SAML attribute mappings.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the SAML attribute.",
							Required:    true,
						},
						"name_format": schema.StringAttribute{
							Description: "The name format of the SAML attribute.",
							Optional:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value expression for the SAML attribute.",
							Optional:    true,
						},
					},
				},
			},

			// URLs
			"signup_url": schema.StringAttribute{
				Description: "Custom signup URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"signin_url": schema.StringAttribute{
				Description: "Custom signin URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"forget_url": schema.StringAttribute{
				Description: "Custom forgot password URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"affiliation_url": schema.StringAttribute{
				Description: "Affiliation URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},

			// UI/Appearance
			"header_html": schema.StringAttribute{
				Description: "Custom HTML for the page header.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"footer_html": schema.StringAttribute{
				Description: "Custom HTML for the page footer.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"signup_html": schema.StringAttribute{
				Description: "Custom HTML for the signup page.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"signin_html": schema.StringAttribute{
				Description: "Custom HTML for the signin page.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"form_css": schema.StringAttribute{
				Description: "Custom CSS for the form.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"form_css_mobile": schema.StringAttribute{
				Description: "Custom CSS for the form on mobile devices.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"form_offset": schema.Int64Attribute{
				Description: "Form offset position.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"form_side_html": schema.StringAttribute{
				Description: "Custom HTML for the form side panel.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"form_background_url": schema.StringAttribute{
				Description: "Background image URL for the form.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"form_background_url_mobile": schema.StringAttribute{
				Description: "Background image URL for the form on mobile devices.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"theme_data": schema.SingleNestedAttribute{
				Description: "Theme configuration for the application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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

			// Security
			"ip_restriction": schema.StringAttribute{
				Description: "IP restriction rules.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ip_whitelist": schema.StringAttribute{
				Description: "IP whitelist.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"failed_signin_limit": schema.Int64Attribute{
				Description: "Maximum number of failed signin attempts before lockout.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"failed_signin_frozen_time": schema.Int64Attribute{
				Description: "Duration in minutes to freeze account after exceeding failed signin limit.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(15),
			},

			// Misc
			"code_resend_timeout": schema.Int64Attribute{
				Description: "The code resend timeout in seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"org_choice_mode": schema.StringAttribute{
				Description: "Organization choice mode for multi-org applications.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"terms_of_use": schema.StringAttribute{
				Description: "Terms of use URL or text.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tags": schema.ListAttribute{
				Description: "Tags for the application.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"cert_public_key": schema.StringAttribute{
				Description: "The public key of the certificate. Computed from the cert field.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"forced_redirect_origin": schema.StringAttribute{
				Description: "Forced redirect origin for OAuth callbacks.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"order": schema.Int64Attribute{
				Description: "Display order of the application.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},

			// Complex nested types
			"providers": schema.ListNestedAttribute{
				Description: "List of identity providers configured for the application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"owner": schema.StringAttribute{
							Description: "The owner of the provider.",
							Optional:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the provider.",
							Required:    true,
						},
						"can_sign_up": schema.BoolAttribute{
							Description: "Whether users can sign up with this provider.",
							Optional:    true,
						},
						"can_sign_in": schema.BoolAttribute{
							Description: "Whether users can sign in with this provider.",
							Optional:    true,
						},
						"can_unlink": schema.BoolAttribute{
							Description: "Whether users can unlink this provider.",
							Optional:    true,
						},
						"prompted": schema.BoolAttribute{
							Description: "Whether this provider is prompted during login.",
							Optional:    true,
						},
						"rule": schema.StringAttribute{
							Description: "Rule for the provider.",
							Optional:    true,
						},
						"signup_group": schema.StringAttribute{
							Description: "The signup group for the provider.",
							Optional:    true,
						},
						"country_codes": schema.ListAttribute{
							Description: "Country codes for the provider.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"signin_methods": schema.ListNestedAttribute{
				Description: "List of signin methods.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the signin method.",
							Required:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The display name of the signin method.",
							Optional:    true,
						},
						"rule": schema.StringAttribute{
							Description: "The rule for the signin method.",
							Optional:    true,
						},
					},
				},
			},
			"signup_items": schema.ListNestedAttribute{
				Description: "List of signup form items.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the signup item.",
							Required:    true,
						},
						"visible": schema.BoolAttribute{
							Description: "Whether the item is visible.",
							Optional:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Whether the item is required.",
							Optional:    true,
						},
						"prompted": schema.BoolAttribute{
							Description: "Whether the item is prompted.",
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the item.",
							Optional:    true,
						},
						"rule": schema.StringAttribute{
							Description: "Validation rule for the item.",
							Optional:    true,
						},
						"label": schema.StringAttribute{
							Description: "Label for the item.",
							Optional:    true,
						},
						"placeholder": schema.StringAttribute{
							Description: "Placeholder text for the item.",
							Optional:    true,
						},
						"regex": schema.StringAttribute{
							Description: "Regex pattern for validation.",
							Optional:    true,
						},
						"custom_css": schema.StringAttribute{
							Description: "Custom CSS for the item.",
							Optional:    true,
						},
						"options": schema.ListAttribute{
							Description: "Options for select-type items.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"signin_items": schema.ListNestedAttribute{
				Description: "List of signin form items.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the signin item.",
							Required:    true,
						},
						"visible": schema.BoolAttribute{
							Description: "Whether the item is visible.",
							Optional:    true,
						},
						"is_custom": schema.BoolAttribute{
							Description: "Whether the item is custom.",
							Optional:    true,
						},
						"label": schema.StringAttribute{
							Description: "Label for the item.",
							Optional:    true,
						},
						"placeholder": schema.StringAttribute{
							Description: "Placeholder text for the item.",
							Optional:    true,
						},
						"rule": schema.StringAttribute{
							Description: "Rule for the item.",
							Optional:    true,
						},
						"custom_css": schema.StringAttribute{
							Description: "Custom CSS for the item.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *ApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert list types to Go slices.
	var redirectURIs, grantTypes, tokenFields, tags []string

	if !plan.RedirectURIs.IsNull() && !plan.RedirectURIs.IsUnknown() {
		resp.Diagnostics.Append(plan.RedirectURIs.ElementsAs(ctx, &redirectURIs, false)...)
	}
	if !plan.GrantTypes.IsNull() && !plan.GrantTypes.IsUnknown() {
		resp.Diagnostics.Append(plan.GrantTypes.ElementsAs(ctx, &grantTypes, false)...)
	}
	if !plan.TokenFields.IsNull() && !plan.TokenFields.IsUnknown() {
		resp.Diagnostics.Append(plan.TokenFields.ElementsAs(ctx, &tokenFields, false)...)
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nested objects.
	var themeData *casdoorsdk.ThemeData
	if !plan.ThemeData.IsNull() && !plan.ThemeData.IsUnknown() {
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

	// Convert providers
	var providers []*casdoorsdk.ProviderItem
	if !plan.Providers.IsNull() && !plan.Providers.IsUnknown() {
		var providerModels []ProviderItemModel
		resp.Diagnostics.Append(plan.Providers.ElementsAs(ctx, &providerModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, p := range providerModels {
			var countryCodes []string
			if !p.CountryCodes.IsNull() {
				resp.Diagnostics.Append(p.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
			}
			providers = append(providers, &casdoorsdk.ProviderItem{
				Owner:        p.Owner.ValueString(),
				Name:         p.Name.ValueString(),
				CanSignUp:    p.CanSignUp.ValueBool(),
				CanSignIn:    p.CanSignIn.ValueBool(),
				CanUnlink:    p.CanUnlink.ValueBool(),
				Prompted:     p.Prompted.ValueBool(),
				Rule:         p.Rule.ValueString(),
				SignupGroup:  p.SignupGroup.ValueString(),
				CountryCodes: countryCodes,
			})
		}
	}

	// Convert signin methods
	var signinMethods []*casdoorsdk.SigninMethod
	if !plan.SigninMethods.IsNull() && !plan.SigninMethods.IsUnknown() {
		var methodModels []SigninMethodModel
		resp.Diagnostics.Append(plan.SigninMethods.ElementsAs(ctx, &methodModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, m := range methodModels {
			signinMethods = append(signinMethods, &casdoorsdk.SigninMethod{
				Name:        m.Name.ValueString(),
				DisplayName: m.DisplayName.ValueString(),
				Rule:        m.Rule.ValueString(),
			})
		}
	}

	// Convert signup items
	var signupItems []*casdoorsdk.SignupItem
	if !plan.SignupItems.IsNull() && !plan.SignupItems.IsUnknown() {
		var itemModels []SignupItemModel
		resp.Diagnostics.Append(plan.SignupItems.ElementsAs(ctx, &itemModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, i := range itemModels {
			var options []string
			if !i.Options.IsNull() {
				resp.Diagnostics.Append(i.Options.ElementsAs(ctx, &options, false)...)
			}
			signupItems = append(signupItems, &casdoorsdk.SignupItem{
				Name:        i.Name.ValueString(),
				Visible:     i.Visible.ValueBool(),
				Required:    i.Required.ValueBool(),
				Prompted:    i.Prompted.ValueBool(),
				Type:        i.Type.ValueString(),
				Rule:        i.Rule.ValueString(),
				Label:       i.Label.ValueString(),
				Placeholder: i.Placeholder.ValueString(),
				Regex:       i.Regex.ValueString(),
				CustomCss:   i.CustomCSS.ValueString(),
				Options:     options,
			})
		}
	}

	// Convert signin items
	var signinItems []*casdoorsdk.SigninItem
	if !plan.SigninItems.IsNull() && !plan.SigninItems.IsUnknown() {
		var itemModels []SigninItemModel
		resp.Diagnostics.Append(plan.SigninItems.ElementsAs(ctx, &itemModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, i := range itemModels {
			signinItems = append(signinItems, &casdoorsdk.SigninItem{
				Name:        i.Name.ValueString(),
				Visible:     i.Visible.ValueBool(),
				IsCustom:    i.IsCustom.ValueBool(),
				Label:       i.Label.ValueString(),
				Placeholder: i.Placeholder.ValueString(),
				Rule:        i.Rule.ValueString(),
				CustomCss:   i.CustomCSS.ValueString(),
			})
		}
	}

	// Convert SAML attributes
	var samlAttributes []*casdoorsdk.SamlItem
	if !plan.SamlAttributes.IsNull() && !plan.SamlAttributes.IsUnknown() {
		var samlModels []SamlItemModel
		resp.Diagnostics.Append(plan.SamlAttributes.ElementsAs(ctx, &samlModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, s := range samlModels {
			samlAttributes = append(samlAttributes, &casdoorsdk.SamlItem{
				Name:       s.Name.ValueString(),
				NameFormat: s.NameFormat.ValueString(),
				Value:      s.Value.ValueString(),
			})
		}
	}

	// Convert token attributes
	var tokenAttributes []*casdoorsdk.JwtItem
	if !plan.TokenAttributes.IsNull() && !plan.TokenAttributes.IsUnknown() {
		var jwtModels []JwtItemModel
		resp.Diagnostics.Append(plan.TokenAttributes.ElementsAs(ctx, &jwtModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, j := range jwtModels {
			tokenAttributes = append(tokenAttributes, &casdoorsdk.JwtItem{
				Name:  j.Name.ValueString(),
				Value: j.Value.ValueString(),
				Type:  j.Type.ValueString(),
			})
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	app := &casdoorsdk.Application{
		Owner:                        plan.Owner.ValueString(),
		Name:                         plan.Name.ValueString(),
		CreatedTime:                  createdTime,
		DisplayName:                  plan.DisplayName.ValueString(),
		Title:                        plan.Title.ValueString(),
		Favicon:                      plan.Favicon.ValueString(),
		Logo:                         plan.Logo.ValueString(),
		HomepageUrl:                  plan.HomepageURL.ValueString(),
		Description:                  plan.Description.ValueString(),
		Organization:                 plan.Organization.ValueString(),
		Cert:                         plan.Cert.ValueString(),
		DefaultGroup:                 plan.DefaultGroup.ValueString(),
		EnablePassword:               plan.EnablePassword.ValueBool(),
		EnableSignUp:                 plan.EnableSignUp.ValueBool(),
		EnableSigninSession:          plan.EnableSigninSession.ValueBool(),
		EnableAutoSignin:             plan.EnableAutoSignin.ValueBool(),
		EnableCodeSignin:             plan.EnableCodeSignin.ValueBool(),
		EnableExclusiveSignin:        plan.EnableExclusiveSignin.ValueBool(),
		EnableSamlCompress:           plan.EnableSamlCompress.ValueBool(),
		EnableSamlC14n10:             plan.EnableSamlC14n10.ValueBool(),
		EnableSamlPostBinding:        plan.EnableSamlPostBinding.ValueBool(),
		EnableSamlAssertionSignature: plan.EnableSamlAssertionSignature.ValueBool(),
		DisableSamlAttributes:        plan.DisableSamlAttributes.ValueBool(),
		UseEmailAsSamlNameId:         plan.UseEmailAsSamlNameId.ValueBool(),
		EnableWebAuthn:               plan.EnableWebAuthn.ValueBool(),
		EnableLinkWithEmail:          plan.EnableLinkWithEmail.ValueBool(),
		DisableSignin:                plan.DisableSignin.ValueBool(),
		IsShared:                     plan.IsShared.ValueBool(),
		RedirectUris:                 redirectURIs,
		TokenFormat:                  plan.TokenFormat.ValueString(),
		TokenSigningMethod:           plan.TokenSigningMethod.ValueString(),
		TokenFields:                  tokenFields,
		TokenAttributes:              tokenAttributes,
		ExpireInHours:                plan.ExpireInHours.ValueFloat64(),
		RefreshExpireInHours:         plan.RefreshExpireInHours.ValueFloat64(),
		CookieExpireInHours:          plan.CookieExpireInHours.ValueInt64(),
		GrantTypes:                   grantTypes,
		SamlReplyUrl:                 plan.SamlReplyUrl.ValueString(),
		SamlHashAlgorithm:            plan.SamlHashAlgorithm.ValueString(),
		SamlAttributes:               samlAttributes,
		SignupUrl:                    plan.SignupUrl.ValueString(),
		SigninUrl:                    plan.SigninUrl.ValueString(),
		ForgetUrl:                    plan.ForgetUrl.ValueString(),
		AffiliationUrl:               plan.AffiliationUrl.ValueString(),
		HeaderHtml:                   plan.HeaderHtml.ValueString(),
		FooterHtml:                   plan.FooterHtml.ValueString(),
		SignupHtml:                   plan.SignupHtml.ValueString(),
		SigninHtml:                   plan.SigninHtml.ValueString(),
		FormCss:                      plan.FormCss.ValueString(),
		FormCssMobile:                plan.FormCssMobile.ValueString(),
		FormOffset:                   int(plan.FormOffset.ValueInt64()),
		FormSideHtml:                 plan.FormSideHtml.ValueString(),
		FormBackgroundUrl:            plan.FormBackgroundUrl.ValueString(),
		FormBackgroundUrlMobile:      plan.FormBackgroundUrlMobile.ValueString(),
		ThemeData:                    themeData,
		IpRestriction:                plan.IpRestriction.ValueString(),
		IpWhitelist:                  plan.IpWhitelist.ValueString(),
		FailedSigninLimit:            int(plan.FailedSigninLimit.ValueInt64()),
		FailedSigninFrozenTime:       int(plan.FailedSigninFrozenTime.ValueInt64()),
		CodeResendTimeout:            int(plan.CodeResendTimeout.ValueInt64()),
		OrgChoiceMode:                plan.OrgChoiceMode.ValueString(),
		TermsOfUse:                   plan.TermsOfUse.ValueString(),
		Tags:                         tags,
		ForcedRedirectOrigin:         plan.ForcedRedirectOrigin.ValueString(),
		Order:                        int(plan.Order.ValueInt64()),
		Providers:                    providers,
		SigninMethods:                signinMethods,
		SignupItems:                  signupItems,
		SigninItems:                  signinItems,
	}

	success, err := r.client.AddApplication(app)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Application",
			fmt.Sprintf("Could not create application %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Application",
			fmt.Sprintf("Casdoor returned failure when creating application %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the created application to get computed fields.
	createdApp, err := r.client.GetApplication(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Application",
			fmt.Sprintf("Could not read application %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdApp == nil {
		resp.Diagnostics.AddError(
			"Error Reading Application",
			fmt.Sprintf("Application %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	// Update computed fields from API response.
	plan.CreatedTime = types.StringValue(createdApp.CreatedTime)
	plan.ClientID = types.StringValue(createdApp.ClientId)
	plan.ClientSecret = types.StringValue(createdApp.ClientSecret)
	plan.CertPublicKey = types.StringValue(createdApp.CertPublicKey)

	// For computed list/object fields that were unknown (not set in config),
	// populate them from the API response. For fields that were set in config,
	// keep the configured values.

	// RedirectURIs
	if plan.RedirectURIs.IsUnknown() {
		if len(createdApp.RedirectUris) > 0 {
			uriList, diags := types.ListValueFrom(ctx, types.StringType, createdApp.RedirectUris)
			resp.Diagnostics.Append(diags...)
			plan.RedirectURIs = uriList
		} else {
			plan.RedirectURIs = types.ListNull(types.StringType)
		}
	}

	// GrantTypes
	if plan.GrantTypes.IsUnknown() {
		if len(createdApp.GrantTypes) > 0 {
			gtList, diags := types.ListValueFrom(ctx, types.StringType, createdApp.GrantTypes)
			resp.Diagnostics.Append(diags...)
			plan.GrantTypes = gtList
		} else {
			plan.GrantTypes = types.ListNull(types.StringType)
		}
	}

	// TokenFields
	if plan.TokenFields.IsUnknown() {
		if len(createdApp.TokenFields) > 0 {
			tfList, diags := types.ListValueFrom(ctx, types.StringType, createdApp.TokenFields)
			resp.Diagnostics.Append(diags...)
			plan.TokenFields = tfList
		} else {
			plan.TokenFields = types.ListNull(types.StringType)
		}
	}

	// Tags
	if plan.Tags.IsUnknown() {
		if len(createdApp.Tags) > 0 {
			tagList, diags := types.ListValueFrom(ctx, types.StringType, createdApp.Tags)
			resp.Diagnostics.Append(diags...)
			plan.Tags = tagList
		} else {
			plan.Tags = types.ListNull(types.StringType)
		}
	}

	// ThemeData
	if plan.ThemeData.IsUnknown() {
		if createdApp.ThemeData != nil {
			themeObj, diags := types.ObjectValue(ThemeDataAttrTypes(), map[string]attr.Value{
				"theme_type":    types.StringValue(createdApp.ThemeData.ThemeType),
				"color_primary": types.StringValue(createdApp.ThemeData.ColorPrimary),
				"border_radius": types.Int64Value(int64(createdApp.ThemeData.BorderRadius)),
				"is_compact":    types.BoolValue(createdApp.ThemeData.IsCompact),
				"is_enabled":    types.BoolValue(createdApp.ThemeData.IsEnabled),
			})
			resp.Diagnostics.Append(diags...)
			plan.ThemeData = themeObj
		} else {
			plan.ThemeData = types.ObjectNull(ThemeDataAttrTypes())
		}
	}

	// Providers
	if plan.Providers.IsUnknown() {
		if len(createdApp.Providers) > 0 {
			providerObjList := make([]attr.Value, 0, len(createdApp.Providers))
			for _, p := range createdApp.Providers {
				providerObj, diags := types.ObjectValue(ProviderItemAttrTypes(), map[string]attr.Value{
					"owner":       types.StringValue(p.Owner),
					"name":        types.StringValue(p.Name),
					"can_sign_up": types.BoolValue(p.CanSignUp),
					"can_sign_in": types.BoolValue(p.CanSignIn),
					"can_unlink":  types.BoolValue(p.CanUnlink),
					"prompted":    types.BoolValue(p.Prompted),
					"rule":        types.StringValue(p.Rule),
				})
				resp.Diagnostics.Append(diags...)
				providerObjList = append(providerObjList, providerObj)
			}
			providerList, diags := types.ListValue(types.ObjectType{AttrTypes: ProviderItemAttrTypes()}, providerObjList)
			resp.Diagnostics.Append(diags...)
			plan.Providers = providerList
		} else {
			plan.Providers = types.ListNull(types.ObjectType{AttrTypes: ProviderItemAttrTypes()})
		}
	}

	// SigninMethods
	if plan.SigninMethods.IsUnknown() {
		if len(createdApp.SigninMethods) > 0 {
			methodObjList := make([]attr.Value, 0, len(createdApp.SigninMethods))
			for _, m := range createdApp.SigninMethods {
				methodObj, diags := types.ObjectValue(SigninMethodAttrTypes(), map[string]attr.Value{
					"name":         types.StringValue(m.Name),
					"display_name": types.StringValue(m.DisplayName),
					"rule":         types.StringValue(m.Rule),
				})
				resp.Diagnostics.Append(diags...)
				methodObjList = append(methodObjList, methodObj)
			}
			methodList, diags := types.ListValue(types.ObjectType{AttrTypes: SigninMethodAttrTypes()}, methodObjList)
			resp.Diagnostics.Append(diags...)
			plan.SigninMethods = methodList
		} else {
			plan.SigninMethods = types.ListNull(types.ObjectType{AttrTypes: SigninMethodAttrTypes()})
		}
	}

	// SignupItems
	if plan.SignupItems.IsUnknown() {
		if len(createdApp.SignupItems) > 0 {
			itemObjList := make([]attr.Value, 0, len(createdApp.SignupItems))
			for _, i := range createdApp.SignupItems {
				var options basetypes.ListValue
				if len(i.Options) > 0 {
					optionsList, optsDiags := types.ListValueFrom(ctx, types.StringType, i.Options)
					resp.Diagnostics.Append(optsDiags...)
					options = optionsList
				} else {
					options = types.ListNull(types.StringType)
				}
				itemObj, diags := types.ObjectValue(SignupItemAttrTypes(), map[string]attr.Value{
					"name":        types.StringValue(i.Name),
					"visible":     types.BoolValue(i.Visible),
					"required":    types.BoolValue(i.Required),
					"prompted":    types.BoolValue(i.Prompted),
					"type":        types.StringValue(i.Type),
					"rule":        types.StringValue(i.Rule),
					"label":       types.StringValue(i.Label),
					"placeholder": types.StringValue(i.Placeholder),
					"regex":       types.StringValue(i.Regex),
					"custom_css":  types.StringValue(i.CustomCss),
					"options":     options,
				})
				resp.Diagnostics.Append(diags...)
				itemObjList = append(itemObjList, itemObj)
			}
			itemList, diags := types.ListValue(types.ObjectType{AttrTypes: SignupItemAttrTypes()}, itemObjList)
			resp.Diagnostics.Append(diags...)
			plan.SignupItems = itemList
		} else {
			plan.SignupItems = types.ListNull(types.ObjectType{AttrTypes: SignupItemAttrTypes()})
		}
	}

	// SigninItems
	if plan.SigninItems.IsUnknown() {
		if len(createdApp.SigninItems) > 0 {
			itemObjList := make([]attr.Value, 0, len(createdApp.SigninItems))
			for _, i := range createdApp.SigninItems {
				itemObj, diags := types.ObjectValue(SigninItemAttrTypes(), map[string]attr.Value{
					"name":        types.StringValue(i.Name),
					"visible":     types.BoolValue(i.Visible),
					"is_custom":   types.BoolValue(i.IsCustom),
					"label":       types.StringValue(i.Label),
					"placeholder": types.StringValue(i.Placeholder),
					"rule":        types.StringValue(i.Rule),
					"custom_css":  types.StringValue(i.CustomCss),
				})
				resp.Diagnostics.Append(diags...)
				itemObjList = append(itemObjList, itemObj)
			}
			itemList, diags := types.ListValue(types.ObjectType{AttrTypes: SigninItemAttrTypes()}, itemObjList)
			resp.Diagnostics.Append(diags...)
			plan.SigninItems = itemList
		} else {
			plan.SigninItems = types.ListNull(types.ObjectType{AttrTypes: SigninItemAttrTypes()})
		}
	}

	// SamlAttributes
	if plan.SamlAttributes.IsUnknown() {
		if len(createdApp.SamlAttributes) > 0 {
			samlObjList := make([]attr.Value, 0, len(createdApp.SamlAttributes))
			for _, s := range createdApp.SamlAttributes {
				samlObj, diags := types.ObjectValue(SamlItemAttrTypes(), map[string]attr.Value{
					"name":        types.StringValue(s.Name),
					"name_format": types.StringValue(s.NameFormat),
					"value":       types.StringValue(s.Value),
				})
				resp.Diagnostics.Append(diags...)
				samlObjList = append(samlObjList, samlObj)
			}
			samlList, diags := types.ListValue(types.ObjectType{AttrTypes: SamlItemAttrTypes()}, samlObjList)
			resp.Diagnostics.Append(diags...)
			plan.SamlAttributes = samlList
		} else {
			plan.SamlAttributes = types.ListNull(types.ObjectType{AttrTypes: SamlItemAttrTypes()})
		}
	}

	// TokenAttributes
	if plan.TokenAttributes.IsUnknown() {
		if len(createdApp.TokenAttributes) > 0 {
			jwtObjList := make([]attr.Value, 0, len(createdApp.TokenAttributes))
			for _, j := range createdApp.TokenAttributes {
				jwtObj, diags := types.ObjectValue(JwtItemAttrTypes(), map[string]attr.Value{
					"name":  types.StringValue(j.Name),
					"value": types.StringValue(j.Value),
					"type":  types.StringValue(j.Type),
				})
				resp.Diagnostics.Append(diags...)
				jwtObjList = append(jwtObjList, jwtObj)
			}
			jwtList, diags := types.ListValue(types.ObjectType{AttrTypes: JwtItemAttrTypes()}, jwtObjList)
			resp.Diagnostics.Append(diags...)
			plan.TokenAttributes = jwtList
		} else {
			plan.TokenAttributes = types.ListNull(types.ObjectType{AttrTypes: JwtItemAttrTypes()})
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetApplication(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Application",
			fmt.Sprintf("Could not read application %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if app == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set scalar fields
	state.Owner = types.StringValue(app.Owner)
	state.Name = types.StringValue(app.Name)
	state.CreatedTime = types.StringValue(app.CreatedTime)
	state.DisplayName = types.StringValue(app.DisplayName)
	state.Title = types.StringValue(app.Title)
	state.Favicon = types.StringValue(app.Favicon)
	state.Logo = types.StringValue(app.Logo)
	state.HomepageURL = types.StringValue(app.HomepageUrl)
	state.Description = types.StringValue(app.Description)
	state.Organization = types.StringValue(app.Organization)
	state.Cert = types.StringValue(app.Cert)
	state.DefaultGroup = types.StringValue(app.DefaultGroup)
	state.EnablePassword = types.BoolValue(app.EnablePassword)
	state.EnableSignUp = types.BoolValue(app.EnableSignUp)
	state.EnableSigninSession = types.BoolValue(app.EnableSigninSession)
	state.EnableAutoSignin = types.BoolValue(app.EnableAutoSignin)
	state.EnableCodeSignin = types.BoolValue(app.EnableCodeSignin)
	state.EnableExclusiveSignin = types.BoolValue(app.EnableExclusiveSignin)
	state.EnableSamlCompress = types.BoolValue(app.EnableSamlCompress)
	state.EnableSamlC14n10 = types.BoolValue(app.EnableSamlC14n10)
	state.EnableSamlPostBinding = types.BoolValue(app.EnableSamlPostBinding)
	state.EnableSamlAssertionSignature = types.BoolValue(app.EnableSamlAssertionSignature)
	state.DisableSamlAttributes = types.BoolValue(app.DisableSamlAttributes)
	state.UseEmailAsSamlNameId = types.BoolValue(app.UseEmailAsSamlNameId)
	state.EnableWebAuthn = types.BoolValue(app.EnableWebAuthn)
	state.EnableLinkWithEmail = types.BoolValue(app.EnableLinkWithEmail)
	state.DisableSignin = types.BoolValue(app.DisableSignin)
	state.IsShared = types.BoolValue(app.IsShared)
	state.ClientID = types.StringValue(app.ClientId)
	state.ClientSecret = types.StringValue(app.ClientSecret)
	state.TokenFormat = types.StringValue(app.TokenFormat)
	state.TokenSigningMethod = types.StringValue(app.TokenSigningMethod)
	state.ExpireInHours = types.Float64Value(app.ExpireInHours)
	state.RefreshExpireInHours = types.Float64Value(app.RefreshExpireInHours)
	state.CookieExpireInHours = types.Int64Value(app.CookieExpireInHours)
	state.SamlReplyUrl = types.StringValue(app.SamlReplyUrl)
	state.SamlHashAlgorithm = types.StringValue(app.SamlHashAlgorithm)
	state.SignupUrl = types.StringValue(app.SignupUrl)
	state.SigninUrl = types.StringValue(app.SigninUrl)
	state.ForgetUrl = types.StringValue(app.ForgetUrl)
	state.AffiliationUrl = types.StringValue(app.AffiliationUrl)
	state.HeaderHtml = types.StringValue(app.HeaderHtml)
	state.FooterHtml = types.StringValue(app.FooterHtml)
	state.SignupHtml = types.StringValue(app.SignupHtml)
	state.SigninHtml = types.StringValue(app.SigninHtml)
	state.FormCss = types.StringValue(app.FormCss)
	state.FormCssMobile = types.StringValue(app.FormCssMobile)
	state.FormOffset = types.Int64Value(int64(app.FormOffset))
	state.FormSideHtml = types.StringValue(app.FormSideHtml)
	state.FormBackgroundUrl = types.StringValue(app.FormBackgroundUrl)
	state.FormBackgroundUrlMobile = types.StringValue(app.FormBackgroundUrlMobile)
	state.IpRestriction = types.StringValue(app.IpRestriction)
	state.IpWhitelist = types.StringValue(app.IpWhitelist)
	state.FailedSigninLimit = types.Int64Value(int64(app.FailedSigninLimit))
	state.FailedSigninFrozenTime = types.Int64Value(int64(app.FailedSigninFrozenTime))
	state.CodeResendTimeout = types.Int64Value(int64(app.CodeResendTimeout))
	state.OrgChoiceMode = types.StringValue(app.OrgChoiceMode)
	state.TermsOfUse = types.StringValue(app.TermsOfUse)
	state.CertPublicKey = types.StringValue(app.CertPublicKey)
	state.ForcedRedirectOrigin = types.StringValue(app.ForcedRedirectOrigin)
	state.Order = types.Int64Value(int64(app.Order))

	// Convert string slices to list types
	if len(app.RedirectUris) > 0 {
		redirectURIs, diags := types.ListValueFrom(ctx, types.StringType, app.RedirectUris)
		resp.Diagnostics.Append(diags...)
		state.RedirectURIs = redirectURIs
	} else {
		state.RedirectURIs = types.ListNull(types.StringType)
	}

	if len(app.GrantTypes) > 0 {
		grantTypes, diags := types.ListValueFrom(ctx, types.StringType, app.GrantTypes)
		resp.Diagnostics.Append(diags...)
		state.GrantTypes = grantTypes
	} else {
		state.GrantTypes = types.ListNull(types.StringType)
	}

	if len(app.TokenFields) > 0 {
		tokenFields, diags := types.ListValueFrom(ctx, types.StringType, app.TokenFields)
		resp.Diagnostics.Append(diags...)
		state.TokenFields = tokenFields
	} else {
		state.TokenFields = types.ListNull(types.StringType)
	}

	if len(app.Tags) > 0 {
		tags, diags := types.ListValueFrom(ctx, types.StringType, app.Tags)
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	} else {
		state.Tags = types.ListNull(types.StringType)
	}

	// Convert ThemeData to object type
	if app.ThemeData != nil {
		themeObj, diags := types.ObjectValue(ThemeDataAttrTypes(), map[string]attr.Value{
			"theme_type":    types.StringValue(app.ThemeData.ThemeType),
			"color_primary": types.StringValue(app.ThemeData.ColorPrimary),
			"border_radius": types.Int64Value(int64(app.ThemeData.BorderRadius)),
			"is_compact":    types.BoolValue(app.ThemeData.IsCompact),
			"is_enabled":    types.BoolValue(app.ThemeData.IsEnabled),
		})
		resp.Diagnostics.Append(diags...)
		state.ThemeData = themeObj
	} else {
		state.ThemeData = types.ObjectNull(ThemeDataAttrTypes())
	}

	// Convert Providers to list of objects
	if len(app.Providers) > 0 {
		providerObjList := make([]attr.Value, 0, len(app.Providers))
		for _, p := range app.Providers {
			var countryCodes basetypes.ListValue
			if len(p.CountryCodes) > 0 {
				ccList, ccDiags := types.ListValueFrom(ctx, types.StringType, p.CountryCodes)
				resp.Diagnostics.Append(ccDiags...)
				countryCodes = ccList
			} else {
				countryCodes = types.ListNull(types.StringType)
			}
			providerObj, diags := types.ObjectValue(ProviderItemAttrTypes(), map[string]attr.Value{
				"owner":         types.StringValue(p.Owner),
				"name":          types.StringValue(p.Name),
				"can_sign_up":   types.BoolValue(p.CanSignUp),
				"can_sign_in":   types.BoolValue(p.CanSignIn),
				"can_unlink":    types.BoolValue(p.CanUnlink),
				"prompted":      types.BoolValue(p.Prompted),
				"rule":          types.StringValue(p.Rule),
				"signup_group":  types.StringValue(p.SignupGroup),
				"country_codes": countryCodes,
			})
			resp.Diagnostics.Append(diags...)
			providerObjList = append(providerObjList, providerObj)
		}
		providerList, diags := types.ListValue(types.ObjectType{AttrTypes: ProviderItemAttrTypes()}, providerObjList)
		resp.Diagnostics.Append(diags...)
		state.Providers = providerList
	} else {
		state.Providers = types.ListNull(types.ObjectType{AttrTypes: ProviderItemAttrTypes()})
	}

	// Convert SigninMethods to list of objects
	if len(app.SigninMethods) > 0 {
		methodObjList := make([]attr.Value, 0, len(app.SigninMethods))
		for _, m := range app.SigninMethods {
			methodObj, diags := types.ObjectValue(SigninMethodAttrTypes(), map[string]attr.Value{
				"name":         types.StringValue(m.Name),
				"display_name": types.StringValue(m.DisplayName),
				"rule":         types.StringValue(m.Rule),
			})
			resp.Diagnostics.Append(diags...)
			methodObjList = append(methodObjList, methodObj)
		}
		methodList, diags := types.ListValue(types.ObjectType{AttrTypes: SigninMethodAttrTypes()}, methodObjList)
		resp.Diagnostics.Append(diags...)
		state.SigninMethods = methodList
	} else {
		state.SigninMethods = types.ListNull(types.ObjectType{AttrTypes: SigninMethodAttrTypes()})
	}

	// Convert SignupItems to list of objects
	if len(app.SignupItems) > 0 {
		itemObjList := make([]attr.Value, 0, len(app.SignupItems))
		for _, i := range app.SignupItems {
			var options basetypes.ListValue
			if len(i.Options) > 0 {
				optionsList, optsDiags := types.ListValueFrom(ctx, types.StringType, i.Options)
				resp.Diagnostics.Append(optsDiags...)
				options = optionsList
			} else {
				options = types.ListNull(types.StringType)
			}
			itemObj, diags := types.ObjectValue(SignupItemAttrTypes(), map[string]attr.Value{
				"name":        types.StringValue(i.Name),
				"visible":     types.BoolValue(i.Visible),
				"required":    types.BoolValue(i.Required),
				"prompted":    types.BoolValue(i.Prompted),
				"type":        types.StringValue(i.Type),
				"rule":        types.StringValue(i.Rule),
				"label":       types.StringValue(i.Label),
				"placeholder": types.StringValue(i.Placeholder),
				"regex":       types.StringValue(i.Regex),
				"custom_css":  types.StringValue(i.CustomCss),
				"options":     options,
			})
			resp.Diagnostics.Append(diags...)
			itemObjList = append(itemObjList, itemObj)
		}
		itemList, diags := types.ListValue(types.ObjectType{AttrTypes: SignupItemAttrTypes()}, itemObjList)
		resp.Diagnostics.Append(diags...)
		state.SignupItems = itemList
	} else {
		state.SignupItems = types.ListNull(types.ObjectType{AttrTypes: SignupItemAttrTypes()})
	}

	// Convert SigninItems to list of objects
	if len(app.SigninItems) > 0 {
		itemObjList := make([]attr.Value, 0, len(app.SigninItems))
		for _, i := range app.SigninItems {
			itemObj, diags := types.ObjectValue(SigninItemAttrTypes(), map[string]attr.Value{
				"name":        types.StringValue(i.Name),
				"visible":     types.BoolValue(i.Visible),
				"is_custom":   types.BoolValue(i.IsCustom),
				"label":       types.StringValue(i.Label),
				"placeholder": types.StringValue(i.Placeholder),
				"rule":        types.StringValue(i.Rule),
				"custom_css":  types.StringValue(i.CustomCss),
			})
			resp.Diagnostics.Append(diags...)
			itemObjList = append(itemObjList, itemObj)
		}
		itemList, diags := types.ListValue(types.ObjectType{AttrTypes: SigninItemAttrTypes()}, itemObjList)
		resp.Diagnostics.Append(diags...)
		state.SigninItems = itemList
	} else {
		state.SigninItems = types.ListNull(types.ObjectType{AttrTypes: SigninItemAttrTypes()})
	}

	// Convert SamlAttributes to list of objects
	if len(app.SamlAttributes) > 0 {
		samlObjList := make([]attr.Value, 0, len(app.SamlAttributes))
		for _, s := range app.SamlAttributes {
			samlObj, diags := types.ObjectValue(SamlItemAttrTypes(), map[string]attr.Value{
				"name":        types.StringValue(s.Name),
				"name_format": types.StringValue(s.NameFormat),
				"value":       types.StringValue(s.Value),
			})
			resp.Diagnostics.Append(diags...)
			samlObjList = append(samlObjList, samlObj)
		}
		samlList, diags := types.ListValue(types.ObjectType{AttrTypes: SamlItemAttrTypes()}, samlObjList)
		resp.Diagnostics.Append(diags...)
		state.SamlAttributes = samlList
	} else {
		state.SamlAttributes = types.ListNull(types.ObjectType{AttrTypes: SamlItemAttrTypes()})
	}

	// Convert TokenAttributes to list of objects
	if len(app.TokenAttributes) > 0 {
		jwtObjList := make([]attr.Value, 0, len(app.TokenAttributes))
		for _, j := range app.TokenAttributes {
			jwtObj, diags := types.ObjectValue(JwtItemAttrTypes(), map[string]attr.Value{
				"name":  types.StringValue(j.Name),
				"value": types.StringValue(j.Value),
				"type":  types.StringValue(j.Type),
			})
			resp.Diagnostics.Append(diags...)
			jwtObjList = append(jwtObjList, jwtObj)
		}
		jwtList, diags := types.ListValue(types.ObjectType{AttrTypes: JwtItemAttrTypes()}, jwtObjList)
		resp.Diagnostics.Append(diags...)
		state.TokenAttributes = jwtList
	} else {
		state.TokenAttributes = types.ListNull(types.ObjectType{AttrTypes: JwtItemAttrTypes()})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve computed fields from state.
	var state ApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert list types to Go slices.
	var redirectURIs, grantTypes, tokenFields, tags []string

	if !plan.RedirectURIs.IsNull() && !plan.RedirectURIs.IsUnknown() {
		resp.Diagnostics.Append(plan.RedirectURIs.ElementsAs(ctx, &redirectURIs, false)...)
	}
	if !plan.GrantTypes.IsNull() && !plan.GrantTypes.IsUnknown() {
		resp.Diagnostics.Append(plan.GrantTypes.ElementsAs(ctx, &grantTypes, false)...)
	}
	if !plan.TokenFields.IsNull() && !plan.TokenFields.IsUnknown() {
		resp.Diagnostics.Append(plan.TokenFields.ElementsAs(ctx, &tokenFields, false)...)
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nested objects.
	var themeData *casdoorsdk.ThemeData
	if !plan.ThemeData.IsNull() && !plan.ThemeData.IsUnknown() {
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

	// Convert providers
	var providers []*casdoorsdk.ProviderItem
	if !plan.Providers.IsNull() && !plan.Providers.IsUnknown() {
		var providerModels []ProviderItemModel
		resp.Diagnostics.Append(plan.Providers.ElementsAs(ctx, &providerModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, p := range providerModels {
			var countryCodes []string
			if !p.CountryCodes.IsNull() {
				resp.Diagnostics.Append(p.CountryCodes.ElementsAs(ctx, &countryCodes, false)...)
			}
			providers = append(providers, &casdoorsdk.ProviderItem{
				Owner:        p.Owner.ValueString(),
				Name:         p.Name.ValueString(),
				CanSignUp:    p.CanSignUp.ValueBool(),
				CanSignIn:    p.CanSignIn.ValueBool(),
				CanUnlink:    p.CanUnlink.ValueBool(),
				Prompted:     p.Prompted.ValueBool(),
				Rule:         p.Rule.ValueString(),
				SignupGroup:  p.SignupGroup.ValueString(),
				CountryCodes: countryCodes,
			})
		}
	}

	// Convert signin methods
	var signinMethods []*casdoorsdk.SigninMethod
	if !plan.SigninMethods.IsNull() && !plan.SigninMethods.IsUnknown() {
		var methodModels []SigninMethodModel
		resp.Diagnostics.Append(plan.SigninMethods.ElementsAs(ctx, &methodModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, m := range methodModels {
			signinMethods = append(signinMethods, &casdoorsdk.SigninMethod{
				Name:        m.Name.ValueString(),
				DisplayName: m.DisplayName.ValueString(),
				Rule:        m.Rule.ValueString(),
			})
		}
	}

	// Convert signup items
	var signupItems []*casdoorsdk.SignupItem
	if !plan.SignupItems.IsNull() && !plan.SignupItems.IsUnknown() {
		var itemModels []SignupItemModel
		resp.Diagnostics.Append(plan.SignupItems.ElementsAs(ctx, &itemModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, i := range itemModels {
			var options []string
			if !i.Options.IsNull() {
				resp.Diagnostics.Append(i.Options.ElementsAs(ctx, &options, false)...)
			}
			signupItems = append(signupItems, &casdoorsdk.SignupItem{
				Name:        i.Name.ValueString(),
				Visible:     i.Visible.ValueBool(),
				Required:    i.Required.ValueBool(),
				Prompted:    i.Prompted.ValueBool(),
				Type:        i.Type.ValueString(),
				Rule:        i.Rule.ValueString(),
				Label:       i.Label.ValueString(),
				Placeholder: i.Placeholder.ValueString(),
				Regex:       i.Regex.ValueString(),
				CustomCss:   i.CustomCSS.ValueString(),
				Options:     options,
			})
		}
	}

	// Convert signin items
	var signinItems []*casdoorsdk.SigninItem
	if !plan.SigninItems.IsNull() && !plan.SigninItems.IsUnknown() {
		var itemModels []SigninItemModel
		resp.Diagnostics.Append(plan.SigninItems.ElementsAs(ctx, &itemModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, i := range itemModels {
			signinItems = append(signinItems, &casdoorsdk.SigninItem{
				Name:        i.Name.ValueString(),
				Visible:     i.Visible.ValueBool(),
				IsCustom:    i.IsCustom.ValueBool(),
				Label:       i.Label.ValueString(),
				Placeholder: i.Placeholder.ValueString(),
				Rule:        i.Rule.ValueString(),
				CustomCss:   i.CustomCSS.ValueString(),
			})
		}
	}

	// Convert SAML attributes
	var samlAttributes []*casdoorsdk.SamlItem
	if !plan.SamlAttributes.IsNull() && !plan.SamlAttributes.IsUnknown() {
		var samlModels []SamlItemModel
		resp.Diagnostics.Append(plan.SamlAttributes.ElementsAs(ctx, &samlModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, s := range samlModels {
			samlAttributes = append(samlAttributes, &casdoorsdk.SamlItem{
				Name:       s.Name.ValueString(),
				NameFormat: s.NameFormat.ValueString(),
				Value:      s.Value.ValueString(),
			})
		}
	}

	// Convert token attributes
	var tokenAttributes []*casdoorsdk.JwtItem
	if !plan.TokenAttributes.IsNull() && !plan.TokenAttributes.IsUnknown() {
		var jwtModels []JwtItemModel
		resp.Diagnostics.Append(plan.TokenAttributes.ElementsAs(ctx, &jwtModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, j := range jwtModels {
			tokenAttributes = append(tokenAttributes, &casdoorsdk.JwtItem{
				Name:  j.Name.ValueString(),
				Value: j.Value.ValueString(),
				Type:  j.Type.ValueString(),
			})
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	app := &casdoorsdk.Application{
		Owner:                        plan.Owner.ValueString(),
		Name:                         plan.Name.ValueString(),
		CreatedTime:                  state.CreatedTime.ValueString(),
		DisplayName:                  plan.DisplayName.ValueString(),
		Title:                        plan.Title.ValueString(),
		Favicon:                      plan.Favicon.ValueString(),
		Logo:                         plan.Logo.ValueString(),
		HomepageUrl:                  plan.HomepageURL.ValueString(),
		Description:                  plan.Description.ValueString(),
		Organization:                 plan.Organization.ValueString(),
		Cert:                         plan.Cert.ValueString(),
		DefaultGroup:                 plan.DefaultGroup.ValueString(),
		EnablePassword:               plan.EnablePassword.ValueBool(),
		EnableSignUp:                 plan.EnableSignUp.ValueBool(),
		EnableSigninSession:          plan.EnableSigninSession.ValueBool(),
		EnableAutoSignin:             plan.EnableAutoSignin.ValueBool(),
		EnableCodeSignin:             plan.EnableCodeSignin.ValueBool(),
		EnableExclusiveSignin:        plan.EnableExclusiveSignin.ValueBool(),
		EnableSamlCompress:           plan.EnableSamlCompress.ValueBool(),
		EnableSamlC14n10:             plan.EnableSamlC14n10.ValueBool(),
		EnableSamlPostBinding:        plan.EnableSamlPostBinding.ValueBool(),
		EnableSamlAssertionSignature: plan.EnableSamlAssertionSignature.ValueBool(),
		DisableSamlAttributes:        plan.DisableSamlAttributes.ValueBool(),
		UseEmailAsSamlNameId:         plan.UseEmailAsSamlNameId.ValueBool(),
		EnableWebAuthn:               plan.EnableWebAuthn.ValueBool(),
		EnableLinkWithEmail:          plan.EnableLinkWithEmail.ValueBool(),
		DisableSignin:                plan.DisableSignin.ValueBool(),
		IsShared:                     plan.IsShared.ValueBool(),
		ClientId:                     state.ClientID.ValueString(),
		ClientSecret:                 state.ClientSecret.ValueString(),
		RedirectUris:                 redirectURIs,
		TokenFormat:                  plan.TokenFormat.ValueString(),
		TokenSigningMethod:           plan.TokenSigningMethod.ValueString(),
		TokenFields:                  tokenFields,
		TokenAttributes:              tokenAttributes,
		ExpireInHours:                plan.ExpireInHours.ValueFloat64(),
		RefreshExpireInHours:         plan.RefreshExpireInHours.ValueFloat64(),
		CookieExpireInHours:          plan.CookieExpireInHours.ValueInt64(),
		GrantTypes:                   grantTypes,
		SamlReplyUrl:                 plan.SamlReplyUrl.ValueString(),
		SamlHashAlgorithm:            plan.SamlHashAlgorithm.ValueString(),
		SamlAttributes:               samlAttributes,
		SignupUrl:                    plan.SignupUrl.ValueString(),
		SigninUrl:                    plan.SigninUrl.ValueString(),
		ForgetUrl:                    plan.ForgetUrl.ValueString(),
		AffiliationUrl:               plan.AffiliationUrl.ValueString(),
		HeaderHtml:                   plan.HeaderHtml.ValueString(),
		FooterHtml:                   plan.FooterHtml.ValueString(),
		SignupHtml:                   plan.SignupHtml.ValueString(),
		SigninHtml:                   plan.SigninHtml.ValueString(),
		FormCss:                      plan.FormCss.ValueString(),
		FormCssMobile:                plan.FormCssMobile.ValueString(),
		FormOffset:                   int(plan.FormOffset.ValueInt64()),
		FormSideHtml:                 plan.FormSideHtml.ValueString(),
		FormBackgroundUrl:            plan.FormBackgroundUrl.ValueString(),
		FormBackgroundUrlMobile:      plan.FormBackgroundUrlMobile.ValueString(),
		ThemeData:                    themeData,
		IpRestriction:                plan.IpRestriction.ValueString(),
		IpWhitelist:                  plan.IpWhitelist.ValueString(),
		FailedSigninLimit:            int(plan.FailedSigninLimit.ValueInt64()),
		FailedSigninFrozenTime:       int(plan.FailedSigninFrozenTime.ValueInt64()),
		CodeResendTimeout:            int(plan.CodeResendTimeout.ValueInt64()),
		OrgChoiceMode:                plan.OrgChoiceMode.ValueString(),
		TermsOfUse:                   plan.TermsOfUse.ValueString(),
		Tags:                         tags,
		ForcedRedirectOrigin:         plan.ForcedRedirectOrigin.ValueString(),
		Order:                        int(plan.Order.ValueInt64()),
		Providers:                    providers,
		SigninMethods:                signinMethods,
		SignupItems:                  signupItems,
		SigninItems:                  signinItems,
	}

	success, err := r.client.UpdateApplication(app)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Application",
			fmt.Sprintf("Could not update application %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Application",
			fmt.Sprintf("Casdoor returned failure when updating application %q", plan.Name.ValueString()),
		)
		return
	}

	// Keep computed values from state.
	plan.CreatedTime = state.CreatedTime
	plan.ClientID = state.ClientID
	plan.ClientSecret = state.ClientSecret
	plan.CertPublicKey = state.CertPublicKey

	// Set list values to null if empty to match plan.
	if len(redirectURIs) == 0 {
		plan.RedirectURIs = types.ListNull(types.StringType)
	}
	if len(grantTypes) == 0 {
		plan.GrantTypes = types.ListNull(types.StringType)
	}
	if len(tokenFields) == 0 {
		plan.TokenFields = types.ListNull(types.StringType)
	}
	if len(tags) == 0 {
		plan.Tags = types.ListNull(types.StringType)
	}
	if len(samlAttributes) == 0 {
		plan.SamlAttributes = types.ListNull(types.ObjectType{AttrTypes: SamlItemAttrTypes()})
	}
	if len(tokenAttributes) == 0 {
		plan.TokenAttributes = types.ListNull(types.ObjectType{AttrTypes: JwtItemAttrTypes()})
	}
	if len(providers) == 0 {
		plan.Providers = types.ListNull(types.ObjectType{AttrTypes: ProviderItemAttrTypes()})
	}
	if len(signinMethods) == 0 {
		plan.SigninMethods = types.ListNull(types.ObjectType{AttrTypes: SigninMethodAttrTypes()})
	}
	if len(signupItems) == 0 {
		plan.SignupItems = types.ListNull(types.ObjectType{AttrTypes: SignupItemAttrTypes()})
	}
	if len(signinItems) == 0 {
		plan.SigninItems = types.ListNull(types.ObjectType{AttrTypes: SigninItemAttrTypes()})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Include organization field as the Casdoor API requires it for deletion.
	app := &casdoorsdk.Application{
		Owner:        state.Owner.ValueString(),
		Name:         state.Name.ValueString(),
		Organization: state.Organization.ValueString(),
	}

	success, err := r.client.DeleteApplication(app)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Application",
			fmt.Sprintf("Could not delete application %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Application",
			fmt.Sprintf("Casdoor returned failure when deleting application %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *ApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
