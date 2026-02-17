// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &SyncerResource{}
	_ resource.ResourceWithConfigure   = &SyncerResource{}
	_ resource.ResourceWithImportState = &SyncerResource{}
)

type SyncerResource struct {
	client *casdoorsdk.Client
}

type TableColumnModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CasdoorName types.String `tfsdk:"casdoor_name"`
	IsKey       types.Bool   `tfsdk:"is_key"`
	IsHashed    types.Bool   `tfsdk:"is_hashed"`
	Values      types.List   `tfsdk:"values"`
}

type SyncerResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Owner            types.String `tfsdk:"owner"`
	Name             types.String `tfsdk:"name"`
	CreatedTime      types.String `tfsdk:"created_time"`
	Organization     types.String `tfsdk:"organization"`
	Type             types.String `tfsdk:"type"`
	Host             types.String `tfsdk:"host"`
	Port             types.Int64  `tfsdk:"port"`
	User             types.String `tfsdk:"user"`
	Password         types.String `tfsdk:"password"`
	DatabaseType     types.String `tfsdk:"database_type"`
	SslMode          types.String `tfsdk:"ssl_mode"`
	SshType          types.String `tfsdk:"ssh_type"`
	SshHost          types.String `tfsdk:"ssh_host"`
	SshPort          types.Int64  `tfsdk:"ssh_port"`
	SshUser          types.String `tfsdk:"ssh_user"`
	SshPassword      types.String `tfsdk:"ssh_password"`
	Cert             types.String `tfsdk:"cert"`
	Database         types.String `tfsdk:"database"`
	Table            types.String `tfsdk:"table"`
	TableColumns     types.List   `tfsdk:"table_columns"`
	AffiliationTable types.String `tfsdk:"affiliation_table"`
	AvatarBaseUrl    types.String `tfsdk:"avatar_base_url"`
	ErrorText        types.String `tfsdk:"error_text"`
	SyncInterval     types.Int64  `tfsdk:"sync_interval"`
	IsReadOnly       types.Bool   `tfsdk:"is_read_only"`
	IsEnabled        types.Bool   `tfsdk:"is_enabled"`
}

func NewSyncerResource() resource.Resource {
	return &SyncerResource{}
}

func (r *SyncerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syncer"
}

var tableColumnAttrTypes = map[string]attr.Type{
	"name":         types.StringType,
	"type":         types.StringType,
	"casdoor_name": types.StringType,
	"is_key":       types.BoolType,
	"is_hashed":    types.BoolType,
	"values":       types.ListType{ElemType: types.StringType},
}

func (r *SyncerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor syncer for external system synchronization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the syncer in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this syncer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the syncer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the syncer was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The organization to sync users to.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description: "The type of the syncer (e.g., 'Database').",
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
			"database_type": schema.StringAttribute{
				Description: "The database type (e.g., 'mysql', 'postgres').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ssl_mode": schema.StringAttribute{
				Description: "The SSL mode for database connections.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ssh_type": schema.StringAttribute{
				Description: "The SSH tunnel type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ssh_host": schema.StringAttribute{
				Description: "The SSH host address.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ssh_port": schema.Int64Attribute{
				Description: "The SSH port number.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"ssh_user": schema.StringAttribute{
				Description: "The SSH username.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ssh_password": schema.StringAttribute{
				Description: "The SSH password.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"cert": schema.StringAttribute{
				Description: "The certificate for database connections.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"database": schema.StringAttribute{
				Description: "The database name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"table": schema.StringAttribute{
				Description: "The table name to sync from.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"table_columns": schema.ListNestedAttribute{
				Description: "The column mappings for synchronization.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The column name in the source table.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "The column type.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"casdoor_name": schema.StringAttribute{
							Description: "The corresponding Casdoor user field name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"is_key": schema.BoolAttribute{
							Description: "Whether this column is a key column.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"is_hashed": schema.BoolAttribute{
							Description: "Whether this column value is hashed.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"values": schema.ListAttribute{
							Description: "Possible values for this column.",
							Optional:    true,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"affiliation_table": schema.StringAttribute{
				Description: "The affiliation table name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"avatar_base_url": schema.StringAttribute{
				Description: "The base URL for user avatars.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"error_text": schema.StringAttribute{
				Description: "Error text from the last sync operation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sync_interval": schema.Int64Attribute{
				Description: "The synchronization interval in seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"is_read_only": schema.BoolAttribute{
				Description: "Whether the syncer is read-only.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the syncer is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *SyncerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func syncerTableColumnsToSDK(ctx context.Context, plan SyncerResourceModel) ([]*casdoorsdk.TableColumn, error) {
	if plan.TableColumns.IsNull() || plan.TableColumns.IsUnknown() {
		return nil, nil
	}

	var columns []TableColumnModel
	diags := plan.TableColumns.ElementsAs(ctx, &columns, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert table columns")
	}

	result := make([]*casdoorsdk.TableColumn, len(columns))
	for i, col := range columns {
		var values []string
		if !col.Values.IsNull() && !col.Values.IsUnknown() {
			diags := col.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to convert column values")
			}
		}

		result[i] = &casdoorsdk.TableColumn{
			Name:        col.Name.ValueString(),
			Type:        col.Type.ValueString(),
			CasdoorName: col.CasdoorName.ValueString(),
			IsKey:       col.IsKey.ValueBool(),
			IsHashed:    col.IsHashed.ValueBool(),
			Values:      values,
		}
	}
	return result, nil
}

func (r *SyncerResource) tableColumnsFromSDK(ctx context.Context, columns []*casdoorsdk.TableColumn) (types.List, error) {
	if len(columns) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: tableColumnAttrTypes}), nil
	}

	objs := make([]attr.Value, len(columns))
	for i, col := range columns {
		values, _ := types.ListValueFrom(ctx, types.StringType, col.Values)
		if col.Values == nil {
			values = types.ListNull(types.StringType)
		}

		obj, diags := types.ObjectValue(tableColumnAttrTypes, map[string]attr.Value{
			"name":         types.StringValue(col.Name),
			"type":         types.StringValue(col.Type),
			"casdoor_name": types.StringValue(col.CasdoorName),
			"is_key":       types.BoolValue(col.IsKey),
			"is_hashed":    types.BoolValue(col.IsHashed),
			"values":       values,
		})
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: tableColumnAttrTypes}), fmt.Errorf("failed to create object")
		}
		objs[i] = obj
	}

	result, diags := types.ListValue(types.ObjectType{AttrTypes: tableColumnAttrTypes}, objs)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: tableColumnAttrTypes}), fmt.Errorf("failed to create list")
	}
	return result, nil
}

func syncerPlanToSDK(ctx context.Context, plan SyncerResourceModel, createdTime string) (*casdoorsdk.Syncer, diag.Diagnostics) {
	var diags diag.Diagnostics

	tableColumns, err := syncerTableColumnsToSDK(ctx, plan)
	if err != nil {
		diags.AddError("Error Converting Table Columns", err.Error())
		return nil, diags
	}

	return &casdoorsdk.Syncer{
		Owner:            plan.Owner.ValueString(),
		Name:             plan.Name.ValueString(),
		CreatedTime:      createdTime,
		Organization:     plan.Organization.ValueString(),
		Type:             plan.Type.ValueString(),
		Host:             plan.Host.ValueString(),
		Port:             int(plan.Port.ValueInt64()),
		User:             plan.User.ValueString(),
		Password:         plan.Password.ValueString(),
		DatabaseType:     plan.DatabaseType.ValueString(),
		SslMode:          plan.SslMode.ValueString(),
		SshType:          plan.SshType.ValueString(),
		SshHost:          plan.SshHost.ValueString(),
		SshPort:          int(plan.SshPort.ValueInt64()),
		SshUser:          plan.SshUser.ValueString(),
		SshPassword:      plan.SshPassword.ValueString(),
		Cert:             plan.Cert.ValueString(),
		Database:         plan.Database.ValueString(),
		Table:            plan.Table.ValueString(),
		TableColumns:     tableColumns,
		AffiliationTable: plan.AffiliationTable.ValueString(),
		AvatarBaseUrl:    plan.AvatarBaseUrl.ValueString(),
		SyncInterval:     int(plan.SyncInterval.ValueInt64()),
		IsReadOnly:       plan.IsReadOnly.ValueBool(),
		IsEnabled:        plan.IsEnabled.ValueBool(),
	}, diags
}

func (r *SyncerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SyncerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	syncer, diags := syncerPlanToSDK(ctx, plan, createdTime)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.AddSyncer(syncer)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating syncer %q", plan.Name.ValueString())) {
		return
	}

	// Read back the syncer to get server-generated values.
	createdSyncer, err := r.client.GetSyncer(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Syncer",
			fmt.Sprintf("Could not read syncer %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdSyncer != nil {
		plan.CreatedTime = types.StringValue(createdSyncer.CreatedTime)
		plan.ErrorText = types.StringValue(createdSyncer.ErrorText)
		tableColumnsList, err := r.tableColumnsFromSDK(ctx, createdSyncer.TableColumns)
		if err != nil {
			resp.Diagnostics.AddError("Error Converting Table Columns", err.Error())
			return
		}
		plan.TableColumns = tableColumnsList
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SyncerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SyncerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	syncer, err := r.client.GetSyncer(state.Owner.ValueString() + "/" + state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Syncer",
			fmt.Sprintf("Could not read syncer %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if syncer == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(syncer.Owner + "/" + syncer.Name)
	state.Owner = types.StringValue(syncer.Owner)
	state.Name = types.StringValue(syncer.Name)
	state.CreatedTime = types.StringValue(syncer.CreatedTime)
	state.Organization = types.StringValue(syncer.Organization)
	state.Type = types.StringValue(syncer.Type)
	state.Host = types.StringValue(syncer.Host)
	state.Port = types.Int64Value(int64(syncer.Port))
	state.User = types.StringValue(syncer.User)
	// Password is always masked by Casdoor API ("***"), preserve from state.
	if syncer.Password != "***" {
		state.Password = types.StringValue(syncer.Password)
	}
	state.DatabaseType = types.StringValue(syncer.DatabaseType)
	state.SslMode = types.StringValue(syncer.SslMode)
	state.SshType = types.StringValue(syncer.SshType)
	state.SshHost = types.StringValue(syncer.SshHost)
	state.SshPort = types.Int64Value(int64(syncer.SshPort))
	state.SshUser = types.StringValue(syncer.SshUser)
	// SshPassword is always masked by Casdoor API ("***"), preserve from state.
	if syncer.SshPassword != "***" {
		state.SshPassword = types.StringValue(syncer.SshPassword)
	}
	state.Cert = types.StringValue(syncer.Cert)
	state.Database = types.StringValue(syncer.Database)
	state.Table = types.StringValue(syncer.Table)
	state.AffiliationTable = types.StringValue(syncer.AffiliationTable)
	state.AvatarBaseUrl = types.StringValue(syncer.AvatarBaseUrl)
	state.ErrorText = types.StringValue(syncer.ErrorText)
	state.SyncInterval = types.Int64Value(int64(syncer.SyncInterval))
	state.IsReadOnly = types.BoolValue(syncer.IsReadOnly)
	state.IsEnabled = types.BoolValue(syncer.IsEnabled)

	tableColumnsList, err := r.tableColumnsFromSDK(ctx, syncer.TableColumns)
	if err != nil {
		resp.Diagnostics.AddError("Error Converting Table Columns", err.Error())
		return
	}
	state.TableColumns = tableColumnsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SyncerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SyncerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	syncer, diags := syncerPlanToSDK(ctx, plan, plan.CreatedTime.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.UpdateSyncer(syncer)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("updating syncer %q", plan.Name.ValueString())) {
		return
	}

	// Read back to get error_text if any.
	updatedSyncer, err := r.client.GetSyncer(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	if err == nil && updatedSyncer != nil {
		plan.ErrorText = types.StringValue(updatedSyncer.ErrorText)
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SyncerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SyncerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	syncer := &casdoorsdk.Syncer{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	_, err := r.client.DeleteSyncer(syncer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Syncer",
			fmt.Sprintf("Could not delete syncer %q: %s", state.Name.ValueString(), err),
		)
		return
	}
}

func (r *SyncerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
