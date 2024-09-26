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
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

// organizationDataSource is the data source implementation.
type organizationDataSource struct {
	client *forgejo.Client
}

// organizationDataSourceModel maps the data source schema data.
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
		MarkdownDescription: "Organization data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            false,
				Required:            true,
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"website": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *organizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organizationDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using tflog package
	tflog.Trace(ctx, "get organization by name", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Use forgejo.Client to get organization by name
	o, r, err := d.client.GetOrg(data.Name.ValueString())
	if err != nil {
		// Log HTTP status
		tflog.Error(ctx, "error", map[string]interface{}{
			"status": r.Status,
		})

		msg := err.Error()
		if r.StatusCode == 404 {
			msg = fmt.Sprintf("Organization with name %s not found.", data.Name.String())
		}
		resp.Diagnostics.AddError("Unable to read Forgejo organization", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(int64(o.ID))
	data.Name = types.StringValue(o.UserName)
	data.FullName = types.StringValue(o.FullName)
	data.AvatarURL = types.StringValue(o.AvatarURL)
	data.Description = types.StringValue(o.Description)
	data.Website = types.StringValue(o.Website)
	data.Location = types.StringValue(o.Location)
	data.Visibility = types.StringValue(o.Visibility)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
