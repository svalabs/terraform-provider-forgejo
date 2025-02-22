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
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

// organizationDataSource is the data source implementation.
type organizationDataSource struct {
	client *forgejo.Client
}

// organizationDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Organization
type organizationDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	FullName    types.String `tfsdk:"full_name"`
	AvatarURL   types.String `tfsdk:"avatar_url"`
	Description types.String `tfsdk:"description"`
	Website     types.String `tfsdk:"website"`
	Location    types.String `tfsdk:"location"`
	Visibility  types.String `tfsdk:"visibility"`
}

// Metadata returns the data source type name.
func (d *organizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the data source.
func (d *organizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo organization data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the organization.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the organization.",
				Required:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the organization.",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the organization.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the organization.",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website of the organization.",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "Location of the organization.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Visibility of the organization.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *organizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read organization data source"))

	var data organizationDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get organization by name", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to get organization by name
	org, res, err := d.client.GetOrg(data.Name.ValueString())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Organization with name %s not found: %s",
				data.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get organization by name", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(org.ID)
	data.Name = types.StringValue(org.UserName)
	data.FullName = types.StringValue(org.FullName)
	data.AvatarURL = types.StringValue(org.AvatarURL)
	data.Description = types.StringValue(org.Description)
	data.Website = types.StringValue(org.Website)
	data.Location = types.StringValue(org.Location)
	data.Visibility = types.StringValue(org.Visibility)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}
