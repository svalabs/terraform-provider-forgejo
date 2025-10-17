package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (personal repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "tftest"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
	website     = "http://localhost:3000"

	has_issues        = false
	has_pull_requests = false
	has_wiki          = false
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000")),
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
	website     = "http://localhost:3000"

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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest1.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("rebase")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 123")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("allow_only_contributors_to_track_time"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_issue_dependencies"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_time_tracker"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000")),
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
	website     = "http://localhost:3000"

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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest2.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("squash")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 456")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_url"), knownvalue.StringExact("https://some.tracker")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_format"), knownvalue.StringExact("https://some.tracker/{user}/{repo}/{index}")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker").AtMapKey("external_tracker_style"), knownvalue.StringExact("alphanumeric")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki").AtMapKey("external_wiki_url"), knownvalue.StringExact("https://some.wiki")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest2")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000")),
				},
			},
			// Recreate and Read testing (clone repo)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name       = "tftest"
	clone_addr = "https://github.com/svalabs/terraform-provider-forgejo"
	mirror     = false
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/svalabs/terraform-provider-forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing... 123")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/acch/test-non-default-branch")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest-mirror-develop.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest-mirror-develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest-mirror-develop")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest-mirror-develop.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("https://github.com/acch/test-non-default-branch")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest-mirror-production.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest-mirror-production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest-mirror-production")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest-mirror-production.git")),
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
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("autodetect_manual_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_addr"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("Go")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/tfadmin/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("Default")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("service"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("ssh://git@localhost:2222/tfadmin/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("collaborator")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
