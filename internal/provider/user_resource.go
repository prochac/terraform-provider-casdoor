// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

type UserResource struct {
	client *casdoorsdk.Client
}

type UserResourceModel struct {
	Owner                  types.String `tfsdk:"owner"`
	Name                   types.String `tfsdk:"name"`
	ID                     types.String `tfsdk:"id"`
	Type                   types.String `tfsdk:"type"`
	Password               types.String `tfsdk:"password"`
	PasswordType           types.String `tfsdk:"password_type"`
	DisplayName            types.String `tfsdk:"display_name"`
	FirstName              types.String `tfsdk:"first_name"`
	LastName               types.String `tfsdk:"last_name"`
	Avatar                 types.String `tfsdk:"avatar"`
	Email                  types.String `tfsdk:"email"`
	EmailVerified          types.Bool   `tfsdk:"email_verified"`
	Phone                  types.String `tfsdk:"phone"`
	CountryCode            types.String `tfsdk:"country_code"`
	Region                 types.String `tfsdk:"region"`
	Location               types.String `tfsdk:"location"`
	Affiliation            types.String `tfsdk:"affiliation"`
	Title                  types.String `tfsdk:"title"`
	Homepage               types.String `tfsdk:"homepage"`
	Bio                    types.String `tfsdk:"bio"`
	Tag                    types.String `tfsdk:"tag"`
	Language               types.String `tfsdk:"language"`
	Gender                 types.String `tfsdk:"gender"`
	Birthday               types.String `tfsdk:"birthday"`
	Education              types.String `tfsdk:"education"`
	Score                  types.Int64  `tfsdk:"score"`
	Karma                  types.Int64  `tfsdk:"karma"`
	Ranking                types.Int64  `tfsdk:"ranking"`
	IsAdmin                types.Bool   `tfsdk:"is_admin"`
	IsForbidden            types.Bool   `tfsdk:"is_forbidden"`
	IsDeleted              types.Bool   `tfsdk:"is_deleted"`
	SignupApplication      types.String `tfsdk:"signup_application"`
	CreatedTime            types.String `tfsdk:"created_time"`
	UpdatedTime            types.String `tfsdk:"updated_time"`
	ExternalId             types.String `tfsdk:"external_id"`
	PasswordSalt           types.String `tfsdk:"password_salt"`
	AvatarType             types.String `tfsdk:"avatar_type"`
	PermanentAvatar        types.String `tfsdk:"permanent_avatar"`
	Address                types.List   `tfsdk:"address"`
	IdCardType             types.String `tfsdk:"id_card_type"`
	IdCard                 types.String `tfsdk:"id_card"`
	IsDefaultAvatar        types.Bool   `tfsdk:"is_default_avatar"`
	IsOnline               types.Bool   `tfsdk:"is_online"`
	Hash                   types.String `tfsdk:"hash"`
	PreHash                types.String `tfsdk:"pre_hash"`
	AccessKey              types.String `tfsdk:"access_key"`
	AccessSecret           types.String `tfsdk:"access_secret"`
	CreatedIp              types.String `tfsdk:"created_ip"`
	LastSigninTime         types.String `tfsdk:"last_signin_time"`
	LastSigninIp           types.String `tfsdk:"last_signin_ip"`
	Invitation             types.String `tfsdk:"invitation"`
	InvitationCode         types.String `tfsdk:"invitation_code"`
	Ldap                   types.String `tfsdk:"ldap"`
	Properties             types.Map    `tfsdk:"properties"`
	NeedUpdatePassword     types.Bool   `tfsdk:"need_update_password"`
	LastChangePasswordTime types.String `tfsdk:"last_change_password_time"`
	LastSigninWrongTime    types.String `tfsdk:"last_signin_wrong_time"`
	SigninWrongTimes       types.Int64  `tfsdk:"signin_wrong_times"`
	PreferredMfaType       types.String `tfsdk:"preferred_mfa_type"`
	RecoveryCodes          types.List   `tfsdk:"recovery_codes"`
	TotpSecret             types.String `tfsdk:"totp_secret"`
	MfaPhoneEnabled        types.Bool   `tfsdk:"mfa_phone_enabled"`
	MfaEmailEnabled        types.Bool   `tfsdk:"mfa_email_enabled"`
	Groups                 types.List   `tfsdk:"groups"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor user.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique username.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the user in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The user type (e.g., 'normal-user').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("normal-user"),
			},
			"password": schema.StringAttribute{
				Description: "The user's password. Note: This is write-only and will not be read back from Casdoor.",
				Optional:    true,
				Sensitive:   true,
			},
			"password_type": schema.StringAttribute{
				Description: "The password hashing type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"first_name": schema.StringAttribute{
				Description: "The user's first name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"last_name": schema.StringAttribute{
				Description: "The user's last name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"avatar": schema.StringAttribute{
				Description: "URL of the user's avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"email_verified": schema.BoolAttribute{
				Description: "Whether the user's email has been verified.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"phone": schema.StringAttribute{
				Description: "The user's phone number.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"country_code": schema.StringAttribute{
				Description: "The country code for the phone number.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"region": schema.StringAttribute{
				Description: "The user's region.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"location": schema.StringAttribute{
				Description: "The user's location.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"affiliation": schema.StringAttribute{
				Description: "The user's affiliation (e.g., company name).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"title": schema.StringAttribute{
				Description: "The user's job title.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"homepage": schema.StringAttribute{
				Description: "The user's homepage URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"bio": schema.StringAttribute{
				Description: "The user's biography.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tag": schema.StringAttribute{
				Description: "A tag for the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"language": schema.StringAttribute{
				Description: "The user's preferred language.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"gender": schema.StringAttribute{
				Description: "The user's gender.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"birthday": schema.StringAttribute{
				Description: "The user's birthday (ISO 8601 format).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"education": schema.StringAttribute{
				Description: "The user's education level.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"score": schema.Int64Attribute{
				Description: "The user's score.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"karma": schema.Int64Attribute{
				Description: "The user's karma points.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"ranking": schema.Int64Attribute{
				Description: "The user's ranking.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"is_admin": schema.BoolAttribute{
				Description: "Whether the user is an administrator.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_forbidden": schema.BoolAttribute{
				Description: "Whether the user is forbidden (disabled).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_deleted": schema.BoolAttribute{
				Description: "Whether the user is soft-deleted.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"signup_application": schema.StringAttribute{
				Description: "The application through which the user signed up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the user was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_time": schema.StringAttribute{
				Description: "The time when the user was last updated.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_id": schema.StringAttribute{
				Description: "External ID for the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_salt": schema.StringAttribute{
				Description: "The password salt.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"avatar_type": schema.StringAttribute{
				Description: "The type of avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"permanent_avatar": schema.StringAttribute{
				Description: "URL of the permanent avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"address": schema.ListAttribute{
				Description: "The user's address lines.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"id_card_type": schema.StringAttribute{
				Description: "The type of ID card.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"id_card": schema.StringAttribute{
				Description: "The ID card number.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"is_default_avatar": schema.BoolAttribute{
				Description: "Whether the user has the default avatar.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_online": schema.BoolAttribute{
				Description: "Whether the user is currently online.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hash": schema.StringAttribute{
				Description: "The user hash.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pre_hash": schema.StringAttribute{
				Description: "The previous user hash.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_key": schema.StringAttribute{
				Description: "The user's access key.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"access_secret": schema.StringAttribute{
				Description: "The user's access secret.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"created_ip": schema.StringAttribute{
				Description: "The IP address the user was created from.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_time": schema.StringAttribute{
				Description: "The last sign-in time.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_ip": schema.StringAttribute{
				Description: "The last sign-in IP address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invitation": schema.StringAttribute{
				Description: "The invitation used to sign up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"invitation_code": schema.StringAttribute{
				Description: "The invitation code used to sign up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ldap": schema.StringAttribute{
				Description: "LDAP identifier.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"properties": schema.MapAttribute{
				Description: "Custom properties for the user.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"need_update_password": schema.BoolAttribute{
				Description: "Whether the user needs to update their password.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"last_change_password_time": schema.StringAttribute{
				Description: "The last time the password was changed.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_wrong_time": schema.StringAttribute{
				Description: "The last time a wrong sign-in attempt was made.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"signin_wrong_times": schema.Int64Attribute{
				Description: "The number of wrong sign-in attempts.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"preferred_mfa_type": schema.StringAttribute{
				Description: "The preferred MFA type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"recovery_codes": schema.ListAttribute{
				Description: "MFA recovery codes.",
				Optional:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"totp_secret": schema.StringAttribute{
				Description: "The TOTP secret for MFA.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_phone_enabled": schema.BoolAttribute{
				Description: "Whether phone-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mfa_email_enabled": schema.BoolAttribute{
				Description: "Whether email-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"groups": schema.ListAttribute{
				Description: "List of groups the user belongs to.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []string
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	}
	var address []string
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		resp.Diagnostics.Append(plan.Address.ElementsAs(ctx, &address, false)...)
	}
	var recoveryCodes []string
	if !plan.RecoveryCodes.IsNull() && !plan.RecoveryCodes.IsUnknown() {
		resp.Diagnostics.Append(plan.RecoveryCodes.ElementsAs(ctx, &recoveryCodes, false)...)
	}
	var properties map[string]string
	if !plan.Properties.IsNull() && !plan.Properties.IsUnknown() {
		properties = make(map[string]string)
		resp.Diagnostics.Append(plan.Properties.ElementsAs(ctx, &properties, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	user := &casdoorsdk.User{
		Owner:              plan.Owner.ValueString(),
		Name:               plan.Name.ValueString(),
		CreatedTime:        createdTime,
		UpdatedTime:        createdTime,
		Type:               plan.Type.ValueString(),
		Password:           plan.Password.ValueString(),
		PasswordType:       plan.PasswordType.ValueString(),
		DisplayName:        plan.DisplayName.ValueString(),
		FirstName:          plan.FirstName.ValueString(),
		LastName:           plan.LastName.ValueString(),
		Avatar:             plan.Avatar.ValueString(),
		Email:              plan.Email.ValueString(),
		EmailVerified:      plan.EmailVerified.ValueBool(),
		Phone:              plan.Phone.ValueString(),
		CountryCode:        plan.CountryCode.ValueString(),
		Region:             plan.Region.ValueString(),
		Location:           plan.Location.ValueString(),
		Affiliation:        plan.Affiliation.ValueString(),
		Title:              plan.Title.ValueString(),
		Homepage:           plan.Homepage.ValueString(),
		Bio:                plan.Bio.ValueString(),
		Tag:                plan.Tag.ValueString(),
		Language:           plan.Language.ValueString(),
		Gender:             plan.Gender.ValueString(),
		Birthday:           plan.Birthday.ValueString(),
		Education:          plan.Education.ValueString(),
		Score:              int(plan.Score.ValueInt64()),
		Karma:              int(plan.Karma.ValueInt64()),
		Ranking:            int(plan.Ranking.ValueInt64()),
		IsAdmin:            plan.IsAdmin.ValueBool(),
		IsForbidden:        plan.IsForbidden.ValueBool(),
		IsDeleted:          plan.IsDeleted.ValueBool(),
		SignupApplication:  plan.SignupApplication.ValueString(),
		ExternalId:         plan.ExternalId.ValueString(),
		PasswordSalt:       plan.PasswordSalt.ValueString(),
		AvatarType:         plan.AvatarType.ValueString(),
		PermanentAvatar:    plan.PermanentAvatar.ValueString(),
		Address:            address,
		IdCardType:         plan.IdCardType.ValueString(),
		IdCard:             plan.IdCard.ValueString(),
		AccessKey:          plan.AccessKey.ValueString(),
		AccessSecret:       plan.AccessSecret.ValueString(),
		Invitation:         plan.Invitation.ValueString(),
		InvitationCode:     plan.InvitationCode.ValueString(),
		Ldap:               plan.Ldap.ValueString(),
		Properties:         properties,
		NeedUpdatePassword: plan.NeedUpdatePassword.ValueBool(),
		PreferredMfaType:   plan.PreferredMfaType.ValueString(),
		RecoveryCodes:      recoveryCodes,
		TotpSecret:         plan.TotpSecret.ValueString(),
		MfaPhoneEnabled:    plan.MfaPhoneEnabled.ValueBool(),
		MfaEmailEnabled:    plan.MfaEmailEnabled.ValueBool(),
		Groups:             groups,
	}

	success, err := r.client.AddUser(user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating User",
			fmt.Sprintf("Could not create user %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating User",
			fmt.Sprintf("Casdoor returned failure when creating user %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the user to get the generated ID.
	createdUser, err := r.client.GetUser(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User After Create",
			fmt.Sprintf("Could not read user %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdUser != nil {
		plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
		plan.CreatedTime = types.StringValue(createdUser.CreatedTime)
		plan.UpdatedTime = types.StringValue(createdUser.UpdatedTime)
		plan.IsDefaultAvatar = types.BoolValue(createdUser.IsDefaultAvatar)
		plan.IsOnline = types.BoolValue(createdUser.IsOnline)
		plan.Hash = types.StringValue(createdUser.Hash)
		plan.PreHash = types.StringValue(createdUser.PreHash)
		plan.CreatedIp = types.StringValue(createdUser.CreatedIp)
		plan.LastSigninTime = types.StringValue(createdUser.LastSigninTime)
		plan.LastSigninIp = types.StringValue(createdUser.LastSigninIp)
		plan.LastChangePasswordTime = types.StringValue(createdUser.LastChangePasswordTime)
		plan.LastSigninWrongTime = types.StringValue(createdUser.LastSigninWrongTime)
		plan.SigninWrongTimes = types.Int64Value(int64(createdUser.SigninWrongTimes))
		// Preserve server-generated sensitive values.
		if createdUser.PasswordSalt != "***" {
			plan.PasswordSalt = types.StringValue(createdUser.PasswordSalt)
		}
		if createdUser.AccessKey != "***" {
			plan.AccessKey = types.StringValue(createdUser.AccessKey)
		}
		if createdUser.AccessSecret != "***" {
			plan.AccessSecret = types.StringValue(createdUser.AccessSecret)
		}
		if createdUser.TotpSecret != "***" {
			plan.TotpSecret = types.StringValue(createdUser.TotpSecret)
		}
	} else {
		// GetUser uses the provider's OrganizationName which may differ from the
		// user's owner. Set computed fields to defaults to avoid unknown values.
		plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
		plan.CreatedTime = types.StringValue(createdTime)
		plan.UpdatedTime = types.StringValue(createdTime)
		plan.IsDefaultAvatar = types.BoolValue(false)
		plan.IsOnline = types.BoolValue(false)
		plan.Hash = types.StringValue("")
		plan.PreHash = types.StringValue("")
		plan.CreatedIp = types.StringValue("")
		plan.LastSigninTime = types.StringValue("")
		plan.LastSigninIp = types.StringValue("")
		plan.LastChangePasswordTime = types.StringValue("")
		plan.LastSigninWrongTime = types.StringValue("")
		plan.SigninWrongTimes = types.Int64Value(0)
	}

	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
	}
	if len(address) == 0 {
		plan.Address = types.ListNull(types.StringType)
	}
	if len(recoveryCodes) == 0 {
		plan.RecoveryCodes = types.ListNull(types.StringType)
	}
	if len(properties) == 0 {
		plan.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User",
			fmt.Sprintf("Could not read user %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(user.Owner)
	state.Name = types.StringValue(user.Name)
	state.ID = types.StringValue(user.Owner + "/" + user.Name)
	state.Type = types.StringValue(user.Type)
	// Password is write-only; never read it back from the server.
	state.PasswordType = types.StringValue(user.PasswordType)
	state.DisplayName = types.StringValue(user.DisplayName)
	state.FirstName = types.StringValue(user.FirstName)
	state.LastName = types.StringValue(user.LastName)
	state.Avatar = types.StringValue(user.Avatar)
	state.Email = types.StringValue(user.Email)
	state.EmailVerified = types.BoolValue(user.EmailVerified)
	state.Phone = types.StringValue(user.Phone)
	state.CountryCode = types.StringValue(user.CountryCode)
	state.Region = types.StringValue(user.Region)
	state.Location = types.StringValue(user.Location)
	state.Affiliation = types.StringValue(user.Affiliation)
	state.Title = types.StringValue(user.Title)
	state.Homepage = types.StringValue(user.Homepage)
	state.Bio = types.StringValue(user.Bio)
	state.Tag = types.StringValue(user.Tag)
	state.Language = types.StringValue(user.Language)
	state.Gender = types.StringValue(user.Gender)
	state.Birthday = types.StringValue(user.Birthday)
	state.Education = types.StringValue(user.Education)
	state.Score = types.Int64Value(int64(user.Score))
	state.Karma = types.Int64Value(int64(user.Karma))
	state.Ranking = types.Int64Value(int64(user.Ranking))
	state.IsAdmin = types.BoolValue(user.IsAdmin)
	state.IsForbidden = types.BoolValue(user.IsForbidden)
	state.IsDeleted = types.BoolValue(user.IsDeleted)
	state.SignupApplication = types.StringValue(user.SignupApplication)
	state.CreatedTime = types.StringValue(user.CreatedTime)
	state.UpdatedTime = types.StringValue(user.UpdatedTime)
	state.ExternalId = types.StringValue(user.ExternalId)
	if user.PasswordSalt != "***" {
		state.PasswordSalt = types.StringValue(user.PasswordSalt)
	}
	state.AvatarType = types.StringValue(user.AvatarType)
	state.PermanentAvatar = types.StringValue(user.PermanentAvatar)
	state.IdCardType = types.StringValue(user.IdCardType)
	state.IdCard = types.StringValue(user.IdCard)
	state.IsDefaultAvatar = types.BoolValue(user.IsDefaultAvatar)
	state.IsOnline = types.BoolValue(user.IsOnline)
	state.Hash = types.StringValue(user.Hash)
	state.PreHash = types.StringValue(user.PreHash)
	if user.AccessKey != "***" {
		state.AccessKey = types.StringValue(user.AccessKey)
	}
	if user.AccessSecret != "***" {
		state.AccessSecret = types.StringValue(user.AccessSecret)
	}
	state.CreatedIp = types.StringValue(user.CreatedIp)
	state.LastSigninTime = types.StringValue(user.LastSigninTime)
	state.LastSigninIp = types.StringValue(user.LastSigninIp)
	state.Invitation = types.StringValue(user.Invitation)
	state.InvitationCode = types.StringValue(user.InvitationCode)
	state.Ldap = types.StringValue(user.Ldap)
	state.NeedUpdatePassword = types.BoolValue(user.NeedUpdatePassword)
	state.LastChangePasswordTime = types.StringValue(user.LastChangePasswordTime)
	state.LastSigninWrongTime = types.StringValue(user.LastSigninWrongTime)
	state.SigninWrongTimes = types.Int64Value(int64(user.SigninWrongTimes))
	state.PreferredMfaType = types.StringValue(user.PreferredMfaType)
	if user.TotpSecret != "***" {
		state.TotpSecret = types.StringValue(user.TotpSecret)
	}
	state.MfaPhoneEnabled = types.BoolValue(user.MfaPhoneEnabled)
	state.MfaEmailEnabled = types.BoolValue(user.MfaEmailEnabled)

	if len(user.Address) > 0 {
		addressList, diags := types.ListValueFrom(ctx, types.StringType, user.Address)
		resp.Diagnostics.Append(diags...)
		state.Address = addressList
	} else {
		state.Address = types.ListNull(types.StringType)
	}

	if len(user.RecoveryCodes) > 0 {
		rcList, diags := types.ListValueFrom(ctx, types.StringType, user.RecoveryCodes)
		resp.Diagnostics.Append(diags...)
		state.RecoveryCodes = rcList
	} else {
		state.RecoveryCodes = types.ListNull(types.StringType)
	}

	if len(user.Properties) > 0 {
		props, diags := types.MapValueFrom(ctx, types.StringType, user.Properties)
		resp.Diagnostics.Append(diags...)
		state.Properties = props
	} else {
		state.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if len(user.Groups) > 0 {
		groups, diags := types.ListValueFrom(ctx, types.StringType, user.Groups)
		resp.Diagnostics.Append(diags...)
		state.Groups = groups
	} else {
		state.Groups = types.ListNull(types.StringType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []string
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	}
	var address []string
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		resp.Diagnostics.Append(plan.Address.ElementsAs(ctx, &address, false)...)
	}
	var recoveryCodes []string
	if !plan.RecoveryCodes.IsNull() && !plan.RecoveryCodes.IsUnknown() {
		resp.Diagnostics.Append(plan.RecoveryCodes.ElementsAs(ctx, &recoveryCodes, false)...)
	}
	var properties map[string]string
	if !plan.Properties.IsNull() && !plan.Properties.IsUnknown() {
		properties = make(map[string]string)
		resp.Diagnostics.Append(plan.Properties.ElementsAs(ctx, &properties, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the existing user to preserve the Casdoor-internal Id field,
	// which is immutable and must not change during updates.
	existingUser, err := r.client.GetUser(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User Before Update",
			fmt.Sprintf("Could not read user %q before update: %s", plan.Name.ValueString(), err),
		)
		return
	}

	var internalID string
	if existingUser != nil {
		internalID = existingUser.Id
	}

	user := &casdoorsdk.User{
		Owner:              plan.Owner.ValueString(),
		Name:               plan.Name.ValueString(),
		Id:                 internalID,
		Type:               plan.Type.ValueString(),
		Password:           plan.Password.ValueString(),
		PasswordType:       plan.PasswordType.ValueString(),
		DisplayName:        plan.DisplayName.ValueString(),
		FirstName:          plan.FirstName.ValueString(),
		LastName:           plan.LastName.ValueString(),
		Avatar:             plan.Avatar.ValueString(),
		Email:              plan.Email.ValueString(),
		EmailVerified:      plan.EmailVerified.ValueBool(),
		Phone:              plan.Phone.ValueString(),
		CountryCode:        plan.CountryCode.ValueString(),
		Region:             plan.Region.ValueString(),
		Location:           plan.Location.ValueString(),
		Affiliation:        plan.Affiliation.ValueString(),
		Title:              plan.Title.ValueString(),
		Homepage:           plan.Homepage.ValueString(),
		Bio:                plan.Bio.ValueString(),
		Tag:                plan.Tag.ValueString(),
		Language:           plan.Language.ValueString(),
		Gender:             plan.Gender.ValueString(),
		Birthday:           plan.Birthday.ValueString(),
		Education:          plan.Education.ValueString(),
		Score:              int(plan.Score.ValueInt64()),
		Karma:              int(plan.Karma.ValueInt64()),
		Ranking:            int(plan.Ranking.ValueInt64()),
		IsAdmin:            plan.IsAdmin.ValueBool(),
		IsForbidden:        plan.IsForbidden.ValueBool(),
		IsDeleted:          plan.IsDeleted.ValueBool(),
		SignupApplication:  plan.SignupApplication.ValueString(),
		CreatedTime:        plan.CreatedTime.ValueString(),
		ExternalId:         plan.ExternalId.ValueString(),
		PasswordSalt:       plan.PasswordSalt.ValueString(),
		AvatarType:         plan.AvatarType.ValueString(),
		PermanentAvatar:    plan.PermanentAvatar.ValueString(),
		Address:            address,
		IdCardType:         plan.IdCardType.ValueString(),
		IdCard:             plan.IdCard.ValueString(),
		AccessKey:          plan.AccessKey.ValueString(),
		AccessSecret:       plan.AccessSecret.ValueString(),
		Invitation:         plan.Invitation.ValueString(),
		InvitationCode:     plan.InvitationCode.ValueString(),
		Ldap:               plan.Ldap.ValueString(),
		Properties:         properties,
		NeedUpdatePassword: plan.NeedUpdatePassword.ValueBool(),
		PreferredMfaType:   plan.PreferredMfaType.ValueString(),
		RecoveryCodes:      recoveryCodes,
		TotpSecret:         plan.TotpSecret.ValueString(),
		MfaPhoneEnabled:    plan.MfaPhoneEnabled.ValueBool(),
		MfaEmailEnabled:    plan.MfaEmailEnabled.ValueBool(),
		Groups:             groups,
	}

	success, err := r.client.UpdateUser(user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User",
			fmt.Sprintf("Could not update user %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating User",
			fmt.Sprintf("Casdoor returned failure when updating user %q", plan.Name.ValueString()),
		)
		return
	}

	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
	}
	if len(address) == 0 {
		plan.Address = types.ListNull(types.StringType)
	}
	if len(recoveryCodes) == 0 {
		plan.RecoveryCodes = types.ListNull(types.StringType)
	}
	if len(properties) == 0 {
		plan.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &casdoorsdk.User{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	_, err := r.client.DeleteUser(user)
	if err != nil {
		// Casdoor returns "session is nil" when deleting users from the built-in
		// organization. This is a known server-side bug; treat it as a warning
		// rather than a hard error.
		if strings.Contains(err.Error(), "session is nil") {
			resp.Diagnostics.AddWarning(
				"User Deletion Warning",
				fmt.Sprintf("Casdoor returned 'session is nil' when deleting user %q. "+
					"This is a known Casdoor issue with the built-in organization.", state.Name.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting User",
			fmt.Sprintf("Could not delete user %q: %s", state.Name.ValueString(), err),
		)
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
