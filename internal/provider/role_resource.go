// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &RoleResource{}
	_ resource.ResourceWithConfigure   = &RoleResource{}
	_ resource.ResourceWithImportState = &RoleResource{}
)

type RoleResource struct {
	client *casdoorsdk.Client
}

type RoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Owner       types.String `tfsdk:"owner"`
	Name        types.String `tfsdk:"name"`
	CreatedTime types.String `tfsdk:"created_time"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Users       types.List   `tfsdk:"users"`
	Groups      types.List   `tfsdk:"groups"`
	Roles       types.List   `tfsdk:"roles"`
	Domains     types.List   `tfsdk:"domains"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
}

func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the role in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this role.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the role.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the role was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the role.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "A description of the role.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"users": schema.ListAttribute{
				Description: "List of users assigned to this role (format: 'organization/username').",
				Optional:    true,
				ElementType: types.StringType,
			},
			"groups": schema.ListAttribute{
				Description: "List of groups assigned to this role.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"roles": schema.ListAttribute{
				Description: "List of sub-roles (for role hierarchy).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"domains": schema.ListAttribute{
				Description: "List of domains where this role applies.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the role is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users, groups, roles, domains []string

	if !plan.Users.IsNull() {
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
	}
	if !plan.Groups.IsNull() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	}
	if !plan.Roles.IsNull() {
		resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, false)...)
	}
	if !plan.Domains.IsNull() {
		resp.Diagnostics.Append(plan.Domains.ElementsAs(ctx, &domains, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	role := &casdoorsdk.Role{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: createdTime,
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Users:       users,
		Groups:      groups,
		Roles:       roles,
		Domains:     domains,
		IsEnabled:   plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.AddRole(role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Role",
			fmt.Sprintf("Could not create role %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Role",
			fmt.Sprintf("Casdoor returned failure when creating role %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the role to get server-generated values like CreatedTime.
	createdRole, err := r.client.GetRole(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role",
			fmt.Sprintf("Could not read role %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdRole != nil {
		plan.CreatedTime = types.StringValue(createdRole.CreatedTime)
	}

	// Set list values to null if empty.
	if len(users) == 0 {
		plan.Users = types.ListNull(types.StringType)
	}
	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
	}
	if len(roles) == 0 {
		plan.Roles = types.ListNull(types.StringType)
	}
	if len(domains) == 0 {
		plan.Domains = types.ListNull(types.StringType)
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.client.GetRole(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role",
			fmt.Sprintf("Could not read role %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if role == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(role.Owner + "/" + role.Name)
	state.Owner = types.StringValue(role.Owner)
	state.Name = types.StringValue(role.Name)
	state.DisplayName = types.StringValue(role.DisplayName)
	state.Description = types.StringValue(role.Description)
	state.CreatedTime = types.StringValue(role.CreatedTime)
	state.IsEnabled = types.BoolValue(role.IsEnabled)

	if len(role.Users) > 0 {
		users, diags := types.ListValueFrom(ctx, types.StringType, role.Users)
		resp.Diagnostics.Append(diags...)
		state.Users = users
	} else {
		state.Users = types.ListNull(types.StringType)
	}

	if len(role.Groups) > 0 {
		groups, diags := types.ListValueFrom(ctx, types.StringType, role.Groups)
		resp.Diagnostics.Append(diags...)
		state.Groups = groups
	} else {
		state.Groups = types.ListNull(types.StringType)
	}

	if len(role.Roles) > 0 {
		roles, diags := types.ListValueFrom(ctx, types.StringType, role.Roles)
		resp.Diagnostics.Append(diags...)
		state.Roles = roles
	} else {
		state.Roles = types.ListNull(types.StringType)
	}

	if len(role.Domains) > 0 {
		domains, diags := types.ListValueFrom(ctx, types.StringType, role.Domains)
		resp.Diagnostics.Append(diags...)
		state.Domains = domains
	} else {
		state.Domains = types.ListNull(types.StringType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users, groups, roles, domains []string

	if !plan.Users.IsNull() {
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
	}
	if !plan.Groups.IsNull() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	}
	if !plan.Roles.IsNull() {
		resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, false)...)
	}
	if !plan.Domains.IsNull() {
		resp.Diagnostics.Append(plan.Domains.ElementsAs(ctx, &domains, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	role := &casdoorsdk.Role{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: plan.CreatedTime.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Users:       users,
		Groups:      groups,
		Roles:       roles,
		Domains:     domains,
		IsEnabled:   plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.UpdateRole(role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			fmt.Sprintf("Could not update role %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			fmt.Sprintf("Casdoor returned failure when updating role %q", plan.Name.ValueString()),
		)
		return
	}

	// Set list values to null if empty.
	if len(users) == 0 {
		plan.Users = types.ListNull(types.StringType)
	}
	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
	}
	if len(roles) == 0 {
		plan.Roles = types.ListNull(types.StringType)
	}
	if len(domains) == 0 {
		plan.Domains = types.ListNull(types.StringType)
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := &casdoorsdk.Role{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteRole(role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Role",
			fmt.Sprintf("Could not delete role %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Role",
			fmt.Sprintf("Casdoor returned failure when deleting role %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
