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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
				Default:     stringdefault.StaticString(""),
			},
			"image": schema.StringAttribute{
				Description: "The image URL for the product.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"detail": schema.StringAttribute{
				Description: "Detailed information about the product.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "A short description of the product.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tag": schema.StringAttribute{
				Description: "A tag for categorizing the product.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"currency": schema.StringAttribute{
				Description: "The currency for the price (e.g., 'USD', 'EUR').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"price": schema.Float64Attribute{
				Description: "The price of the product.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"quantity": schema.Int64Attribute{
				Description: "The available quantity of the product.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"sold": schema.Int64Attribute{
				Description: "The number of products sold.",
				Computed:    true,
			},
			"is_recharge": schema.BoolAttribute{
				Description: "Whether this is a recharge product.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"recharge_options": schema.ListAttribute{
				Description: "List of recharge amount options.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Float64Type,
			},
			"disable_custom_recharge": schema.BoolAttribute{
				Description: "Whether to disable custom recharge amounts.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"success_url": schema.StringAttribute{
				Description: "The URL to redirect to after successful payment.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"providers": schema.ListAttribute{
				Description: "List of payment provider names for this product.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"state": schema.StringAttribute{
				Description: "The current state of the product.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
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

func (r *ProductResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProductResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var providers []string
	if !plan.Providers.IsNull() && !plan.Providers.IsUnknown() {
		resp.Diagnostics.Append(plan.Providers.ElementsAs(ctx, &providers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	product := &casdoorsdk.Product{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: createdTime,
		DisplayName: plan.DisplayName.ValueString(),
		Image:       plan.Image.ValueString(),
		Detail:      plan.Detail.ValueString(),
		Description: plan.Description.ValueString(),
		Tag:         plan.Tag.ValueString(),
		Currency:    plan.Currency.ValueString(),
		Price:       plan.Price.ValueFloat64(),
		Quantity:    int(plan.Quantity.ValueInt64()),
		Providers:   providers,
		State:       plan.State.ValueString(),
	}

	success, err := r.client.AddProduct(product)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Product",
			fmt.Sprintf("Could not create product %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Product",
			fmt.Sprintf("Casdoor returned failure when creating product %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the product to get server-generated values.
	createdProduct, err := r.client.GetProduct(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Product",
			fmt.Sprintf("Could not read product %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdProduct != nil {
		plan.CreatedTime = types.StringValue(createdProduct.CreatedTime)
		plan.Sold = types.Int64Value(int64(createdProduct.Sold))
		providersList, _ := types.ListValueFrom(ctx, types.StringType, createdProduct.Providers)
		plan.Providers = providersList
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

	product, err := r.client.GetProduct(state.Name.ValueString())
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProductResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProductResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var providers []string
	if !plan.Providers.IsNull() && !plan.Providers.IsUnknown() {
		resp.Diagnostics.Append(plan.Providers.ElementsAs(ctx, &providers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	product := &casdoorsdk.Product{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		CreatedTime: plan.CreatedTime.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Image:       plan.Image.ValueString(),
		Detail:      plan.Detail.ValueString(),
		Description: plan.Description.ValueString(),
		Tag:         plan.Tag.ValueString(),
		Currency:    plan.Currency.ValueString(),
		Price:       plan.Price.ValueFloat64(),
		Quantity:    int(plan.Quantity.ValueInt64()),
		Providers:   providers,
		State:       plan.State.ValueString(),
	}

	success, err := r.client.UpdateProduct(product)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Product",
			fmt.Sprintf("Could not update product %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Product",
			fmt.Sprintf("Casdoor returned failure when updating product %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back to get server-updated fields.
	updatedProduct, err := r.client.GetProduct(plan.Name.ValueString())
	if err == nil && updatedProduct != nil {
		plan.Sold = types.Int64Value(int64(updatedProduct.Sold))
		plan.IsRecharge = types.BoolValue(updatedProduct.IsRecharge)
		plan.DisableCustomRecharge = types.BoolValue(updatedProduct.DisableCustomRecharge)
		plan.SuccessUrl = types.StringValue(updatedProduct.SuccessUrl)
		plan.State = types.StringValue(updatedProduct.State)
		providersList, _ := types.ListValueFrom(ctx, types.StringType, updatedProduct.Providers)
		plan.Providers = providersList
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

	product := &casdoorsdk.Product{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteProduct(product)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Product",
			fmt.Sprintf("Could not delete product %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Product",
			fmt.Sprintf("Casdoor returned failure when deleting product %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *ProductResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
