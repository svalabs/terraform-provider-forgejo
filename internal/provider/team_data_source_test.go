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

func TestAccTeamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent org by ID)
			{
				Config: providerConfig + `
data "forgejo_team" "test" {
	name            = "tftest"
	organization_id = 1011
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Read testing (non-existent org by name)
			{
				Config: providerConfig + `
data "forgejo_team" "test" {
	name         = "tftest"
	organization = "non-existent"
}`,
				ExpectError: regexp.MustCompile("Organization with name 'non-existent' not found"),
			},
			// Read testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
data "forgejo_team" "test_by_id" {
	name            = "non-existent"
	organization_id = forgejo_organization.test.id
}
data "forgejo_team" "test_by_name" {
	name         = "non-existent"
	organization = forgejo_organization.test.name
}`,
				ExpectError: regexp.MustCompile("Team with name 'non-existent' not found"),
			},
			// Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_team" "test_by_id" {
	name                      = "test_team_by_id"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
data "forgejo_team" "test_by_id" {
	name            = forgejo_team.test_by_id.name
	organization_id = forgejo_team.test_by_id.organization_id
}
resource "forgejo_team" "test_by_name" {
	name                      = "test_team_by_name"
	organization              = forgejo_organization.test.name
	can_create_org_repo       = true
	includes_all_repositories = true
	permission                = "read"
	units_map                 = {
		"repo.code" = "read"
	}
}
data "forgejo_team" "test_by_name" {
	name         = forgejo_team.test_by_name.name
	organization = forgejo_team.test_by_name.organization
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("data.forgejo_team.test_by_id", plancheck.ResourceActionRead),
						plancheck.ExpectResourceAction("data.forgejo_team.test_by_name", plancheck.ResourceActionRead),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_id")),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("data.forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code": knownvalue.StringExact("read"),
					})),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_name")),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					// organization_id not populated when fetching team by name
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("data.forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code": knownvalue.StringExact("read"),
					})),
				},
			},
		},
	})
}
