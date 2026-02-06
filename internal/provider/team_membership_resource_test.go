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

func TestAccTeamMembershipResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_membership" "test" {
	team_id = 1010
	user_id = forgejo_user.test.id
}`,
				ExpectError: regexp.MustCompile("Either User test_user or Team with ID 1010 not found"),
			},
			// Create testing (non-existent user)
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
	units                     = ["repo.code"]
}
resource "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = 1011
}`,
				ExpectError: regexp.MustCompile("user not found with id 1011"),
			},
			// Create testing
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
	units                     = ["repo.code"]
}
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = forgejo_user.test.id
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("user_id"), "forgejo_user.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team_membership.test", tfjsonpath.New("user_name"), knownvalue.StringExact("test_user")),
				},
			},
			// Update testing (updating any value recreates the resource -- user)
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
	units                     = ["repo.code"]
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
resource "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = forgejo_user.test2.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_membership.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("user_id"), "forgejo_user.test2", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team_membership.test", tfjsonpath.New("user_name"), knownvalue.StringExact("second_test_user")),
				},
			},
			// Update testing (updating any value recreates the resource -- team)
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
	units                     = ["repo.code"]
}
resource "forgejo_team" "test2" {
	name                      = "second_test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units                     = ["repo.code"]
}
resource "forgejo_user" "test2" {
	login    = "second_test_user"
	email    = "second_test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
resource "forgejo_team_membership" "test" {
	team_id = forgejo_team.test2.id
	user_id = forgejo_user.test2.id
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_membership.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("team_id"), "forgejo_team.test2", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("forgejo_team_membership.test", tfjsonpath.New("user_id"), "forgejo_user.test2", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team_membership.test", tfjsonpath.New("user_name"), knownvalue.StringExact("second_test_user")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
