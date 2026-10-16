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
	user   = "non_existing_user"
	name   = "tftest"
	scopes = ["all"]
}`,
				ExpectError: regexp.MustCompile("Personal access token for user \"non_existing_user\" not found"),
			},
			// Create and Read testing
			{
				Config: providerBasicAuthConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
resource "forgejo_personal_access_token" "test" {
	user   = forgejo_user.test.login
	name   = "tftest"
	scopes = ["all"]
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
				Config: providerBasicAuthConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
resource "forgejo_personal_access_token" "test" {
	user   = forgejo_user.test.login
	name   = "tftest"
	scopes = ["all"]
}
resource "forgejo_personal_access_token" "test1" {
	user   = forgejo_user.test.login
	name   = "tftest"
	scopes = ["all"]
}`,
				ExpectError: regexp.MustCompile("access token name has been used already"),
			},
			// Changing scope recreates the token.
			{
				Config: providerBasicAuthConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
resource "forgejo_personal_access_token" "test" {
	user   = forgejo_user.test.login
	name   = "tftest"
	scopes = [
		"read:organization",
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
						knownvalue.StringExact("read:organization"),
						knownvalue.StringExact("read:repository"),
					})),
				},
			},
			// Changing name recreates the token.
			{
				Config: providerBasicAuthConfig + `
resource "forgejo_user" "test" {
	login    = "test_user"
	password = "password"
	email    = "test_user@example.com"
}
resource "forgejo_personal_access_token" "test" {
	user   = forgejo_user.test.login
	name   = "tftest1"
	scopes = [
		"read:organization",
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
						knownvalue.StringExact("read:organization"),
						knownvalue.StringExact("read:repository"),
					})),
				},
			},
		},
	})
}

func TestAccPersonalAccessTokenRequiresBasicAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: providerConfig + `
resource "forgejo_personal_access_token" "token_auth" {
	user   = "tfadmin"
	name   = "test_token_auth_not_allowed"
	scopes = ["read:repository"]
}
`,
			ExpectError: regexp.MustCompile("Authentication method not allowed, use basic-auth"),
		}},
	})
}
