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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &PermissionResource{}
	_ resource.ResourceWithConfigure   = &PermissionResource{}
	_ resource.ResourceWithImportState = &PermissionResource{}
)

type PermissionResource struct {
	client *casdoorsdk.Client
}

type PermissionResourceModel struct {
	Owner        types.String `tfsdk:"owner"`
	Name         types.String `tfsdk:"name"`
	CreatedTime  types.String `tfsdk:"created_time"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	Users        types.List   `tfsdk:"users"`
	Groups       types.List   `tfsdk:"groups"`
	Roles        types.List   `tfsdk:"roles"`
	Domains      types.List   `tfsdk:"domains"`
	Model        types.String `tfsdk:"model"`
	Adapter      types.String `tfsdk:"adapter"`
	ResourceType types.String `tfsdk:"resource_type"`
	Resources    types.List   `tfsdk:"resources"`
	Actions      types.List   `tfsdk:"actions"`
	Effect       types.String `tfsdk:"effect"`
	IsEnabled    types.Bool   `tfsdk:"is_enabled"`
	Submitter    types.String `tfsdk:"submitter"`
	Approver     types.String `tfsdk:"approver"`
	ApproveTime  types.String `tfsdk:"approve_time"`
	State        types.String `tfsdk:"state"`
}

func NewPermissionResource() resource.Resource {
	return &PermissionResource{}
}

func (r *PermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *PermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor permission.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this permission.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the permission.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the permission was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "A description of the permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"users": schema.ListAttribute{
				Description: "List of users this permission applies to (format: 'organization/username').",
				Optional:    true,
				ElementType: types.StringType,
			},
			"groups": schema.ListAttribute{
				Description: "List of groups this permission applies to.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"roles": schema.ListAttribute{
				Description: "List of roles this permission applies to.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"domains": schema.ListAttribute{
				Description: "List of domains where this permission applies.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"model": schema.StringAttribute{
				Description: "The Casbin model for this permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"adapter": schema.StringAttribute{
				Description: "The Casbin adapter for this permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"resource_type": schema.StringAttribute{
				Description: "The type of resource this permission controls.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"resources": schema.ListAttribute{
				Description: "List of resources this permission controls.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"actions": schema.ListAttribute{
				Description: "List of actions allowed by this permission (e.g., 'Read', 'Write', 'Admin').",
				Optional:    true,
				ElementType: types.StringType,
			},
			"effect": schema.StringAttribute{
				Description: "The effect of this permission ('Allow' or 'Deny').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Allow"),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the permission is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"submitter": schema.StringAttribute{
				Description: "The user who submitted this permission for approval.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"approver": schema.StringAttribute{
				Description: "The user who approved this permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"approve_time": schema.StringAttribute{
				Description: "The time when this permission was approved.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"state": schema.StringAttribute{
				Description: "The approval state of this permission.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *PermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PermissionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users, groups, roles, domains, resources, actions []string

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
	if !plan.Resources.IsNull() {
		resp.Diagnostics.Append(plan.Resources.ElementsAs(ctx, &resources, false)...)
	}
	if !plan.Actions.IsNull() {
		resp.Diagnostics.Append(plan.Actions.ElementsAs(ctx, &actions, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	permission := &casdoorsdk.Permission{
		Owner:        plan.Owner.ValueString(),
		Name:         plan.Name.ValueString(),
		CreatedTime:  createdTime,
		DisplayName:  plan.DisplayName.ValueString(),
		Description:  plan.Description.ValueString(),
		Users:        users,
		Groups:       groups,
		Roles:        roles,
		Domains:      domains,
		Model:        plan.Model.ValueString(),
		Adapter:      plan.Adapter.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		Resources:    resources,
		Actions:      actions,
		Effect:       plan.Effect.ValueString(),
		IsEnabled:    plan.IsEnabled.ValueBool(),
		Submitter:    plan.Submitter.ValueString(),
		Approver:     plan.Approver.ValueString(),
		ApproveTime:  plan.ApproveTime.ValueString(),
		State:        plan.State.ValueString(),
	}

	success, err := r.client.AddPermission(permission)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Permission",
			fmt.Sprintf("Could not create permission %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Permission",
			fmt.Sprintf("Casdoor returned failure when creating permission %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the permission to get server-generated values like CreatedTime.
	createdPermission, err := r.client.GetPermission(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Permission",
			fmt.Sprintf("Could not read permission %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdPermission != nil {
		plan.CreatedTime = types.StringValue(createdPermission.CreatedTime)
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
	if len(resources) == 0 {
		plan.Resources = types.ListNull(types.StringType)
	}
	if len(actions) == 0 {
		plan.Actions = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PermissionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.GetPermission(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Permission",
			fmt.Sprintf("Could not read permission %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if permission == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(permission.Owner)
	state.Name = types.StringValue(permission.Name)
	state.CreatedTime = types.StringValue(permission.CreatedTime)
	state.DisplayName = types.StringValue(permission.DisplayName)
	state.Description = types.StringValue(permission.Description)
	state.Model = types.StringValue(permission.Model)
	state.Adapter = types.StringValue(permission.Adapter)
	state.ResourceType = types.StringValue(permission.ResourceType)
	state.Effect = types.StringValue(permission.Effect)
	state.IsEnabled = types.BoolValue(permission.IsEnabled)
	state.Submitter = types.StringValue(permission.Submitter)
	state.Approver = types.StringValue(permission.Approver)
	state.ApproveTime = types.StringValue(permission.ApproveTime)
	state.State = types.StringValue(permission.State)

	if len(permission.Users) > 0 {
		users, diags := types.ListValueFrom(ctx, types.StringType, permission.Users)
		resp.Diagnostics.Append(diags...)
		state.Users = users
	} else {
		state.Users = types.ListNull(types.StringType)
	}

	if len(permission.Groups) > 0 {
		groups, diags := types.ListValueFrom(ctx, types.StringType, permission.Groups)
		resp.Diagnostics.Append(diags...)
		state.Groups = groups
	} else {
		state.Groups = types.ListNull(types.StringType)
	}

	if len(permission.Roles) > 0 {
		roles, diags := types.ListValueFrom(ctx, types.StringType, permission.Roles)
		resp.Diagnostics.Append(diags...)
		state.Roles = roles
	} else {
		state.Roles = types.ListNull(types.StringType)
	}

	if len(permission.Domains) > 0 {
		domains, diags := types.ListValueFrom(ctx, types.StringType, permission.Domains)
		resp.Diagnostics.Append(diags...)
		state.Domains = domains
	} else {
		state.Domains = types.ListNull(types.StringType)
	}

	if len(permission.Resources) > 0 {
		resources, diags := types.ListValueFrom(ctx, types.StringType, permission.Resources)
		resp.Diagnostics.Append(diags...)
		state.Resources = resources
	} else {
		state.Resources = types.ListNull(types.StringType)
	}

	if len(permission.Actions) > 0 {
		actions, diags := types.ListValueFrom(ctx, types.StringType, permission.Actions)
		resp.Diagnostics.Append(diags...)
		state.Actions = actions
	} else {
		state.Actions = types.ListNull(types.StringType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PermissionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users, groups, roles, domains, resources, actions []string

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
	if !plan.Resources.IsNull() {
		resp.Diagnostics.Append(plan.Resources.ElementsAs(ctx, &resources, false)...)
	}
	if !plan.Actions.IsNull() {
		resp.Diagnostics.Append(plan.Actions.ElementsAs(ctx, &actions, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	permission := &casdoorsdk.Permission{
		Owner:        plan.Owner.ValueString(),
		Name:         plan.Name.ValueString(),
		CreatedTime:  plan.CreatedTime.ValueString(),
		DisplayName:  plan.DisplayName.ValueString(),
		Description:  plan.Description.ValueString(),
		Users:        users,
		Groups:       groups,
		Roles:        roles,
		Domains:      domains,
		Model:        plan.Model.ValueString(),
		Adapter:      plan.Adapter.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		Resources:    resources,
		Actions:      actions,
		Effect:       plan.Effect.ValueString(),
		IsEnabled:    plan.IsEnabled.ValueBool(),
		Submitter:    plan.Submitter.ValueString(),
		Approver:     plan.Approver.ValueString(),
		ApproveTime:  plan.ApproveTime.ValueString(),
		State:        plan.State.ValueString(),
	}

	success, err := r.client.UpdatePermission(permission)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Permission",
			fmt.Sprintf("Could not update permission %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Permission",
			fmt.Sprintf("Casdoor returned failure when updating permission %q", plan.Name.ValueString()),
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
	if len(resources) == 0 {
		plan.Resources = types.ListNull(types.StringType)
	}
	if len(actions) == 0 {
		plan.Actions = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PermissionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := &casdoorsdk.Permission{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeletePermission(permission)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Permission",
			fmt.Sprintf("Could not delete permission %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Permission",
			fmt.Sprintf("Casdoor returned failure when deleting permission %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *PermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
