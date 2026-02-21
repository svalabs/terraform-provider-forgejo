package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &gpgKeyResource{}
	_ resource.ResourceWithConfigure = &gpgKeyResource{}
)

// Nested resource types.
var (
	gpgKeyEmailAttrTypes = map[string]attr.Type{
		"email":    types.StringType,
		"verified": types.BoolType,
	}
	gpgKeyEmailType = types.ObjectType{
		AttrTypes: gpgKeyEmailAttrTypes,
	}

	gpgKeySubkeyAttrTypes = map[string]attr.Type{
		"id":                  types.Int64Type,
		"primary_key_id":      types.StringType,
		"key_id":              types.StringType,
		"public_key":          types.StringType,
		"can_sign":            types.BoolType,
		"can_encrypt_comms":   types.BoolType,
		"can_encrypt_storage": types.BoolType,
		"can_certify":         types.BoolType,
		"created_at":          types.StringType,
		"expires_at":          types.StringType,
	}
	gpgKeySubkeyType = types.ObjectType{
		AttrTypes: gpgKeySubkeyAttrTypes,
	}
)

// gpgKeyResource is the resource implementation.
type gpgKeyResource struct {
	client *forgejo.Client
}

// gpgKeyResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#GPGKey
type gpgKeyResourceModel struct {
	ArmoredPublicKey  types.String `tfsdk:"armored_public_key"`
	ID                types.Int64  `tfsdk:"id"`
	KeyID             types.String `tfsdk:"key_id"`
	PrimaryKeyID      types.String `tfsdk:"primary_key_id"`
	PublicKey         types.String `tfsdk:"public_key"`
	CanSign           types.Bool   `tfsdk:"can_sign"`
	CanEncryptComms   types.Bool   `tfsdk:"can_encrypt_comms"`
	CanEncryptStorage types.Bool   `tfsdk:"can_encrypt_storage"`
	CanCertify        types.Bool   `tfsdk:"can_certify"`
	Created           types.String `tfsdk:"created_at"`
	Expires           types.String `tfsdk:"expires_at"`
	Emails            types.List   `tfsdk:"emails"`
	SubKeys           types.List   `tfsdk:"subkeys"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *gpgKeyResourceModel) from(k *forgejo.GPGKey) (diags diag.Diagnostics) {
	if k == nil {
		return diags
	}

	m.ID = types.Int64Value(k.ID)
	m.KeyID = types.StringValue(k.KeyID)
	m.PrimaryKeyID = types.StringValue(k.PrimaryKeyID)
	m.PublicKey = types.StringValue(k.PublicKey)
	m.CanSign = types.BoolValue(k.CanSign)
	m.CanEncryptComms = types.BoolValue(k.CanEncryptComms)
	m.CanEncryptStorage = types.BoolValue(k.CanEncryptStorage)
	m.CanCertify = types.BoolValue(k.CanCertify)
	m.Created = types.StringValue(k.Created.Format(time.RFC3339))
	m.Expires = types.StringValue(k.Expires.Format(time.RFC3339))

	emails, diags := getEmails(k)
	if diags.HasError() {
		return diags
	}
	m.Emails = emails

	subkeys, diags := getSubkeys(k)
	if diags.HasError() {
		return diags
	}
	m.SubKeys = subkeys

	return diags
}

// to is a helper function to save Terraform data model into an API struct.
func (m *gpgKeyResourceModel) to(o *forgejo.CreateGPGKeyOption) {
	if o == nil {
		o = new(forgejo.CreateGPGKeyOption)
	}

	o.ArmoredKey = m.ArmoredPublicKey.ValueString()
}

// getEmails is a helper function to convert a list of emails into a Terraform list.
func getEmails(k *forgejo.GPGKey) (types.List, diag.Diagnostics) {
	emailElements := make([]attr.Value, 0, len(k.Emails))

	for _, e := range k.Emails {
		values := map[string]attr.Value{
			"email":    types.StringValue(e.Email),
			"verified": types.BoolValue(e.Verified),
		}
		elem, diags := types.ObjectValue(gpgKeyEmailAttrTypes, values)

		if diags.HasError() {
			return types.List{}, diags
		}

		emailElements = append(emailElements, elem)
	}

	emails, diags := types.ListValue(gpgKeyEmailType, emailElements)
	if diags.HasError() {
		return types.List{}, diags
	}

	return emails, diags
}

// getSubkeys is a helper function to convert a list of subkeys into a Terraform list.
func getSubkeys(k *forgejo.GPGKey) (types.List, diag.Diagnostics) {
	subkeyElements := make([]attr.Value, 0, len(k.SubsKey))

	for _, e := range k.SubsKey {
		values := map[string]attr.Value{
			"id":                  types.Int64Value(e.ID),
			"primary_key_id":      types.StringValue(e.PrimaryKeyID),
			"key_id":              types.StringValue(e.KeyID),
			"public_key":          types.StringValue(e.PublicKey),
			"can_sign":            types.BoolValue(e.CanSign),
			"can_encrypt_comms":   types.BoolValue(e.CanEncryptComms),
			"can_encrypt_storage": types.BoolValue(e.CanEncryptStorage),
			"can_certify":         types.BoolValue(e.CanCertify),
			"created_at":          types.StringValue(e.Created.Format(time.RFC3339)),
			"expires_at":          types.StringValue(e.Expires.Format(time.RFC3339)),
		}
		elem, diags := types.ObjectValue(gpgKeySubkeyAttrTypes, values)

		if diags.HasError() {
			return types.List{}, diags
		}

		subkeyElements = append(subkeyElements, elem)
	}

	subkeys, diags := types.ListValue(gpgKeySubkeyType, subkeyElements)
	if diags.HasError() {
		return types.List{}, diags
	}

	return subkeys, diags
}

// Metadata returns the resource type name.
func (r *gpgKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpg_key"
}

// Schema defines the schema for the resource.
func (r *gpgKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo user GPG key resource.",

		Attributes: map[string]schema.Attribute{
			"armored_public_key": schema.StringAttribute{
				Description: "Armored GPG Public key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the GPG key.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"key_id": schema.StringAttribute{
				Description: "ID of the GPG key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_key_id": schema.StringAttribute{
				Description: "Primary ID of the GPG key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_key": schema.StringAttribute{
				Description: "The public key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"can_sign": schema.BoolAttribute{
				Description: "Can this key sign.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"can_encrypt_comms": schema.BoolAttribute{
				Description: "Can this key encrypt communications.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"can_encrypt_storage": schema.BoolAttribute{
				Description: "Can this key encrypt storage.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"can_certify": schema.BoolAttribute{
				Description: "Can this key certify.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the GPG key was created.",
				Computed:    true,
				// 6b66d9e: standardize on formatting temporal data in RFC3339 format
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"expires_at": schema.StringAttribute{
				Description: "Time at which the GPG key expires.",
				Computed:    true,
				// 6b66d9e: standardize on formatting temporal data in RFC3339 format
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"emails": schema.ListAttribute{
				Description: "Emails associated with the GPG key.",
				Computed:    true,
				ElementType: gpgKeyEmailType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"subkeys": schema.ListAttribute{
				Description: "Subkeys of the GPG key.",
				Computed:    true,
				ElementType: gpgKeySubkeyType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *gpgKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *gpgKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create GPG key resource"))

	var data gpgKeyResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create GPG key", map[string]any{
		"armored_public_key": data.ArmoredPublicKey.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateGPGKeyOption{}
	data.to(&opts)

	// Use Forgejo client to create new GPG key
	key, res, err := r.client.CreateGPGKey(opts)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf(
					"GPG key creation forbidden: %s",
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"GPG key creation not found: %s",
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create GPG key", msg)

		return
	}

	// Map response body to model
	diags = data.from(key)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *gpgKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read GPG key resource"))

	var data gpgKeyResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get GPG key by id", map[string]any{
		"key_id": data.KeyID.ValueString(),
	})

	// Use Forgejo client to get GPG key
	key, res, err := r.client.GetGPGKey(data.ID.ValueInt64())
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf(
					"GPG key with id %s forbidden: %s",
					data.ID.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"GPG key with id %s not found: %s",
					data.ID.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to get GPG key by id", msg)

		return
	}

	// Map response body to model
	diags = data.from(key)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *gpgKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update GPG key resource"))

	/*
	 * GPG keys can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *gpgKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete GPG key resource"))

	var data gpgKeyResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete GPG key", map[string]any{
		"key_id": data.KeyID.ValueString(),
	})

	// Use Forgejo client to delete existing GPG key
	res, err := r.client.DeleteGPGKey(data.ID.ValueInt64())
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf(
					"GPG key with id %s forbidden: %s",
					data.ID.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"GPG key with id %s not found: %s",
					data.ID.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete GPG key", msg)

		return
	}
}

// NewGPGKeyResource is a helper function to simplify the provider implementation.
func NewGPGKeyResource() resource.Resource {
	return &gpgKeyResource{}
}
