package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v4"
)

var (
	_ resource.Resource                = &teamRepositoryResource{}
	_ resource.ResourceWithConfigure   = &teamRepositoryResource{}
	_ resource.ResourceWithImportState = &teamRepositoryResource{}
)

type teamRepositoryResource struct {
	client *forgejo.Client
}

type teamRepositoryResourceModel struct {
	ID           types.String `tfsdk:"id"`
	TeamID       types.Int64  `tfsdk:"team_id"`
	RepositoryID types.Int64  `tfsdk:"repository_id"`
}

func (r *teamRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_repository"
}

func (r *teamRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns one repository to an organization team. Permissions remain managed by ` + "`forgejo_team`" + `.

The team and repository must belong to the same organization. Teams with
` + "`includes_all_repositories = true`" + ` already have access and cannot be used to create
or import an individual assignment. Refreshing an existing assignment does not change
that setting. Destroy removes only the assignment, not the team or repository, and does
not revoke access granted by other teams, collaborators, or all-repositories access.

Import uses ` + "`org/team/repo`" + ` and only reads the existing configuration.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the assignment in team_id/repository_id format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team. Changing this forces a new resource to be created.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *teamRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *teamRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create team repository resource"))

	var data teamRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, diags := getOrgTeamByID(ctx, r.client, data.TeamID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if team.IncludesAllRepositories {
		resp.Diagnostics.AddError("Unable to create team repository", "The team includes all repositories. Set includes_all_repositories = false before managing individual assignments.")

		return
	}

	repo, diags := getRepositoryByID(ctx, r.client, data.RepositoryID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if team.Organization == nil || repo.Owner == nil {
		resp.Diagnostics.AddError("Unable to create team repository", "The API response is missing the team organization or repository owner.")

		return
	}
	if team.Organization.ID != repo.Owner.ID {
		resp.Diagnostics.AddError("Unable to create team repository", "The team and repository must belong to the same organization.")

		return
	}

	resp.Diagnostics.Append(setTeamRepository(ctx, r.client, team.ID, repo.Owner.UserName, repo.Name)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue(teamRepositoryID(team.ID, repo.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *teamRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read team repository resource"))

	var data teamRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := checkTeamRepository(ctx, r.client, data.TeamID.ValueInt64(), data.RepositoryID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *teamRepositoryResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team repository resource"))
	// Both writable attributes require replacement.
}

func (r *teamRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team repository resource"))

	var data teamRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := checkTeamRepository(ctx, r.client, data.TeamID.ValueInt64(), data.RepositoryID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || !found {
		return
	}

	// Removal uses the current owner/name so repository renames do not target another repository.
	repo, diags := getRepositoryByID(ctx, r.client, data.RepositoryID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if repo.Owner == nil {
		resp.Diagnostics.AddError("Unable to delete team repository", "The API response is missing the repository owner.")

		return
	}
	resp.Diagnostics.Append(deleteTeamRepository(ctx, r.client, data.TeamID.ValueInt64(), repo.Owner.UserName, repo.Name)...)
}

func (r *teamRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import team repository resource"))

	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" || strings.TrimSpace(parts[2]) == "" {
		resp.Diagnostics.AddError("Unable to parse import identifier", "Expected import identifier with format: 'org/team/repo'.")

		return
	}
	team, diags := getOrgTeamByName(ctx, r.client, parts[0], parts[1])
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if team.IncludesAllRepositories {
		resp.Diagnostics.AddError("Unable to import team repository", "The team includes all repositories; its access cannot be imported as an individual assignment.")

		return
	}
	repo, diags := getRepositoryByName(ctx, r.client, parts[0], parts[2])
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := checkTeamRepository(ctx, r.client, team.ID, repo.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Unable to import team repository", fmt.Sprintf("Repository '%s/%s' is not assigned to team '%s'.", parts[0], parts[2], parts[1]))

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &teamRepositoryResourceModel{
		ID:           types.StringValue(teamRepositoryID(team.ID, repo.ID)),
		TeamID:       types.Int64Value(team.ID),
		RepositoryID: types.Int64Value(repo.ID),
	})...)
}

func NewTeamRepositoryResource() resource.Resource {
	return &teamRepositoryResource{}
}

func teamRepositoryID(teamID, repositoryID int64) string {
	return fmt.Sprintf("%d/%d", teamID, repositoryID)
}

func checkTeamRepository(ctx context.Context, client *forgejo.Client, teamID, repositoryID int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Info(ctx, "List team repositories", map[string]any{
		"team_id":       teamID,
		"repository_id": repositoryID,
	})

	// Forgejo may omit Link headers and cap page size, so continue until an empty page.
	for page := 1; ; page++ {
		repos, res, err := client.ListTeamRepositories(teamID, forgejo.ListTeamRepositoriesOptions{
			ListOptions: forgejo.ListOptions{Page: page, PageSize: 50},
		})
		if err != nil {
			if res != nil && res.StatusCode == 404 {
				return false, diags
			}
			var msg string
			if res == nil {
				msg = fmt.Sprintf("Unknown error with nil response: %s", err)
			} else {
				msg = fmt.Sprintf("Unable to list repositories of team %d (status %d): %s", teamID, res.StatusCode, err)
			}
			diags.AddError("Unable to read team repository", msg)

			return false, diags
		}
		if len(repos) == 0 {
			return false, diags
		}
		for _, repo := range repos {
			if repo != nil && repo.ID == repositoryID {
				return true, diags
			}
		}
	}
}

func setTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, name string) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Info(ctx, "Create team repository", map[string]any{"team_id": teamID, "owner": owner, "name": name})
	res, err := client.AddTeamRepository(teamID, owner, name)
	if err == nil {
		return diags
	}
	var msg string
	if res == nil {
		msg = fmt.Sprintf("Unknown error with nil response: %s", err)
	} else {
		msg = fmt.Sprintf("Unable to assign repository '%s/%s' to team %d (status %d): %s", owner, name, teamID, res.StatusCode, err)
	}
	diags.AddError("Unable to create team repository", msg)

	return diags
}

func deleteTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, name string) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Info(ctx, "Delete team repository", map[string]any{"team_id": teamID, "owner": owner, "name": name})
	res, err := client.RemoveTeamRepository(teamID, owner, name)
	if err == nil || (res != nil && res.StatusCode == 404) {
		return diags
	}
	var msg string
	if res == nil {
		msg = fmt.Sprintf("Unknown error with nil response: %s", err)
	} else {
		msg = fmt.Sprintf("Unable to remove repository '%s/%s' from team %d (status %d): %s", owner, name, teamID, res.StatusCode, err)
	}
	diags.AddError("Unable to delete team repository", msg)

	return diags
}
