// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	ID           types.String `tfsdk:"id"`
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
			"id": schema.StringAttribute{
				Description: "The ID of the permission in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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

	users, diags := stringListToSDK(ctx, plan.Users)
	resp.Diagnostics.Append(diags...)
	groups, diags := stringListToSDK(ctx, plan.Groups)
	resp.Diagnostics.Append(diags...)
	roles, diags := stringListToSDK(ctx, plan.Roles)
	resp.Diagnostics.Append(diags...)
	domains, diags := stringListToSDK(ctx, plan.Domains)
	resp.Diagnostics.Append(diags...)
	resources, diags := stringListToSDK(ctx, plan.Resources)
	resp.Diagnostics.Append(diags...)
	actions, diags := stringListToSDK(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
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

	ok, err := r.client.AddPermission(permission)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating permission %q", plan.Name.ValueString())) {
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
	plan.Users, diags = stringListFromSDK(ctx, users)
	resp.Diagnostics.Append(diags...)
	plan.Groups, diags = stringListFromSDK(ctx, groups)
	resp.Diagnostics.Append(diags...)
	plan.Roles, diags = stringListFromSDK(ctx, roles)
	resp.Diagnostics.Append(diags...)
	plan.Domains, diags = stringListFromSDK(ctx, domains)
	resp.Diagnostics.Append(diags...)
	plan.Resources, diags = stringListFromSDK(ctx, resources)
	resp.Diagnostics.Append(diags...)
	plan.Actions, diags = stringListFromSDK(ctx, actions)
	resp.Diagnostics.Append(diags...)

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PermissionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := getByOwnerName[casdoorsdk.Permission](r.client, "get-permission", state.Owner.ValueString(), state.Name.ValueString())
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

	state.ID = types.StringValue(permission.Owner + "/" + permission.Name)
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

	var diags diag.Diagnostics

	state.Users, diags = stringListFromSDK(ctx, permission.Users)
	resp.Diagnostics.Append(diags...)
	state.Groups, diags = stringListFromSDK(ctx, permission.Groups)
	resp.Diagnostics.Append(diags...)
	state.Roles, diags = stringListFromSDK(ctx, permission.Roles)
	resp.Diagnostics.Append(diags...)
	state.Domains, diags = stringListFromSDK(ctx, permission.Domains)
	resp.Diagnostics.Append(diags...)
	state.Resources, diags = stringListFromSDK(ctx, permission.Resources)
	resp.Diagnostics.Append(diags...)
	state.Actions, diags = stringListFromSDK(ctx, permission.Actions)
	resp.Diagnostics.Append(diags...)
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

	users, diags := stringListToSDK(ctx, plan.Users)
	resp.Diagnostics.Append(diags...)
	groups, diags := stringListToSDK(ctx, plan.Groups)
	resp.Diagnostics.Append(diags...)
	roles, diags := stringListToSDK(ctx, plan.Roles)
	resp.Diagnostics.Append(diags...)
	domains, diags := stringListToSDK(ctx, plan.Domains)
	resp.Diagnostics.Append(diags...)
	resources, diags := stringListToSDK(ctx, plan.Resources)
	resp.Diagnostics.Append(diags...)
	actions, diags := stringListToSDK(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
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

	ok, err := r.client.UpdatePermission(permission)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("updating permission %q", plan.Name.ValueString())) {
		return
	}

	// Set list values to null if empty.
	plan.Users, diags = stringListFromSDK(ctx, users)
	resp.Diagnostics.Append(diags...)
	plan.Groups, diags = stringListFromSDK(ctx, groups)
	resp.Diagnostics.Append(diags...)
	plan.Roles, diags = stringListFromSDK(ctx, roles)
	resp.Diagnostics.Append(diags...)
	plan.Domains, diags = stringListFromSDK(ctx, domains)
	resp.Diagnostics.Append(diags...)
	plan.Resources, diags = stringListFromSDK(ctx, resources)
	resp.Diagnostics.Append(diags...)
	plan.Actions, diags = stringListFromSDK(ctx, actions)
	resp.Diagnostics.Append(diags...)

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
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

	ok, err := r.client.DeletePermission(permission)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("deleting permission %q", state.Name.ValueString())) {
		return
	}
}

func (r *PermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
