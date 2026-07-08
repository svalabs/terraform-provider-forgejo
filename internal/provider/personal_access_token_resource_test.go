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

func TestAccPersonalAccessTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent user)
			{
				Config: providerConfig + `
resource "forgejo_personal_access_token" "test" {
	user_id = 1111
	name    = "tftest"
	scopes  = ["all"]
}`,
				ExpectError: regexp.MustCompile("user not found with id 1111"),
			},
			// Create and Read testing
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
	scopes  = ["all"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_personal_access_token.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token_last_eight"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("all"),
					})),
				},
			},
			// Duplicate token names are not allowed.
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
		"all"
	]
}
resource "forgejo_personal_access_token" "test1" {
	provider = forgejo.basicAuth

	user_id = forgejo_user.test.id
	name    = "tftest"
	scopes  = [
		"all"
	]
}`,
				ExpectError: regexp.MustCompile("Input validation error: tftest has been used as an application name already. Please use a new one."),
			},
			// Changing scope recreates the token.
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
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_personal_access_token.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token_last_eight"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("all"),
						knownvalue.StringExact("read:repository"),
					})),
				},
			},
			// Changing name recreates the token.
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
	name    = "tftest1"
	scopes  = [
		"all",
		"read:repository"
	]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_personal_access_token.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("token_last_eight"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_personal_access_token.test", tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("all"),
						knownvalue.StringExact("read:repository"),
					})),
				},
			},
		},
	})
}
