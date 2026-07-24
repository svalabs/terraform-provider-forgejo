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

func TestAccPersonalAccessTokenDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent user)
			{
				Config: providerConfig + `
data "forgejo_personal_access_token" "test" {
	user_id = 1111
	name    = "tftest"
}`,
				ExpectError: regexp.MustCompile("user not found with id 1111"),
			},
			// Read testing (non-existent resource)
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
data "forgejo_personal_access_token" "test" {
	user_id = forgejo_user.test.id
	name    = "non_existent"
}`,
				ExpectError: regexp.MustCompile("Personal access token from user \"test_user\" and name non_existent not found"),
			},
			// Read testing
			{
				Config: providerConfig + `
variable "FORGEJO_BASIC_AUTH_USERNAME" {
	type = string
}
variable "FORGEJO_BASIC_AUTH_PASSWORD" {
	type = string
}
provider "forgejo" {
	alias    = "basicAuth"
	host     = "` + forgejoTestHost + `"
	username = var.FORGEJO_BASIC_AUTH_USERNAME
	password = var.FORGEJO_BASIC_AUTH_PASSWORD
}
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
resource "forgejo_personal_access_token" "test" {
	provider = forgejo.basicAuth

	user_id = forgejo_user.test.id
	name    = "tftest"
	scopes  = [
		"all",
		"read:repository"
	]
}
data "forgejo_personal_access_token" "test" {
	user_id = forgejo_user.test.id
	name    = forgejo_personal_access_token.test.name
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("data.forgejo_personal_access_token.test", plancheck.ResourceActionRead),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_personal_access_token.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_personal_access_token.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_personal_access_token.test", tfjsonpath.New("token_last_eight"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_personal_access_token.test", tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("all"),
						knownvalue.StringExact("read:repository"),
					})),
				},
			},
		},
	})
}
