package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &teamMembershipDataSource{}
	_ datasource.DataSourceWithConfigure = &teamMembershipDataSource{}
)

// teamMembershipDataSource is the data source implementation.
type teamMembershipDataSource struct {
	client *forgejo.Client
}

// teamMembershipDataSourceModel maps the data source schema data.
type teamMembershipDataSourceModel struct {
	TeamID   types.Int64  `tfsdk:"team_id"`
	UserID   types.Int64  `tfsdk:"user_id"`
	UserName types.String `tfsdk:"user_name"`
}

// Metadata returns the data source type name.
func (d *teamMembershipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_membership"
}

// Schema defines the schema for the data source.
func (d *teamMembershipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo team membership data source.",

		Attributes: map[string]schema.Attribute{
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team.",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Numeric identifier of the user.",
				Required:    true,
			},
			"user_name": schema.StringAttribute{
				Description: "Name of the user.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *teamMembershipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *teamMembershipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read team membership data source"))

	var data teamMembershipDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, diags := getUserByID(ctx, d.client, data.UserID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, "Read user's team membership", map[string]any{
		"team_id":   data.TeamID.ValueInt64(),
		"user_id":   data.UserID.ValueInt64(),
		"user_name": user.UserName,
	})

	user, diags = getTeamMembership(ctx, d.client,
		data.TeamID.ValueInt64(),
		user.UserName,
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.UserName = types.StringValue(user.UserName)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewTeamMembershipDataSource is a helper function to simplify the provider implementation.
func NewTeamMembershipDataSource() datasource.DataSource {
	return &teamMembershipDataSource{}
}

// getTeamMembership fetches a team membership and handles errors consistently.
func getTeamMembership(ctx context.Context, client *forgejo.Client, teamID int64, userName string) (*forgejo.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Get team membership", map[string]any{
		"team_id":   teamID,
		"user_name": userName,
	})

	// Use Forgejo client to get team membership
	user, res, err := client.GetTeamMember(teamID, userName)
	if err == nil {
		return user, diags
	}

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
				"User %s is not a member of team with ID %d: %s",
				userName,
				teamID,
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
	}
	diags.AddError("Unable to get team membership", msg)
	return nil, diags
}
