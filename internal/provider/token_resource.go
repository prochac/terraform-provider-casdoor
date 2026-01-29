// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

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
	_ resource.Resource                = &TokenResource{}
	_ resource.ResourceWithConfigure   = &TokenResource{}
	_ resource.ResourceWithImportState = &TokenResource{}
)

type TokenResource struct {
	client *casdoorsdk.Client
}

type TokenResourceModel struct {
	Owner            types.String `tfsdk:"owner"`
	Name             types.String `tfsdk:"name"`
	CreatedTime      types.String `tfsdk:"created_time"`
	Application      types.String `tfsdk:"application"`
	Organization     types.String `tfsdk:"organization"`
	User             types.String `tfsdk:"user"`
	Code             types.String `tfsdk:"code"`
	AccessToken      types.String `tfsdk:"access_token"`
	RefreshToken     types.String `tfsdk:"refresh_token"`
	AccessTokenHash  types.String `tfsdk:"access_token_hash"`
	RefreshTokenHash types.String `tfsdk:"refresh_token_hash"`
	ExpiresIn        types.Int64  `tfsdk:"expires_in"`
	Scope            types.String `tfsdk:"scope"`
	TokenType        types.String `tfsdk:"token_type"`
	CodeChallenge    types.String `tfsdk:"code_challenge"`
	CodeIsUsed       types.Bool   `tfsdk:"code_is_used"`
	CodeExpireIn     types.Int64  `tfsdk:"code_expire_in"`
}

func NewTokenResource() resource.Resource {
	return &TokenResource{}
}

func (r *TokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token"
}

func (r *TokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor token.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the token was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application": schema.StringAttribute{
				Description: "The application this token belongs to.",
				Required:    true,
			},
			"organization": schema.StringAttribute{
				Description: "The organization this token belongs to.",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "The user this token belongs to.",
				Required:    true,
			},
			"code": schema.StringAttribute{
				Description: "The authorization code.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"access_token": schema.StringAttribute{
				Description: "The access token.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"refresh_token": schema.StringAttribute{
				Description: "The refresh token.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"access_token_hash": schema.StringAttribute{
				Description: "The hash of the access token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"refresh_token_hash": schema.StringAttribute{
				Description: "The hash of the refresh token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_in": schema.Int64Attribute{
				Description: "Token expiration time in seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(7200),
			},
			"scope": schema.StringAttribute{
				Description: "The scope of the token.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"token_type": schema.StringAttribute{
				Description: "The type of the token (e.g., 'Bearer').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Bearer"),
			},
			"code_challenge": schema.StringAttribute{
				Description: "The PKCE code challenge.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"code_is_used": schema.BoolAttribute{
				Description: "Whether the authorization code has been used.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"code_expire_in": schema.Int64Attribute{
				Description: "Code expiration time in seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
		},
	}
}

func (r *TokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TokenResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	token := &casdoorsdk.Token{
		Owner:         plan.Owner.ValueString(),
		Name:          plan.Name.ValueString(),
		CreatedTime:   createdTime,
		Application:   plan.Application.ValueString(),
		Organization:  plan.Organization.ValueString(),
		User:          plan.User.ValueString(),
		Code:          plan.Code.ValueString(),
		AccessToken:   plan.AccessToken.ValueString(),
		RefreshToken:  plan.RefreshToken.ValueString(),
		ExpiresIn:     int(plan.ExpiresIn.ValueInt64()),
		Scope:         plan.Scope.ValueString(),
		TokenType:     plan.TokenType.ValueString(),
		CodeChallenge: plan.CodeChallenge.ValueString(),
		CodeIsUsed:    plan.CodeIsUsed.ValueBool(),
		CodeExpireIn:  plan.CodeExpireIn.ValueInt64(),
	}

	success, err := r.client.AddToken(token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Token",
			fmt.Sprintf("Could not create token %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Token",
			fmt.Sprintf("Casdoor returned failure when creating token %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the token to get generated values.
	createdToken, err := r.client.GetToken(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Token After Create",
			fmt.Sprintf("Could not read token %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdToken != nil {
		plan.CreatedTime = types.StringValue(createdToken.CreatedTime)
		plan.AccessToken = types.StringValue(createdToken.AccessToken)
		plan.RefreshToken = types.StringValue(createdToken.RefreshToken)
		plan.AccessTokenHash = types.StringValue(createdToken.AccessTokenHash)
		plan.RefreshTokenHash = types.StringValue(createdToken.RefreshTokenHash)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TokenResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetToken(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Token",
			fmt.Sprintf("Could not read token %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if token == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(token.Owner)
	state.Name = types.StringValue(token.Name)
	state.CreatedTime = types.StringValue(token.CreatedTime)
	state.Application = types.StringValue(token.Application)
	state.Organization = types.StringValue(token.Organization)
	state.User = types.StringValue(token.User)
	state.Code = types.StringValue(token.Code)
	state.AccessToken = types.StringValue(token.AccessToken)
	state.RefreshToken = types.StringValue(token.RefreshToken)
	state.AccessTokenHash = types.StringValue(token.AccessTokenHash)
	state.RefreshTokenHash = types.StringValue(token.RefreshTokenHash)
	state.ExpiresIn = types.Int64Value(int64(token.ExpiresIn))
	state.Scope = types.StringValue(token.Scope)
	state.TokenType = types.StringValue(token.TokenType)
	state.CodeChallenge = types.StringValue(token.CodeChallenge)
	state.CodeIsUsed = types.BoolValue(token.CodeIsUsed)
	state.CodeExpireIn = types.Int64Value(token.CodeExpireIn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TokenResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := &casdoorsdk.Token{
		Owner:         plan.Owner.ValueString(),
		Name:          plan.Name.ValueString(),
		CreatedTime:   plan.CreatedTime.ValueString(),
		Application:   plan.Application.ValueString(),
		Organization:  plan.Organization.ValueString(),
		User:          plan.User.ValueString(),
		Code:          plan.Code.ValueString(),
		AccessToken:   plan.AccessToken.ValueString(),
		RefreshToken:  plan.RefreshToken.ValueString(),
		ExpiresIn:     int(plan.ExpiresIn.ValueInt64()),
		Scope:         plan.Scope.ValueString(),
		TokenType:     plan.TokenType.ValueString(),
		CodeChallenge: plan.CodeChallenge.ValueString(),
		CodeIsUsed:    plan.CodeIsUsed.ValueBool(),
		CodeExpireIn:  plan.CodeExpireIn.ValueInt64(),
	}

	success, err := r.client.UpdateToken(token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Token",
			fmt.Sprintf("Could not update token %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Token",
			fmt.Sprintf("Casdoor returned failure when updating token %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TokenResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := &casdoorsdk.Token{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	_, err := r.client.DeleteToken(token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Token",
			fmt.Sprintf("Could not delete token %q: %s", state.Name.ValueString(), err),
		)
		return
	}
}

func (r *TokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
