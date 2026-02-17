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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &GroupResource{}
	_ resource.ResourceWithConfigure   = &GroupResource{}
	_ resource.ResourceWithImportState = &GroupResource{}
)

type GroupResource struct {
	client *casdoorsdk.Client
}

type GroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Owner        types.String `tfsdk:"owner"`
	Name         types.String `tfsdk:"name"`
	CreatedTime  types.String `tfsdk:"created_time"`
	UpdatedTime  types.String `tfsdk:"updated_time"`
	DisplayName  types.String `tfsdk:"display_name"`
	Manager      types.String `tfsdk:"manager"`
	ContactEmail types.String `tfsdk:"contact_email"`
	Type         types.String `tfsdk:"type"`
	ParentId     types.String `tfsdk:"parent_id"`
	ParentName   types.String `tfsdk:"parent_name"`
	Title        types.String `tfsdk:"title"`
	Key          types.String `tfsdk:"key"`
	HaveChildren types.Bool   `tfsdk:"have_children"`
	IsTopGroup   types.Bool   `tfsdk:"is_top_group"`
	Users        types.List   `tfsdk:"users"`
	IsEnabled    types.Bool   `tfsdk:"is_enabled"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor user group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the group in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the group was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_time": schema.StringAttribute{
				Description: "The time when the group was last updated.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"manager": schema.StringAttribute{
				Description: "The manager of the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"contact_email": schema.StringAttribute{
				Description: "The contact email for the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description: "The type of the group (e.g., 'Physical', 'Virtual').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"parent_id": schema.StringAttribute{
				Description: "The parent group ID for hierarchical groups.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"parent_name": schema.StringAttribute{
				Description: "The parent group name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"title": schema.StringAttribute{
				Description: "The title of the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"key": schema.StringAttribute{
				Description: "The key identifier of the group.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"have_children": schema.BoolAttribute{
				Description: "Whether this group has child groups.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_top_group": schema.BoolAttribute{
				Description: "Whether this is a top-level group.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"users": schema.ListAttribute{
				Description: "List of users in this group.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the group is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users := make([]string, 0)
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	group := &casdoorsdk.Group{
		Owner:        plan.Owner.ValueString(),
		Name:         plan.Name.ValueString(),
		CreatedTime:  createdTime,
		UpdatedTime:  createdTime,
		DisplayName:  plan.DisplayName.ValueString(),
		Manager:      plan.Manager.ValueString(),
		ContactEmail: plan.ContactEmail.ValueString(),
		Type:         plan.Type.ValueString(),
		ParentId:     plan.ParentId.ValueString(),
		ParentName:   plan.ParentName.ValueString(),
		Title:        plan.Title.ValueString(),
		Key:          plan.Key.ValueString(),
		IsTopGroup:   plan.IsTopGroup.ValueBool(),
		Users:        users,
		IsEnabled:    plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.AddGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Group",
			fmt.Sprintf("Could not create group %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Group",
			fmt.Sprintf("Casdoor returned failure when creating group %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the group to get server-generated values like CreatedTime.
	createdGroup, err := r.client.GetGroup(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Group",
			fmt.Sprintf("Could not read group %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdGroup == nil {
		resp.Diagnostics.AddError(
			"Error Reading Group",
			fmt.Sprintf("Group %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	plan.CreatedTime = types.StringValue(createdGroup.CreatedTime)
	plan.UpdatedTime = types.StringValue(createdGroup.UpdatedTime)
	plan.ParentName = types.StringValue(createdGroup.ParentName)
	plan.Title = types.StringValue(createdGroup.Title)
	plan.Key = types.StringValue(createdGroup.Key)
	plan.HaveChildren = types.BoolValue(createdGroup.HaveChildren)
	usersList, diags := types.ListValueFrom(ctx, types.StringType, createdGroup.Users)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Users = usersList

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := getByOwnerName[casdoorsdk.Group](r.client, "get-group", state.Owner.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Group",
			fmt.Sprintf("Could not read group %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if group == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(group.Owner + "/" + group.Name)
	state.Owner = types.StringValue(group.Owner)
	state.Name = types.StringValue(group.Name)
	state.CreatedTime = types.StringValue(group.CreatedTime)
	state.UpdatedTime = types.StringValue(group.UpdatedTime)
	state.DisplayName = types.StringValue(group.DisplayName)
	state.Manager = types.StringValue(group.Manager)
	state.ContactEmail = types.StringValue(group.ContactEmail)
	state.Type = types.StringValue(group.Type)
	state.ParentId = types.StringValue(group.ParentId)
	state.ParentName = types.StringValue(group.ParentName)
	state.Title = types.StringValue(group.Title)
	state.Key = types.StringValue(group.Key)
	state.HaveChildren = types.BoolValue(group.HaveChildren)
	state.IsTopGroup = types.BoolValue(group.IsTopGroup)
	state.IsEnabled = types.BoolValue(group.IsEnabled)

	usersList, diags := types.ListValueFrom(ctx, types.StringType, group.Users)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Users = usersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users := make([]string, 0)
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	group := &casdoorsdk.Group{
		Owner:        plan.Owner.ValueString(),
		Name:         plan.Name.ValueString(),
		CreatedTime:  plan.CreatedTime.ValueString(),
		UpdatedTime:  plan.UpdatedTime.ValueString(),
		DisplayName:  plan.DisplayName.ValueString(),
		Manager:      plan.Manager.ValueString(),
		ContactEmail: plan.ContactEmail.ValueString(),
		Type:         plan.Type.ValueString(),
		ParentId:     plan.ParentId.ValueString(),
		IsTopGroup:   plan.IsTopGroup.ValueBool(),
		Users:        users,
		IsEnabled:    plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.UpdateGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Group",
			fmt.Sprintf("Could not update group %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Group",
			fmt.Sprintf("Casdoor returned failure when updating group %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back to get updated fields.
	updatedGroup, err := r.client.GetGroup(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Group",
			fmt.Sprintf("Could not read group %q after update: %s", plan.Name.ValueString(), err),
		)
		return
	}
	if updatedGroup == nil {
		resp.Diagnostics.AddError(
			"Error Reading Group",
			fmt.Sprintf("Group %q not found after update", plan.Name.ValueString()),
		)
		return
	}

	plan.UpdatedTime = types.StringValue(updatedGroup.UpdatedTime)
	plan.ParentName = types.StringValue(updatedGroup.ParentName)
	plan.Title = types.StringValue(updatedGroup.Title)
	plan.Key = types.StringValue(updatedGroup.Key)
	plan.HaveChildren = types.BoolValue(updatedGroup.HaveChildren)
	usersList, diags := types.ListValueFrom(ctx, types.StringType, updatedGroup.Users)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Users = usersList

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &casdoorsdk.Group{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Group",
			fmt.Sprintf("Could not delete group %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Group",
			fmt.Sprintf("Casdoor returned failure when deleting group %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
