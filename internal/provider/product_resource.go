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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ProductResource{}
	_ resource.ResourceWithConfigure   = &ProductResource{}
	_ resource.ResourceWithImportState = &ProductResource{}
)

type ProductResource struct {
	client *casdoorsdk.Client
}

type ProductResourceModel struct {
	ID                    types.String  `tfsdk:"id"`
	Owner                 types.String  `tfsdk:"owner"`
	Name                  types.String  `tfsdk:"name"`
	CreatedTime           types.String  `tfsdk:"created_time"`
	DisplayName           types.String  `tfsdk:"display_name"`
	Image                 types.String  `tfsdk:"image"`
	Detail                types.String  `tfsdk:"detail"`
	Description           types.String  `tfsdk:"description"`
	Tag                   types.String  `tfsdk:"tag"`
	Currency              types.String  `tfsdk:"currency"`
	Price                 types.Float64 `tfsdk:"price"`
	Quantity              types.Int64   `tfsdk:"quantity"`
	Sold                  types.Int64   `tfsdk:"sold"`
	IsRecharge            types.Bool    `tfsdk:"is_recharge"`
	RechargeOptions       types.List    `tfsdk:"recharge_options"`
	DisableCustomRecharge types.Bool    `tfsdk:"disable_custom_recharge"`
	SuccessUrl            types.String  `tfsdk:"success_url"`
	Providers             types.List    `tfsdk:"providers"`
	State                 types.String  `tfsdk:"state"`
	ManagedByPlan         types.Bool    `tfsdk:"managed_by_plan"`
}

func NewProductResource() resource.Resource {
	return &ProductResource{}
}

func (r *ProductResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

func (r *ProductResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor product for the SaaS product catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the product in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this product.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the product.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the product was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image": schema.StringAttribute{
				Description: "The image URL for the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"detail": schema.StringAttribute{
				Description: "Detailed information about the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A short description of the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tag": schema.StringAttribute{
				Description: "A tag for categorizing the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"currency": schema.StringAttribute{
				Description: "The currency for the price (e.g., 'USD', 'EUR').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"price": schema.Float64Attribute{
				Description: "The price of the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"quantity": schema.Int64Attribute{
				Description: "The available quantity of the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sold": schema.Int64Attribute{
				Description: "The number of products sold.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"is_recharge": schema.BoolAttribute{
				Description: "Whether this is a recharge product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"recharge_options": schema.ListAttribute{
				Description: "List of recharge amount options.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Float64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_custom_recharge": schema.BoolAttribute{
				Description: "Whether to disable custom recharge amounts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"success_url": schema.StringAttribute{
				Description: "The URL to redirect to after successful payment.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"providers": schema.ListAttribute{
				Description: "List of payment provider names for this product.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The current state of the product.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"managed_by_plan": schema.BoolAttribute{
				Description: "True when this product was auto-created by a casdoor_plan. Deletion is a no-op for plan-managed products.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ProductResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func productPlanToSDK(ctx context.Context, plan ProductResourceModel, createdTime string) (*casdoorsdk.Product, diag.Diagnostics) {
	var diags diag.Diagnostics

	providers, d := stringListToSDK(ctx, plan.Providers)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	rechargeOptions := make([]float64, 0)
	if !plan.RechargeOptions.IsNull() && !plan.RechargeOptions.IsUnknown() {
		diags.Append(plan.RechargeOptions.ElementsAs(ctx, &rechargeOptions, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	return &casdoorsdk.Product{
		Owner:                 plan.Owner.ValueString(),
		Name:                  plan.Name.ValueString(),
		CreatedTime:           createdTime,
		DisplayName:           plan.DisplayName.ValueString(),
		Image:                 plan.Image.ValueString(),
		Detail:                plan.Detail.ValueString(),
		Description:           plan.Description.ValueString(),
		Tag:                   plan.Tag.ValueString(),
		Currency:              plan.Currency.ValueString(),
		Price:                 plan.Price.ValueFloat64(),
		Quantity:              int(plan.Quantity.ValueInt64()),
		Sold:                  int(plan.Sold.ValueInt64()),
		IsRecharge:            plan.IsRecharge.ValueBool(),
		RechargeOptions:       rechargeOptions,
		DisableCustomRecharge: plan.DisableCustomRecharge.ValueBool(),
		SuccessUrl:            plan.SuccessUrl.ValueString(),
		Providers:             providers,
		State:                 plan.State.ValueString(),
	}, diags
}

func (r *ProductResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProductResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	productID := plan.Owner.ValueString() + "/" + plan.Name.ValueString()

	// Check if the product already exists (e.g. auto-created by a casdoor_plan).
	existing, err := r.client.GetProduct(productID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Could not check for existing product %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if existing != nil {
		// Adopt the existing plan-created product: overlay user-specified
		// values onto the existing product, preserving server-set fields
		// (like currency) that the user didn't configure.
		plan.ManagedByPlan = types.BoolValue(true)

		if v := plan.DisplayName.ValueString(); v != "" {
			existing.DisplayName = v
		}
		if v := plan.Image.ValueString(); v != "" {
			existing.Image = v
		}
		if v := plan.Detail.ValueString(); v != "" {
			existing.Detail = v
		}
		if v := plan.Description.ValueString(); v != "" {
			existing.Description = v
		}
		if v := plan.Tag.ValueString(); v != "" {
			existing.Tag = v
		}
		if v := plan.Currency.ValueString(); v != "" {
			existing.Currency = v
		}
		if v := plan.Price.ValueFloat64(); v != 0 {
			existing.Price = v
		}
		if v := int(plan.Quantity.ValueInt64()); v != 0 {
			existing.Quantity = v
		}
		if plan.IsRecharge.ValueBool() {
			existing.IsRecharge = true
		}
		if plan.DisableCustomRecharge.ValueBool() {
			existing.DisableCustomRecharge = true
		}
		if v := plan.SuccessUrl.ValueString(); v != "" {
			existing.SuccessUrl = v
		}
		if v := plan.State.ValueString(); v != "" {
			existing.State = v
		}
		if !plan.Providers.IsNull() && !plan.Providers.IsUnknown() {
			providers, d := stringListToSDK(ctx, plan.Providers)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(providers) > 0 {
				existing.Providers = providers
			}
		}
		if !plan.RechargeOptions.IsNull() && !plan.RechargeOptions.IsUnknown() {
			var rechargeOptions []float64
			resp.Diagnostics.Append(plan.RechargeOptions.ElementsAs(ctx, &rechargeOptions, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(rechargeOptions) > 0 {
				existing.RechargeOptions = rechargeOptions
			}
		}

		ok, err := r.client.UpdateProduct(existing)
		if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("adopting product %q", plan.Name.ValueString())) {
			return
		}
	} else {
		// Normal creation flow.
		plan.ManagedByPlan = types.BoolValue(false)

		createdTime := plan.CreatedTime.ValueString()
		if createdTime == "" {
			createdTime = time.Now().UTC().Format(time.RFC3339)
		}

		product, diags := productPlanToSDK(ctx, plan, createdTime)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ok, err := r.client.AddProduct(product)
		if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating product %q", plan.Name.ValueString())) {
			return
		}
	}

	// Read back the product to get server-generated values.
	createdProduct, err := r.client.GetProduct(productID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Could not read product %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdProduct == nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Product %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringValue(productID)
	plan.Owner = types.StringValue(createdProduct.Owner)
	plan.Name = types.StringValue(createdProduct.Name)
	plan.CreatedTime = types.StringValue(createdProduct.CreatedTime)
	plan.DisplayName = types.StringValue(createdProduct.DisplayName)
	plan.Image = types.StringValue(createdProduct.Image)
	plan.Detail = types.StringValue(createdProduct.Detail)
	plan.Description = types.StringValue(createdProduct.Description)
	plan.Tag = types.StringValue(createdProduct.Tag)
	plan.Currency = types.StringValue(createdProduct.Currency)
	plan.Price = types.Float64Value(createdProduct.Price)
	plan.Quantity = types.Int64Value(int64(createdProduct.Quantity))
	plan.Sold = types.Int64Value(int64(createdProduct.Sold))
	plan.IsRecharge = types.BoolValue(createdProduct.IsRecharge)
	plan.DisableCustomRecharge = types.BoolValue(createdProduct.DisableCustomRecharge)
	plan.SuccessUrl = types.StringValue(createdProduct.SuccessUrl)
	plan.State = types.StringValue(createdProduct.State)
	if len(createdProduct.Providers) > 0 || !plan.Providers.IsNull() {
		providersList, _ := types.ListValueFrom(ctx, types.StringType, createdProduct.Providers)
		plan.Providers = providersList
	}
	if len(createdProduct.RechargeOptions) > 0 || !plan.RechargeOptions.IsNull() {
		rechargeOptionsList, _ := types.ListValueFrom(ctx, types.Float64Type, createdProduct.RechargeOptions)
		plan.RechargeOptions = rechargeOptionsList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProductResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProductResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	product, err := r.client.GetProduct(state.Owner.ValueString() + "/" + state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Could not read product %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if product == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(product.Owner + "/" + product.Name)
	state.Owner = types.StringValue(product.Owner)
	state.Name = types.StringValue(product.Name)
	state.CreatedTime = types.StringValue(product.CreatedTime)
	state.DisplayName = types.StringValue(product.DisplayName)
	state.Image = types.StringValue(product.Image)
	state.Detail = types.StringValue(product.Detail)
	state.Description = types.StringValue(product.Description)
	state.Tag = types.StringValue(product.Tag)
	state.Currency = types.StringValue(product.Currency)
	state.Price = types.Float64Value(product.Price)
	state.Quantity = types.Int64Value(int64(product.Quantity))
	state.Sold = types.Int64Value(int64(product.Sold))
	state.IsRecharge = types.BoolValue(product.IsRecharge)
	state.DisableCustomRecharge = types.BoolValue(product.DisableCustomRecharge)
	state.SuccessUrl = types.StringValue(product.SuccessUrl)
	state.State = types.StringValue(product.State)

	providersList, _ := types.ListValueFrom(ctx, types.StringType, product.Providers)
	state.Providers = providersList
	rechargeOptionsList, _ := types.ListValueFrom(ctx, types.Float64Type, product.RechargeOptions)
	state.RechargeOptions = rechargeOptionsList

	// ManagedByPlan is Terraform-only state; default to false on import.
	if state.ManagedByPlan.IsNull() || state.ManagedByPlan.IsUnknown() {
		state.ManagedByPlan = types.BoolValue(false)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProductResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProductResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	product, diags := productPlanToSDK(ctx, plan, plan.CreatedTime.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.UpdateProduct(product)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("updating product %q", plan.Name.ValueString())) {
		return
	}

	// Read back to get server-updated fields.
	updatedProduct, err := r.client.GetProduct(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Could not read product %q after update: %s", plan.Name.ValueString(), err),
		)
		return
	}
	if updatedProduct == nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Product %q not found after update", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringValue(updatedProduct.Owner + "/" + updatedProduct.Name)
	plan.Owner = types.StringValue(updatedProduct.Owner)
	plan.Name = types.StringValue(updatedProduct.Name)
	plan.CreatedTime = types.StringValue(updatedProduct.CreatedTime)
	plan.DisplayName = types.StringValue(updatedProduct.DisplayName)
	plan.Image = types.StringValue(updatedProduct.Image)
	plan.Detail = types.StringValue(updatedProduct.Detail)
	plan.Description = types.StringValue(updatedProduct.Description)
	plan.Tag = types.StringValue(updatedProduct.Tag)
	plan.Currency = types.StringValue(updatedProduct.Currency)
	plan.Price = types.Float64Value(updatedProduct.Price)
	plan.Quantity = types.Int64Value(int64(updatedProduct.Quantity))
	plan.Sold = types.Int64Value(int64(updatedProduct.Sold))
	plan.IsRecharge = types.BoolValue(updatedProduct.IsRecharge)
	plan.DisableCustomRecharge = types.BoolValue(updatedProduct.DisableCustomRecharge)
	plan.SuccessUrl = types.StringValue(updatedProduct.SuccessUrl)
	plan.State = types.StringValue(updatedProduct.State)
	if len(updatedProduct.Providers) > 0 || !plan.Providers.IsNull() {
		providersList, _ := types.ListValueFrom(ctx, types.StringType, updatedProduct.Providers)
		plan.Providers = providersList
	}
	if len(updatedProduct.RechargeOptions) > 0 || !plan.RechargeOptions.IsNull() {
		rechargeOptionsList, _ := types.ListValueFrom(ctx, types.Float64Type, updatedProduct.RechargeOptions)
		plan.RechargeOptions = rechargeOptionsList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProductResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProductResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Plan-managed products are deleted by Casdoor when the plan is destroyed.
	if state.ManagedByPlan.ValueBool() {
		return
	}

	product := &casdoorsdk.Product{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	ok, err := r.client.DeleteProduct(product)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("deleting product %q", state.Name.ValueString())) {
		return
	}
}

func (r *ProductResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
