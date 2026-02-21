package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &branchProtectionResource{}
	_ resource.ResourceWithConfigure   = &branchProtectionResource{}
	_ resource.ResourceWithImportState = &branchProtectionResource{}
)

// branchProtectionResource is the resource implementation.
type branchProtectionResource struct {
	client *forgejo.Client
}

// branchProtectionResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#CreateBranchProtectionOption
type branchProtectionResourceModel struct {
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

// Metadata returns the resource type name.
func (r *branchProtectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_protection"
}

// Schema defines the schema for the resource.
func (r *branchProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo branch protection resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
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
func (r *branchProtectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *branchProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create branch protection resource"))

	var data branchProtectionResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryId.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	opts := r.toCreateOption(ctx, &data)

	tflog.Info(ctx, "Create branch protection", map[string]any{
		"repository_id":                     data.RepositoryId.ValueInt64(),
		"repository_name":                   repo.Name,
		"repository_owner":                  repo.Owner.UserName,
		"branch_name":                       data.BranchName.ValueString(),
		"enable_push":                       opts.EnablePush,
		"enable_push_whitelist":             opts.EnablePushWhitelist,
		"push_whitelist_deploy_keys":        opts.PushWhitelistDeployKeys,
		"enable_status_check":               opts.EnableStatusCheck,
		"require_signed_commits":            opts.RequireSignedCommits,
		"enable_merge_whitelist":            opts.EnableMergeWhitelist,
		"enable_approvals_whitelist":        opts.EnableApprovalsWhitelist,
		"required_approvals":                opts.RequiredApprovals,
		"block_on_rejected_reviews":         opts.BlockOnRejectedReviews,
		"block_on_official_review_requests": opts.BlockOnOfficialReviewRequests,
		"block_on_outdated_branch":          opts.BlockOnOutdatedBranch,
		"dismiss_stale_approvals":           opts.DismissStaleApprovals,
	})

	// Use Forgejo client to create branch protection
	protection, res, err := r.client.CreateBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		opts,
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"error": err.Error(),
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf("Repository forbidden: %s", err)
			case 404:
				msg = fmt.Sprintf("Repository not found: %s", err)
			case 422:
				msg = fmt.Sprintf("Validation error: %s", err)
			case 423:
				msg = fmt.Sprintf("Repository is already archived: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create branch protection", msg)

		return
	}

	// Update model with response data to ensure computed fields are correctly populated
	diags = r.from(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *branchProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read branch protection resource"))

	var data branchProtectionResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryId.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read branch protection", map[string]any{
		"repository_id":    data.RepositoryId.ValueInt64(),
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Use Forgejo client to get branch protection
	protection, res, err := r.client.GetBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
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
					"Branch with name %s not found: %s",
					data.BranchName.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read branch protection", msg)

		return
	}

	// Update model with response data
	diags = r.from(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *branchProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update branch protection resource"))

	var data branchProtectionResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryId.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	opts := r.toEditOption(ctx, &data)

	tflog.Info(ctx, "Update branch protection", map[string]any{
		"repository_id":                     data.RepositoryId.ValueInt64(),
		"repository_name":                   repo.Name,
		"repository_owner":                  repo.Owner.UserName,
		"branch_name":                       data.BranchName.ValueString(),
		"enable_push":                       opts.EnablePush,
		"enable_push_whitelist":             opts.EnablePushWhitelist,
		"push_whitelist_deploy_keys":        opts.PushWhitelistDeployKeys,
		"enable_status_check":               opts.EnableStatusCheck,
		"require_signed_commits":            opts.RequireSignedCommits,
		"enable_merge_whitelist":            opts.EnableMergeWhitelist,
		"enable_approvals_whitelist":        opts.EnableApprovalsWhitelist,
		"required_approvals":                opts.RequiredApprovals,
		"block_on_rejected_reviews":         opts.BlockOnRejectedReviews,
		"block_on_official_review_requests": opts.BlockOnOfficialReviewRequests,
		"block_on_outdated_branch":          opts.BlockOnOutdatedBranch,
		"dismiss_stale_approvals":           opts.DismissStaleApprovals,
	})

	// Use Forgejo client to update branch protection
	protection, res, err := r.client.EditBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
		opts,
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
			case 403:
				msg = fmt.Sprintf("Repository forbidden: %s", err)
			case 404:
				msg = fmt.Sprintf("Repository not found: %s", err)
			case 422:
				msg = fmt.Sprintf("Validation error: %s", err)
			case 423:
				msg = fmt.Sprintf("Repository is already archived: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update branch protection", msg)

		return
	}

	// Update model with response data
	diags = r.from(protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *branchProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete branch protection resource"))

	var data branchProtectionResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	repo, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryId.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete branch protection", map[string]any{
		"repository_id":    data.RepositoryId.ValueInt64(),
		"repository_name":  repo.Name,
		"repository_owner": repo.Owner.UserName,
		"branch_name":      data.BranchName.ValueString(),
	})

	// Use Forgejo client to delete branch protection
	res, err := r.client.DeleteBranchProtection(
		repo.Owner.UserName,
		repo.Name,
		data.BranchName.ValueString(),
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
				msg = fmt.Sprintf("Branch protection not found: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete branch protection", msg)

		return
	}
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *branchProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, response *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import branch protection resource"))

	var data branchProtectionResourceModel

	// Parse import identifier
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		response.Diagnostics.AddError(req.ID, "Import ID must be in format 'owner/repo/branch'")

		return
	}
	owner, repo, branchName := parts[0], parts[1], parts[2]

	tflog.Info(ctx, "Read branch protection", map[string]any{
		"owner":       owner,
		"repo":        repo,
		"branch_name": branchName,
	})

	// Use Forgejo client to get branch protection
	protection, res, err := r.client.GetBranchProtection(
		owner,
		repo,
		branchName,
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
					"Branch protection not found for %s/%s/%s: %s",
					owner,
					repo,
					branchName,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		response.Diagnostics.AddError("Unable to read branch protection", msg)

		return
	}

	tflog.Info(ctx, "Read repository", map[string]any{
		"owner": owner,
		"repo":  repo,
	})

	// Use Forgejo client to get repository by owner and name
	repository, res, err := r.client.GetRepo(owner, repo)
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
					"Repository with owner '%s' and name '%s' not found: %s",
					owner,
					repo,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		response.Diagnostics.AddError("Unable to read repository", msg)

		return
	}

	// Map response to model
	data.BranchName = types.StringValue(branchName)
	data.RepositoryId = types.Int64Value(repository.ID)
	diags := r.from(protection, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = response.State.Set(ctx, &data)
	response.Diagnostics.Append(diags...)
}

// NewBranchProtectionResource is a helper function to simplify the provider implementation.
func NewBranchProtectionResource() resource.Resource {
	return &branchProtectionResource{}
}

// Helper function to convert model to CreateBranchProtectionOption.
func (r *branchProtectionResource) toCreateOption(ctx context.Context, data *branchProtectionResourceModel) forgejo.CreateBranchProtectionOption {
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

// Helper function to convert model to EditBranchProtectionOption.
func (r *branchProtectionResource) toEditOption(ctx context.Context, data *branchProtectionResourceModel) forgejo.EditBranchProtectionOption {
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

// from is a helper function to load an API struct into Terraform data model.
func (r *branchProtectionResource) from(protection *forgejo.BranchProtection, data *branchProtectionResourceModel) (diags diag.Diagnostics) {
	if protection == nil || data == nil {
		return diags
	}

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

func (r *branchProtectionResource) stringSliceToList(slice []string) (types.List, diag.Diagnostics) {
	if len(slice) == 0 {
		return types.ListNull(types.StringType), diag.Diagnostics{}
	}
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
		elements[i] = types.StringValue(v)
	}

	return types.ListValue(types.StringType, elements)
}
