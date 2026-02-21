package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &gpgKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &gpgKeyDataSource{}
)

// gpgKeyDataSource is the data source implementation.
type gpgKeyDataSource struct {
	client *forgejo.Client
}

// gpgKeyDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#GPGKey
type gpgKeyDataSourceModel struct {
	User              types.String `tfsdk:"user"`
	KeyID             types.String `tfsdk:"key_id"`
	ID                types.Int64  `tfsdk:"id"`
	PrimaryKeyID      types.String `tfsdk:"primary_key_id"`
	PublicKey         types.String `tfsdk:"public_key"`
	CanSign           types.Bool   `tfsdk:"can_sign"`
	CanEncryptComms   types.Bool   `tfsdk:"can_encrypt_comms"`
	CanEncryptStorage types.Bool   `tfsdk:"can_encrypt_storage"`
	CanCertify        types.Bool   `tfsdk:"can_certify"`
	Created           types.String `tfsdk:"created_at"`
	Expires           types.String `tfsdk:"expires_at"`
	Emails            types.List   `tfsdk:"emails"`
	Subkeys           types.List   `tfsdk:"subkeys"`
}

// Metadata returns the data source type name.
func (d *gpgKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpg_key"
}

// Schema defines the schema for the data source.
func (d *gpgKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo user GPG key data source.",

		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Description: "ID of the GPG key.",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "Name of the user.",
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the GPG key.",
				Computed:    true,
			},
			"primary_key_id": schema.StringAttribute{
				Description: "Primary ID of the GPG key.",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "The public key.",
				Computed:    true,
			},
			"can_sign": schema.BoolAttribute{
				Description: "Can this key sign.",
				Computed:    true,
			},
			"can_encrypt_comms": schema.BoolAttribute{
				Description: "Can this key encrypt communications.",
				Computed:    true,
			},
			"can_encrypt_storage": schema.BoolAttribute{
				Description: "Can this key encrypt storage.",
				Computed:    true,
			},
			"can_certify": schema.BoolAttribute{
				Description: "Can this key certify.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the GPG key was created.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Time at which the GPG key expires.",
				Computed:    true,
			},
			"emails": schema.ListAttribute{
				Description: "Emails associated with the GPG key.",
				Computed:    true,
				ElementType: gpgKeyEmailType,
			},
			"subkeys": schema.ListAttribute{
				Description: "Subkeys of the GPG key.",
				Computed:    true,
				ElementType: gpgKeySubkeyType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *gpgKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *gpgKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read GPG key data source"))

	var data gpgKeyDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "List GPG keys", map[string]any{
		"user": data.User.ValueString(),
	})

	var (
		keys []*forgejo.GPGKey
		res  *forgejo.Response
		err  error
	)

	// Use Forgejo client to list GPG keys
	if data.User.ValueString() != "" {
		keys, res, err = d.client.ListGPGKeys(
			data.User.ValueString(),
			forgejo.ListGPGKeysOptions{},
		)
	} else {
		keys, res, err = d.client.ListMyGPGKeys(
			&forgejo.ListGPGKeysOptions{},
		)
	}
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 404:
				// If the user was not provided, we should never get a 404, so the message here should always have a user.
				msg = fmt.Sprintf(
					"GPG keys for user %s not found: %s",
					data.User.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to list GPG keys", msg)

		return
	}

	// Search for GPG key with given title
	idx := slices.IndexFunc(keys, func(k *forgejo.GPGKey) bool {
		return strings.EqualFold(k.KeyID, data.KeyID.ValueString())
	})
	if idx == -1 {
		var msg string
		if data.User.ValueString() != "" {
			msg = fmt.Sprintf(
				"GPG key with user %s and key_id %s not found",
				data.User.String(),
				data.KeyID.String(),
			)
		} else {
			msg = fmt.Sprintf(
				"GPG key with key_id %s not found",
				data.KeyID.String(),
			)
		}
		resp.Diagnostics.AddError("Unable to find GPG key by ID", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(keys[idx].ID)
	data.KeyID = types.StringValue(keys[idx].KeyID)
	data.PrimaryKeyID = types.StringValue(keys[idx].PrimaryKeyID)
	data.PublicKey = types.StringValue(keys[idx].PublicKey)
	data.CanSign = types.BoolValue(keys[idx].CanSign)
	data.CanEncryptComms = types.BoolValue(keys[idx].CanEncryptComms)
	data.CanEncryptStorage = types.BoolValue(keys[idx].CanEncryptStorage)
	data.CanCertify = types.BoolValue(keys[idx].CanCertify)
	data.Created = types.StringValue(keys[idx].Created.Format(time.RFC3339))
	data.Expires = types.StringValue(keys[idx].Expires.Format(time.RFC3339))

	data.Emails, diags = getEmails(keys[idx])
	resp.Diagnostics.Append(diags...)

	data.Subkeys, diags = getSubkeys(keys[idx])
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewGPGKeyDataSource is a helper function to simplify the provider implementation.
func NewGPGKeyDataSource() datasource.DataSource {
	return &gpgKeyDataSource{}
}
