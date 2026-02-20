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
	_ datasource.DataSource              = &teamDataSource{}
	_ datasource.DataSourceWithConfigure = &teamDataSource{}
)

// teamDataSource is the data source implementation.
type teamDataSource struct {
	client *forgejo.Client
}

// teamDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Team
type teamDataSourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	OrganizationID          types.Int64  `tfsdk:"organization_id"`
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Permission              types.String `tfsdk:"permission"`
	Units                   types.Set    `tfsdk:"units"`
}

// Metadata returns the data source type name.
func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the data source.
func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo team data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the team.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the team.",
				Required:    true,
			},
			"organization_id": schema.Int64Attribute{
				Description: "ID of the owning organization.",
				Required:    true,
			},
			"can_create_org_repo": schema.BoolAttribute{
				Description: "Can create repositories?",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the team.",
				Computed:    true,
			},
			"includes_all_repositories": schema.BoolAttribute{
				Description: "Has access to all repositories?",
				Computed:    true,
			},
			"permission": schema.StringAttribute{
				Description: "Permissions within the owning organization.",
				Computed:    true,
			},
			"units": schema.SetAttribute{
				Description: "Set of units.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read team data source"))

	var data teamDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read team", map[string]any{
		"name":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	// Use Forgejo client to get team by name
	team, diags := getOrgTeamByName(
		ctx,
		d.client,
		data.OrganizationID,
		data.Name,
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Map response body to model
	data.ID = types.Int64Value(team.ID)
	data.Name = types.StringValue(team.Name)
	data.Description = types.StringValue(team.Description)
	if team.Organization != nil {
		data.OrganizationID = types.Int64Value(team.Organization.ID)
	}
	data.Permission = types.StringValue(string(team.Permission))
	data.CanCreateOrgRepo = types.BoolValue(team.CanCreateOrgRepo)
	data.IncludesAllRepositories = types.BoolValue(team.IncludesAllRepositories)
	data.Units, diags = types.SetValueFrom(ctx, types.StringType, team.Units)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewTeamDataSource is a helper function to simplify the provider implementation.
func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

// getOrgTeamByName fetches a team by its name and handles errors consistently.
func getOrgTeamByName(ctx context.Context, client *forgejo.Client, orgID types.Int64, teamName types.String) (*forgejo.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use Forgejo client to get organization by ID
	organization, diags := getOrganizationByID(
		ctx,
		client,
		orgID,
	)
	if diags.HasError() {
		return nil, diags
	}

	tflog.Info(ctx, "List teams", map[string]any{
		"name":            teamName,
		"organization_id": orgID,
	})

	// Use Forgejo client to list teams in organization
	teams, res, err := client.ListOrgTeams(organization.UserName, forgejo.ListTeamsOptions{})
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
				msg = fmt.Sprintf(
					"Organization with name '%s' not found: %s",
					organization.UserName,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		diags.AddError("Unable to list teams", msg)

		return nil, diags
	}

	// Find team by name
	var team *forgejo.Team
	for _, potentialTeam := range teams {
		if teamName.Equal(types.StringValue(potentialTeam.Name)) {
			team = potentialTeam
			break
		}
	}

	if team == nil {
		diags.AddError(
			"Unable to find team",
			fmt.Sprintf("Team with name %s not found", teamName.String()),
		)

		return nil, diags
	}

	return team, diags
}
