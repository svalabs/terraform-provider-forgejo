package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &teamMemberDataSource{}
	_ datasource.DataSourceWithConfigure = &teamMemberDataSource{}
)

// teamMemberDataSource is the data source implementation.
type teamMemberDataSource struct {
	client *forgejo.Client
}

// teamMemberDataSourceModel maps the data source schema data.
type teamMemberDataSourceModel struct {
	TeamID types.Int64  `tfsdk:"team_id"`
	User   types.String `tfsdk:"user"`
}

// Metadata returns the data source type name.
func (d *teamMemberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

// Schema defines the schema for the data source.
func (d *teamMemberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo team member data source.",

		Attributes: map[string]schema.Attribute{
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team.",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "Username of the team member.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *teamMemberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *teamMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read team member data source"))

	var data teamMemberDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to check a team member
	diags = checkTeamMember(
		ctx,
		d.client,
		data.TeamID.ValueInt64(),
		data.User.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewTeamMemberDataSource is a helper function to simplify the provider implementation.
func NewTeamMemberDataSource() datasource.DataSource {
	return &teamMemberDataSource{}
}

// checkTeamMember fetches a team member and handles errors consistently.
func checkTeamMember(ctx context.Context, client *forgejo.Client, teamID int64, userName string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Read team member", map[string]any{
		"team_id": teamID,
		"user":    userName,
	})

	// Use Forgejo client to get team member
	_, res, err := client.GetTeamMember(teamID, userName)
	if err == nil {
		return diags
	}

	// Handle errors
	var msg string
	if res == nil {
		msg = fmt.Sprintf("Unknown error with nil response: %s", err)
	} else {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"User '%s' in team with ID %d not found: %s",
				userName,
				teamID,
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
	diags.AddError("Unable to read team member", msg)

	return diags
}
