package provider_test

import (
	"regexp"
	"testing"

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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Null()),
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
				ExpectError: regexp.MustCompile(`Repository with owner "tfadmin" and name "test_repo_branch_protection"
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
				ImportStateId: "tfadmin/non-existent/main",
				ExpectError:   regexp.MustCompile("Branch protection not found for 'tfadmin/non-existent/main'"),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_branch_protection.test",
				ImportState:   true,
				ImportStateId: "tfadmin/test_repo_branch_protection/non-existent",
				ExpectError:   regexp.MustCompile("Branch protection not found for\n'tfadmin/test_repo_branch_protection/non-existent'"),
			},
			// Import testing
			{
				ResourceName:                         "forgejo_branch_protection.test",
				ImportState:                          true,
				ImportStateId:                        "tfadmin/test_repo_branch_protection/main",
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("*.tf")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("*.log")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("dev")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_rejected_reviews"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_official_review_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("block_on_outdated_branch"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("dismiss_stale_approvals"), knownvalue.Bool(false)),
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Null()),
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
	push_whitelist_usernames          = ["tfadmin"]
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.StringExact("*.tf")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.StringExact("*.log")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("release/*")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("protected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("unprotected_file_patterns"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_teams"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Null()),
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
				ExpectError: regexp.MustCompile("Repository with owner \"tfadmin\" and name \"test_repo_archived\" is archived"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
