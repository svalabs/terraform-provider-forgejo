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

func TestAccRepositoryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent user)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	owner = "non_existent"
	name  = "tftest"
}`,
				ExpectError: regexp.MustCompile("Neither organization nor user with name \"non_existent\" exists"),
			},
			// Create and Read testing (personal repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Create and Read testing (duplicate name)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "tftest"
}
resource "forgejo_repository" "duplicate" {
	name = "tftest"
}`,
				ExpectError: regexp.MustCompile("Repository with name \"tftest\" already exists"),
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_repository.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile("Expected import identifier with format: 'owner/name', got: 'invalid'"),
			},
			// Import testing (non-existent user)
			{
				ResourceName:  "forgejo_repository.test",
				ImportState:   true,
				ImportStateId: "non-existent/tftest",
				ExpectError:   regexp.MustCompile("Repository with owner 'non-existent' and name 'tftest' not found"),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_repository.test",
				ImportState:   true,
				ImportStateId: forgejoTestUser + "/non-existent",
				ExpectError:   regexp.MustCompile("Repository with owner '" + forgejoTestUser + "' and name 'non-existent' not found"),
			},
			// Import testing
			{
				ResourceName:      "forgejo_repository.test",
				ImportState:       true,
				ImportStateId:     forgejoTestUser + "/tftest",
				ImportStateVerify: true,
			},
			// Recreate and Read testing (org repo)
			{
				Config: providerConfig + `
resource "forgejo_organization" "owner" {
	name = "test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.owner.name
	name  = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest"

	description = "Purely for testing..."
	website     = "` + forgejoTestHost + `"

	has_issues        = false
	has_pull_requests = false
	has_wiki          = false
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact(forgejoTestHost)),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest1"

	description = "Purely for testing... 123"
	website     = "` + forgejoTestHost + `"

	has_issues = true
	internal_tracker = {
		enable_time_tracker                   = false
		allow_only_contributors_to_track_time = false
		enable_issue_dependencies             = false
	}

	has_pull_requests           = true
	allow_manual_merge          = true
	allow_merge_commits         = false
	allow_rebase                = true
	allow_rebase_explicit       = true
	allow_squash_merge          = false
	autodetect_manual_merge     = true
	default_merge_style         = "rebase"
	ignore_whitespace_conflicts = true

	has_wiki = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest1.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("rebase")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 123")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("allow_only_contributors_to_track_time"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_issue_dependencies"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_time_tracker"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest1.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact(forgejoTestHost)),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest2"

	description = "Purely for testing... 456"
	website     = "` + forgejoTestHost + `"

	archived                    = true
	has_actions                 = false
	has_packages                = false
	has_projects                = false
	has_releases                = false
	issue_labels                = "Advanced"
	private                     = true
	template                    = true

	has_issues = true
	external_tracker = {
		external_tracker_url    = "https://some.tracker"
		external_tracker_format = "https://some.tracker/{user}/{repo}/{index}"
		external_tracker_style  = "alphanumeric"
	}

	has_pull_requests           = true
	allow_manual_merge          = true
	allow_merge_commits         = false
	allow_rebase                = false
	allow_rebase_explicit       = false
	allow_squash_merge          = true
	autodetect_manual_merge     = true
	default_merge_style         = "squash"
	ignore_whitespace_conflicts = true

	has_wiki = true
	external_wiki = {
		external_wiki_url = "https://some.wiki"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest2.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("squash")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 456")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_url"), knownvalue.StringExact("https://some.tracker")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_format"), knownvalue.StringExact("https://some.tracker/{user}/{repo}/{index}")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_style"), knownvalue.StringExact("alphanumeric")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_regexp_pattern"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki").AtMapKey("external_wiki_url"), knownvalue.StringExact("https://some.wiki")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest2")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest2")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("Advanced")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest2")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest2.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact(forgejoTestHost)),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest3"

	description = "Purely for testing... 789"
	website     = "` + forgejoTestHost + `"

	archived                    = false
	has_actions                 = true
	has_packages                = true
	has_projects                = true
	has_releases                = true
	issue_labels                = "Default"
	private                     = false
	template                    = false

	has_issues = true
	external_tracker = {
		external_tracker_url            = "https://another.tracker"
		external_tracker_format         = "https://another.tracker/{user}/{repo}/{index}"
		external_tracker_style          = "regexp"
		external_tracker_regexp_pattern = "issue-[0-9]+"
	}

	has_pull_requests           = true
	allow_manual_merge          = false
	allow_merge_commits         = true
	allow_rebase                = true
	allow_rebase_explicit       = true
	allow_squash_merge          = false
	autodetect_manual_merge     = false
	default_merge_style         = "rebase"
	ignore_whitespace_conflicts = false

	has_wiki = true
	external_wiki = {
		external_wiki_url = "https://another.wiki"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest3.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("rebase")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 789")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_url"), knownvalue.StringExact("https://another.tracker")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_format"), knownvalue.StringExact("https://another.tracker/{user}/{repo}/{index}")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_style"), knownvalue.StringExact("regexp")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_regexp_pattern"), knownvalue.StringExact("issue-[0-9]+")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki").AtMapKey("external_wiki_url"), knownvalue.StringExact("https://another.wiki")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("Default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest3.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact(forgejoTestHost)),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest3"

	description = "Purely for testing... abc"
	website     = "` + forgejoTestHost + `"

	archived                    = false
	has_actions                 = true
	has_packages                = true
	has_projects                = true
	has_releases                = true
	issue_labels                = "Default"
	private                     = false
	template                    = false

	has_issues = true
	external_tracker = {
		external_tracker_url            = "https://yet.another.tracker"
		external_tracker_format         = "https://yet.another.tracker/{user}/{repo}/{index}"
		external_tracker_style          = "regexp"
		external_tracker_regexp_pattern = "issue-[0-9]+"
	}

	has_pull_requests           = true
	allow_manual_merge          = false
	allow_merge_commits         = true
	allow_rebase                = true
	allow_rebase_explicit       = true
	allow_squash_merge          = false
	autodetect_manual_merge     = false
	default_merge_style         = "fast-forward-only"
	ignore_whitespace_conflicts = false

	has_wiki = true
	external_wiki = {
		external_wiki_url = "https://yet.another.wiki"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest3.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("fast-forward-only")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... abc")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_url"), knownvalue.StringExact("https://yet.another.tracker")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_format"), knownvalue.StringExact("https://yet.another.tracker/{user}/{repo}/{index}")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_style"), knownvalue.StringExact("regexp")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_regexp_pattern"), knownvalue.StringExact("issue-[0-9]+")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki").AtMapKey("external_wiki_url"), knownvalue.StringExact("https://yet.another.wiki")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("Default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest3")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest3.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact(forgejoTestHost)),
				},
			},
			// Update and Read testing (user repo)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (user repo, default_delete_branch_after_merge)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest"

	has_pull_requests                 = true
	default_delete_branch_after_merge = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (user repo, default_delete_branch_after_merge reset)
			{
				Config: providerConfig + `
resource "forgejo_user" "owner" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.owner.login
	name  = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (clone repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name       = "tftest"
	clone_addr = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror     = false
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (clone repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name        = "tftest"
	clone_addr  = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror      = false
	archived    = true
	description = "Purely for testing..."
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (mirror repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name       = "tftest"
	clone_addr = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror     = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("8h0m0s")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (mirror repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name            = "tftest"
	clone_addr      = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror          = true
	mirror_interval = "12h0m0s"
	description     = "Purely for testing..."
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("12h0m0s")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (mirror repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name            = "tftest"
	clone_addr      = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror          = true
	mirror_interval = "24h0m0s"
	description     = "Purely for testing... 123"
	lfs             = true
	milestones      = true
	labels          = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 123")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("24h0m0s")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (mirror repo with non-default branch)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name            = "tftest-mirror-develop"
	clone_addr      = "https://github.com/acch/test-non-default-branch"
	mirror          = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/acch/test-non-default-branch")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest-mirror-develop.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest-mirror-develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest-mirror-develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("8h0m0s")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest-mirror-develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest-mirror-develop.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Update and Read testing (mirror repo with non-default branch)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name            = "tftest-mirror-production"
	clone_addr      = "https://github.com/acch/test-non-default-branch"
	default_branch  = "production"
	mirror          = true
	mirror_interval = "12h0m0s"
	description     = "Purely for testing..."
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/acch/test-non-default-branch")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest-mirror-production.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest-mirror-production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest-mirror-production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("12h0m0s")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest-mirror-production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest-mirror-production.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (auto-init)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name         = "tftest"
	auto_init    = false
	gitignores   = "Go"
	issue_labels = "Advanced"
	license      = "MIT"
	readme       = "Default"
	trust_model  = "collaborator"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact(forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("Go")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/"+forgejoTestUser+"/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("Advanced")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("MIT")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("Default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/"+forgejoTestUser+"/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("collaborator")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Archive on destroy testing (organisation repo; creation)
			{
				Config: providerConfig + `
resource "forgejo_organization" "owner" {
	name = "test_org"
}
resource "forgejo_repository" "test" {
	owner              = forgejo_organization.owner.name
	name               = "tftest"
	archive_on_destroy = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archive_on_destroy"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Archive on destroy testing (organisation repo; destroying)
			{
				Config: providerConfig + `
resource "forgejo_organization" "owner" {
	name = "test_org"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionDestroy),
					},
				},
			},
			// Archive on destroy testing (organisation repo; reading after destroy)
			{
				Config: providerConfig + `
resource "forgejo_organization" "owner" {
	name = "test_org"
}
data "forgejo_repository" "test" {
	owner = forgejo_organization.owner.name
	name  = "tftest"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest.git")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_org/tftest")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_org/tftest.git")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Archive on destroy testing (organisation repo; importing -- unarchives)
			{
				Config: providerConfig + `
resource "forgejo_organization" "owner" {
	name = "test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.owner.name
	name  = "tftest"
}
data "forgejo_repository" "import_test" {
	owner = forgejo_organization.owner.name
	name  = "tftest"
}
import {
	id = "test_org/tftest"
	to = forgejo_repository.test
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_delete_branch_after_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact(forgejoTestHost+"/test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("labels"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs_endpoint"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("lfs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("milestones"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
