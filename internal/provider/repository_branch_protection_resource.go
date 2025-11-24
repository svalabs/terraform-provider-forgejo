package provider

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	ID                            types.String `tfsdk:"id"`
	Owner                         types.String `tfsdk:"owner"`
	Repo                          types.String `tfsdk:"repo"`
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
			"id": schema.StringAttribute{
				Description: "Computed ID of the branch protection (format: owner/repo/branch_name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // avoiding unnecessarily showing changes, i.e. '(known after apply)', on update/refresh
				},
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo": schema.StringAttribute{
				Description: "Name of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Default:     booldefault.StaticBool(true),
			},
			"enable_push_whitelist": schema.BoolAttribute{
				Description: "Enable push whitelist.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"push_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames allowed to push.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"push_whitelist_teams": schema.ListAttribute{
				Description: "Teams allowed to push.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"push_whitelist_deploy_keys": schema.BoolAttribute{
				Description: "Allow deploy keys to push.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_status_check": schema.BoolAttribute{
				Description: "Enable status checks.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"status_check_contexts": schema.ListAttribute{
				Description: "Status check contexts that must pass.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"require_signed_commits": schema.BoolAttribute{
				Description: "Require signed commits.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"protected_file_patterns": schema.StringAttribute{
				Description: "Patterns for protected files.",
				Optional:    true,
			},
			"unprotected_file_patterns": schema.StringAttribute{
				Description: "Patterns for unprotected files.",
				Optional:    true,
			},
			"enable_merge_whitelist": schema.BoolAttribute{
				Description: "Enable merge whitelist.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"merge_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames allowed to merge.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"merge_whitelist_teams": schema.ListAttribute{
				Description: "Teams allowed to merge.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enable_approvals_whitelist": schema.BoolAttribute{
				Description: "Enable approvals whitelist.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"approvals_whitelist_usernames": schema.ListAttribute{
				Description: "Usernames that can approve.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"approvals_whitelist_teams": schema.ListAttribute{
				Description: "Teams that can approve.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"required_approvals": schema.Int64Attribute{
				Description: "Number of required approvals.",
				Optional:    true,
			},
			"block_on_rejected_reviews": schema.BoolAttribute{
				Description: "Block merge on rejected reviews.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"block_on_official_review_requests": schema.BoolAttribute{
				Description: "Block merge on official review requests.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"block_on_outdated_branch": schema.BoolAttribute{
				Description: "Block merge on outdated branch.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dismiss_stale_approvals": schema.BoolAttribute{
				Description: "Dismiss stale approvals.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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

	tflog.Info(ctx, "Create repository branch protection", map[string]any{
		"owner":       data.Owner.ValueString(),
		"repo":        data.Repo.ValueString(),
		"branch_name": data.BranchName.ValueString(),
	})

	// Convert model to API request
	opts := r.modelToCreateOption(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create branch protection
	_, res, err := r.client.CreateBranchProtection(
		data.Owner.ValueString(),
		data.Repo.ValueString(),
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
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		} else {
			msg = fmt.Sprintf("Failed to create branch protection: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create branch protection", msg)
		return
	}

	// Set the ID
	data.ID = types.StringValue(fmt.Sprintf("%s/%s/%s",
		data.Owner.ValueString(),
		data.Repo.ValueString(),
		data.BranchName.ValueString()))

	// Update model with response data
	tflog.Trace(ctx, "created branch protection resource")

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

	// Get branch protection
	_, res, err := r.client.GetBranchProtection(
		data.Owner.ValueString(),
		data.Repo.ValueString(),
		data.BranchName.ValueString(),
	)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			// Branch protection no longer exists
			tflog.Trace(ctx, "branch protection not found, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to read branch protection",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// Set the ID (in case it wasn't set or needs to be refreshed)
	data.ID = types.StringValue(fmt.Sprintf("%s/%s/%s",
		data.Owner.ValueString(),
		data.Repo.ValueString(),
		data.BranchName.ValueString()))

	// Update model with response data
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

	tflog.Info(ctx, "Update repository branch protection", map[string]any{
		"owner":       data.Owner.ValueString(),
		"repo":        data.Repo.ValueString(),
		"branch_name": data.BranchName.ValueString(),
	})

	// Convert model to API request
	opts := r.modelToEditOption(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update branch protection
	_, res, err := r.client.EditBranchProtection(
		data.Owner.ValueString(),
		data.Repo.ValueString(),
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
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		} else {
			msg = fmt.Sprintf("Failed to update branch protection: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update branch protection", msg)
		return
	}

	// Set the ID
	data.ID = types.StringValue(fmt.Sprintf("%s/%s/%s",
		data.Owner.ValueString(),
		data.Repo.ValueString(),
		data.BranchName.ValueString()))

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

	tflog.Info(ctx, "Delete repository branch protection", map[string]any{
		"owner":       data.Owner.ValueString(),
		"repo":        data.Repo.ValueString(),
		"branch_name": data.BranchName.ValueString(),
	})

	// Delete branch protection
	res, err := r.client.DeleteBranchProtection(
		data.Owner.ValueString(),
		data.Repo.ValueString(),
		data.BranchName.ValueString(),
	)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			// Already deleted
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

	// Map response to model
	var data repositoryBranchProtectionResourceModel
	data.ID = types.StringValue(request.ID)
	data.Owner = types.StringValue(owner)
	data.Repo = types.StringValue(repo)
	data.BranchName = types.StringValue(branchName)

	// Map protection settings
	data.EnablePush = types.BoolValue(protection.EnablePush)
	data.EnablePushWhitelist = types.BoolValue(protection.EnablePushWhitelist)
	data.PushWhitelistDeployKeys = types.BoolValue(protection.PushWhitelistDeployKeys)
	data.EnableStatusCheck = types.BoolValue(protection.EnableStatusCheck)
	data.RequireSignedCommits = types.BoolValue(protection.RequireSignedCommits)

	// Handle optional string attributes - use null if empty
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

	// Handle required_approvals - use null if zero
	if protection.RequiredApprovals != 0 {
		data.RequiredApprovals = types.Int64Value(protection.RequiredApprovals)
	} else {
		data.RequiredApprovals = types.Int64Null()
	}

	data.BlockOnRejectedReviews = types.BoolValue(protection.BlockOnRejectedReviews)
	data.BlockOnOfficialReviewRequests = types.BoolValue(protection.BlockOnOfficialReviewRequests)
	data.BlockOnOutdatedBranch = types.BoolValue(protection.BlockOnOutdatedBranch)
	data.DismissStaleApprovals = types.BoolValue(protection.DismissStaleApprovals)
	data.PushWhitelistUsernames = types.ListNull(types.StringType)
	data.PushWhitelistTeams = types.ListNull(types.StringType)
	data.StatusCheckContexts = types.ListNull(types.StringType)
	data.MergeWhitelistUsernames = types.ListNull(types.StringType)
	data.MergeWhitelistTeams = types.ListNull(types.StringType)
	data.ApprovalsWhitelistUsernames = types.ListNull(types.StringType)
	data.ApprovalsWhitelistTeams = types.ListNull(types.StringType)

	if len(protection.PushWhitelistUsernames) > 0 {
		elements := make([]attr.Value, len(protection.PushWhitelistUsernames))
		for i, username := range protection.PushWhitelistUsernames {
			elements[i] = types.StringValue(username)
		}
		data.PushWhitelistUsernames, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.PushWhitelistTeams) > 0 {
		elements := make([]attr.Value, len(protection.PushWhitelistTeams))
		for i, team := range protection.PushWhitelistTeams {
			elements[i] = types.StringValue(team)
		}
		data.PushWhitelistTeams, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.StatusCheckContexts) > 0 {
		elements := make([]attr.Value, len(protection.StatusCheckContexts))
		for i, checkContext := range protection.StatusCheckContexts {
			elements[i] = types.StringValue(checkContext)
		}
		data.StatusCheckContexts, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.MergeWhitelistUsernames) > 0 {
		elements := make([]attr.Value, len(protection.MergeWhitelistUsernames))
		for i, username := range protection.MergeWhitelistUsernames {
			elements[i] = types.StringValue(username)
		}
		data.MergeWhitelistUsernames, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.MergeWhitelistTeams) > 0 {
		elements := make([]attr.Value, len(protection.MergeWhitelistTeams))
		for i, team := range protection.MergeWhitelistTeams {
			elements[i] = types.StringValue(team)
		}
		data.MergeWhitelistTeams, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.ApprovalsWhitelistUsernames) > 0 {
		elements := make([]attr.Value, len(protection.ApprovalsWhitelistUsernames))
		for i, username := range protection.ApprovalsWhitelistUsernames {
			elements[i] = types.StringValue(username)
		}
		data.ApprovalsWhitelistUsernames, _ = types.ListValue(types.StringType, elements)
	}

	if len(protection.ApprovalsWhitelistTeams) > 0 {
		elements := make([]attr.Value, len(protection.ApprovalsWhitelistTeams))
		for i, team := range protection.ApprovalsWhitelistTeams {
			elements[i] = types.StringValue(team)
		}
		data.ApprovalsWhitelistTeams, _ = types.ListValue(types.StringType, elements)
	}

	// Save data into Terraform state
	diags := response.State.Set(ctx, &data)
	response.Diagnostics.Append(diags...)
}

// Helper function to convert model to CreateBranchProtectionOption.
func (r *repositoryBranchProtectionResource) modelToCreateOption(ctx context.Context, data *repositoryBranchProtectionResourceModel) forgejo.CreateBranchProtectionOption {

	opts := forgejo.CreateBranchProtectionOption{
		BranchName: data.BranchName.ValueString(),
	}

	if !data.EnablePush.IsNull() {
		opts.EnablePush = *data.EnablePush.ValueBoolPointer()
	}

	if !data.EnablePushWhitelist.IsNull() {
		opts.EnablePushWhitelist = *data.EnablePushWhitelist.ValueBoolPointer()
	}

	if !data.PushWhitelistUsernames.IsNull() {
		var usernames []string
		data.PushWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.PushWhitelistUsernames = usernames
	}

	if !data.PushWhitelistTeams.IsNull() {
		var teams []string
		data.PushWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.PushWhitelistTeams = teams
	}

	if !data.PushWhitelistDeployKeys.IsNull() {
		opts.PushWhitelistDeployKeys = *data.PushWhitelistDeployKeys.ValueBoolPointer()
	}

	if !data.EnableStatusCheck.IsNull() {
		opts.EnableStatusCheck = *data.EnableStatusCheck.ValueBoolPointer()
	}

	if !data.StatusCheckContexts.IsNull() {
		var contexts []string
		data.StatusCheckContexts.ElementsAs(ctx, &contexts, false)
		opts.StatusCheckContexts = contexts
	}

	if !data.RequireSignedCommits.IsNull() {
		opts.RequireSignedCommits = *data.RequireSignedCommits.ValueBoolPointer()
	}

	if !data.ProtectedFilePatterns.IsNull() {
		opts.ProtectedFilePatterns = data.ProtectedFilePatterns.ValueString()
	}

	if !data.UnprotectedFilePatterns.IsNull() {
		opts.UnprotectedFilePatterns = data.UnprotectedFilePatterns.ValueString()
	}

	if !data.EnableMergeWhitelist.IsNull() {
		opts.EnableMergeWhitelist = *data.EnableMergeWhitelist.ValueBoolPointer()
	}

	if !data.MergeWhitelistUsernames.IsNull() {
		var usernames []string
		data.MergeWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.MergeWhitelistUsernames = usernames
	}

	if !data.MergeWhitelistTeams.IsNull() {
		var teams []string
		data.MergeWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.MergeWhitelistTeams = teams
	}

	if !data.EnableApprovalsWhitelist.IsNull() {
		opts.EnableApprovalsWhitelist = *data.EnableApprovalsWhitelist.ValueBoolPointer()
	}

	if !data.ApprovalsWhitelistUsernames.IsNull() {
		var usernames []string
		data.ApprovalsWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.ApprovalsWhitelistUsernames = usernames
	}

	if !data.ApprovalsWhitelistTeams.IsNull() {
		var teams []string
		data.ApprovalsWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.ApprovalsWhitelistTeams = teams
	}

	if !data.RequiredApprovals.IsNull() {
		approvals := data.RequiredApprovals.ValueInt64()
		opts.RequiredApprovals = approvals
	}

	if !data.BlockOnRejectedReviews.IsNull() {
		opts.BlockOnRejectedReviews = *data.BlockOnRejectedReviews.ValueBoolPointer()
	}

	if !data.BlockOnOfficialReviewRequests.IsNull() {
		opts.BlockOnOfficialReviewRequests = *data.BlockOnOfficialReviewRequests.ValueBoolPointer()
	}

	if !data.BlockOnOutdatedBranch.IsNull() {
		opts.BlockOnOutdatedBranch = *data.BlockOnOutdatedBranch.ValueBoolPointer()
	}

	if !data.DismissStaleApprovals.IsNull() {
		opts.DismissStaleApprovals = *data.DismissStaleApprovals.ValueBoolPointer()
	}

	return opts
}

// Helper function to convert model to EditBranchProtectionOption.
func (r *repositoryBranchProtectionResource) modelToEditOption(ctx context.Context, data *repositoryBranchProtectionResourceModel) forgejo.EditBranchProtectionOption {
	opts := forgejo.EditBranchProtectionOption{}

	if !data.EnablePush.IsNull() {
		opts.EnablePush = data.EnablePush.ValueBoolPointer()
	}

	if !data.EnablePushWhitelist.IsNull() {
		opts.EnablePushWhitelist = data.EnablePushWhitelist.ValueBoolPointer()
	}

	if !data.PushWhitelistUsernames.IsNull() {
		var usernames []string
		data.PushWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.PushWhitelistUsernames = usernames
	}

	if !data.PushWhitelistTeams.IsNull() {
		var teams []string
		data.PushWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.PushWhitelistTeams = teams
	}

	if !data.PushWhitelistDeployKeys.IsNull() {
		opts.PushWhitelistDeployKeys = data.PushWhitelistDeployKeys.ValueBoolPointer()
	}

	if !data.EnableStatusCheck.IsNull() {
		opts.EnableStatusCheck = data.EnableStatusCheck.ValueBoolPointer()
	}

	if !data.StatusCheckContexts.IsNull() {
		var contexts []string
		data.StatusCheckContexts.ElementsAs(ctx, &contexts, false)
		opts.StatusCheckContexts = contexts
	}

	if !data.RequireSignedCommits.IsNull() {
		opts.RequireSignedCommits = data.RequireSignedCommits.ValueBoolPointer()
	}

	if !data.ProtectedFilePatterns.IsNull() {
		patterns := data.ProtectedFilePatterns.ValueString()
		opts.ProtectedFilePatterns = &patterns
	}

	if !data.UnprotectedFilePatterns.IsNull() {
		patterns := data.UnprotectedFilePatterns.ValueString()
		opts.UnprotectedFilePatterns = &patterns
	}

	if !data.EnableMergeWhitelist.IsNull() {
		opts.EnableMergeWhitelist = data.EnableMergeWhitelist.ValueBoolPointer()
	}

	if !data.MergeWhitelistUsernames.IsNull() {
		var usernames []string
		data.MergeWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.MergeWhitelistUsernames = usernames
	}

	if !data.MergeWhitelistTeams.IsNull() {
		var teams []string
		data.MergeWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.MergeWhitelistTeams = teams
	}

	if !data.EnableApprovalsWhitelist.IsNull() {
		opts.EnableApprovalsWhitelist = data.EnableApprovalsWhitelist.ValueBoolPointer()
	}

	if !data.ApprovalsWhitelistUsernames.IsNull() {
		var usernames []string
		data.ApprovalsWhitelistUsernames.ElementsAs(ctx, &usernames, false)
		opts.ApprovalsWhitelistUsernames = usernames
	}

	if !data.ApprovalsWhitelistTeams.IsNull() {
		var teams []string
		data.ApprovalsWhitelistTeams.ElementsAs(ctx, &teams, false)
		opts.ApprovalsWhitelistTeams = teams
	}

	if !data.RequiredApprovals.IsNull() {
		approvals := data.RequiredApprovals.ValueInt64()
		opts.RequiredApprovals = &approvals
	}

	if !data.BlockOnRejectedReviews.IsNull() {
		opts.BlockOnRejectedReviews = data.BlockOnRejectedReviews.ValueBoolPointer()
	}

	if !data.BlockOnOfficialReviewRequests.IsNull() {
		opts.BlockOnOfficialReviewRequests = data.BlockOnOfficialReviewRequests.ValueBoolPointer()
	}

	if !data.BlockOnOutdatedBranch.IsNull() {
		opts.BlockOnOutdatedBranch = data.BlockOnOutdatedBranch.ValueBoolPointer()
	}

	if !data.DismissStaleApprovals.IsNull() {
		opts.DismissStaleApprovals = data.DismissStaleApprovals.ValueBoolPointer()
	}

	return opts
}
