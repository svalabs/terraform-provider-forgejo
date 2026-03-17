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
	name            = "tftest"
	organization_id = 1011
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("read"),
					})),
				},
			},
			// Not allowed to create the same team twice.
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	}
}
resource "forgejo_team" "test2" {
	# Make sure this second team is created later.
	name            = forgejo_team.test.name
	organization_id = forgejo_organization.test.id
	permission      = "write"
}`,
				ExpectError: regexp.MustCompile("team already exists"),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name                      = "test_team"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = false
	description               = "Updated test team."
	includes_all_repositories = true
	permission                = "write"
	units_map                 = {
		"repo.issues" = "write"
	}
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Update and Read testing (rename)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test" {
	name            = "renamed_test_team"
	organization_id = forgejo_organization.test.id
	description     = "Updated test team."
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Changing the parent organization recreates the resource.
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test" {
	name            = "renamed_test_team"
	organization_id = forgejo_organization.new_test.id
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Admin permission gives all units.
			{
				Config: providerConfig + `
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test" {
	name            = "renamed_test_team"
	organization_id = forgejo_organization.new_test.id
	permission      = "admin"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team")),
					statecheck.CompareValuePairs("forgejo_team.test", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code":       knownvalue.StringExact("admin"),
						"repo.issues":     knownvalue.StringExact("admin"),
						"repo.pulls":      knownvalue.StringExact("admin"),
						"repo.ext_issues": knownvalue.StringExact("read"),
						"repo.wiki":       knownvalue.StringExact("admin"),
						"repo.ext_wiki":   knownvalue.StringExact("read"),
						"repo.releases":   knownvalue.StringExact("admin"),
						"repo.projects":   knownvalue.StringExact("admin"),
						"repo.packages":   knownvalue.StringExact("admin"),
						"repo.actions":    knownvalue.StringExact("admin"),
					})),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
