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

func TestAccTeamRepositoryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo"
}
resource "forgejo_team_repository" "test" {
	team_id       = 1010
	repository_id = forgejo_repository.test.id
}`,
				ExpectError: regexp.MustCompile("Team with ID 1010 not found"),
			},
			// Create and Read testing (non-existent repository)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test.id
	repository_id = 1010
}`,
				ExpectError: regexp.MustCompile("Repository with ID 1010 not found"),
			},
			// Create and Read testing (team includes all repositories)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test.id
	repository_id = forgejo_repository.test.id
}`,
				ExpectError: regexp.MustCompile(`Team\s+with\s+ID\s+[0-9]+\s+includes\s+all\s+repositories`),
			},
			// Create and Read testing (repository outside the organization)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo"
}
resource "forgejo_repository" "personal" {
	owner = "` + forgejoTestUser + `"
	name  = "personal_test_repo"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test.id
	repository_id = forgejo_repository.personal.id
}`,
				ExpectError: regexp.MustCompile(`belongs\s+to\s+organization\s+'team_repo_test_org',\s+but\s+repository\s+with\s+ID\s+[0-9]+\s+is\s+owned\s+by\s+'` + forgejoTestUser + `'`),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo"
}
resource "forgejo_repository" "test2" {
	owner = forgejo_organization.test.name
	name  = "second_test_repo"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test.id
	repository_id = forgejo_repository.test.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team_repository.test", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`^[0-9]+/[0-9]+$`))),
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile(`Expected\s+import\s+identifier\s+with\s+format:\s+'org/team/repo',\s+got:\s+'invalid'`),
			},
			// Import testing (non-existent org)
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "non-existent/test_team/test_repo",
				ExpectError:   regexp.MustCompile(`Organization\s+with\s+name\s+'non-existent'\s+not\s+found`),
			},
			// Import testing (non-existent team)
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "team_repo_test_org/non-existent/test_repo",
				ExpectError:   regexp.MustCompile(`Team\s+with\s+name\s+'non-existent'\s+not\s+found`),
			},
			// Import testing (non-existent repository)
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "team_repo_test_org/test_team/non-existent",
				ExpectError:   regexp.MustCompile(`Repository\s+with\s+owner\s+'team_repo_test_org'\s+and\s+name\s+'non-existent'\s+not\s+found`),
			},
			// Import testing (non-existent association)
			{
				ResourceName:  "forgejo_team_repository.test",
				ImportState:   true,
				ImportStateId: "team_repo_test_org/test_team/second_test_repo",
				ExpectError:   regexp.MustCompile(`Team\s+repository\s+with\s+organization\s+'team_repo_test_org',\s+team\s+'test_team'\s+and\s+repository\s+'second_test_repo'\s+not\s+found`),
			},
			// Import testing
			{
				ResourceName:      "forgejo_team_repository.test",
				ImportState:       true,
				ImportStateId:     "team_repo_test_org/test_team/test_repo",
				ImportStateVerify: true,
			},
			// Recreate and Read testing (updating any value recreates the resource -- repository)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo"
}
resource "forgejo_repository" "test2" {
	owner = forgejo_organization.test.name
	name  = "second_test_repo"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test.id
	repository_id = forgejo_repository.test2.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("repository_id"), "forgejo_repository.test2", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			// Recreate and Read testing (updating any value recreates the resource -- team)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_repo_test_org"
}
resource "forgejo_repository" "test2" {
	owner = forgejo_organization.test.name
	name  = "second_test_repo"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team" "test2" {
	name                      = "second_test_team"
	organization_id           = forgejo_organization.test.id
	includes_all_repositories = false
	permission                = "write"
	units_map                 = {
		"repo.code" = "write"
	}
}
resource "forgejo_team_repository" "test" {
	team_id       = forgejo_team.test2.id
	repository_id = forgejo_repository.test2.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_repository.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("team_id"), "forgejo_team.test2", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_repository.test", tfjsonpath.New("repository_id"), "forgejo_repository.test2", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
