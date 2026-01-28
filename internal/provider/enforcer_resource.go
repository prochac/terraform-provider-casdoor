// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &EnforcerResource{}
	_ resource.ResourceWithConfigure   = &EnforcerResource{}
	_ resource.ResourceWithImportState = &EnforcerResource{}
)

type EnforcerResource struct {
	client *casdoorsdk.Client
}

type EnforcerResourceModel struct {
	Owner       types.String `tfsdk:"owner"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Model       types.String `tfsdk:"model"`
	Adapter     types.String `tfsdk:"adapter"`
}

func NewEnforcerResource() resource.Resource {
	return &EnforcerResource{}
}

func (r *EnforcerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcer"
}

func (r *EnforcerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor Casbin enforcer.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this enforcer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the enforcer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the enforcer.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "A description of the enforcer.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"model": schema.StringAttribute{
				Description: "The Casbin model name to use (format: 'organization/model-name').",
				Required:    true,
			},
			"adapter": schema.StringAttribute{
				Description: "The Casbin adapter name to use (format: 'organization/adapter-name').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *EnforcerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnforcerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnforcerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enforcer := &casdoorsdk.Enforcer{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Model:       plan.Model.ValueString(),
		Adapter:     plan.Adapter.ValueString(),
	}

	success, err := r.client.AddEnforcer(enforcer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Enforcer",
			fmt.Sprintf("Could not create enforcer %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Enforcer",
			fmt.Sprintf("Casdoor returned failure when creating enforcer %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *EnforcerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnforcerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enforcer, err := r.client.GetEnforcer(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Enforcer",
			fmt.Sprintf("Could not read enforcer %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if enforcer == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(enforcer.Owner)
	state.Name = types.StringValue(enforcer.Name)
	state.DisplayName = types.StringValue(enforcer.DisplayName)
	state.Description = types.StringValue(enforcer.Description)
	state.Model = types.StringValue(enforcer.Model)
	state.Adapter = types.StringValue(enforcer.Adapter)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnforcerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnforcerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enforcer := &casdoorsdk.Enforcer{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Model:       plan.Model.ValueString(),
		Adapter:     plan.Adapter.ValueString(),
	}

	success, err := r.client.UpdateEnforcer(enforcer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Enforcer",
			fmt.Sprintf("Could not update enforcer %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Enforcer",
			fmt.Sprintf("Casdoor returned failure when updating enforcer %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *EnforcerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnforcerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enforcer := &casdoorsdk.Enforcer{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteEnforcer(enforcer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Enforcer",
			fmt.Sprintf("Could not delete enforcer %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Enforcer",
			fmt.Sprintf("Casdoor returned failure when deleting enforcer %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *EnforcerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
