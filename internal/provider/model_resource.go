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
	_ resource.Resource                = &ModelResource{}
	_ resource.ResourceWithConfigure   = &ModelResource{}
	_ resource.ResourceWithImportState = &ModelResource{}
)

type ModelResource struct {
	client *casdoorsdk.Client
}

type ModelResourceModel struct {
	Owner       types.String `tfsdk:"owner"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	ModelText   types.String `tfsdk:"model_text"`
}

func NewModelResource() resource.Resource {
	return &ModelResource{}
}

func (r *ModelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

func (r *ModelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor Casbin model.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this model.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the model.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"model_text": schema.StringAttribute{
				Description: "The Casbin model definition text (PERM format).",
				Required:    true,
			},
		},
	}
}

func (r *ModelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model := &casdoorsdk.Model{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		ModelText:   plan.ModelText.ValueString(),
	}

	success, err := r.client.AddModel(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Model",
			fmt.Sprintf("Could not create model %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Model",
			fmt.Sprintf("Casdoor returned failure when creating model %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetModel(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Model",
			fmt.Sprintf("Could not read model %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if model == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(model.Owner)
	state.Name = types.StringValue(model.Name)
	state.DisplayName = types.StringValue(model.DisplayName)
	state.ModelText = types.StringValue(model.ModelText)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model := &casdoorsdk.Model{
		Owner:       plan.Owner.ValueString(),
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		ModelText:   plan.ModelText.ValueString(),
	}

	success, err := r.client.UpdateModel(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Model",
			fmt.Sprintf("Could not update model %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Model",
			fmt.Sprintf("Casdoor returned failure when updating model %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model := &casdoorsdk.Model{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteModel(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Model",
			fmt.Sprintf("Could not delete model %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Model",
			fmt.Sprintf("Casdoor returned failure when deleting model %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *ModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
