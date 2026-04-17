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

func TestAccTeamMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_member" "test" {
	team_id = 1010
	user    = forgejo_user.test.login
}`,
				ExpectError: regexp.MustCompile("Either user 'test_user' or team with ID 1010 not found"),
			},
			// Create and Read testing (non-existent user)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = "non-existing-user"
}`,
				ExpectError: regexp.MustCompile("Either user 'non-existing-user' or team with ID [0-9]+ not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = forgejo_user.test.login
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_member.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("user"), "forgejo_user.test", tfjsonpath.New("login"), compare.ValuesSame()),
				},
			},
			// Recreate and Read testing (updating any value recreates the resource -- user)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_user" "test2" {
	login    = "second_test_user"
	email    = "second_test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = forgejo_user.test2.login
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_member.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("user"), "forgejo_user.test2", tfjsonpath.New("login"), compare.ValuesSame()),
				},
			},
			// Recreate and Read testing (updating any value recreates the resource -- team)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_team" "test2" {
	name                      = "second_test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
resource "forgejo_user" "test2" {
	login    = "second_test_user"
	email    = "second_test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_member" "test" {
	team_id = forgejo_team.test2.id
	user    = forgejo_user.test2.login
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_member.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("team_id"), "forgejo_team.test2", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_member.test", tfjsonpath.New("user"), "forgejo_user.test2", tfjsonpath.New("login"), compare.ValuesSame()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
