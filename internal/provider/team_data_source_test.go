package provider_test

import (
	"regexp"
	"testing"

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
	name         = "tftest"
	organization = "test_org"
}`,
				ExpectError: regexp.MustCompile("Organization with name 'test_org' not found"),
			},
			// Read testing (non-existent team)
			{
				Config: providerConfig + `
data "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
}`,
				ExpectError: regexp.MustCompile("Team with name 'test_team' not found"),
			},
			// Create and Read testing (import if exists)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name                     = "test_team"
	organization             = "tftest"
	can_create_org_repo      = true
	include_all_repositories = true
	permission               = "read"
	units                    = [ "repo.code" ]
}
data "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("include_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
		},
	})
}
