package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamRepositoryResource(t *testing.T) {
	fixtures := providerConfig + `
resource "forgejo_organization" "test" {
	name = "assignment_import_org"
}
resource "forgejo_repository" "first" {
	owner   = forgejo_organization.test.name
	name    = "first"
	private = true
}
resource "forgejo_repository" "second" {
	owner   = forgejo_organization.test.name
	name    = "second"
	private = true
}
resource "forgejo_team" "first" {
	name                      = "developers"
	organization              = forgejo_organization.test.name
	includes_all_repositories = false
	units_map                 = { "repo.code" = "read" }
}
resource "forgejo_team" "second" {
	name                      = "reviewers"
	organization              = forgejo_organization.test.name
	includes_all_repositories = false
	units_map                 = { "repo.code" = "read" }
	depends_on                = [forgejo_team.first]
}`
	config := func(team, repo string) string {
		return fixtures + `
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.` + team + `.id
	repository_id = forgejo_repository.` + repo + `.id
}`
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("first", "first"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("team_id"), "forgejo_team.first", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("repository_id"), "forgejo_repository.first", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile("Unable to parse import identifier"),
			},
			{
				ResourceName:      "forgejo_team_repository.test",
				ImportState:       true,
				ImportStateId:     "assignment_import_org/developers/first",
				ImportStateVerify: true,
			},
			{
				Config: config("first", "first"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: config("first", "second"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				Config: config("second", "second"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				Config: fixtures,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionDestroy),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository.second", tfjsonpath.New("owner"), "forgejo_organization.test", tfjsonpath.New("name"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team.second", tfjsonpath.New("organization"), "forgejo_organization.test", tfjsonpath.New("name"), compare.ValuesSame()),
				},
			},
		},
	})
}
