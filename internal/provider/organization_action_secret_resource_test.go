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

func TestAccOrganizationActionSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent org by ID)
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization_id = 1011
	name            = "my_secret"
	data            = "my_secret_value"
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Create and Read testing (non-existent org by name)
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization = "non-existent"
	name         = "my_secret"
	data         = "my_secret_value"
}`,
				ExpectError: regexp.MustCompile("Organization with name \"non-existent\" not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_secret_by_id"
	data            = "my_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_secret_by_name"
	data         = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_secret_by_id")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret_by_id"
	data            = "my_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_new_secret_by_name"
	data         = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_id")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret_by_id"
	data            = "my_new_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_new_secret_by_name"
	data         = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_id")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Recreate and Read testing (long name)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret_by_id_with_a_very_long_name_that_is_over_30_characters_long"
	data            = "my_new_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long"
	data         = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_id_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Changing the parent organization recreates the resource
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization" "new_test" {
	name = "new_test_org"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.new_test.id
	name            = "my_new_secret_by_id_with_a_very_long_name_that_is_over_30_characters_long"
	data            = "my_new_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.new_test.name
	name         = "my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long"
	data         = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("new_test_org")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.new_test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_id_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("new_test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
