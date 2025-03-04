package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "forgejo_repository" "test" {
  owner = {
	  login = "achim"
	}
  name = "user_test_repo_1"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_merge_commits"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_rebase"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_rebase_explicit"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("allow_squash_merge"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("archived"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("avatar_url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("clone_url"), knownvalue.StringExact("http://localhost:3000/achim/user_test_repo_1.git")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("default_branch"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("default_merge_style"), knownvalue.StringExact("merge")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("description"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("empty"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("external_tracker"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("external_wiki"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("fork"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("forks_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("full_name"), knownvalue.StringExact("achim/user_test_repo_1")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_actions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_issues"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_packages"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_projects"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_pull_requests"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_releases"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("has_wiki"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("html_url"), knownvalue.StringExact("http://localhost:3000/achim/user_test_repo_1")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("ignore_whitespace_conflicts"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("internal"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("internal_tracker"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror_interval"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("mirror_updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("name"), knownvalue.StringExact("user_test_repo_1")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("open_issues_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("open_pr_counter"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("original_url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("parent_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("permissions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("private"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("release_counter"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("ssh_url"), knownvalue.StringExact("git@localhost:achim/user_test_repo_1.git")),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("stars_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("watchers_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000/")),
				},
			},
		},
	})
}
