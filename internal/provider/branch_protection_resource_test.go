package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccBranchProtectionResource1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (personal repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
			// Create and Read testing (duplicate branch_name)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id
}
resource "forgejo_branch_protection" "duplicate" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id
}`,
				ExpectError: regexp.MustCompile(`Repository with owner "` + forgejoTestUser + `" and name "test_repo_branch_protection"
forbidden: Branch protection already exist`),
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_branch_protection.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile("Expected import identifier with format: 'owner/repo/branch', got: 'invalid'"),
			},
			// Import testing (non-existent repo)
			{
				ResourceName:  "forgejo_branch_protection.test",
				ImportState:   true,
				ImportStateId: "" + forgejoTestUser + "/non-existent/main",
				ExpectError: regexp.MustCompile(`Branch protection with owner '` + forgejoTestUser + `', repo 'non-existent' and name 'main'
not found`),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_branch_protection.test",
				ImportState:   true,
				ImportStateId: forgejoTestUser + "/test_repo_branch_protection/non-existent",
				ExpectError: regexp.MustCompile(`Branch protection with owner '` + forgejoTestUser + `', repo 'test_repo_branch_protection'
and name 'non-existent' not found`),
			},
			// Import testing
			{
				ResourceName:                         "forgejo_branch_protection.test",
				ImportState:                          true,
				ImportStateId:                        forgejoTestUser + "/test_repo_branch_protection/main",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "branch_name",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id

	enable_push               = true
	require_signed_commits    = true
	required_approvals        = 1
	protected_file_patterns   = "*.tf"
	unprotected_file_patterns = "*.log"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("*.tf")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("*.log")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
			// Recreate and Read testing since the branch name has changed
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "dev"
	repository_id = forgejo_repository.test.id

	enable_push            = false
	require_signed_commits = false
	required_approvals     = null

	enable_status_check   = true
	status_check_contexts = ["ci/on-submit"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("dev")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
			// Recreate and Read testing (organization repo)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_branch_protection"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name    = "main"
	repository_id  = forgejo_repository.test.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
			// Update and Read testing (custom sort order)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_branch_protection"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo_branch_protection"
}
resource "forgejo_user" "test" {
  for_each = toset(["alice", "bob"])
  login    = each.key
  email    = "${each.key}@localhost"
  password = "12345678"
}
data "forgejo_team" "test" {
  organization_id = forgejo_organization.test.id
  name            = "Owners"
}
resource "forgejo_team_member" "test" {
  for_each = toset(["alice", "bob"])
  team_id  = data.forgejo_team.test.id
  user     = forgejo_user.test[each.key].login
}
resource "forgejo_branch_protection" "test" {
  depends_on = [forgejo_team_member.test]
  branch_name   = "main"
  repository_id = forgejo_repository.test.id

  block_on_outdated_branch  = true
  block_on_rejected_reviews = true
  dismiss_stale_approvals   = true
  required_approvals        = 1

  enable_approvals_whitelist = true
  approvals_whitelist_usernames = [
    "bob",
    "alice",
  ]

  enable_merge_whitelist = true
  merge_whitelist_usernames = [
    "bob",
    "alice",
  ]

  enable_push           = true
  enable_push_whitelist = true
  push_whitelist_usernames = [
    "bob",
    "alice",
  ]

  enable_status_check = true
  status_check_contexts = [
    "foo",
    "bar",
  ]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(true)),
				},
			},
			// Recreate and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user_branch_protection"
	email    = "test_user_branch_protection@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.test.login
	name  = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name		= "main"
	repository_id	= forgejo_repository.test.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccBranchProtectionResource2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (with optional attributes)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                       = "main"
	repository_id                     = forgejo_repository.test.id
	enable_push                       = true
	enable_push_whitelist             = true
	push_whitelist_usernames          = ["` + forgejoTestUser + `"]
	require_signed_commits            = true
	required_approvals                = 2
	block_on_rejected_reviews         = true
	block_on_official_review_requests = true
	dismiss_stale_approvals           = true
	protected_file_patterns           = "*.tf"
	unprotected_file_patterns         = "*.log"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("*.tf")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("*.log")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(2)),
				},
			},
			// Recreate and Read testing (branch pattern)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "release/*"
	repository_id = forgejo_repository.test.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_branch_protection.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("release/*")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
				},
			},
			// Create and Read testing (invalid repo)
			{
				Config: providerConfig + `
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = 123
}`,
				ExpectError: regexp.MustCompile("Repository with ID 123 not found"),
			},
			// Create and Read testing (archived repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name     = "test_repo_archived"
	archived = true
}
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id
}`,
				ExpectError: regexp.MustCompile("Repository with owner \"" + forgejoTestUser + "\" and name \"test_repo_archived\" is archived"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBranchProtectionValidationPushConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check push_enabled and push_whitelist_enabled attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_push           = false
	enable_push_whitelist = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If enable_push_whitelist is 'true', enable_push must be 'true'"),
			},
			// Check push_whitelist_usernames attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name              = "main"
	repository_id            = forgejo_repository.test.id
	enable_push              = false
	enable_push_whitelist    = false
	push_whitelist_usernames = ["user1", "user2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If push_whitelist_usernames is not empty, enable_push_whitelist must be\n'true'"),
			},
			// Check push_whitelist_teams attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_push           = false
	enable_push_whitelist = false
	push_whitelist_teams  = ["team1", "team2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If push_whitelist_teams is not empty, enable_push_whitelist must be 'true'"),
			},
			// Check push_whitelist_deploy_keys attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = false
	enable_push_whitelist      = false
	push_whitelist_deploy_keys = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If push_whitelist_deploy_keys is 'true', enable_push_whitelist must be 'true'"),
			},
			// enabling push, but not push_whitelist while configuring whitelist attributes ought to also error
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = true
	enable_push_whitelist      = false
	push_whitelist_deploy_keys = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If push_whitelist_deploy_keys is 'true', enable_push_whitelist must be 'true'"),
			},
			// successful push configuration with push_whitelist
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = true
	enable_push_whitelist      = true
	push_whitelist_deploy_keys = true
	push_whitelist_usernames   = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
		},
	})
}

func TestAccBranchProtectionValidationStatusCheckConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check status_check_enabled and status_check_contexts attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = false
	status_check_contexts = ["context1", "context2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If status_check_contexts is not empty, enable_status_check must be 'true'"),
			},
			// Valid. status_check_enabled is false and contexts are empty
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = false
	status_check_contexts = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
			// Valid. status_check_enabled is enabled and contexts are set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = true
	status_check_contexts = ["*"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
		},
	})
}

func TestAccBranchProtectionValidationMergeConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check merge_settings attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = false
	merge_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If merge_whitelist_usernames is not empty, enable_merge_whitelist must be\n'true'"),
			},
			// Valid. merge_whitelist_usernames is set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = true
	merge_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
			// Valid. merge configuration is disabled
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = false
	merge_whitelist_usernames = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
		},
	})
}

func TestAccBranchProtectionValidationApprovalsConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check approvals attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = false
	approvals_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("If approvals_whitelist_usernames is not empty, enable_approvals_whitelist\nmust be 'true'"),
			},
			// Valid. approvals_whitelist_usernames is set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = true
	approvals_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
			// Valid. merge configuration is disabled
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = false
	approvals_whitelist_usernames = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false))},
			},
		},
	})
}
