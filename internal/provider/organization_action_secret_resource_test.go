package provider_test

import (
	"regexp"
	"testing"

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
			// Create and Read testing (non-existent organization)
			{
				Config: providerConfig + `
resource "forgejo_organization_action_secret" "test" {
	organization_id = 1011
	name            = "my_secret"
	data            = "my_secret_value"
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_secret"
	data            = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_secret")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret"
	data            = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret"
	data            = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Recreate and Read testing (long name)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_secret_with_a_very_long_name_that_is_over_30_characters_long"
	data            = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_organization_action_secret.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
