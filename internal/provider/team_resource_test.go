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

func TestAccTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent org)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "tftest"
	organization = "test_org"
}`,
				ExpectError: regexp.MustCompile("Organization with name 'test_org' not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
	permission   = "read"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
				},
			},
			// Create and Read testing (import if exists)
			{
				Config: providerConfig + `
data "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
}
resource "forgejo_team" "test2" {
	name         = "test_team"
	organization = "tftest"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection("forgejo_team.test2", []tfjsonpath.Path{tfjsonpath.New("id")}, "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
	permission   = "write"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
				},
			},
			// Update and Read testing (rename)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "test_team_2"
	organization = "tftest"
	permission   = "write"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team_2")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
