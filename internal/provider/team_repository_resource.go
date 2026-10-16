package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &teamRepositoryResource{}
	_ resource.ResourceWithConfigure   = &teamRepositoryResource{}
	_ resource.ResourceWithImportState = &teamRepositoryResource{}
)

// teamRepositoryResource is the resource implementation.
type teamRepositoryResource struct {
	client *forgejo.Client
}

// teamRepositoryResourceModel maps the resource schema data.
// The association between a team and a repository has no counterpart in the
// SDK, so there are no from()/to() helpers. The composite 'id' attribute is
// kept because neither 'team_id' nor 'repository_id' alone identifies the
// association, and acceptance tests need a unique identifier for
// 'ImportStateVerify'.
type teamRepositoryResourceModel struct {
	ID           types.String `tfsdk:"id"`
	TeamID       types.Int64  `tfsdk:"team_id"`
	RepositoryID types.Int64  `tfsdk:"repository_id"`
}

// Metadata returns the resource type name.
func (r *teamRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_repository"
}

// Schema defines the schema for the resource.
func (r *teamRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo team repository resource. Assigns a repository to an organization team.

**Note**: The team and repository must belong to the same organization. Teams with ` + "`includes_all_repositories = true`" + ` already have access to every repository of the organization and can not be managed by this resource. A team that is switched to ` + "`includes_all_repositories = true`" + ` afterwards keeps the explicit association. Destroying this resource removes the association without deleting the team or repository.`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the association in team_id/repository_id format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *teamRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *teamRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create team repository resource"))

	var data teamRepositoryResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get team
	team, diags := getOrgTeamByID(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Teams with access to all repositories can not hold explicit associations
	if team.IncludesAllRepositories {
		resp.Diagnostics.AddError(
			"Unable to create team repository",
			fmt.Sprintf(
				"Team with ID %d includes all repositories. Set includes_all_repositories = false on the team before managing explicit repository associations.",
				team.ID,
			),
		)

		return
	}

	// Use Forgejo client to get repository
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure team and repository belong to the same organization
	if team.Organization == nil {
		resp.Diagnostics.AddError(
			"Unable to create team repository",
			fmt.Sprintf(
				"Team with ID %d has no organization in the API response. Please report this issue to the provider developers.",
				team.ID,
			),
		)

		return
	}
	if repo.Owner == nil {
		resp.Diagnostics.AddError(
			"Unable to create team repository",
			fmt.Sprintf(
				"Repository with ID %d has no owner in the API response. Please report this issue to the provider developers.",
				repo.ID,
			),
		)

		return
	}
	if team.Organization.ID != repo.Owner.ID {
		resp.Diagnostics.AddError(
			"Unable to create team repository",
			fmt.Sprintf(
				"Team with ID %d belongs to organization '%s', but repository with ID %d is owned by '%s'. The team and repository must belong to the same organization.",
				team.ID,
				team.Organization.UserName,
				repo.ID,
				repo.Owner.UserName,
			),
		)

		return
	}

	// Use Forgejo client to set a team repository
	diags = setTeamRepository(
		ctx,
		r.client,
		team.ID,
		repo.Owner.UserName,
		repo.Name,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.ID = types.StringValue(teamRepositoryID(team.ID, repo.ID))

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *teamRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read team repository resource"))

	var data teamRepositoryResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to check a team repository
	found, diags := checkTeamRepository(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove the resource from state if the association no longer exists
	if !found {
		tflog.Warn(ctx, "Team repository not found, removing from state", map[string]any{
			"team_id":       data.TeamID.ValueInt64(),
			"repository_id": data.RepositoryID.ValueInt64(),
		})
		resp.State.RemoveResource(ctx)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team repository resource"))

	/*
	 * Team repositories can not be updated in-place. All writable attributes
	 * have 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team repository resource"))

	var data teamRepositoryResourceModel

	// Read Terraform prior state into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository (team repositories are removed by name)
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if repo.Owner == nil {
		resp.Diagnostics.AddError(
			"Unable to delete team repository",
			fmt.Sprintf(
				"Repository with ID %d has no owner in the API response. Please report this issue to the provider developers.",
				repo.ID,
			),
		)

		return
	}

	// Use Forgejo client to delete a team repository
	diags = deleteTeamRepository(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		repo.Owner.UserName,
		repo.Name,
	)
	resp.Diagnostics.Append(diags...)
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *teamRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import team repository resource"))

	var state teamRepositoryResourceModel

	// Parse import identifier
	cmp := strings.Split(req.ID, "/")
	if len(cmp) != 3 {
		resp.Diagnostics.AddError(
			"Unable to parse import identifier",
			fmt.Sprintf(
				"Expected import identifier with format: 'org/team/repo', got: '%s'",
				req.ID,
			),
		)

		return
	}
	orgName, teamName, repoName := cmp[0], cmp[1], cmp[2]

	// Use Forgejo client to get team
	team, diags := getOrgTeamByName(
		ctx,
		r.client,
		orgName,
		teamName,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository (the owner is the organization)
	repo, diags := getRepositoryByName(
		ctx,
		r.client,
		orgName,
		repoName,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to check a team repository
	found, diags := checkTeamRepository(
		ctx,
		r.client,
		team.ID,
		repo.ID,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError(
			"Unable to read team repository",
			fmt.Sprintf(
				"Team repository with organization '%s', team '%s' and repository '%s' not found",
				orgName,
				teamName,
				repoName,
			),
		)

		return
	}

	// Map response body to model
	state.ID = types.StringValue(teamRepositoryID(team.ID, repo.ID))
	state.TeamID = types.Int64Value(team.ID)
	state.RepositoryID = types.Int64Value(repo.ID)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// NewTeamRepositoryResource is a helper function to simplify the provider implementation.
func NewTeamRepositoryResource() resource.Resource {
	return &teamRepositoryResource{}
}

// teamRepositoryID is a helper function to build the composite identifier of
// a team repository association.
func teamRepositoryID(teamID, repositoryID int64) string {
	return fmt.Sprintf("%d/%d", teamID, repositoryID)
}

// teamRepositoryPageSize is the page size requested when listing team
// repositories. It matches Forgejo's default MAX_RESPONSE_ITEMS; the server
// may still cap the actual page size.
const teamRepositoryPageSize = 50

// checkTeamRepository is a helper function to check whether a repository is
// assigned to a team and handles errors consistently. A missing team is
// reported as a missing association rather than an error.
func checkTeamRepository(ctx context.Context, client *forgejo.Client, teamID, repositoryID int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "List team repositories", map[string]any{
		"team_id":       teamID,
		"repository_id": repositoryID,
	})

	/*
	 * Forgejo's 'GET /teams/{id}/repos' endpoint does not set a 'Link'
	 * header, so the SDK can not report the next page. Request pages until
	 * one comes back empty; the server may cap the page size, so a short
	 * page does not mean it was the last one.
	 */
	for page := 1; ; page++ {
		// Use Forgejo client to list team repositories
		repos, res, err := client.ListTeamRepositories(
			teamID,
			forgejo.ListTeamRepositoriesOptions{
				ListOptions: forgejo.ListOptions{
					Page:     page,
					PageSize: teamRepositoryPageSize,
				},
			},
		)
		if err != nil {
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
					tflog.Warn(ctx, "Team not found", map[string]any{
						"team_id": teamID,
					})

					return false, diags
				case 403:
					msg = fmt.Sprintf(
						"Access to repositories of team with ID %d forbidden: %s",
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
			diags.AddError("Unable to read team repository", msg)

			return false, diags
		}

		// An empty page means all repositories have been listed
		if len(repos) == 0 {
			return false, diags
		}

		// Search for repository with given ID
		for _, repo := range repos {
			if repo != nil && repo.ID == repositoryID {
				return true, diags
			}
		}
	}
}

// setTeamRepository is a helper function to add a repository to a team.
func setTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, name string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Create team repository", map[string]any{
		"team_id": teamID,
		"owner":   owner,
		"name":    name,
	})

	// Use Forgejo client to set a team repository
	res, err := client.AddTeamRepository(teamID, owner, name)
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
		case 403:
			msg = fmt.Sprintf(
				"Administrative access to repository '%s/%s' required: %s",
				owner,
				name,
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Either repository '%s/%s' or team with ID %d not found: %s",
				owner,
				name,
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
	diags.AddError("Unable to create team repository", msg)

	return diags
}

// deleteTeamRepository is a helper function to remove a repository from a team.
func deleteTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, name string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Delete team repository", map[string]any{
		"team_id": teamID,
		"owner":   owner,
		"name":    name,
	})

	// Use Forgejo client to delete a team repository
	res, err := client.RemoveTeamRepository(teamID, owner, name)
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
		case 403:
			msg = fmt.Sprintf(
				"Administrative access to repository '%s/%s' required: %s",
				owner,
				name,
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Either repository '%s/%s' or team with ID %d not found: %s",
				owner,
				name,
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
	diags.AddError("Unable to delete team repository", msg)

	return diags
}
