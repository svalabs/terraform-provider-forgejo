package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
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
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#Team
type teamDataSourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Organization            types.String `tfsdk:"organization"`
	OrganizationID          types.Int64  `tfsdk:"organization_id"`
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Permission              types.String `tfsdk:"permission"`
	UnitsMap                types.Map    `tfsdk:"units_map"`
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
			"units_map": schema.MapAttribute{
				Description: "Map of access units.",
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

	tflog.Info(ctx, "Read team", map[string]any{
		"organization": data.Organization.ValueString(),
		"name":         data.Name.ValueString(),
	})

	// Use Forgejo client to get team
	team, diags := getOrgTeamByName(
		ctx,
		d.client,
		data.Organization.ValueString(),
		data.Name.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.ID = types.Int64Value(team.ID)
	data.Name = types.StringValue(team.Name)
	data.Description = types.StringValue(team.Description)
	if team.Organization != nil {
		data.Organization = types.StringValue(team.Organization.UserName)
		data.OrganizationID = types.Int64Value(team.Organization.ID)
	}
	data.Permission = types.StringValue(string(team.Permission))
	data.CanCreateOrgRepo = types.BoolValue(team.CanCreateOrgRepo)
	data.IncludesAllRepositories = types.BoolValue(team.IncludesAllRepositories)
	data.UnitsMap, diags = types.MapValueFrom(ctx, types.StringType, team.UnitsMap)
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

// getOrgTeamByID fetches a team by its ID and handles errors consistently.
func getOrgTeamByID(ctx context.Context, client *forgejo.Client, id int64) (*forgejo.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Read team", map[string]any{
		"id": id,
	})

	// Use Forgejo client to get team
	team, res, err := client.GetTeam(id)
	if err == nil {
		return team, diags
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
				"Team with ID %d not found: %s",
				id,
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
	diags.AddError("Unable to read team", msg)

	return nil, diags
}

// getOrgTeamByName fetches a team by its name and handles errors consistently.
func getOrgTeamByName(ctx context.Context, client *forgejo.Client, org, name string) (*forgejo.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "List teams", map[string]any{
		"organization": org,
		"name":         name,
	})

	// Use Forgejo client to list teams in organization
	teams, res, err := client.ListOrgTeams(
		org,
		forgejo.ListTeamsOptions{
			ListOptions: forgejo.ListOptions{
				Page: -1,
			},
		},
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
			case 404:
				msg = fmt.Sprintf(
					"Organization with name '%s' not found: %s",
					org,
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
		diags.AddError("Unable to list teams", msg)

		return nil, diags
	}

	// Search for team with given name
	idx := slices.IndexFunc(teams, func(t *forgejo.Team) bool {
		return t.Name == name
	})
	if idx == -1 {
		diags.AddError(
			"Unable to find team by name",
			fmt.Sprintf("Team with name '%s' not found", name),
		)

		return nil, diags
	}

	return teams[idx], diags
}
