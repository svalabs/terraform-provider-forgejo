package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
  password = "passw0rd"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("active"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("admin"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("email"), knownvalue.StringRegexp(regexp.MustCompile(`^[0-9a-z]+@[a-z]+\.[a-z]+$`))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("followers_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("following_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("language"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("last_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("must_change_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("prohibit_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("restricted"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("send_notify"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("source_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("starred_repos_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
					statecheck.ExpectSensitiveValue("forgejo_user.test", tfjsonpath.New("password")),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "tftest1"
	email    = "tftest1@localhost.localdomain"
  password = "passw1rd"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("active"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("admin"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("email"), knownvalue.StringRegexp(regexp.MustCompile(`^[0-9a-z]+@[a-z]+\.[a-z]+$`))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("followers_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("following_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("language"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("last_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("must_change_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("prohibit_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("restricted"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("send_notify"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("source_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("starred_repos_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
					statecheck.ExpectSensitiveValue("forgejo_user.test", tfjsonpath.New("password")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login                = "tftest1"
	email                = "tftest1@localhost.localdomain"
  password             = "passw1rd"
	visibility           = "limited"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("active"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("admin"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("email"), knownvalue.StringRegexp(regexp.MustCompile(`^[0-9a-z]+@[a-z]+\.[a-z]+$`))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("followers_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("following_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("language"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("last_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("must_change_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("prohibit_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("restricted"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("send_notify"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("source_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("starred_repos_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("visibility"), knownvalue.StringExact("limited")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("website"), knownvalue.StringExact("")),
					statecheck.ExpectSensitiveValue("forgejo_user.test", tfjsonpath.New("password")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login                = "tftest1"
	email                = "tftest1@localhost.localdomain"
  password             = "passw1rd"
	active               = false
	admin                = true
	description          = "Purely for testing..."
	location             = "Mêlée Island"
	must_change_password = false
	prohibit_login       = true
	restricted           = true
	send_notify          = false
	visibility           = "private"
	website              = "http://localhost:3000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("active"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("admin"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("avatar_url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/avatars/[0-9a-z]{32}$"))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("description"), knownvalue.StringExact("Purely for testing...")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("email"), knownvalue.StringRegexp(regexp.MustCompile(`^[0-9a-z]+@[a-z]+\.[a-z]+$`))),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("followers_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("following_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("full_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("language"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("last_login"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("location"), knownvalue.StringExact("Mêlée Island")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("login"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("must_change_password"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("prohibit_login"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("restricted"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("send_notify"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("source_id"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("starred_repos_count"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("visibility"), knownvalue.StringExact("private")),
					statecheck.ExpectKnownValue("forgejo_user.test", tfjsonpath.New("website"), knownvalue.StringExact("http://localhost:3000")),
					statecheck.ExpectSensitiveValue("forgejo_user.test", tfjsonpath.New("password")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
