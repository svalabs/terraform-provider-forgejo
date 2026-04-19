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

func TestAccOrganizationActionVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent org by ID)
			{
				Config: providerConfig + `
data "forgejo_organization_action_variable" "test" {
	organization_id = 1011
	name            = "my_variable"
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Read testing (non-existent org by name)
			{
				Config: providerConfig + `
data "forgejo_organization_action_variable" "test" {
	organization = "non-existent"
	name         = "my_variable"
}`,
				ExpectError: regexp.MustCompile("Action variable with organization \"non-existent\" and name \"my_variable\" not found"),
			},
			// Read testing (non-existent variable by ID)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
data "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_variable"
}`,
				ExpectError: regexp.MustCompile(`Action variable with organization "test_org" and name "my_variable" not found`),
			},
			// Read testing (non-existent variable by name)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
data "forgejo_organization_action_variable" "test" {
	organization = forgejo_organization.test.name
	name         = "my_variable"
}`,
				ExpectError: regexp.MustCompile(`Action variable with organization "test_org" and name "my_variable" not found`),
			},
			// Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_variable"
	data            = "my_variable_value"
}
data "forgejo_organization_action_variable" "test_by_id" {
	organization_id = forgejo_organization.test.id
	name            = forgejo_organization_action_variable.test.name
}
data "forgejo_organization_action_variable" "test_by_name" {
	organization = forgejo_organization.test.name
	name         = forgejo_organization_action_variable.test.name
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("data.forgejo_organization_action_variable.test_by_id", plancheck.ResourceActionRead),
						plancheck.ExpectResourceAction("data.forgejo_organization_action_variable.test_by_name", plancheck.ResourceActionRead),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_id", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.CompareValuePairs("data.forgejo_organization_action_variable.test_by_id", tfjsonpath.New("organization_id"), "forgejo_organization.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_id", tfjsonpath.New("name"), knownvalue.StringExact("my_variable")),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_id", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_name", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_name", tfjsonpath.New("organization_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_name", tfjsonpath.New("name"), knownvalue.StringExact("my_variable")),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test_by_name", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
				},
			},
		},
	})
}
