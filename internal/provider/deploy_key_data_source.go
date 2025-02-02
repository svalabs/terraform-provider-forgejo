package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deployKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &deployKeyDataSource{}
)

// deployKeyDataSource is the data source implementation.
type deployKeyDataSource struct {
	client *forgejo.Client
}

// deployKeyDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#CreateKeyOption
type deployKeyDataSourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	KeyID        types.Int64  `tfsdk:"key_id"`
	Key          types.String `tfsdk:"key"`
	URL          types.String `tfsdk:"url"`
	Title        types.String `tfsdk:"title"`
	Fingerprint  types.String `tfsdk:"fingerprint"`
	Created      types.String `tfsdk:"created_at"`
	ReadOnly     types.Bool   `tfsdk:"read_only"`
}

// Metadata returns the data source type name.
func (d *deployKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy_key"
}

// Schema defines the schema for the data source.
func (d *deployKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo deploy key data source",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
			},
			"key_id": schema.Int64Attribute{
				Description: "Numeric identifier of the deploy key.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "Armored SSH key",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL of the deploy key.",
				Computed:    true,
			},
			"title": schema.StringAttribute{
				Description: "Title of the deploy key.",
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "Fingerprint of the deploy key.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the deploy key was created.",
				Computed:    true,
			},
			"read_only": schema.BoolAttribute{
				Description: "Does the key have only read access?",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *deployKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *deployKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read deploy key data source"))

	var (
		repo repositoryResourceModel
		data deployKeyDataSourceModel
	)

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := d.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Get deploy key by id", map[string]any{
		"user":   repo.Owner.ValueString(),
		"repo":   repo.Name.ValueString(),
		"key_id": data.KeyID.ValueInt64(),
	})

	// Use Forgejo client to get deploy key
	key, res, err := d.client.GetDeployKey(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.KeyID.ValueInt64(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Deploy key with user %s repo %s and id %d not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.KeyID.ValueInt64(),
				err,
			)
		// TODO: check other error codes
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get deploy key by id", msg)

		return
	}

	// Map response body to model
	data.KeyID = types.Int64Value(key.KeyID)
	data.Key = types.StringValue(key.Key)
	data.URL = types.StringValue(key.URL)
	data.Title = types.StringValue(key.Title)
	data.Fingerprint = types.StringValue(key.Fingerprint)
	data.Created = types.StringValue(key.Created.String())
	data.ReadOnly = types.BoolValue(key.ReadOnly)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewDeployKeyDataSource is a helper function to simplify the provider implementation.
func NewDeployKeyDataSource() datasource.DataSource {
	return &deployKeyDataSource{}
}
