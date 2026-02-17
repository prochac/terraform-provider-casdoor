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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &AdapterResource{}
	_ resource.ResourceWithConfigure   = &AdapterResource{}
	_ resource.ResourceWithImportState = &AdapterResource{}
)

type AdapterResource struct {
	client *casdoorsdk.Client
}

type AdapterResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Owner           types.String `tfsdk:"owner"`
	Name            types.String `tfsdk:"name"`
	CreatedTime     types.String `tfsdk:"created_time"`
	UseSameDb       types.Bool   `tfsdk:"use_same_db"`
	Type            types.String `tfsdk:"type"`
	DatabaseType    types.String `tfsdk:"database_type"`
	Host            types.String `tfsdk:"host"`
	Port            types.Int64  `tfsdk:"port"`
	User            types.String `tfsdk:"user"`
	Password        types.String `tfsdk:"password"`
	Database        types.String `tfsdk:"database"`
	Table           types.String `tfsdk:"table"`
	TableNamePrefix types.String `tfsdk:"table_name_prefix"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
}

func NewAdapterResource() resource.Resource {
	return &AdapterResource{}
}

func (r *AdapterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adapter"
}

func (r *AdapterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor Casbin adapter for policy storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the adapter in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this adapter.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the adapter.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the adapter was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"table": schema.StringAttribute{
				Description: "The table name for storing policies.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"use_same_db": schema.BoolAttribute{
				Description: "Whether to use the same database as Casdoor.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"type": schema.StringAttribute{
				Description: "The type of the adapter (e.g., 'Database').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"database_type": schema.StringAttribute{
				Description: "The database type (e.g., 'mysql', 'postgres', 'sqlite3').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"host": schema.StringAttribute{
				Description: "The database host address.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"port": schema.Int64Attribute{
				Description: "The database port number.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"user": schema.StringAttribute{
				Description: "The database username.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				Description: "The database password.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"database": schema.StringAttribute{
				Description: "The database name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"table_name_prefix": schema.StringAttribute{
				Description: "The table name prefix for policy storage.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether this adapter is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *AdapterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AdapterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	adapter := &casdoorsdk.Adapter{
		Owner:           plan.Owner.ValueString(),
		Name:            plan.Name.ValueString(),
		CreatedTime:     createdTime,
		UseSameDb:       plan.UseSameDb.ValueBool(),
		Type:            plan.Type.ValueString(),
		DatabaseType:    plan.DatabaseType.ValueString(),
		Host:            plan.Host.ValueString(),
		Port:            int(plan.Port.ValueInt64()),
		User:            plan.User.ValueString(),
		Password:        plan.Password.ValueString(),
		Database:        plan.Database.ValueString(),
		Table:           plan.Table.ValueString(),
		TableNamePrefix: plan.TableNamePrefix.ValueString(),
		IsEnabled:       plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.AddAdapter(adapter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Adapter",
			fmt.Sprintf("Could not create adapter %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Adapter",
			fmt.Sprintf("Casdoor returned failure when creating adapter %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the adapter to get server-generated values like CreatedTime.
	createdAdapter, err := r.client.GetAdapter(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Adapter",
			fmt.Sprintf("Could not read adapter %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdAdapter != nil {
		plan.CreatedTime = types.StringValue(createdAdapter.CreatedTime)
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AdapterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapter, err := getByOwnerName[casdoorsdk.Adapter](r.client, "get-adapter", state.Owner.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Adapter",
			fmt.Sprintf("Could not read adapter %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if adapter == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(adapter.Owner + "/" + adapter.Name)
	state.Owner = types.StringValue(adapter.Owner)
	state.Name = types.StringValue(adapter.Name)
	state.CreatedTime = types.StringValue(adapter.CreatedTime)
	state.Table = types.StringValue(adapter.Table)
	state.UseSameDb = types.BoolValue(adapter.UseSameDb)
	state.Type = types.StringValue(adapter.Type)
	state.DatabaseType = types.StringValue(adapter.DatabaseType)
	state.Host = types.StringValue(adapter.Host)
	state.Port = types.Int64Value(int64(adapter.Port))
	state.User = types.StringValue(adapter.User)
	state.Password = types.StringValue(adapter.Password)
	state.Database = types.StringValue(adapter.Database)
	state.TableNamePrefix = types.StringValue(adapter.TableNamePrefix)
	state.IsEnabled = types.BoolValue(adapter.IsEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AdapterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapter := &casdoorsdk.Adapter{
		Owner:           plan.Owner.ValueString(),
		Name:            plan.Name.ValueString(),
		CreatedTime:     plan.CreatedTime.ValueString(),
		UseSameDb:       plan.UseSameDb.ValueBool(),
		Type:            plan.Type.ValueString(),
		DatabaseType:    plan.DatabaseType.ValueString(),
		Host:            plan.Host.ValueString(),
		Port:            int(plan.Port.ValueInt64()),
		User:            plan.User.ValueString(),
		Password:        plan.Password.ValueString(),
		Database:        plan.Database.ValueString(),
		Table:           plan.Table.ValueString(),
		TableNamePrefix: plan.TableNamePrefix.ValueString(),
		IsEnabled:       plan.IsEnabled.ValueBool(),
	}

	success, err := r.client.UpdateAdapter(adapter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adapter",
			fmt.Sprintf("Could not update adapter %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Adapter",
			fmt.Sprintf("Casdoor returned failure when updating adapter %q", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AdapterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapter := &casdoorsdk.Adapter{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteAdapter(adapter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Adapter",
			fmt.Sprintf("Could not delete adapter %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Adapter",
			fmt.Sprintf("Casdoor returned failure when deleting adapter %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *AdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
