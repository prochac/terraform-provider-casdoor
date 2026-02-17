// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &LdapResource{}
	_ resource.ResourceWithConfigure   = &LdapResource{}
	_ resource.ResourceWithImportState = &LdapResource{}
)

type LdapResource struct {
	client *casdoorsdk.Client
}

type LdapResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Owner               types.String `tfsdk:"owner"`
	CreatedTime         types.String `tfsdk:"created_time"`
	ServerName          types.String `tfsdk:"server_name"`
	Host                types.String `tfsdk:"host"`
	Port                types.Int64  `tfsdk:"port"`
	EnableSsl           types.Bool   `tfsdk:"enable_ssl"`
	AllowSelfSignedCert types.Bool   `tfsdk:"allow_self_signed_cert"`
	Username            types.String `tfsdk:"username"`
	Password            types.String `tfsdk:"password"`
	BaseDn              types.String `tfsdk:"base_dn"`
	Filter              types.String `tfsdk:"filter"`
	FilterFields        types.List   `tfsdk:"filter_fields"`
	DefaultGroup        types.String `tfsdk:"default_group"`
	PasswordType        types.String `tfsdk:"password_type"`
	CustomAttributes    types.Map    `tfsdk:"custom_attributes"`
	AutoSync            types.Int64  `tfsdk:"auto_sync"`
	LastSync            types.String `tfsdk:"last_sync"`
}

func NewLdapResource() resource.Resource {
	return &LdapResource{}
}

func (r *LdapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap"
}

func (r *LdapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor LDAP configuration for user synchronization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the LDAP configuration.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this LDAP configuration.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the LDAP configuration was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_name": schema.StringAttribute{
				Description: "A friendly name for the LDAP server.",
				Required:    true,
			},
			"host": schema.StringAttribute{
				Description: "The LDAP server hostname or IP address.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The LDAP server port (typically 389 for LDAP, 636 for LDAPS).",
				Required:    true,
			},
			"enable_ssl": schema.BoolAttribute{
				Description: "Whether to use SSL/TLS for the LDAP connection.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_self_signed_cert": schema.BoolAttribute{
				Description: "Whether to allow self-signed certificates when using SSL.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"username": schema.StringAttribute{
				Description: "The bind DN (Distinguished Name) for authenticating to the LDAP server.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the bind DN.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"base_dn": schema.StringAttribute{
				Description: "The base DN for LDAP searches.",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "The LDAP filter for searching users (e.g., '(objectClass=posixAccount)').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"filter_fields": schema.ListAttribute{
				Description: "List of LDAP attributes to use as filter fields.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"default_group": schema.StringAttribute{
				Description: "The default group to assign to synchronized users.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_type": schema.StringAttribute{
				Description: "The password hashing algorithm used by LDAP (e.g., 'plain', 'md5', 'sha256').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"custom_attributes": schema.MapAttribute{
				Description: "Custom attribute mappings from LDAP to Casdoor user fields.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"auto_sync": schema.Int64Attribute{
				Description: "Auto-sync interval in minutes. 0 means no auto-sync.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"last_sync": schema.StringAttribute{
				Description: "The timestamp of the last synchronization.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *LdapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LdapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LdapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert filter_fields list to Go slice.
	filterFields := make([]string, 0)
	if !plan.FilterFields.IsNull() {
		resp.Diagnostics.Append(plan.FilterFields.ElementsAs(ctx, &filterFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert custom_attributes map to Go map.
	var customAttributes map[string]string
	if !plan.CustomAttributes.IsNull() {
		customAttributes = make(map[string]string)
		resp.Diagnostics.Append(plan.CustomAttributes.ElementsAs(ctx, &customAttributes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	ldap := &casdoorsdk.Ldap{
		Id:                  plan.Id.ValueString(),
		Owner:               plan.Owner.ValueString(),
		CreatedTime:         createdTime,
		ServerName:          plan.ServerName.ValueString(),
		Host:                plan.Host.ValueString(),
		Port:                int(plan.Port.ValueInt64()),
		EnableSsl:           plan.EnableSsl.ValueBool(),
		AllowSelfSignedCert: plan.AllowSelfSignedCert.ValueBool(),
		Username:            plan.Username.ValueString(),
		Password:            plan.Password.ValueString(),
		BaseDn:              plan.BaseDn.ValueString(),
		Filter:              plan.Filter.ValueString(),
		FilterFields:        filterFields,
		DefaultGroup:        plan.DefaultGroup.ValueString(),
		PasswordType:        plan.PasswordType.ValueString(),
		CustomAttributes:    customAttributes,
		AutoSync:            int(plan.AutoSync.ValueInt64()),
	}

	success, err := r.client.AddLdap(ldap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating LDAP",
			fmt.Sprintf("Could not create LDAP %q: %s", plan.Id.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating LDAP",
			fmt.Sprintf("Casdoor returned failure when creating LDAP %q", plan.Id.ValueString()),
		)
		return
	}

	// Read back the LDAP to get server-generated values.
	createdLdap, err := r.client.GetLdap(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading LDAP",
			fmt.Sprintf("Could not read LDAP %q after creation: %s", plan.Id.ValueString(), err),
		)
		return
	}

	if createdLdap != nil {
		plan.CreatedTime = types.StringValue(createdLdap.CreatedTime)
		plan.LastSync = types.StringValue(createdLdap.LastSync)
	}

	// Set list values to null if empty to match plan.
	if len(filterFields) == 0 {
		plan.FilterFields = types.ListNull(types.StringType)
	}
	if len(customAttributes) == 0 {
		plan.CustomAttributes = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LdapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LdapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ldap, err := getByOwnerName[casdoorsdk.Ldap](r.client, "get-ldap", state.Owner.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading LDAP",
			fmt.Sprintf("Could not read LDAP %q: %s", state.Id.ValueString(), err),
		)
		return
	}

	if ldap == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(ldap.Id)
	state.Owner = types.StringValue(ldap.Owner)
	state.CreatedTime = types.StringValue(ldap.CreatedTime)
	state.ServerName = types.StringValue(ldap.ServerName)
	state.Host = types.StringValue(ldap.Host)
	state.Port = types.Int64Value(int64(ldap.Port))
	state.EnableSsl = types.BoolValue(ldap.EnableSsl)
	state.AllowSelfSignedCert = types.BoolValue(ldap.AllowSelfSignedCert)
	state.Username = types.StringValue(ldap.Username)
	// Password is always masked by Casdoor API ("***"), preserve from state.
	// On import (when state is null) fall back to empty string.
	if ldap.Password == "***" {
		if state.Password.IsNull() {
			state.Password = types.StringValue("")
		}
	} else {
		state.Password = types.StringValue(ldap.Password)
	}
	state.BaseDn = types.StringValue(ldap.BaseDn)
	state.Filter = types.StringValue(ldap.Filter)
	state.DefaultGroup = types.StringValue(ldap.DefaultGroup)
	state.PasswordType = types.StringValue(ldap.PasswordType)
	state.AutoSync = types.Int64Value(int64(ldap.AutoSync))
	state.LastSync = types.StringValue(ldap.LastSync)

	// Convert filter_fields to list type.
	if len(ldap.FilterFields) > 0 {
		filterFields, diags := types.ListValueFrom(ctx, types.StringType, ldap.FilterFields)
		resp.Diagnostics.Append(diags...)
		state.FilterFields = filterFields
	} else {
		state.FilterFields = types.ListNull(types.StringType)
	}

	// Convert custom_attributes to map type.
	if len(ldap.CustomAttributes) > 0 {
		attrValues := make(map[string]attr.Value)
		for k, v := range ldap.CustomAttributes {
			attrValues[k] = types.StringValue(v)
		}
		customAttrs, diags := types.MapValue(types.StringType, attrValues)
		resp.Diagnostics.Append(diags...)
		state.CustomAttributes = customAttrs
	} else {
		state.CustomAttributes = types.MapNull(types.StringType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LdapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LdapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert filter_fields list to Go slice.
	filterFields := make([]string, 0)
	if !plan.FilterFields.IsNull() {
		resp.Diagnostics.Append(plan.FilterFields.ElementsAs(ctx, &filterFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert custom_attributes map to Go map.
	var customAttributes map[string]string
	if !plan.CustomAttributes.IsNull() {
		customAttributes = make(map[string]string)
		resp.Diagnostics.Append(plan.CustomAttributes.ElementsAs(ctx, &customAttributes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	ldap := &casdoorsdk.Ldap{
		Id:                  plan.Id.ValueString(),
		Owner:               plan.Owner.ValueString(),
		CreatedTime:         plan.CreatedTime.ValueString(),
		ServerName:          plan.ServerName.ValueString(),
		Host:                plan.Host.ValueString(),
		Port:                int(plan.Port.ValueInt64()),
		EnableSsl:           plan.EnableSsl.ValueBool(),
		AllowSelfSignedCert: plan.AllowSelfSignedCert.ValueBool(),
		Username:            plan.Username.ValueString(),
		Password:            plan.Password.ValueString(),
		BaseDn:              plan.BaseDn.ValueString(),
		Filter:              plan.Filter.ValueString(),
		FilterFields:        filterFields,
		DefaultGroup:        plan.DefaultGroup.ValueString(),
		PasswordType:        plan.PasswordType.ValueString(),
		CustomAttributes:    customAttributes,
		AutoSync:            int(plan.AutoSync.ValueInt64()),
		LastSync:            plan.LastSync.ValueString(),
	}

	success, err := r.client.UpdateLdap(ldap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating LDAP",
			fmt.Sprintf("Could not update LDAP %q: %s", plan.Id.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating LDAP",
			fmt.Sprintf("Casdoor returned failure when updating LDAP %q", plan.Id.ValueString()),
		)
		return
	}

	// Set list values to null if empty to match plan.
	if len(filterFields) == 0 {
		plan.FilterFields = types.ListNull(types.StringType)
	}
	if len(customAttributes) == 0 {
		plan.CustomAttributes = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LdapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LdapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ldap := &casdoorsdk.Ldap{
		Id:    state.Id.ValueString(),
		Owner: state.Owner.ValueString(),
	}

	success, err := r.client.DeleteLdap(ldap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting LDAP",
			fmt.Sprintf("Could not delete LDAP %q: %s", state.Id.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting LDAP",
			fmt.Sprintf("Casdoor returned failure when deleting LDAP %q", state.Id.ValueString()),
		)
		return
	}
}

func (r *LdapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
