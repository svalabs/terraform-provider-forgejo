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

func TestAccSSHKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent user)
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "non_existent"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}`,
				ExpectError: regexp.MustCompile("SSH key with user \"non_existent\" not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_ssh_key.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^"+forgejoTestHost+"/api/v1/user/keys/[0-9]+$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
			// Create and Read testing (duplicate key)
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}
resource "forgejo_ssh_key" "duplicate" {
	user  = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}`,
				ExpectError: regexp.MustCompile("Input validation error: Key content has been used as non-deploy key"),
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest1"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_ssh_key.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^"+forgejoTestHost+"/api/v1/user/keys/[0-9]+$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_ssh_key" "test" {
	user  = forgejo_user.test.login
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_ssh_key.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^"+forgejoTestHost+"/api/v1/user/keys/[0-9]+$"))),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
