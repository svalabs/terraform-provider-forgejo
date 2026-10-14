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
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("new_test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_by_name_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccOrganizationActionSecretValidationDataConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check data and data_wo are mutually exclusive
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization = "test_org_wo"
	name         = "my_wo_secret"
	data         = "my_secret_value"
	data_wo      = "my_wo_secret_value"
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of`),
			},
			// Check one of data or data_wo is required
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization = "test_org_wo"
	name         = "my_wo_secret"
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of`),
			},
			// Check data_wo_version requires data_wo
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization    = "test_org_wo"
	name            = "my_wo_secret"
	data            = "my_secret_value"
	data_wo_version = 1
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute "data_wo" must be specified when`),
			},
		},
	})
}

func TestAccOrganizationActionSecretResourceWriteOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (write-only data without version)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data_wo         = "my_wo_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_wo_secret_by_name"
	data_wo      = "my_wo_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (adding data_wo_version)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data_wo         = "my_wo_value"
	data_wo_version = 1
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization    = forgejo_organization.test.name
	name            = "my_wo_secret_by_name"
	data_wo         = "my_wo_value"
	data_wo_version = 1
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Changing the write-only value alone is not detected (by design)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 1
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization    = forgejo_organization.test.name
	name            = "my_wo_secret_by_name"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 1
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (bumping data_wo_version)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 2
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization    = forgejo_organization.test.name
	name            = "my_wo_secret_by_name"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 2
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (switching to data)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data            = "my_secret_value"
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = "my_wo_secret_by_name"
	data         = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (switching back to write-only data)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}
resource "forgejo_organization_action_secret" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = "my_wo_secret_by_id"
	data_wo         = "my_wo_value"
	data_wo_version = 3
}
resource "forgejo_organization_action_secret" "test_by_name" {
	organization    = forgejo_organization.test.name
	name            = "my_wo_secret_by_name"
	data_wo         = "my_wo_value"
	data_wo_version = 3
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_id", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test_by_name", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.CompareValuePairs("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_id")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_id", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org_wo")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret_by_name")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test_by_name", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
