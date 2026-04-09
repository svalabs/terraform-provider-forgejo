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
			// Create and Read testing (non-existent org by ID)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name            = "tftest"
	organization_id = 1011
	units_map       = {
		"repo.issues" = "read"
	}
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Create and Read testing (non-existent org by name)
			{
				Config: providerConfig + `
resource "forgejo_team" "test" {
	name         = "tftest"
	organization = "non-existent"
	units_map    = {
		"repo.issues" = "read"
	}
}`,
				ExpectError: regexp.MustCompile("Organization with name 'non-existent' not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test_by_id" {
	name                      = "test_team_by_id"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	}
}
resource "forgejo_team" "test_by_name" {
	name                      = "test_team_by_name"
	organization              = forgejo_organization.test.name
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_id")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("read"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_name")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("description"), knownvalue.StringExact("Test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("read"),
					})),
				},
			},
			// Not allowed to create the same team twice
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test_by_id" {
	name                      = "test_team_by_id"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	}
}
resource "forgejo_team" "test_by_id2" {
	# Make sure this second team is created later.
	name            = forgejo_team.test_by_id.name
	organization_id = forgejo_organization.test.id
	permission      = "write"
	units_map                 = {
		"repo.issues" = "read"
	}
}
resource "forgejo_team" "test_by_name" {
	name                      = "test_team_by_name"
	organization              = forgejo_organization.test.name
	can_create_org_repo       = true
	description               = "Test team."
	includes_all_repositories = false
	permission                = "read"
	units_map                 = {
		"repo.issues" = "read"
	}
}
resource "forgejo_team" "test_by_name2" {
	# Make sure this second team is created later.
	name            = forgejo_team.test_by_name.name
	organization    = forgejo_organization.test.name
	permission      = "write"
	units_map                 = {
		"repo.issues" = "read"
	}
}`,
				ExpectError: regexp.MustCompile("Input validation error: team already exists"),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test_by_id" {
	name                      = "test_team_by_id"
	organization_id           = forgejo_organization.test.id
	can_create_org_repo       = false
	description               = "Updated test team."
	includes_all_repositories = true
	permission                = "write"
	units_map                 = {
		"repo.issues" = "write"
	}
}
resource "forgejo_team" "test_by_name" {
	name                      = "test_team_by_name"
	organization              = forgejo_organization.test.name
	can_create_org_repo       = false
	description               = "Updated test team."
	includes_all_repositories = true
	permission                = "write"
	units_map                 = {
		"repo.issues" = "write"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_id")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("test_team_by_name")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test_by_id" {
	name            = "renamed_test_team_by_id"
	organization_id = forgejo_organization.test.id
	description     = "Updated test team."
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}
resource "forgejo_team" "test_by_name" {
	name            = "renamed_test_team_by_name"
	organization_id = forgejo_organization.test.id
	description     = "Updated test team."
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_id")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_name")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("description"), knownvalue.StringExact("Updated test team.")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Changing the parent organization recreates the resource
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test_by_id" {
	name            = "renamed_test_team_by_id"
	organization_id = forgejo_organization.new_test.id
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}
resource "forgejo_team" "test_by_name" {
	name            = "renamed_test_team_by_name"
	organization_id = forgejo_organization.new_test.id
	permission      = "write"
	units_map       = {
		"repo.issues" = "write"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_id")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("new_test")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_name")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("new_test")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.issues": knownvalue.StringExact("write"),
					})),
				},
			},
			// Admin permission gives all units
			{
				Config: providerConfig + `
resource "forgejo_organization" "new_test" {
	name = "new_test"
}
resource "forgejo_team" "test_by_id" {
	name            = "renamed_test_team_by_id"
	organization_id = forgejo_organization.new_test.id
	permission      = "admin"
	units_map       = {
		"repo.code"       = "admin"
		"repo.issues"     = "admin"
		"repo.pulls"      = "admin"
		"repo.ext_issues" = "admin"
		"repo.wiki"       = "admin"
		"repo.ext_wiki"   = "admin"
		"repo.releases"   = "admin"
		"repo.projects"   = "admin"
		"repo.packages"   = "admin"
		"repo.actions"    = "admin"
	}
}
resource "forgejo_team" "test_by_name" {
	name            = "renamed_test_team_by_name"
	organization_id = forgejo_organization.new_test.id
	permission      = "admin"
	units_map       = {
		"repo.code"       = "admin"
		"repo.issues"     = "admin"
		"repo.pulls"      = "admin"
		"repo.ext_issues" = "admin"
		"repo.wiki"       = "admin"
		"repo.ext_wiki"   = "admin"
		"repo.releases"   = "admin"
		"repo.projects"   = "admin"
		"repo.packages"   = "admin"
		"repo.actions"    = "admin"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_id")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("new_test")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code":       knownvalue.StringExact("admin"),
						"repo.issues":     knownvalue.StringExact("admin"),
						"repo.pulls":      knownvalue.StringExact("admin"),
						"repo.ext_issues": knownvalue.StringExact("admin"),
						"repo.wiki":       knownvalue.StringExact("admin"),
						"repo.ext_wiki":   knownvalue.StringExact("admin"),
						"repo.releases":   knownvalue.StringExact("admin"),
						"repo.projects":   knownvalue.StringExact("admin"),
						"repo.packages":   knownvalue.StringExact("admin"),
						"repo.actions":    knownvalue.StringExact("admin"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("renamed_test_team_by_name")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("new_test")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code":       knownvalue.StringExact("admin"),
						"repo.issues":     knownvalue.StringExact("admin"),
						"repo.pulls":      knownvalue.StringExact("admin"),
						"repo.ext_issues": knownvalue.StringExact("admin"),
						"repo.wiki":       knownvalue.StringExact("admin"),
						"repo.ext_wiki":   knownvalue.StringExact("admin"),
						"repo.releases":   knownvalue.StringExact("admin"),
						"repo.projects":   knownvalue.StringExact("admin"),
						"repo.packages":   knownvalue.StringExact("admin"),
						"repo.actions":    knownvalue.StringExact("admin"),
					})),
				},
			},
			// Create and Read testing (many teams)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "team_test_org"
}
resource "forgejo_team" "test_by_id" {
  count = 100

	organization_id = forgejo_organization.test.id
	name            = "test_team_${count.index}_by_id"

	units_map = {
		"repo.code" = "read"
	}
}
resource "forgejo_team" "test_by_name" {
  count = 100

	organization_id = forgejo_organization.test.id
	name            = "test_team_${count.index}_by_name"

	units_map = {
		"repo.code" = "read"
	}
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team.test_by_id[99]", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("forgejo_team.test_by_name[99]", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("name"), knownvalue.StringRegexp(regexp.MustCompile("test_team_[0-9]+_by_id"))),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_id[99]", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_id[99]", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code": knownvalue.StringExact("read"),
					})),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("name"), knownvalue.StringRegexp(regexp.MustCompile("test_team_[0-9]+_by_name"))),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("organization"), knownvalue.StringExact("team_test_org")),
					statecheck.CompareValuePairs("forgejo_team.test_by_name[99]", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("can_create_org_repo"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("includes_all_repositories"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
					statecheck.ExpectKnownValue("forgejo_team.test_by_name[99]", tfjsonpath.New("units_map"), knownvalue.MapExact(map[string]knownvalue.Check{
						"repo.code": knownvalue.StringExact("read"),
					})),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
