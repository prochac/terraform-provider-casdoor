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
	_ resource.Resource                = &PricingResource{}
	_ resource.ResourceWithConfigure   = &PricingResource{}
	_ resource.ResourceWithImportState = &PricingResource{}
)

type PricingResource struct {
	client *casdoorsdk.Client
}

type PricingResourceModel struct {
	Owner         types.String `tfsdk:"owner"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
	Plans         types.List   `tfsdk:"plans"`
	IsEnabled     types.Bool   `tfsdk:"is_enabled"`
	TrialDuration types.Int64  `tfsdk:"trial_duration"`
	Application   types.String `tfsdk:"application"`
	// Submitter     types.String `tfsdk:"submitter"`
	// Approver      types.String `tfsdk:"approver"`
	// ApproveTime   types.String `tfsdk:"approve_time"`
	// State         types.String `tfsdk:"state"`
}

func NewPricingResource() resource.Resource {
	return &PricingResource{}
}

func (r *PricingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pricing"
}

func (r *PricingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor pricing configuration for SaaS products.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this pricing.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the pricing.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the pricing was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the pricing.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "The description of the pricing.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"plans": schema.ListAttribute{
				Description: "List of plan names included in this pricing.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the pricing is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"trial_duration": schema.Int64Attribute{
				Description: "The trial duration in days.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"application": schema.StringAttribute{
				Description: "The application this pricing is for.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *PricingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PricingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PricingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plans []string
	if !plan.Plans.IsNull() && !plan.Plans.IsUnknown() {
		resp.Diagnostics.Append(plan.Plans.ElementsAs(ctx, &plans, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	pricing := &casdoorsdk.Pricing{
		Owner:         plan.Owner.ValueString(),
		Name:          plan.Name.ValueString(),
		DisplayName:   plan.DisplayName.ValueString(),
		Description:   plan.Description.ValueString(),
		Plans:         plans,
		IsEnabled:     plan.IsEnabled.ValueBool(),
		TrialDuration: int(plan.TrialDuration.ValueInt64()),
		Application:   plan.Application.ValueString(),
	}

	success, err := r.client.AddPricing(pricing)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Pricing",
			fmt.Sprintf("Could not create pricing %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Pricing",
			fmt.Sprintf("Casdoor returned failure when creating pricing %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the pricing to get server-generated values.
	createdPricing, err := r.client.GetPricing(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pricing",
			fmt.Sprintf("Could not read pricing %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdPricing != nil {
		plan.CreatedTime = types.StringValue(createdPricing.CreatedTime)
		plansList, _ := types.ListValueFrom(ctx, types.StringType, createdPricing.Plans)
		plan.Plans = plansList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PricingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PricingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pricing, err := r.client.GetPricing(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pricing",
			fmt.Sprintf("Could not read pricing %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if pricing == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(pricing.Owner)
	state.Name = types.StringValue(pricing.Name)
	state.CreatedTime = types.StringValue(pricing.CreatedTime)
	state.DisplayName = types.StringValue(pricing.DisplayName)
	state.Description = types.StringValue(pricing.Description)
	state.IsEnabled = types.BoolValue(pricing.IsEnabled)
	state.TrialDuration = types.Int64Value(int64(pricing.TrialDuration))
	state.Application = types.StringValue(pricing.Application)

	plansList, _ := types.ListValueFrom(ctx, types.StringType, pricing.Plans)
	state.Plans = plansList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PricingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PricingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plans []string
	if !plan.Plans.IsNull() && !plan.Plans.IsUnknown() {
		resp.Diagnostics.Append(plan.Plans.ElementsAs(ctx, &plans, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	pricing := &casdoorsdk.Pricing{
		Owner:         plan.Owner.ValueString(),
		Name:          plan.Name.ValueString(),
		CreatedTime:   plan.CreatedTime.ValueString(),
		DisplayName:   plan.DisplayName.ValueString(),
		Description:   plan.Description.ValueString(),
		Plans:         plans,
		IsEnabled:     plan.IsEnabled.ValueBool(),
		TrialDuration: int(plan.TrialDuration.ValueInt64()),
		Application:   plan.Application.ValueString(),
	}

	success, err := r.client.UpdatePricing(pricing)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Pricing",
			fmt.Sprintf("Could not update pricing %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Pricing",
			fmt.Sprintf("Casdoor returned failure when updating pricing %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back to get updated computed fields.
	updatedPricing, err := r.client.GetPricing(plan.Name.ValueString())
	if err == nil && updatedPricing != nil {
		// TODO update attributes
		// plan.SomeAttr = types.StringValue(updatedPricing.SomeAttr)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PricingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PricingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pricing := &casdoorsdk.Pricing{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeletePricing(pricing)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Pricing",
			fmt.Sprintf("Could not delete pricing %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Pricing",
			fmt.Sprintf("Casdoor returned failure when deleting pricing %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *PricingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
