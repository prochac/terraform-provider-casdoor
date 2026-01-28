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
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

type UserResource struct {
	client *casdoorsdk.Client
}

type UserResourceModel struct {
	Owner             types.String `tfsdk:"owner"`
	Name              types.String `tfsdk:"name"`
	ID                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Password          types.String `tfsdk:"password"`
	PasswordType      types.String `tfsdk:"password_type"`
	DisplayName       types.String `tfsdk:"display_name"`
	FirstName         types.String `tfsdk:"first_name"`
	LastName          types.String `tfsdk:"last_name"`
	Avatar            types.String `tfsdk:"avatar"`
	Email             types.String `tfsdk:"email"`
	EmailVerified     types.Bool   `tfsdk:"email_verified"`
	Phone             types.String `tfsdk:"phone"`
	CountryCode       types.String `tfsdk:"country_code"`
	Region            types.String `tfsdk:"region"`
	Location          types.String `tfsdk:"location"`
	Affiliation       types.String `tfsdk:"affiliation"`
	Title             types.String `tfsdk:"title"`
	Homepage          types.String `tfsdk:"homepage"`
	Bio               types.String `tfsdk:"bio"`
	Tag               types.String `tfsdk:"tag"`
	Language          types.String `tfsdk:"language"`
	Gender            types.String `tfsdk:"gender"`
	Birthday          types.String `tfsdk:"birthday"`
	Education         types.String `tfsdk:"education"`
	Score             types.Int64  `tfsdk:"score"`
	Karma             types.Int64  `tfsdk:"karma"`
	Ranking           types.Int64  `tfsdk:"ranking"`
	IsAdmin           types.Bool   `tfsdk:"is_admin"`
	IsForbidden       types.Bool   `tfsdk:"is_forbidden"`
	IsDeleted         types.Bool   `tfsdk:"is_deleted"`
	SignupApplication types.String `tfsdk:"signup_application"`
	Groups            types.List   `tfsdk:"groups"`
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
				Description: "The user ID (auto-generated if not provided).",
				Optional:    true,
				Computed:    true,
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
	if !plan.Groups.IsNull() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	user := &casdoorsdk.User{
		Owner:             plan.Owner.ValueString(),
		Name:              plan.Name.ValueString(),
		Id:                plan.ID.ValueString(),
		Type:              plan.Type.ValueString(),
		Password:          plan.Password.ValueString(),
		PasswordType:      plan.PasswordType.ValueString(),
		DisplayName:       plan.DisplayName.ValueString(),
		FirstName:         plan.FirstName.ValueString(),
		LastName:          plan.LastName.ValueString(),
		Avatar:            plan.Avatar.ValueString(),
		Email:             plan.Email.ValueString(),
		EmailVerified:     plan.EmailVerified.ValueBool(),
		Phone:             plan.Phone.ValueString(),
		CountryCode:       plan.CountryCode.ValueString(),
		Region:            plan.Region.ValueString(),
		Location:          plan.Location.ValueString(),
		Affiliation:       plan.Affiliation.ValueString(),
		Title:             plan.Title.ValueString(),
		Homepage:          plan.Homepage.ValueString(),
		Bio:               plan.Bio.ValueString(),
		Tag:               plan.Tag.ValueString(),
		Language:          plan.Language.ValueString(),
		Gender:            plan.Gender.ValueString(),
		Birthday:          plan.Birthday.ValueString(),
		Education:         plan.Education.ValueString(),
		Score:             int(plan.Score.ValueInt64()),
		Karma:             int(plan.Karma.ValueInt64()),
		Ranking:           int(plan.Ranking.ValueInt64()),
		IsAdmin:           plan.IsAdmin.ValueBool(),
		IsForbidden:       plan.IsForbidden.ValueBool(),
		IsDeleted:         plan.IsDeleted.ValueBool(),
		SignupApplication: plan.SignupApplication.ValueString(),
		Groups:            groups,
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
		plan.ID = types.StringValue(createdUser.Id)
	}

	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
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
	state.ID = types.StringValue(user.Id)
	state.Type = types.StringValue(user.Type)
	// Password is always masked by Casdoor API ("***"), preserve from state.
	if user.Password != "***" {
		state.Password = types.StringValue(user.Password)
	}
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
	if !plan.Groups.IsNull() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	user := &casdoorsdk.User{
		Owner:             plan.Owner.ValueString(),
		Name:              plan.Name.ValueString(),
		Id:                plan.ID.ValueString(),
		Type:              plan.Type.ValueString(),
		Password:          plan.Password.ValueString(),
		PasswordType:      plan.PasswordType.ValueString(),
		DisplayName:       plan.DisplayName.ValueString(),
		FirstName:         plan.FirstName.ValueString(),
		LastName:          plan.LastName.ValueString(),
		Avatar:            plan.Avatar.ValueString(),
		Email:             plan.Email.ValueString(),
		EmailVerified:     plan.EmailVerified.ValueBool(),
		Phone:             plan.Phone.ValueString(),
		CountryCode:       plan.CountryCode.ValueString(),
		Region:            plan.Region.ValueString(),
		Location:          plan.Location.ValueString(),
		Affiliation:       plan.Affiliation.ValueString(),
		Title:             plan.Title.ValueString(),
		Homepage:          plan.Homepage.ValueString(),
		Bio:               plan.Bio.ValueString(),
		Tag:               plan.Tag.ValueString(),
		Language:          plan.Language.ValueString(),
		Gender:            plan.Gender.ValueString(),
		Birthday:          plan.Birthday.ValueString(),
		Education:         plan.Education.ValueString(),
		Score:             int(plan.Score.ValueInt64()),
		Karma:             int(plan.Karma.ValueInt64()),
		Ranking:           int(plan.Ranking.ValueInt64()),
		IsAdmin:           plan.IsAdmin.ValueBool(),
		IsForbidden:       plan.IsForbidden.ValueBool(),
		IsDeleted:         plan.IsDeleted.ValueBool(),
		SignupApplication: plan.SignupApplication.ValueString(),
		Groups:            groups,
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

	success, err := r.client.DeleteUser(user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting User",
			fmt.Sprintf("Could not delete user %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting User",
			fmt.Sprintf("Casdoor returned failure when deleting user %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
