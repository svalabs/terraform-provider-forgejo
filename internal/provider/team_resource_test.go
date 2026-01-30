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
	units = ["repo.code"]
}`,
				ExpectError: regexp.MustCompile("Organization with name 'test_org' not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
	can_create_org_repo = true
	description = "Test team."
	includes_all_repositories = false
	permission   = "read"
	units = ["repo.code"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
			// Not allowed to create the same team twice.
			{
				Config: providerConfig + `
resource "forgejo_team" "test2" {
	name         = "test_team"
	organization = "tftest"
	units = ["repo.code"]
}`,
				ExpectError: regexp.MustCompile("Team 'test_team' already exists"),
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
	import_if_exists = true
	permission   = "write"
	units = ["repo.code"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection("forgejo_team.test2", []tfjsonpath.Path{tfjsonpath.New("id")}, "forgejo_team.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test2", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "test_team"
	organization = "tftest"
	can_create_org_repo = true
	description = "Test team."
	includes_all_repositories = false
	permission   = "admin"
	units = ["repo.code"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
			// Update and Read testing (rename)
			{
				Config: providerConfig + `
data "forgejo_team" "test_data" {
	name         = "test_team"
	organization = "tftest"
}
resource "forgejo_team" "test" {
	name         = "renamed_test_team"
	organization = "tftest"
	units = ["repo.code"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection("forgejo_team.test", []tfjsonpath.Path{tfjsonpath.New("id")}, "forgejo_team.test_data", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
			// Changing the parent organization recreates the resource.
			{
				Config: providerConfig + `
resource "forgejo_organization" "test_org" {
	name         = "test-org"
}
resource "forgejo_team" "test" {
	name         = "renamed_test_team"
	organization = "test-org"
	units = ["repo.code"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("organization"), knownvalue.StringExact("test-org")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("none")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("repo.code"),
					})),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
