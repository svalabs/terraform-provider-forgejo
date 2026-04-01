package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &repositoryActionVariableDataSource{}
	_ datasource.DataSourceWithConfigure = &repositoryActionVariableDataSource{}
)

// repositoryActionVariableDataSource is the data source implementation.
type repositoryActionVariableDataSource struct {
	client *forgejo.Client
}

// repositoryActionVariableDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#ActionVariable
type repositoryActionVariableDataSourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	Name         types.String `tfsdk:"name"`
	Data         types.String `tfsdk:"data"`
}

// Metadata returns the data source type name.
func (d *repositoryActionVariableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_action_variable"
}

// Schema defines the schema for the data source.
func (d *repositoryActionVariableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository action variable data source.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the variable.",
				Required:    true,
			},
			"data": schema.StringAttribute{
				Description: "Data of the variable.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *repositoryActionVariableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *repositoryActionVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read repository action variable data source"))

	var (
		repo repositoryResourceModel
		data repositoryActionVariableDataSourceModel
	)

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByID(
		ctx,
		d.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Read repository action variable", map[string]any{
		"repository_id": data.RepositoryID.ValueInt64(),
		"user":          repo.Owner.ValueString(),
		"repo":          repo.Name.ValueString(),
		"name":          data.Name.ValueString(),
	})

	// Use Forgejo client to get repository action variable
	variable, res, err := d.client.GetRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.Name.ValueString(),
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 400:
				msg = fmt.Sprintf("Bad request: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Action variable with owner %s, repo %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read repository action variable", msg)

		return
	}

	// Map response body to model
	// name is omitted here, to maintain the user's configuration casing
	data.Data = types.StringValue(variable.Data)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewRepositoryActionVariableDataSource is a helper function to simplify the provider implementation.
func NewRepositoryActionVariableDataSource() datasource.DataSource {
	return &repositoryActionVariableDataSource{}
}
