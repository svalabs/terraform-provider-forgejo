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

func TestAccTeamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent org)
			{
				Config: providerConfig + `
data "forgejo_team" "test" {
	name            = "tftest"
	organization_id = 1011
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Read testing (non-existent team)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
data "forgejo_team" "test" {
	name            = "test_team"
	organization_id = forgejo_organization.test.id
}`,
				ExpectError: regexp.MustCompile("Team with name \"test_team\" not found"),
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
data "forgejo_team" "test" {
	name            = forgejo_team.test.name
	organization_id = forgejo_team.test.organization_id
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.CompareValuePairs("data.forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("data.forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("data.forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
		},
	})
}
