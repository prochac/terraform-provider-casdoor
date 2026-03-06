// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ResourceResource{}
	_ resource.ResourceWithConfigure   = &ResourceResource{}
	_ resource.ResourceWithImportState = &ResourceResource{}
)

type ResourceResource struct {
	client *casdoorsdk.Client
}

type ResourceResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Owner           types.String `tfsdk:"owner"`
	Name            types.String `tfsdk:"name"`
	CreatedTime     types.String `tfsdk:"created_time"`
	User            types.String `tfsdk:"user"`
	Tag             types.String `tfsdk:"tag"`
	Parent          types.String `tfsdk:"parent"`
	FileName        types.String `tfsdk:"file_name"`
	ContentBase64   types.String `tfsdk:"content_base64"`
	Description     types.String `tfsdk:"description"`
	URL             types.String `tfsdk:"url"`
	FileType        types.String `tfsdk:"file_type"`
	FileFormat      types.String `tfsdk:"file_format"`
	FileSize        types.Int64  `tfsdk:"file_size"`
	StorageProvider types.String `tfsdk:"storage_provider"`
	Application     types.String `tfsdk:"application"`
}

func NewResourceResource() resource.Resource {
	return &ResourceResource{}
}

func (r *ResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *ResourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor resource (uploaded file).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this resource. Determined by the provider's organization_name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The server-generated name of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the resource was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user": schema.StringAttribute{
				Description: "The Casdoor user performing the upload.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tag": schema.StringAttribute{
				Description: "The resource tag/path.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent": schema.StringAttribute{
				Description: "The parent path of the resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Description: "The full file path for upload.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_base64": schema.StringAttribute{
				Description: "Base64-encoded file content. Use filebase64() to read from disk.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description of the resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The generated download URL.",
				Computed:    true,
			},
			"file_type": schema.StringAttribute{
				Description: "The MIME type of the resource.",
				Computed:    true,
			},
			"file_format": schema.StringAttribute{
				Description: "The file format of the resource.",
				Computed:    true,
			},
			"file_size": schema.Int64Attribute{
				Description: "The size of the resource in bytes.",
				Computed:    true,
			},
			"storage_provider": schema.StringAttribute{
				Description: "The storage provider.",
				Computed:    true,
			},
			"application": schema.StringAttribute{
				Description: "The application from the client configuration.",
				Computed:    true,
			},
		},
	}
}

func (r *ResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(plan.ContentBase64.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Base64 Content",
			fmt.Sprintf("Could not decode content_base64: %s", err),
		)
		return
	}

	fileURL, name, err := r.client.UploadResourceEx(
		plan.User.ValueString(),
		plan.Tag.ValueString(),
		plan.Parent.ValueString(),
		plan.FileName.ValueString(),
		fileBytes,
		"",
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Uploading Resource",
			fmt.Sprintf("Could not upload resource: %s", err),
		)
		return
	}

	if fileURL == "" || name == "" {
		resp.Diagnostics.AddError(
			"Error Uploading Resource",
			"Casdoor returned empty URL or name after upload",
		)
		return
	}

	// The owner comes from the client's OrganizationName.
	id := r.client.OrganizationName + "/" + name

	created, err := r.client.GetResource(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource After Create",
			fmt.Sprintf("Could not read resource %q after creation: %s", id, err),
		)
		return
	}

	if created == nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource",
			fmt.Sprintf("Resource %q not found after creation", id),
		)
		return
	}

	plan.ID = types.StringValue(created.Owner + "/" + created.Name)
	plan.Owner = types.StringValue(created.Owner)
	plan.Name = types.StringValue(created.Name)
	plan.CreatedTime = types.StringValue(created.CreatedTime)
	plan.URL = types.StringValue(created.Url)
	plan.FileType = types.StringValue(created.FileType)
	plan.FileFormat = types.StringValue(created.FileFormat)
	plan.FileSize = types.Int64Value(int64(created.FileSize))
	plan.StorageProvider = types.StringValue(created.Provider)
	plan.Application = types.StringValue(created.Application)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetResource(state.Owner.ValueString() + "/" + state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource",
			fmt.Sprintf("Could not read resource %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if res == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(res.Owner + "/" + res.Name)
	state.Owner = types.StringValue(res.Owner)
	state.Name = types.StringValue(res.Name)
	state.CreatedTime = types.StringValue(res.CreatedTime)
	state.User = types.StringValue(res.User)
	state.Tag = types.StringValue(res.Tag)
	state.Parent = types.StringValue(res.Parent)
	state.Description = types.StringValue(res.Description)
	state.URL = types.StringValue(res.Url)
	state.FileType = types.StringValue(res.FileType)
	state.FileFormat = types.StringValue(res.FileFormat)
	state.FileSize = types.Int64Value(int64(res.FileSize))
	state.StorageProvider = types.StringValue(res.Provider)
	state.Application = types.StringValue(res.Application)

	// Only set file_name from server during import (when state value is unknown/null).
	// The server may strip path prefixes from the upload file_name.
	if state.FileName.IsNull() || state.FileName.IsUnknown() {
		state.FileName = types.StringValue(res.FileName)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Casdoor resources do not support updates. All attributes use RequiresReplace.",
	)
}

func (r *ResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := &casdoorsdk.Resource{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	ok, err := r.client.DeleteResource(res)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("deleting resource %q", state.Name.ValueString())) {
		return
	}
}

func (r *ResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
