package provider

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &repositoryBranchProtectionResource{}
	_ resource.ResourceWithConfigure   = &repositoryBranchProtectionResource{}
	_ resource.ResourceWithImportState = &repositoryBranchProtectionResource{}
)

type repositoryBranchProtectionResource struct {
	client *forgejo.Client
}

// repositoryBranchProtectionResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#CreateBranchProtectionOption
type repositoryBranchProtectionResourceModel struct {
	RepositoryId                  types.Int64  `tfsdk:"repository_id"`
	BranchName                    types.String `tfsdk:"branch_name"`
	EnablePush                    types.Bool   `tfsdk:"enable_push"`
	EnablePushWhitelist           types.Bool   `tfsdk:"enable_push_whitelist"`
	PushWhitelistUsernames        types.List   `tfsdk:"push_whitelist_usernames"`
	PushWhitelistTeams            types.List   `tfsdk:"push_whitelist_teams"`
	PushWhitelistDeployKeys       types.Bool   `tfsdk:"push_whitelist_deploy_keys"`
	EnableStatusCheck             types.Bool   `tfsdk:"enable_status_check"`
	StatusCheckContexts           types.List   `tfsdk:"status_check_contexts"`
	RequireSignedCommits          types.Bool   `tfsdk:"require_signed_commits"`
	ProtectedFilePatterns         types.String `tfsdk:"protected_file_patterns"`
	UnprotectedFilePatterns       types.String `tfsdk:"unprotected_file_patterns"`
	EnableMergeWhitelist          types.Bool   `tfsdk:"enable_merge_whitelist"`
	MergeWhitelistUsernames       types.List   `tfsdk:"merge_whitelist_usernames"`
	MergeWhitelistTeams           types.List   `tfsdk:"merge_whitelist_teams"`
	EnableApprovalsWhitelist      types.Bool   `tfsdk:"enable_approvals_whitelist"`
	ApprovalsWhitelistUsernames   types.List   `tfsdk:"approvals_whitelist_usernames"`
	ApprovalsWhitelistTeams       types.List   `tfsdk:"approvals_whitelist_teams"`
	RequiredApprovals             types.Int64  `tfsdk:"required_approvals"`
	BlockOnRejectedReviews        types.Bool   `tfsdk:"block_on_rejected_reviews"`
	BlockOnOfficialReviewRequests types.Bool   `tfsdk:"block_on_official_review_requests"`
	BlockOnOutdatedBranch         types.Bool   `tfsdk:"block_on_outdated_branch"`
	DismissStaleApprovals         types.Bool   `tfsdk:"dismiss_stale_approvals"`
}

func NewRepositoryBranchProtectionResource() resource.Resource {
	return &repositoryBranchProtectionResource{}
}

func (r *repositoryBranchProtectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_branch_protection"
}

func (r *repositoryBranchProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository branch protection resource.",
		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "The ID of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"branch_name": schema.StringAttribute{
				Description: "Name of the branch to protect (can be a pattern).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable_push": schema.BoolAttribute{
				Description: "Enable push to the branch.",
				Optional:    true,
				Computed:    true,
			},
			"enable_push_whitelist": schema.BoolAttribute{
				Description: "Enable push whitelist.",
				Optional:    true,
				Computed:    true,
			},
			"push_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames allowed to push.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"push_whitelist_teams": schema.ListAttribute{
				Description: "Teams allowed to push.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"push_whitelist_deploy_keys": schema.BoolAttribute{
				Description: "Allow deploy keys to push.",
				Optional:    true,
				Computed:    true,
			},
			"enable_status_check": schema.BoolAttribute{
				Description: "Enable status checks.",
				Optional:    true,
				Computed:    true,
			},
			"status_check_contexts": schema.ListAttribute{
				Description: "Status check contexts that must pass.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_signed_commits": schema.BoolAttribute{
				Description: "Require signed commits.",
				Optional:    true,
				Computed:    true,
			},
			"protected_file_patterns": schema.StringAttribute{
				Description: "Patterns for protected files.",
				Optional:    true,
				Computed:    true,
			},
			"unprotected_file_patterns": schema.StringAttribute{
				Description: "Patterns for unprotected files.",
				Optional:    true,
				Computed:    true,
			},
			"enable_merge_whitelist": schema.BoolAttribute{
				Description: "Enable merge whitelist.",
				Optional:    true,
				Computed:    true,
			},
			"merge_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames allowed to merge.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"merge_whitelist_teams": schema.ListAttribute{
				Description: "Teams allowed to merge.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enable_approvals_whitelist": schema.BoolAttribute{
				Description: "Enable approvals whitelist.",
				Optional:    true,
				Computed:    true,
			},
			"approvals_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames that can approve.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"approvals_whitelist_teams": schema.ListAttribute{
				Description: "Teams that can approve.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"required_approvals": schema.Int64Attribute{
				Description: "Number of required approvals.",
				Computed:    true,
				Optional:    true,
			},
			"block_on_rejected_reviews": schema.BoolAttribute{
				Description: "Block merge on rejected reviews.",
				Optional:    true,
				Computed:    true,
			},
			"block_on_official_review_requests": schema.BoolAttribute{
				Description: "Block merge on official review requests.",
				Optional:    true,
				Computed:    true,
			},
			"block_on_outdated_branch": schema.BoolAttribute{
				Description: "Block merge on outdated branch.",
				Optional:    true,
				Computed:    true,
			},
			"dismiss_stale_approvals": schema.BoolAttribute{
				Description: "Dismiss stale approvals.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *repositoryBranchProtectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *repositoryBranchProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository branch protection resource"))

	var data repositoryBranchProtectionResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryID := data.RepositoryId.ValueInt64()

	repo, diags := r.fetchRepository(ctx, repositoryID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create repository branch protection", map[string]any{
		"repository_id":    repositoryID,
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Convert model to API request
	opts := r.modelToCreateOption(ctx, &data)

	tflog.Debug(ctx, "Create branch protection request parameters", map[string]any{
		"options": opts,
	})

	// Create branch protection
	protection, res, err := r.client.CreateBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"error": err.Error(),
		})

		var msg string
		if res != nil {
			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf("Forbidden: %s", err)
			case 404:
				msg = fmt.Sprintf("Repository not found: %s", err)
			case 422:
				msg = fmt.Sprintf("Validation error: %s", err)
			case 423:
				msg = fmt.Sprintf("Repository is already archived: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		} else {
			msg = fmt.Sprintf("Failed to create branch protection: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create branch protection", msg)
		return
	}

	// Update model with response data to ensure Computed fields are correctly populated
	diags = r.mapResponseToModel(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *repositoryBranchProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository branch protection resource"))

	var data repositoryBranchProtectionResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryID := data.RepositoryId.ValueInt64()

	repo, diags := r.fetchRepository(ctx, repositoryID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read repository branch protection", map[string]any{
		"repository_id":    repositoryID,
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Get branch protection
	protection, res, err := r.client.GetBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Branch with name %s not found: %s",
				data.BranchName.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get branch", msg)

		return
	}

	// Update model with response data
	diags = r.mapResponseToModel(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *repositoryBranchProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository branch protection resource"))

	var data repositoryBranchProtectionResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryID := data.RepositoryId.ValueInt64()

	repo, diags := r.fetchRepository(ctx, repositoryID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update repository branch protection", map[string]any{
		"repository_id":    repositoryID,
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Convert model to API request
	opts := r.modelToEditOption(ctx, &data)

	tflog.Debug(ctx, "Update branch protection request parameters", map[string]any{
		"options": opts,
	})

	// Update branch protection
	protection, res, err := r.client.EditBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
		opts,
	)
	if err != nil {
		var msg string
		if res != nil {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf("Forbidden: %s", err)
			case 404:
				msg = fmt.Sprintf("Branch protection not found: %s", err)
			case 422:
				msg = fmt.Sprintf("Validation error: %s", err)
			case 423:
				msg = fmt.Sprintf("Repository is already archived: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		} else {
			msg = fmt.Sprintf("Failed to update branch protection: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update branch protection", msg)
		return
	}

	// Update model with response data
	diags = r.mapResponseToModel(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update model with response data
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *repositoryBranchProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository branch protection resource"))

	var data repositoryBranchProtectionResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryID := data.RepositoryId.ValueInt64()

	repo, diags := r.fetchRepository(ctx, repositoryID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete repository branch protection", map[string]any{
		"repository_id":    repositoryID,
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Delete branch protection
	res, err := r.client.DeleteBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
	)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			// Already deleted
			resp.Diagnostics.AddError(
				"Branch protection is already deleted",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete branch protection",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}
}

func (r *repositoryBranchProtectionResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import repository branch protection resource"))

	parts := strings.Split(request.ID, "/")
	if len(parts) != 3 {
		response.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'owner/repo/branch_name', got: %s", request.ID),
		)
		return
	}

	owner := parts[0]
	repo := parts[1]
	branchName := parts[2]

	tflog.Info(ctx, "Import repository branch protection", map[string]any{
		"owner":       owner,
		"repo":        repo,
		"branch_name": branchName,
	})

	// Get branch protection
	protection, res, err := r.client.GetBranchProtection(owner, repo, branchName)
	if err != nil {
		var msg string
		if res != nil {
			switch res.StatusCode {
			case 404:
				msg = fmt.Sprintf("Branch protection not found for %s/%s/%s: %s", owner, repo, branchName, err)
			default:
				msg = fmt.Sprintf("Error getting branch protection: %s", err)
			}
		} else {
			msg = fmt.Sprintf("Failed to get branch protection: %s", err)
		}
		response.Diagnostics.AddError("Unable to import branch protection", msg)
		return
	}

	repository, _, err := r.client.GetRepo(owner, repo)
	if err != nil {
		response.Diagnostics.AddError("Unable to get repo", err.Error())
		return
	}

	// Map response to model
	var data repositoryBranchProtectionResourceModel
	data.BranchName = types.StringValue(branchName)
	data.RepositoryId = types.Int64Value(repository.ID)

	diags := r.mapResponseToModel(protection, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = response.State.Set(ctx, &data)
	response.Diagnostics.Append(diags...)
}

// Helper function to convert model to CreateBranchProtectionOption.
func (r *repositoryBranchProtectionResource) modelToCreateOption(ctx context.Context, data *repositoryBranchProtectionResourceModel) forgejo.CreateBranchProtectionOption {
	opts := forgejo.CreateBranchProtectionOption{
		BranchName: data.BranchName.ValueString(),
	}

	opts.EnablePush = data.EnablePush.ValueBool()
	opts.EnablePushWhitelist = data.EnablePushWhitelist.ValueBool()

	var usernames []string
	data.PushWhitelistUsernames.ElementsAs(ctx, &usernames, false)
	opts.PushWhitelistUsernames = usernames

	var teams []string
	data.PushWhitelistTeams.ElementsAs(ctx, &teams, false)
	opts.PushWhitelistTeams = teams

	opts.PushWhitelistDeployKeys = data.PushWhitelistDeployKeys.ValueBool()
	opts.EnableStatusCheck = data.EnableStatusCheck.ValueBool()

	var contexts []string
	data.StatusCheckContexts.ElementsAs(ctx, &contexts, false)
	opts.StatusCheckContexts = contexts

	opts.RequireSignedCommits = data.RequireSignedCommits.ValueBool()
	opts.ProtectedFilePatterns = data.ProtectedFilePatterns.ValueString()
	opts.UnprotectedFilePatterns = data.UnprotectedFilePatterns.ValueString()
	opts.EnableMergeWhitelist = data.EnableMergeWhitelist.ValueBool()

	var mergeWhitelistUsernames []string
	data.MergeWhitelistUsernames.ElementsAs(ctx, &mergeWhitelistUsernames, false)
	opts.MergeWhitelistUsernames = mergeWhitelistUsernames

	var mergeWhitelistTeams []string
	data.MergeWhitelistTeams.ElementsAs(ctx, &mergeWhitelistTeams, false)
	opts.MergeWhitelistTeams = mergeWhitelistTeams

	opts.EnableApprovalsWhitelist = data.EnableApprovalsWhitelist.ValueBool()

	var approvalsWhitelistUsernames []string
	data.ApprovalsWhitelistUsernames.ElementsAs(ctx, &approvalsWhitelistUsernames, false)
	opts.ApprovalsWhitelistUsernames = approvalsWhitelistUsernames

	var approvalsWhitelistTeams []string
	data.ApprovalsWhitelistTeams.ElementsAs(ctx, &approvalsWhitelistTeams, false)
	opts.ApprovalsWhitelistTeams = approvalsWhitelistTeams

	opts.RequiredApprovals = data.RequiredApprovals.ValueInt64()
	opts.BlockOnRejectedReviews = data.BlockOnRejectedReviews.ValueBool()
	opts.BlockOnOfficialReviewRequests = data.BlockOnOfficialReviewRequests.ValueBool()
	opts.BlockOnOutdatedBranch = data.BlockOnOutdatedBranch.ValueBool()
	opts.DismissStaleApprovals = data.DismissStaleApprovals.ValueBool()

	return opts
}

func (r *repositoryBranchProtectionResource) fetchRepository(ctx context.Context, id int64) (*forgejo.Repository, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(id)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				id,
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		diags.AddError("Unable to get repository by id", msg)

		return nil, diags
	}
	return rep, nil
}

// Helper function to convert model to EditBranchProtectionOption.
func (r *repositoryBranchProtectionResource) modelToEditOption(ctx context.Context, data *repositoryBranchProtectionResourceModel) forgejo.EditBranchProtectionOption {
	opts := forgejo.EditBranchProtectionOption{}

	opts.EnablePush = data.EnablePush.ValueBoolPointer()

	opts.EnablePushWhitelist = data.EnablePushWhitelist.ValueBoolPointer()

	var usernames []string
	data.PushWhitelistUsernames.ElementsAs(ctx, &usernames, false)
	opts.PushWhitelistUsernames = usernames

	var teams []string
	data.PushWhitelistTeams.ElementsAs(ctx, &teams, false)
	opts.PushWhitelistTeams = teams

	opts.PushWhitelistDeployKeys = data.PushWhitelistDeployKeys.ValueBoolPointer()

	opts.EnableStatusCheck = data.EnableStatusCheck.ValueBoolPointer()

	var contexts []string
	data.StatusCheckContexts.ElementsAs(ctx, &contexts, false)
	opts.StatusCheckContexts = contexts

	opts.RequireSignedCommits = data.RequireSignedCommits.ValueBoolPointer()

	patterns := data.ProtectedFilePatterns.ValueString()
	opts.ProtectedFilePatterns = &patterns

	unprotectedPatterns := data.UnprotectedFilePatterns.ValueString()
	opts.UnprotectedFilePatterns = &unprotectedPatterns

	opts.EnableMergeWhitelist = data.EnableMergeWhitelist.ValueBoolPointer()

	var mergeWhitelistUsernames []string
	data.MergeWhitelistUsernames.ElementsAs(ctx, &mergeWhitelistUsernames, false)
	opts.MergeWhitelistUsernames = mergeWhitelistUsernames

	var mergeWhitelistTeams []string
	data.MergeWhitelistTeams.ElementsAs(ctx, &mergeWhitelistTeams, false)
	opts.MergeWhitelistTeams = mergeWhitelistTeams

	opts.EnableApprovalsWhitelist = data.EnableApprovalsWhitelist.ValueBoolPointer()

	var approvalsWhitelistUsernames []string
	data.ApprovalsWhitelistUsernames.ElementsAs(ctx, &approvalsWhitelistUsernames, false)
	opts.ApprovalsWhitelistUsernames = approvalsWhitelistUsernames

	var approvalsWhitelistTeams []string
	data.ApprovalsWhitelistTeams.ElementsAs(ctx, &approvalsWhitelistTeams, false)
	opts.ApprovalsWhitelistTeams = approvalsWhitelistTeams

	approvals := data.RequiredApprovals.ValueInt64()
	opts.RequiredApprovals = &approvals

	opts.BlockOnRejectedReviews = data.BlockOnRejectedReviews.ValueBoolPointer()

	opts.BlockOnOfficialReviewRequests = data.BlockOnOfficialReviewRequests.ValueBoolPointer()

	opts.BlockOnOutdatedBranch = data.BlockOnOutdatedBranch.ValueBoolPointer()

	opts.DismissStaleApprovals = data.DismissStaleApprovals.ValueBoolPointer()

	return opts
}

func (r *repositoryBranchProtectionResource) mapResponseToModel(protection *forgejo.BranchProtection, data *repositoryBranchProtectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.EnablePush = types.BoolValue(protection.EnablePush)
	data.EnablePushWhitelist = types.BoolValue(protection.EnablePushWhitelist)
	data.PushWhitelistDeployKeys = types.BoolValue(protection.PushWhitelistDeployKeys)
	data.EnableStatusCheck = types.BoolValue(protection.EnableStatusCheck)
	data.RequireSignedCommits = types.BoolValue(protection.RequireSignedCommits)

	if protection.ProtectedFilePatterns != "" {
		data.ProtectedFilePatterns = types.StringValue(protection.ProtectedFilePatterns)
	} else {
		data.ProtectedFilePatterns = types.StringNull()
	}

	if protection.UnprotectedFilePatterns != "" {
		data.UnprotectedFilePatterns = types.StringValue(protection.UnprotectedFilePatterns)
	} else {
		data.UnprotectedFilePatterns = types.StringNull()
	}

	data.EnableMergeWhitelist = types.BoolValue(protection.EnableMergeWhitelist)
	data.EnableApprovalsWhitelist = types.BoolValue(protection.EnableApprovalsWhitelist)

	if protection.RequiredApprovals != 0 {
		data.RequiredApprovals = types.Int64Value(protection.RequiredApprovals)
	} else {
		data.RequiredApprovals = types.Int64Null()
	}

	data.BlockOnRejectedReviews = types.BoolValue(protection.BlockOnRejectedReviews)
	data.BlockOnOfficialReviewRequests = types.BoolValue(protection.BlockOnOfficialReviewRequests)
	data.BlockOnOutdatedBranch = types.BoolValue(protection.BlockOnOutdatedBranch)
	data.DismissStaleApprovals = types.BoolValue(protection.DismissStaleApprovals)

	// Handle Lists
	var d diag.Diagnostics
	data.PushWhitelistUsernames, d = r.stringSliceToList(protection.PushWhitelistUsernames)
	diags.Append(d...)

	data.PushWhitelistTeams, d = r.stringSliceToList(protection.PushWhitelistTeams)
	diags.Append(d...)

	data.StatusCheckContexts, d = r.stringSliceToList(protection.StatusCheckContexts)
	diags.Append(d...)

	data.MergeWhitelistUsernames, d = r.stringSliceToList(protection.MergeWhitelistUsernames)
	diags.Append(d...)

	data.MergeWhitelistTeams, d = r.stringSliceToList(protection.MergeWhitelistTeams)
	diags.Append(d...)

	data.ApprovalsWhitelistUsernames, d = r.stringSliceToList(protection.ApprovalsWhitelistUsernames)
	diags.Append(d...)

	data.ApprovalsWhitelistTeams, d = r.stringSliceToList(protection.ApprovalsWhitelistTeams)
	diags.Append(d...)

	return diags
}

func (r *repositoryBranchProtectionResource) stringSliceToList(slice []string) (types.List, diag.Diagnostics) {
	if len(slice) == 0 {
		return types.ListNull(types.StringType), nil
	}
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
		elements[i] = types.StringValue(v)
	}
	return types.ListValue(types.StringType, elements)
}
