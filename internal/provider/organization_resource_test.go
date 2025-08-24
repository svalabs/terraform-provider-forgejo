package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "tftest"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("location"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "tftest1"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("location"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("website"), knownvalue.StringExact(""))},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name        = "tftest1"
	description = "Purely for testing..."
	location    = "Mêlée Island"
	visibility  = "limited"
	website     = "http://localhost:3000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("location"), knownvalue.StringExact("Mêlée Island")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("visibility"), knownvalue.StringExact("limited")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000"))},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name        = "tftest1"
	description = "Purely for testing..."
	location    = "Mêlée Island"
	visibility  = "private"
	website     = "http://localhost:3000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("location"), knownvalue.StringExact("Mêlée Island")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("visibility"), knownvalue.StringExact("private")),
					statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000"))},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
