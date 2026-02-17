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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &PlanResource{}
	_ resource.ResourceWithConfigure   = &PlanResource{}
	_ resource.ResourceWithImportState = &PlanResource{}
)

type PlanResource struct {
	client *casdoorsdk.Client
}

type PlanResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Owner            types.String  `tfsdk:"owner"`
	Name             types.String  `tfsdk:"name"`
	CreatedTime      types.String  `tfsdk:"created_time"`
	DisplayName      types.String  `tfsdk:"display_name"`
	Description      types.String  `tfsdk:"description"`
	Price            types.Float64 `tfsdk:"price"`
	Currency         types.String  `tfsdk:"currency"`
	Period           types.String  `tfsdk:"period"`
	Product          types.String  `tfsdk:"product"`
	PaymentProviders types.List    `tfsdk:"payment_providers"`
	IsEnabled        types.Bool    `tfsdk:"is_enabled"`
	Role             types.String  `tfsdk:"role"`
	Options          types.List    `tfsdk:"options"`
}

func NewPlanResource() resource.Resource {
	return &PlanResource{}
}

func (r *PlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plan"
}

func (r *PlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor subscription plan for SaaS products.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the plan in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this plan.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the plan.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the plan was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the plan.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "The description of the plan.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"price": schema.Float64Attribute{
				Description: "The price of the plan.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"currency": schema.StringAttribute{
				Description: "The currency for the price (e.g., 'USD', 'EUR').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"period": schema.StringAttribute{
				Description: "The billing period (e.g., 'Monthly', 'Yearly').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"product": schema.StringAttribute{
				Description: "The product this plan belongs to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"payment_providers": schema.ListAttribute{
				Description: "List of payment providers for this plan.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the plan is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"role": schema.StringAttribute{
				Description: "The role granted by this plan.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"options": schema.ListAttribute{
				Description: "Additional options for the plan.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PlanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func planPlanToSDK(ctx context.Context, plan PlanResourceModel, createdTime string) (*casdoorsdk.Plan, diag.Diagnostics) {
	var diags diag.Diagnostics

	paymentProviders, d := stringListToSDK(ctx, plan.PaymentProviders)
	diags.Append(d...)
	options, d := stringListToSDK(ctx, plan.Options)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &casdoorsdk.Plan{
		Owner:            plan.Owner.ValueString(),
		Name:             plan.Name.ValueString(),
		CreatedTime:      createdTime,
		DisplayName:      plan.DisplayName.ValueString(),
		Description:      plan.Description.ValueString(),
		Price:            plan.Price.ValueFloat64(),
		Currency:         plan.Currency.ValueString(),
		Period:           plan.Period.ValueString(),
		Product:          plan.Product.ValueString(),
		PaymentProviders: paymentProviders,
		IsEnabled:        plan.IsEnabled.ValueBool(),
		Role:             plan.Role.ValueString(),
		Options:          options,
	}, diags
}

func (r *PlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PlanResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	planObj, diags := planPlanToSDK(ctx, plan, createdTime)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.AddPlan(planObj)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating plan %q", plan.Name.ValueString())) {
		return
	}

	// Read back the plan to get server-generated values.
	createdPlan, err := r.client.GetPlan(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Plan",
			fmt.Sprintf("Could not read plan %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdPlan == nil {
		resp.Diagnostics.AddError(
			"Error Reading Plan",
			fmt.Sprintf("Plan %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	plan.CreatedTime = types.StringValue(createdPlan.CreatedTime)
	plan.Product = types.StringValue(createdPlan.Product)
	plan.Role = types.StringValue(createdPlan.Role)
	providersList, _ := types.ListValueFrom(ctx, types.StringType, createdPlan.PaymentProviders)
	plan.PaymentProviders = providersList
	optionsList, _ := types.ListValueFrom(ctx, types.StringType, createdPlan.Options)
	plan.Options = optionsList

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PlanResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planObj, err := getByOwnerName[casdoorsdk.Plan](r.client, "get-plan", state.Owner.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Plan",
			fmt.Sprintf("Could not read plan %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if planObj == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(planObj.Owner + "/" + planObj.Name)
	state.Owner = types.StringValue(planObj.Owner)
	state.Name = types.StringValue(planObj.Name)
	state.CreatedTime = types.StringValue(planObj.CreatedTime)
	state.DisplayName = types.StringValue(planObj.DisplayName)
	state.Description = types.StringValue(planObj.Description)
	state.Price = types.Float64Value(planObj.Price)
	state.Currency = types.StringValue(planObj.Currency)
	state.Period = types.StringValue(planObj.Period)
	state.Product = types.StringValue(planObj.Product)
	state.IsEnabled = types.BoolValue(planObj.IsEnabled)
	state.Role = types.StringValue(planObj.Role)

	providersList, _ := types.ListValueFrom(ctx, types.StringType, planObj.PaymentProviders)
	state.PaymentProviders = providersList
	optionsList, _ := types.ListValueFrom(ctx, types.StringType, planObj.Options)
	state.Options = optionsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PlanResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planObj, diags := planPlanToSDK(ctx, plan, plan.CreatedTime.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.UpdatePlan(planObj)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("updating plan %q", plan.Name.ValueString())) {
		return
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *PlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PlanResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planObj := &casdoorsdk.Plan{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	ok, err := r.client.DeletePlan(planObj)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("deleting plan %q", state.Name.ValueString())) {
		return
	}
}

func (r *PlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
