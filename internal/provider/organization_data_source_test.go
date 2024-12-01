package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "forgejo_organization" "test" {
  name = "test1"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("avatar_url"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"),
						),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("full_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("location"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test1"),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("visibility"),
						knownvalue.StringRegexp(
							regexp.MustCompile("^(public)|(limited)|(private)$"),
						),
					),
					statecheck.ExpectKnownValue(
						"data.forgejo_organization.test",
						tfjsonpath.New("website"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
