package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationActionSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org"
}
resource "forgejo_organization_action_secret" "test" {
	organization = forgejo_organization.test.name
	name         = "my_secret"
	data         = "my_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
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
	organization = forgejo_organization.test.name
	name         = "my_new_secret"
	data         = "my_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
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
	organization = forgejo_organization.test.name
	name         = "my_new_secret"
	data         = "my_new_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("organization"), knownvalue.StringExact("test_org")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_organization_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_organization_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
