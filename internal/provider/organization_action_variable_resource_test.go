package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationActionVariableResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent organization)
			{
				Config: providerConfig + `
resource "forgejo_organization_action_variable" "test" {
	organization_id = 1011
	name            = "my_variable"
	data            = "my_variable_value"
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_variable"
	data            = "my_variable_value"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_variable")),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
				},
			},
			// Update and Read testing (rename variable)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_variable"
	data            = "my_variable_value"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_variable")),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
				},
			},
			// Update and Read testing (update value)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_variable"
	data            = "my_new_variable_value"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_variable")),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_new_variable_value")),
				},
			},
			// Update and Read testing (long name)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_new_variable_with_a_very_long_name_that_is_over_30_characters_long"
	data            = "my_new_variable_value"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_variable_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectKnownValue("forgejo_organization_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_new_variable_value")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
