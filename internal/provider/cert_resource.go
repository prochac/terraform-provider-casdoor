// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CertResource{}
	_ resource.ResourceWithConfigure   = &CertResource{}
	_ resource.ResourceWithImportState = &CertResource{}
)

type CertResource struct {
	client *casdoorsdk.Client
}

type CertResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Owner                  types.String `tfsdk:"owner"`
	Name                   types.String `tfsdk:"name"`
	CreatedTime            types.String `tfsdk:"created_time"`
	DisplayName            types.String `tfsdk:"display_name"`
	Scope                  types.String `tfsdk:"scope"`
	Type                   types.String `tfsdk:"type"`
	CryptoAlgorithm        types.String `tfsdk:"crypto_algorithm"`
	BitSize                types.Int64  `tfsdk:"bit_size"`
	ExpireInYears          types.Int64  `tfsdk:"expire_in_years"`
	Certificate            types.String `tfsdk:"certificate"`
	PrivateKey             types.String `tfsdk:"private_key"`
	AuthorityPublicKey     types.String `tfsdk:"authority_public_key"`
	AuthorityRootPublicKey types.String `tfsdk:"authority_root_public_key"`
}

func NewCertResource() resource.Resource {
	return &CertResource{}
}

func (r *CertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cert"
}

func (r *CertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the certificate in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The organization that owns this certificate.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the certificate.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the certificate was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the certificate.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"scope": schema.StringAttribute{
				Description: "The scope of the certificate (e.g., 'JWT').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("JWT"),
			},
			"type": schema.StringAttribute{
				Description: "The type of the certificate (e.g., 'x509').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("x509"),
			},
			"crypto_algorithm": schema.StringAttribute{
				Description: "The cryptographic algorithm (e.g., 'RS256').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("RS256"),
			},
			"bit_size": schema.Int64Attribute{
				Description: "The key bit size (e.g., 2048, 4096).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(4096),
			},
			"expire_in_years": schema.Int64Attribute{
				Description: "The certificate expiration in years.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(20),
			},
			"certificate": schema.StringAttribute{
				Description: "The X.509 certificate (PEM format).",
				Optional:    true,
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The private key (PEM format).",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"authority_public_key": schema.StringAttribute{
				Description: "The authority public key (PEM format).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"authority_root_public_key": schema.StringAttribute{
				Description: "The authority root public key (PEM format).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *CertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Certificate.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Certificate",
			"The certificate attribute must be provided when creating a cert resource.",
		)
	}
	if plan.PrivateKey.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Private Key",
			"The private_key attribute must be provided when creating a cert resource.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	cert := &casdoorsdk.Cert{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		CreatedTime:            createdTime,
		DisplayName:            plan.DisplayName.ValueString(),
		Scope:                  plan.Scope.ValueString(),
		Type:                   plan.Type.ValueString(),
		CryptoAlgorithm:        plan.CryptoAlgorithm.ValueString(),
		BitSize:                int(plan.BitSize.ValueInt64()),
		ExpireInYears:          int(plan.ExpireInYears.ValueInt64()),
		Certificate:            plan.Certificate.ValueString(),
		PrivateKey:             plan.PrivateKey.ValueString(),
		AuthorityPublicKey:     plan.AuthorityPublicKey.ValueString(),
		AuthorityRootPublicKey: plan.AuthorityRootPublicKey.ValueString(),
	}

	success, err := r.client.AddCert(cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Certificate",
			fmt.Sprintf("Could not create certificate %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Creating Certificate",
			fmt.Sprintf("Casdoor returned failure when creating certificate %q", plan.Name.ValueString()),
		)
		return
	}

	// Read back the cert to get generated values.
	createdCert, err := r.client.GetCert(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Certificate After Create",
			fmt.Sprintf("Could not read certificate %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdCert == nil {
		resp.Diagnostics.AddError(
			"Error Reading Certificate",
			fmt.Sprintf("Certificate %q not found after creation", plan.Name.ValueString()),
		)
		return
	}

	createdTime = createdCert.CreatedTime
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}
	plan.CreatedTime = types.StringValue(createdTime)
	plan.Certificate = types.StringValue(createdCert.Certificate)
	plan.PrivateKey = types.StringValue(createdCert.PrivateKey)

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := getByOwnerName[casdoorsdk.Cert](r.client, "get-cert", state.Owner.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Certificate",
			fmt.Sprintf("Could not read certificate %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if cert == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(cert.Owner + "/" + cert.Name)
	state.Owner = types.StringValue(cert.Owner)
	state.Name = types.StringValue(cert.Name)
	state.CreatedTime = types.StringValue(cert.CreatedTime)
	state.DisplayName = types.StringValue(cert.DisplayName)
	state.Scope = types.StringValue(cert.Scope)
	state.Type = types.StringValue(cert.Type)
	state.CryptoAlgorithm = types.StringValue(cert.CryptoAlgorithm)
	state.BitSize = types.Int64Value(int64(cert.BitSize))
	state.ExpireInYears = types.Int64Value(int64(cert.ExpireInYears))
	state.Certificate = types.StringValue(cert.Certificate)
	state.PrivateKey = types.StringValue(cert.PrivateKey)
	state.AuthorityPublicKey = types.StringValue(cert.AuthorityPublicKey)
	state.AuthorityRootPublicKey = types.StringValue(cert.AuthorityRootPublicKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert := &casdoorsdk.Cert{
		Owner:                  plan.Owner.ValueString(),
		Name:                   plan.Name.ValueString(),
		CreatedTime:            plan.CreatedTime.ValueString(),
		DisplayName:            plan.DisplayName.ValueString(),
		Scope:                  plan.Scope.ValueString(),
		Type:                   plan.Type.ValueString(),
		CryptoAlgorithm:        plan.CryptoAlgorithm.ValueString(),
		BitSize:                int(plan.BitSize.ValueInt64()),
		ExpireInYears:          int(plan.ExpireInYears.ValueInt64()),
		Certificate:            plan.Certificate.ValueString(),
		PrivateKey:             plan.PrivateKey.ValueString(),
		AuthorityPublicKey:     plan.AuthorityPublicKey.ValueString(),
		AuthorityRootPublicKey: plan.AuthorityRootPublicKey.ValueString(),
	}

	success, err := r.client.UpdateCert(cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Certificate",
			fmt.Sprintf("Could not update certificate %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Updating Certificate",
			fmt.Sprintf("Casdoor returned failure when updating certificate %q", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert := &casdoorsdk.Cert{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	success, err := r.client.DeleteCert(cert)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Certificate",
			fmt.Sprintf("Could not delete certificate %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Certificate",
			fmt.Sprintf("Casdoor returned failure when deleting certificate %q", state.Name.ValueString()),
		)
		return
	}
}

func (r *CertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
