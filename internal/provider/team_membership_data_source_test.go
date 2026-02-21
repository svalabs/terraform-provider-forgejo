package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamMembershipDataSource(t *testing.T) {
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
data "forgejo_team_membership" "test" {
	team_id = 1010
	user_id = forgejo_user.test.id
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
data "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = 1011
}`,
				ExpectError: regexp.MustCompile("user not found with id 1011"),
			},
			// Read testing (non-existent membership)
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
data "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = forgejo_user.test.id
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
resource "forgejo_team_membership" "test" {
	team_id = forgejo_team.test.id
	user_id = forgejo_user.test.id
}
data "forgejo_team_membership" "test" {
	team_id = forgejo_team_membership.test.team_id
	user_id = forgejo_team_membership.test.user_id
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("data.forgejo_team_membership.test", tfjsonpath.New("team_id"), "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs("data.forgejo_team_membership.test", tfjsonpath.New("user_id"), "forgejo_user.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("data.forgejo_team_membership.test", tfjsonpath.New("user_name"), knownvalue.StringExact("test_user")),
				},
			},
		},
	})
}
