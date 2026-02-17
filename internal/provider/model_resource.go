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
	ID           types.String `tfsdk:"id"`
	Owner        types.String `tfsdk:"owner"`
	Name         types.String `tfsdk:"name"`
	CreatedTime  types.String `tfsdk:"created_time"`
	UpdatedTime  types.String `tfsdk:"updated_time"`
	Description  types.String `tfsdk:"description"`
	DisplayName  types.String `tfsdk:"display_name"`
	ModelText    types.String `tfsdk:"model_text"`
	Manager      types.String `tfsdk:"manager"`
	ContactEmail types.String `tfsdk:"contact_email"`
	Type         types.String `tfsdk:"type"`
	ParentId     types.String `tfsdk:"parent_id"`
	IsTopModel   types.Bool   `tfsdk:"is_top_model"`
	IsEnabled    types.Bool   `tfsdk:"is_enabled"`
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
			"id": schema.StringAttribute{
				Description: "The ID of the model in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
			"created_time": schema.StringAttribute{
				Description: "The time when the model was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "A description of the model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"model_text": schema.StringAttribute{
				Description: "The Casbin model definition text (PERM format).",
				Required:    true,
			},
			"updated_time": schema.StringAttribute{
				Description: "The time when the model was last updated.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"manager": schema.StringAttribute{
				Description: "The manager of this model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"contact_email": schema.StringAttribute{
				Description: "The contact email for this model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description: "The type of the model.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"parent_id": schema.StringAttribute{
				Description: "The parent model ID.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"is_top_model": schema.BoolAttribute{
				Description: "Whether this is a top-level model.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether this model is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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

func modelPlanToSDK(plan ModelResourceModel, createdTime string) *casdoorsdk.Model {
	return &casdoorsdk.Model{
		Owner:        plan.Owner.ValueString(),
		Name:         plan.Name.ValueString(),
		CreatedTime:  createdTime,
		Description:  plan.Description.ValueString(),
		DisplayName:  plan.DisplayName.ValueString(),
		ModelText:    plan.ModelText.ValueString(),
		Manager:      plan.Manager.ValueString(),
		ContactEmail: plan.ContactEmail.ValueString(),
		Type:         plan.Type.ValueString(),
		ParentId:     plan.ParentId.ValueString(),
		IsTopModel:   plan.IsTopModel.ValueBool(),
		IsEnabled:    plan.IsEnabled.ValueBool(),
	}
}

func (r *ModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	model := modelPlanToSDK(plan, createdTime)

	ok, err := r.client.AddModel(model)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating model %q", plan.Name.ValueString())) {
		return
	}

	// Read back the model to get server-generated values like CreatedTime.
	createdModel, err := r.client.GetModel(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Model",
			fmt.Sprintf("Could not read model %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdModel == nil {
		resp.Diagnostics.AddError(
			"Error Reading Model",
			fmt.Sprintf("Model %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	plan.CreatedTime = types.StringValue(createdModel.CreatedTime)
	plan.UpdatedTime = types.StringValue(createdModel.UpdatedTime)
	plan.ModelText = types.StringValue(createdModel.ModelText)

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetModel(state.Owner.ValueString() + "/" + state.Name.ValueString())
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

	state.ID = types.StringValue(model.Owner + "/" + model.Name)
	state.Owner = types.StringValue(model.Owner)
	state.Name = types.StringValue(model.Name)
	state.CreatedTime = types.StringValue(model.CreatedTime)
	state.UpdatedTime = types.StringValue(model.UpdatedTime)
	state.Description = types.StringValue(model.Description)
	state.DisplayName = types.StringValue(model.DisplayName)
	state.ModelText = types.StringValue(model.ModelText)
	state.Manager = types.StringValue(model.Manager)
	state.ContactEmail = types.StringValue(model.ContactEmail)
	state.Type = types.StringValue(model.Type)
	state.ParentId = types.StringValue(model.ParentId)
	state.IsTopModel = types.BoolValue(model.IsTopModel)
	state.IsEnabled = types.BoolValue(model.IsEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model := modelPlanToSDK(plan, plan.CreatedTime.ValueString())

	_, err := r.client.UpdateModel(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Model",
			fmt.Sprintf("Could not update model %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	// Read back to get server-normalized values.
	updatedModel, err := r.client.GetModel(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Model",
			fmt.Sprintf("Could not read model %q after update: %s", plan.Name.ValueString(), err),
		)
		return
	}
	if updatedModel == nil {
		resp.Diagnostics.AddError(
			"Error Reading Model",
			fmt.Sprintf("Model %q not found after update", plan.Name.ValueString()),
		)
		return
	}

	plan.ModelText = types.StringValue(updatedModel.ModelText)

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
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

	ok, err := r.client.DeleteModel(model)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("deleting model %q", state.Name.ValueString())) {
		return
	}
}

func (r *ModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
