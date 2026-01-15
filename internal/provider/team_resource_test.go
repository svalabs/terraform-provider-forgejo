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

func TestAccTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent org)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name            = "tftest"
	organization_id = 1011
}`,
				ExpectError: regexp.MustCompile("no Organization with id '1011' was found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units                     = ["repo.issues"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.issues"),
					})),
				},
			},
			// Not allowed to create the same team twice.
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units                     = ["repo.issues"]
}
resource "forgejo_team" "test2" {
	# Make sure this second team is created later.
	name            = forgejo_team.test.name
	organization_id = forgejo_organization.test.id
	permission      = "write"
	units           = ["repo.code"]
}`,
				ExpectError: regexp.MustCompile("team already exists"),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = false
	description               = "Updated test team."
	includes_all_repositories = true
	permission                = "write"
	units                     = ["repo.issues"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.issues"),
					})),
				},
			},
			// Update and Read testing (rename)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "renamed_test_team"
	organization_id           = forgejo_organization.test.id
	description               = "Updated test team."
	permission                = "write"
	units                     = ["repo.pulls"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.pulls"),
					})),
				},
			},
			// Changing the parent organization recreates the resource.
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test" {
	name                      = "renamed_test_team"
	organization_id           = forgejo_organization.new_test.id
	permission                = "write"
	units                     = ["repo.pulls"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.pulls"),
					})),
				},
			},
			// Admin permission needs all units.
			{
				Config: providerConfig + `
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test" {
	name            = "renamed_test_team"
	organization_id = forgejo_organization.new_test.id
	permission      = "admin"
	units           = ["repo.code", "repo.issues", "repo.pulls", "repo.ext_issues", "repo.wiki", "repo.ext_wiki", "repo.releases", "repo.projects", "repo.packages", "repo.actions"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
						knownvalue.StringExact("repo.issues"),
						knownvalue.StringExact("repo.pulls"),
						knownvalue.StringExact("repo.ext_issues"),
						knownvalue.StringExact("repo.wiki"),
						knownvalue.StringExact("repo.ext_wiki"),
						knownvalue.StringExact("repo.releases"),
						knownvalue.StringExact("repo.projects"),
						knownvalue.StringExact("repo.packages"),
						knownvalue.StringExact("repo.actions"),
					})),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
