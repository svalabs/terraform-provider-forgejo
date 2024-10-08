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
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("avatar_url"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"),
						),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(""),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("full_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("location"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("tftest"),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("visibility"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^(public)|(limited)|(private)$"),
						),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("website"),
						knownvalue.NotNull(),
					),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name        = "tftest"
	description = "Purely for testing..."
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("avatar_url"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"),
						),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Purely for testing..."),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("full_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("location"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("tftest"),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("visibility"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^(public)|(limited)|(private)$"),
						),
					),
					statecheck.ExpectKnownValue(
						"forgejo_organization.test",
						tfjsonpath.New("website"),
						knownvalue.NotNull(),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
