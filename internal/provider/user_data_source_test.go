package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent resource)
			{
				Config: providerConfig + `
data "forgejo_user" "test" {
	login = "non_existent"
}`,
				ExpectError: regexp.MustCompile("User with name \"non_existent\" not found"),
			},
			// Read testing
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
data "forgejo_user" "test" {
	login = forgejo_user.test.login
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("active"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("description"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("email"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("followers_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("following_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("full_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("is_admin"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("language"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("last_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("login_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("login"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("prohibit_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("restricted"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("source_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("starred_repos_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("visibility"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_user.test", tfjsonpath.New("website"), knownvalue.NotNull()),
				},
			},
		},
	})
}
