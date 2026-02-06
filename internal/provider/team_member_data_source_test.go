package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamMemberDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	email    = "test_user@localhost.localdomain"
	password = "P@s$w0rd!"
}
data "forgejo_team_member" "test" {
	team_id = 1010
	user    = forgejo_user.test.login
}`,
				ExpectError: regexp.MustCompile("User test_user is not a member of team with ID 1010"),
			},
			// Read testing (non-existent user)
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
data "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = "non-existing-login"
}`,
				ExpectError: regexp.MustCompile("User non-existing-login is not a member of team with ID"),
			},
			// Read testing (non-existent member)
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
data "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = forgejo_user.test.login
}`,
				ExpectError: regexp.MustCompile("User test_user is not a member of team with ID "),
			},
			// Read testing
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
resource "forgejo_team_member" "test" {
	team_id = forgejo_team.test.id
	user    = forgejo_user.test.login
}
data "forgejo_team_member" "test" {
	team_id = forgejo_team_member.test.team_id
	user    = forgejo_team_member.test.user
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("data.forgejo_team_member.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("data.forgejo_team_member.test", tfjsonpath.New("user"), "forgejo_user.test", tfjsonpath.New("login"), compare.ValuesSame()),
				},
			},
		},
	})
}
