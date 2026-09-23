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
type teamRepositoryResourceModel struct {
	TeamID     types.Int64  `tfsdk:"team_id"`
	Owner      types.String `tfsdk:"owner"`
	Repository types.String `tfsdk:"repository"`
}

// Metadata returns the resource type name.
func (r *teamRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_repository"
}

// Schema defines the schema for the resource.
func (r *teamRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo team repository resource.

**Note**: The authenticated user must be a member of the managed organization(s) or have administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "Owner (organization) of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository": schema.StringAttribute{
				Description: "Name of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

	// Use Forgejo client to add a repository to a team
	diags = addTeamRepository(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.Owner.ValueString(),
		data.Repository.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	diags = checkTeamRepository(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.Owner.ValueString(),
		data.Repository.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
	 * Team repositories can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team repository resource"))

	var data teamRepositoryResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to remove a repository from a team
	diags = removeTeamRepository(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.Owner.ValueString(),
		data.Repository.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *teamRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import team repository resource"))

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

	// Use Forgejo client to get team by name
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

	// Verify the repository belongs to the team
	diags = checkTeamRepository(
		ctx,
		r.client,
		team.ID,
		orgName,
		repoName,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &teamRepositoryResourceModel{
		TeamID:     types.Int64Value(team.ID),
		Owner:      types.StringValue(orgName),
		Repository: types.StringValue(repoName),
	})
	resp.Diagnostics.Append(diags...)
}

// NewTeamRepositoryResource is a helper function to simplify the provider implementation.
func NewTeamRepositoryResource() resource.Resource {
	return &teamRepositoryResource{}
}

// checkTeamRepository fetches a team repository and handles errors consistently.
func checkTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, repo string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Read team repository", map[string]any{
		"team_id": teamID,
		"owner":   owner,
		"repo":    repo,
	})

	// Use Forgejo client to check a team repository
	repos, res, err := client.ListTeamRepositories(teamID, forgejo.ListTeamRepositoriesOptions{
		ListOptions: forgejo.ListOptions{
			Page: -1,
		},
	})
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			msg = fmt.Sprintf(
				"Unknown error (status %d): %s",
				res.StatusCode,
				err,
			)
		}
		diags.AddError("Unable to list team repositories", msg)

		return diags
	}

	// Search for repository with given name
	for _, r := range repos {
		if r.Name == repo {
			return diags
		}
	}

	diags.AddError(
		"Unable to find team repository",
		fmt.Sprintf(
			"Repository '%s/%s' not found in team with ID %d",
			owner,
			repo,
			teamID,
		),
	)

	return diags
}

// addTeamRepository is a helper function to add a repository to a team.
func addTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, repo string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Create team repository", map[string]any{
		"team_id": teamID,
		"owner":   owner,
		"repo":    repo,
	})

	// Use Forgejo client to add a repository to a team
	res, err := client.AddTeamRepository(teamID, owner, repo)
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
				"Team with ID %d does not have permission to access '%s/%s': %s",
				teamID,
				owner,
				repo,
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Either repository '%s/%s' or team with ID %d not found: %s",
				owner,
				repo,
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

// removeTeamRepository is a helper function to remove a repository from a team.
func removeTeamRepository(ctx context.Context, client *forgejo.Client, teamID int64, owner, repo string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Delete team repository", map[string]any{
		"team_id": teamID,
		"owner":   owner,
		"repo":    repo,
	})

	// Use Forgejo client to remove a repository from a team
	res, err := client.RemoveTeamRepository(teamID, owner, repo)
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
				"Either repository '%s/%s' or team with ID %d not found: %s",
				owner,
				repo,
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
