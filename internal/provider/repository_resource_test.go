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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/achim/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("achim/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/achim/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:achim/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("")),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_org/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:test_org/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing (user)
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("")),
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

	internal_tracker = {
	  enable_time_tracker                   = false
		allow_only_contributors_to_track_time = false
		enable_issue_dependencies             = false
	}
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
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
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_time_tracker"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("allow_only_contributors_to_track_time"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker").AtMapKey("enable_issue_dependencies"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:test_user/tftest.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("")),
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
  owner                       = forgejo_user.owner.login
  name                        = "tftest1"
	allow_merge_commits         = false
	allow_rebase                = false
	allow_rebase_explicit       = false
	allow_squash_merge          = false
	archived                    = true
	description                 = "Purely for testing..."
	has_actions                 = false
  has_issues                  = false
  has_packages                = false
  has_projects                = false
  has_pull_requests           = true
  has_releases                = false
  has_wiki                    = false
	ignore_whitespace_conflicts = true
	private                     = true
	template                    = true
	website                     = "http://localhost:3000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("auto_init"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest1.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("test_user/tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("gitignores"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/test_user/tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("issue_labels"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("license"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("private"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("readme"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:test_user/tftest1.git")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("template"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("trust_model"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
