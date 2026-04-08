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

func TestAccOrganizationActionVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent organization)
			{
				Config: providerConfig + `
data "forgejo_organization_action_variable" "test" {
	organization_id = 1011
	name            = "my_variable"
}`,
				ExpectError: regexp.MustCompile("Organization with ID 1011 not found"),
			},
			// Read testing (non-existent variable)
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
data "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = "my_variable"
}`,
				ExpectError: regexp.MustCompile("Action variable with org \"test_org\" and name \"my_variable\" not found"),
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
data "forgejo_organization_action_variable" "test" {
	organization_id = forgejo_organization.test.id
	name            = forgejo_organization_action_variable.test.name
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("data.forgejo_organization_action_variable.test", plancheck.ResourceActionRead),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test", tfjsonpath.New("organization_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_variable")),
					statecheck.ExpectKnownValue("data.forgejo_organization_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
				},
			},
		},
	})
}
