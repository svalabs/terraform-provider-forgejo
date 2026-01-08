package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sshKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &sshKeyDataSource{}
)

// sshKeyDataSource is the data source implementation.
type sshKeyDataSource struct {
	client *forgejo.Client
}

// sshKeyDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#PublicKey
type sshKeyDataSourceModel struct {
	User        types.String `tfsdk:"user"`
	Title       types.String `tfsdk:"title"`
	KeyID       types.Int64  `tfsdk:"key_id"`
	Key         types.String `tfsdk:"key"`
	URL         types.String `tfsdk:"url"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created_at"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	KeyType     types.String `tfsdk:"key_type"`
}

// Metadata returns the data source type name.
func (d *sshKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

// Schema defines the schema for the data source.
func (d *sshKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo user SSH key data source.",

		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description: "Name of the user.",
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: "Title of the SSH key.",
				Required:    true,
			},
			"key_id": schema.Int64Attribute{
				Description: "Numeric identifier of the SSH key.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "Armored SSH key.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL of the SSH key.",
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "Fingerprint of the SSH key.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the SSH key was created.",
				Computed:    true,
			},
			"read_only": schema.BoolAttribute{
				Description: "Does the key have only read access?",
				Computed:    true,
			},
			"key_type": schema.StringAttribute{
				Description: "Type of the SSH key.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *sshKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *sshKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read SSH key data source"))

	var data sshKeyDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "List SSH keys", map[string]any{
		"user": data.User.ValueString(),
	})

	// Use Forgejo client to list SSH keys
	keys, res, err := d.client.ListPublicKeys(
		data.User.ValueString(),
		forgejo.ListPublicKeysOptions{},
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"SSH keys for user %s not found: %s",
				data.User.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to list SSH keys", msg)

		return
	}

	// Search for SSH key with given title
	idx := slices.IndexFunc(keys, func(k *forgejo.PublicKey) bool {
		return k.Title == data.Title.ValueString()
	})
	if idx == -1 {
		resp.Diagnostics.AddError(
			"Unable to get SSH key by title",
			fmt.Sprintf(
				"SSH key with user %s and title %s not found.",
				data.User.String(),
				data.Title.String(),
			),
		)

		return
	}

	// Map response body to model
	data.KeyID = types.Int64Value(keys[idx].ID)
	data.Key = types.StringValue(keys[idx].Key)
	data.URL = types.StringValue(keys[idx].URL)
	data.Title = types.StringValue(keys[idx].Title)
	data.Fingerprint = types.StringValue(keys[idx].Fingerprint)
	data.Created = types.StringValue(keys[idx].Created.String())
	data.ReadOnly = types.BoolValue(keys[idx].ReadOnly)
	data.KeyType = types.StringValue(keys[idx].KeyType)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewSSHKeyDataSource is a helper function to simplify the provider implementation.
func NewSSHKeyDataSource() datasource.DataSource {
	return &sshKeyDataSource{}
}
