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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
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
	ID          types.String `tfsdk:"id"`
	Owner       types.String `tfsdk:"owner"`
	Name        types.String `tfsdk:"name"`
	CreatedTime types.String `tfsdk:"created_time"`
	UpdatedTime types.String `tfsdk:"updated_time"`
	ModelCfg    types.Map    `tfsdk:"model_cfg"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Model       types.String `tfsdk:"model"`
	Adapter     types.String `tfsdk:"adapter"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
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
			"id": schema.StringAttribute{
				Description: "The ID of the enforcer in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
			"created_time": schema.StringAttribute{
				Description: "The time when the enforcer was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_time": schema.StringAttribute{
				Description: "The time when the enforcer was last updated.",
				Computed:    true,
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
			"model_cfg": schema.MapAttribute{
				Description: "The model configuration key-value pairs.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether this enforcer is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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

	var modelCfg map[string]string
	if !plan.ModelCfg.IsNull() {
		modelCfg = make(map[string]string)
		resp.Diagnostics.Append(plan.ModelCfg.ElementsAs(ctx, &modelCfg, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	enforcer := &casdoorsdk.Enforcer{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: createdTime,
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Model:       plan.Model.ValueString(),
		Adapter:     plan.Adapter.ValueString(),
		ModelCfg:    modelCfg,
		IsEnabled:   plan.IsEnabled.ValueBool(),
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

	// Read back the enforcer to get server-generated values.
	createdEnforcer, err := r.client.GetEnforcer(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Enforcer",
			fmt.Sprintf("Could not read enforcer %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdEnforcer != nil {
		plan.CreatedTime = types.StringValue(createdEnforcer.CreatedTime)
		plan.UpdatedTime = types.StringValue(createdEnforcer.UpdatedTime)
		if len(createdEnforcer.ModelCfg) > 0 {
			modelCfgMap, diags := types.MapValueFrom(ctx, types.StringType, createdEnforcer.ModelCfg)
			resp.Diagnostics.Append(diags...)
			plan.ModelCfg = modelCfgMap
		} else {
			plan.ModelCfg = types.MapValueMust(types.StringType, map[string]attr.Value{})
		}
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
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

	state.ID = types.StringValue(enforcer.Owner + "/" + enforcer.Name)
	state.Owner = types.StringValue(enforcer.Owner)
	state.Name = types.StringValue(enforcer.Name)
	state.DisplayName = types.StringValue(enforcer.DisplayName)
	state.Description = types.StringValue(enforcer.Description)
	state.Model = types.StringValue(enforcer.Model)
	state.Adapter = types.StringValue(enforcer.Adapter)
	state.IsEnabled = types.BoolValue(enforcer.IsEnabled)

	state.CreatedTime = types.StringValue(enforcer.CreatedTime)
	state.UpdatedTime = types.StringValue(enforcer.UpdatedTime)

	if len(enforcer.ModelCfg) > 0 {
		modelCfgMap, diags := types.MapValueFrom(ctx, types.StringType, enforcer.ModelCfg)
		resp.Diagnostics.Append(diags...)
		state.ModelCfg = modelCfgMap
	} else {
		state.ModelCfg = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnforcerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnforcerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var modelCfg map[string]string
	if !plan.ModelCfg.IsNull() {
		modelCfg = make(map[string]string)
		resp.Diagnostics.Append(plan.ModelCfg.ElementsAs(ctx, &modelCfg, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	enforcer := &casdoorsdk.Enforcer{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: plan.CreatedTime.ValueString(),
		UpdatedTime: plan.UpdatedTime.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		Model:       plan.Model.ValueString(),
		Adapter:     plan.Adapter.ValueString(),
		ModelCfg:    modelCfg,
		IsEnabled:   plan.IsEnabled.ValueBool(),
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

	// Read back to get server-updated fields.
	updatedEnforcer, err := r.client.GetEnforcer(plan.Name.ValueString())
	if err == nil && updatedEnforcer != nil {
		plan.UpdatedTime = types.StringValue(updatedEnforcer.UpdatedTime)
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
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
	importStateOwnerName(ctx, req, resp)
}
