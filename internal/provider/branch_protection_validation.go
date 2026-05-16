package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.ConfigValidator = &branchProtectionResourcePushConfigValidator{}
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
