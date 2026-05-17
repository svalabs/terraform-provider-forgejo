package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.ConfigValidator = &branchProtectionResourcePushConfigValidator{}
	_ resource.ConfigValidator = &branchProtectionResourceStatusCheckConfigValidator{}
	_ resource.ConfigValidator = &branchProtectionResourceMergeConfigValidator{}
	_ resource.ConfigValidator = &branchProtectionResourceApprovalConfigValidator{}
)

// branchProtectionResourcePushConfigValidator validates the configuration of branch protection push settings.
type branchProtectionResourcePushConfigValidator struct {
}

// Description describes the validation in plain text formatting.
func (b branchProtectionResourcePushConfigValidator) Description(_ context.Context) string {
	return `
Validates the push whitelist configuration of a branch protection resource.
If 'enable_push_whitelist' is false, both 'push_whitelist_usernames' and 'push_whitelist_teams' must be empty.
If 'enable_push_whitelist' is false, 'push_whitelist_deploy_keys' must be false.
`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (b branchProtectionResourcePushConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Validates that `push_whitelist_deploy_keys` is not true when `enable_push_whitelist` is false.\n\n" +
		"- If `enable_push_whitelist` is false, both `push_whitelist_usernames` and `push_whitelist_teams` must be empty.\n" +
		"- If `enable_push_whitelist` is false, `push_whitelist_deploy_keys` must be false."
}

// ValidateResource validates the configuration of branch protection push settings.
// Decision Matrix:
// - If 'enable_push' is false, 'enable_push_whitelist' must be false.
// - If 'enable_push_whitelist' is false, both 'push_whitelist_usernames' and 'push_whitelist_teams' must be empty.
// - If 'enable_push_whitelist' is false, 'push_whitelist_deploy_keys' must be false.
func (b branchProtectionResourcePushConfigValidator) ValidateResource(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config branchProtectionResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	enablePush := config.EnablePush
	enablePushWhitelist := config.EnablePushWhitelist
	pushWhitelistUsernames := config.PushWhitelistUsernames
	pushWhitelistTeams := config.PushWhitelistTeams
	pushWhitelistDeployKeys := config.PushWhitelistDeployKeys

	if enablePush.IsUnknown() || enablePushWhitelist.IsUnknown() || pushWhitelistDeployKeys.IsUnknown() || pushWhitelistUsernames.IsUnknown() || pushWhitelistTeams.IsUnknown() {
		return
	}

	if !enablePush.ValueBool() && enablePushWhitelist.ValueBool() {
		response.Diagnostics.AddError(
			"Cannot enable push whitelist without enabling push protection",
			"Set 'enable_push' to true if 'enable_push_whitelist' is true",
		)
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}

	pushWhitelistUsernamesNotEmpty := pushWhitelistUsernames.Length(opts) > 0
	pushWhitelistTeamsNotEmpty := pushWhitelistTeams.Length(opts) > 0

	if !enablePushWhitelist.ValueBool() && (pushWhitelistDeployKeys.ValueBool() || pushWhitelistUsernamesNotEmpty || pushWhitelistTeamsNotEmpty) {
		response.Diagnostics.AddError(
			"Push Whitelist configuration is not valid when 'enable_push_whitelist' is false",
			"Set 'enable_push_whitelist' to true if 'push_whitelist_teams', 'push_whitelist_teams' or 'push_whitelist_deploy_keys' are used",
		)
	}
}

// branchProtectionResourceStatusCheckConfigValidator validates the configuration of branch protection status check settings.
type branchProtectionResourceStatusCheckConfigValidator struct {
}

// Description describes the validation in plain text formatting.
func (b branchProtectionResourceStatusCheckConfigValidator) Description(_ context.Context) string {
	return "Validates that 'enable_status_check' is true when 'status_check_contexts' is set."
}

// MarkdownDescription describes the validation in Markdown formatting.
func (b branchProtectionResourceStatusCheckConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Validates that `enable_status_check` is true when `status_check_contexts` is set."
}

// ValidateResource validates the configuration of branch protection status check settings.
// Decision Matrix:
// - If 'enable_status_check' is false, 'status_check_contexts' must be empty.
func (b branchProtectionResourceStatusCheckConfigValidator) ValidateResource(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config branchProtectionResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	enableStatusCheck := config.EnableStatusCheck
	statusCheckContexts := config.StatusCheckContexts

	if enableStatusCheck.IsUnknown() || statusCheckContexts.IsUnknown() {
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}

	if !enableStatusCheck.ValueBool() && statusCheckContexts.Length(opts) > 0 {
		response.Diagnostics.AddError(
			"Cannot specify status check contexts without enabling status check",
			"Set 'enable_status_check' to true if 'status_check_contexts' are used",
		)
	}
}

// branchProtectionResourceMergeConfigValidator validates the configuration of branch protection merge settings.
type branchProtectionResourceMergeConfigValidator struct {
}

// Description describes the validation in plain text formatting.
func (b branchProtectionResourceMergeConfigValidator) Description(_ context.Context) string {
	return `
Validates the merge configuration of a branch protection resource.
If 'enable_merge_whitelist' is false, both 'merge_whitelist_usernames' and 'merge_whitelist_teams' must be empty.
`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (b branchProtectionResourceMergeConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Validates the merge configuration of a branch protection resource.\n\n" +
		"- If `enable_merge_whitelist` is false, both `merge_whitelist_usernames` and `merge_whitelist_teams` must be empty.\n"
}

// ValidateResource validates the configuration of the merge settings for a branch protection resource.
// Decision Matrix:
// - If 'enable_merge_whitelist' is false, both 'merge_whitelist_usernames' and 'merge_whitelist_teams' must be empty.
func (b branchProtectionResourceMergeConfigValidator) ValidateResource(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config branchProtectionResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	enableWhitelist := config.EnableMergeWhitelist
	whitelistedUsernames := config.MergeWhitelistUsernames
	whitelistedTeams := config.MergeWhitelistTeams

	if enableWhitelist.IsUnknown() || whitelistedUsernames.IsUnknown() || whitelistedTeams.IsUnknown() {
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}

	areWhitelistedUsernamesNotEmpty := whitelistedUsernames.Length(opts) > 0
	areWhitelistedTeamsNotEmpty := whitelistedTeams.Length(opts) > 0

	if !enableWhitelist.ValueBool() && (areWhitelistedUsernamesNotEmpty || areWhitelistedTeamsNotEmpty) {
		response.Diagnostics.AddError(
			"Merge configuration is not valid when 'enable_merge_whitelist' is false",
			"Set 'enable_merge_whitelist' to true if 'merge_whitelist_usernames' or 'merge_whitelist_teams' are used",
		)
	}
}

// branchProtectionResourceApprovalConfigValidator validates the configuration of branch protection approval settings.
type branchProtectionResourceApprovalConfigValidator struct {
}

// Description describes the validation in plain text formatting.
func (b branchProtectionResourceApprovalConfigValidator) Description(_ context.Context) string {
	return `
Validates the approvals configuration of a branch protection resource.
If 'enable_approvals_whitelist' is false, both 'approvals_whitelist_usernames' and 'approvals_whitelist_teams' must be empty.
`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (b branchProtectionResourceApprovalConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Validates the approvals configuration of a branch protection resource.\n\n" +
		"- If `enable_approvals_whitelist` is false, both `approvals_whitelist_usernames` and `approvals_whitelist_teams` must be empty.\n"
}

// ValidateResource validates the configuration of the approval settings for a branch protection resource.
// Decision Matrix:
// - If 'enable_approvals_whitelist' is false, both 'approvals_whitelist_usernames' and 'approvals_whitelist_teams' must be empty.
func (b branchProtectionResourceApprovalConfigValidator) ValidateResource(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config branchProtectionResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	enableWhitelist := config.EnableApprovalsWhitelist
	whitelistedUsernames := config.ApprovalsWhitelistUsernames
	whitelistedTeams := config.ApprovalsWhitelistTeams

	if enableWhitelist.IsUnknown() || whitelistedUsernames.IsUnknown() || whitelistedTeams.IsUnknown() {
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}

	areWhitelistedUsernamesNotEmpty := whitelistedUsernames.Length(opts) > 0
	areWhitelistedTeamsNotEmpty := whitelistedTeams.Length(opts) > 0

	if !enableWhitelist.ValueBool() && (areWhitelistedUsernamesNotEmpty || areWhitelistedTeamsNotEmpty) {
		response.Diagnostics.AddError(
			"Approvals configuration is not valid when 'enable_approvals_whitelist' is false",
			"Set 'enable_approvals_whitelist' to true if 'approvals_whitelist_usernames' or 'approvals_whitelist_teams' are used",
		)
	}
}
