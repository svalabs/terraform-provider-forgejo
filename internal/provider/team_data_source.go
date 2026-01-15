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

	// Read Terraform configuration data into model.
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get team by name", map[string]any{
		"name":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	// Use Forgejo client to get team by name.
	team, err := getOrgTeamByName(ctx, d.client, data.OrganizationID, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get team by name", err.Error())
		return
	}

	if team == nil {
		err = fmt.Errorf(
			"no Team with name '%s' was found",
			data.Name.String(),
		)
		resp.Diagnostics.AddError("Unable to get team by name", err.Error())
		return
	}

	// Map response body to model.
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

	// Save data into Terraform state.
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewTeamDataSource is a helper function to simplify the provider implementation.
func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

// Use Forgejo client to get team by name.
func getOrgTeamByName(ctx context.Context, client *forgejo.Client, orgID types.Int64, teamName types.String) (team *forgejo.Team, err error) {
	tflog.Info(ctx, "Getting team in org", map[string]any{
		"team":            teamName,
		"organization_id": orgID,
	})

	organization, err := getOrganizationByID(ctx, client, orgID)
	if err != nil {
		return nil, err
	}

	if organization == nil {
		err = fmt.Errorf(
			"no Organization with id '%d' was found",
			orgID.ValueInt64(),
		)
		return nil, err
	}

	teams, resp, err := client.ListOrgTeams(organization.UserName, forgejo.ListTeamsOptions{})
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": resp.Status,
		})

		switch resp.StatusCode {
		case 404:
			err = fmt.Errorf(
				"the Organization with name '%s' was not found: %s",
				organization.UserName,
				err,
			)
		default:
			err = fmt.Errorf("unknown error: %s", err)
		}
		return nil, err
	}

	for _, potentialTeam := range teams {
		if teamName.Equal(types.StringValue(potentialTeam.Name)) {
			team = potentialTeam
			break
		}
	}

	return team, nil
}
