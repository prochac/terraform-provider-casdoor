// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &WebhookResource{}
	_ resource.ResourceWithConfigure   = &WebhookResource{}
	_ resource.ResourceWithImportState = &WebhookResource{}
)

type WebhookResource struct {
	client *casdoorsdk.Client
}

type HeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type WebhookResourceModel struct {
	Owner          types.String `tfsdk:"owner"`
	Name           types.String `tfsdk:"name"`
	CreatedTime    types.String `tfsdk:"created_time"`
	Organization   types.String `tfsdk:"organization"`
	Url            types.String `tfsdk:"url"`
	Method         types.String `tfsdk:"method"`
	ContentType    types.String `tfsdk:"content_type"`
	Headers        types.List   `tfsdk:"headers"`
	Events         types.List   `tfsdk:"events"`
	TokenFields    types.List   `tfsdk:"token_fields"`
	ObjectFields   types.List   `tfsdk:"object_fields"`
	IsUserExtended types.Bool   `tfsdk:"is_user_extended"`
	SingleOrgOnly  types.Bool   `tfsdk:"single_org_only"`
	IsEnabled      types.Bool   `tfsdk:"is_enabled"`
}

func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

func (r *WebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

var headerAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

func (r *WebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor webhook for event notifications.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this webhook.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the webhook.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the webhook was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The organization this webhook belongs to.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"url": schema.StringAttribute{
				Description: "The URL to send webhook requests to.",
				Required:    true,
			},
			"method": schema.StringAttribute{
				Description: "The HTTP method to use (e.g., 'POST', 'GET').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("POST"),
			},
			"content_type": schema.StringAttribute{
				Description: "The content type of the webhook request.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("application/json"),
			},
			"headers": schema.ListNestedAttribute{
				Description: "Custom headers to include in webhook requests.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The header name.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The header value.",
							Required:    true,
						},
					},
				},
			},
			"events": schema.ListAttribute{
				Description: "List of events that trigger this webhook.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"token_fields": schema.ListAttribute{
				Description: "Token fields to include in the webhook payload.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"object_fields": schema.ListAttribute{
				Description: "Object fields to include in the webhook payload.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_user_extended": schema.BoolAttribute{
				Description: "Whether to include extended user information.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"single_org_only": schema.BoolAttribute{
				Description: "Whether the webhook is limited to a single organization.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the webhook is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *WebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebhookResource) headersToSDK(ctx context.Context, plan WebhookResourceModel) ([]*casdoorsdk.Header, error) {
	if plan.Headers.IsNull() || plan.Headers.IsUnknown() {
		return nil, nil
	}

	var headers []HeaderModel
	diags := plan.Headers.ElementsAs(ctx, &headers, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert headers")
	}

	result := make([]*casdoorsdk.Header, len(headers))
	for i, h := range headers {
		result[i] = &casdoorsdk.Header{
			Name:  h.Name.ValueString(),
			Value: h.Value.ValueString(),
		}
	}
	return result, nil
}

func (r *WebhookResource) headersFromSDK(_ context.Context, headers []*casdoorsdk.Header) (types.List, error) {
	if len(headers) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: headerAttrTypes}), nil
	}

	objs := make([]attr.Value, len(headers))
	for i, h := range headers {
		obj, diags := types.ObjectValue(headerAttrTypes, map[string]attr.Value{
			"name":  types.StringValue(h.Name),
			"value": types.StringValue(h.Value),
		})
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: headerAttrTypes}), fmt.Errorf("failed to create object")
		}
		objs[i] = obj
	}

	result, diags := types.ListValue(types.ObjectType{AttrTypes: headerAttrTypes}, objs)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: headerAttrTypes}), fmt.Errorf("failed to create list")
	}
	return result, nil
}

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebhookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	headers, err := r.headersToSDK(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error Converting Headers", err.Error())
		return
	}

	var events, tokenFields, objectFields []string
	if !plan.Events.IsNull() && !plan.Events.IsUnknown() {
		resp.Diagnostics.Append(plan.Events.ElementsAs(ctx, &events, false)...)
	}
	if !plan.TokenFields.IsNull() && !plan.TokenFields.IsUnknown() {
		resp.Diagnostics.Append(plan.TokenFields.ElementsAs(ctx, &tokenFields, false)...)
	}
	if !plan.ObjectFields.IsNull() && !plan.ObjectFields.IsUnknown() {
		resp.Diagnostics.Append(plan.ObjectFields.ElementsAs(ctx, &objectFields, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := &casdoorsdk.Webhook{
		Owner:          plan.Owner.ValueString(),
		Name:           plan.Name.ValueString(),
		Organization:   plan.Organization.ValueString(),
		Url:            plan.Url.ValueString(),
		Method:         plan.Method.ValueString(),
		ContentType:    plan.ContentType.ValueString(),
		Headers:        headers,
		Events:         events,
		TokenFields:    tokenFields,
		ObjectFields:   objectFields,
		IsUserExtended: plan.IsUserExtended.ValueBool(),
		SingleOrgOnly:  plan.SingleOrgOnly.ValueBool(),
		IsEnabled:      plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.AddWebhook(webhook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Webhook",
			fmt.Sprintf("Could not create webhook %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Webhook",
			fmt.Sprintf("Casdoor returned failure when creating webhook %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the webhook to get server-generated values.
	createdWebhook, err := r.client.GetWebhook(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webhook",
			fmt.Sprintf("Could not read webhook %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdWebhook != nil {
		plan.CreatedTime = types.StringValue(createdWebhook.CreatedTime)
		headersList, err := r.headersFromSDK(ctx, createdWebhook.Headers)
		if err != nil {
			resp.Diagnostics.AddError("Error Converting Headers", err.Error())
			return
		}
		plan.Headers = headersList

		eventsList, _ := types.ListValueFrom(ctx, types.StringType, createdWebhook.Events)
		plan.Events = eventsList
		tokenFieldsList, _ := types.ListValueFrom(ctx, types.StringType, createdWebhook.TokenFields)
		plan.TokenFields = tokenFieldsList
		objectFieldsList, _ := types.ListValueFrom(ctx, types.StringType, createdWebhook.ObjectFields)
		plan.ObjectFields = objectFieldsList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, err := r.client.GetWebhook(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webhook",
			fmt.Sprintf("Could not read webhook %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if webhook == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(webhook.Owner)
	state.Name = types.StringValue(webhook.Name)
	state.CreatedTime = types.StringValue(webhook.CreatedTime)
	state.Organization = types.StringValue(webhook.Organization)
	state.Url = types.StringValue(webhook.Url)
	state.Method = types.StringValue(webhook.Method)
	state.ContentType = types.StringValue(webhook.ContentType)
	state.IsUserExtended = types.BoolValue(webhook.IsUserExtended)
	state.SingleOrgOnly = types.BoolValue(webhook.SingleOrgOnly)
	state.IsEnabled = types.BoolValue(webhook.IsEnabled)

	headersList, err := r.headersFromSDK(ctx, webhook.Headers)
	if err != nil {
		resp.Diagnostics.AddError("Error Converting Headers", err.Error())
		return
	}
	state.Headers = headersList

	eventsList, _ := types.ListValueFrom(ctx, types.StringType, webhook.Events)
	state.Events = eventsList
	tokenFieldsList, _ := types.ListValueFrom(ctx, types.StringType, webhook.TokenFields)
	state.TokenFields = tokenFieldsList
	objectFieldsList, _ := types.ListValueFrom(ctx, types.StringType, webhook.ObjectFields)
	state.ObjectFields = objectFieldsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebhookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	headers, err := r.headersToSDK(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error Converting Headers", err.Error())
		return
	}

	var events, tokenFields, objectFields []string
	if !plan.Events.IsNull() && !plan.Events.IsUnknown() {
		resp.Diagnostics.Append(plan.Events.ElementsAs(ctx, &events, false)...)
	}
	if !plan.TokenFields.IsNull() && !plan.TokenFields.IsUnknown() {
		resp.Diagnostics.Append(plan.TokenFields.ElementsAs(ctx, &tokenFields, false)...)
	}
	if !plan.ObjectFields.IsNull() && !plan.ObjectFields.IsUnknown() {
		resp.Diagnostics.Append(plan.ObjectFields.ElementsAs(ctx, &objectFields, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := &casdoorsdk.Webhook{
		Owner:          plan.Owner.ValueString(),
		Name:           plan.Name.ValueString(),
		CreatedTime:    plan.CreatedTime.ValueString(),
		Organization:   plan.Organization.ValueString(),
		Url:            plan.Url.ValueString(),
		Method:         plan.Method.ValueString(),
		ContentType:    plan.ContentType.ValueString(),
		Headers:        headers,
		Events:         events,
		TokenFields:    tokenFields,
		ObjectFields:   objectFields,
		IsUserExtended: plan.IsUserExtended.ValueBool(),
		SingleOrgOnly:  plan.SingleOrgOnly.ValueBool(),
		IsEnabled:      plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.UpdateWebhook(webhook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Webhook",
			fmt.Sprintf("Could not update webhook %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Webhook",
			fmt.Sprintf("Casdoor returned failure when updating webhook %q", plan.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := &casdoorsdk.Webhook{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteWebhook(webhook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Webhook",
			fmt.Sprintf("Could not delete webhook %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Webhook",
			fmt.Sprintf("Casdoor returned failure when deleting webhook %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
