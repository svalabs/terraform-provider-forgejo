package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &collaboratorDataSource{}
	_ datasource.DataSourceWithConfigure = &collaboratorDataSource{}
)

// collaboratorDataSource is the data source implementation.
type collaboratorDataSource struct {
	client *forgejo.Client
}

// collaboratorDataSourceModel maps the data source schema data.
type collaboratorDataSourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	User         types.String `tfsdk:"user"`
	Permission   types.String `tfsdk:"permission"`
}

// Metadata returns the data source type name.
func (d *collaboratorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collaborator"
}

// Schema defines the schema for the data source.
func (d *collaboratorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository collaborator data source.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "Username of the collaborator.",
				Required:    true,
			},
			"permission": schema.StringAttribute{
				Description: "Repository permissions of the collaborator.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *collaboratorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *collaboratorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read collaborator data source"))

	var (
		repo repositoryResourceModel
		data collaboratorDataSourceModel
	)

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	rep, diags := getRepositoryByID(ctx, d.client, data.RepositoryID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Get collaborator by username", map[string]any{
		"owner":        repo.Owner.ValueString(),
		"repo":         repo.Name.ValueString(),
		"collaborator": data.User.ValueString(),
	})

	// Use Forgejo client to get collaborator permission
	perms, res, err := d.client.CollaboratorPermission(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.User.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s forbidden: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get collaborator permission", msg)

		return
	}

	// Map response body to model
	data.Permission = types.StringValue(string(perms.Permission))

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewCollaboratorDataSource is a helper function to simplify the provider implementation.
func NewCollaboratorDataSource() datasource.DataSource {
	return &collaboratorDataSource{}
}
