package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &organizationActionVariableDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationActionVariableDataSource{}
)

// organizationActionVariableDataSource is the data source implementation.
type organizationActionVariableDataSource struct {
	client *forgejo.Client
}

// organizationActionVariableDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#ActionVariable
type organizationActionVariableDataSourceModel struct {
	Organization   types.String `tfsdk:"organization"`
	OrganizationID types.Int64  `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Data           types.String `tfsdk:"data"`
}

// Metadata returns the data source type name.
func (d *organizationActionVariableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_action_variable"
}

// Schema defines the schema for the data source.
func (d *organizationActionVariableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo organization action variable data source.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the owning organization. **Note**: One of `organization` or `organization_id` must be specified.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("organization_id"),
					}...),
				},
			},
			"organization_id": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier of the owning organization. **Note**: One of `organization` or `organization_id` must be specified.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("organization"),
					}...),
				},
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
func (d *organizationActionVariableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *organizationActionVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read organization action variable data source"))

	var data organizationActionVariableDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get organization name from ID if not provided
	if data.Organization.IsNull() || data.Organization.IsUnknown() {
		org, diags := getOrganizationByID(
			ctx,
			d.client,
			data.OrganizationID.ValueInt64(),
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Map response body to model
		data.Organization = types.StringValue(org.UserName)
	}

	tflog.Info(ctx, "Read organization action variable", map[string]any{
		"organization": data.Organization.ValueString(),
		"name":         data.Name.ValueString(),
	})

	// Use Forgejo client to get organization action variable
	variable, res, err := d.client.GetOrgActionVariable(
		data.Organization.ValueString(),
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
					"Action variable with organization %s and name %s not found: %s",
					data.Organization.String(),
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf(
					"Unknown error (status %d): %s",
					res.StatusCode,
					err,
				)
			}
		}
		resp.Diagnostics.AddError("Unable to read organization action variable", msg)

		return
	}

	// Map response body to model
	// name is omitted here, to maintain the user's configuration casing
	data.Data = types.StringValue(variable.Data)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewOrganizationActionVariableDataSource is a helper function to simplify the provider implementation.
func NewOrganizationActionVariableDataSource() datasource.DataSource {
	return &organizationActionVariableDataSource{}
}
