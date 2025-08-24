package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDeployKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
  algorithm = "ED25519"
}
resource "forgejo_repository" "test" {
  name = "test_repo"
}
resource "forgejo_deploy_key" "test" {
  repository_id = forgejo_repository.test.id
  key           = trimspace(tls_private_key.test.public_key_openssh)
  title         = "tftest"
  read_only     = false
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/api/v1/repos/tfadmin/test_repo/keys/[0-9]+$"))),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
  algorithm = "ED25519"
}
resource "forgejo_repository" "test" {
  name = "test_repo"
}
resource "forgejo_deploy_key" "test" {
  repository_id = forgejo_repository.test.id
  key           = trimspace(tls_private_key.test.public_key_openssh)
  title         = "tftest"
  read_only     = true
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/api/v1/repos/tfadmin/test_repo/keys/[0-9]+$"))),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
  algorithm = "ED25519"
}
resource "forgejo_repository" "test" {
  name = "test_repo"
}
resource "forgejo_deploy_key" "test" {
  repository_id = forgejo_repository.test.id
  key           = trimspace(tls_private_key.test.public_key_openssh)
  title         = "tftest1"
  read_only     = false
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("^ssh-ed25519 .{68}$"))),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("read_only"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest1")),
					statecheck.ExpectKnownValue("forgejo_deploy_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/api/v1/repos/tfadmin/test_repo/keys/[0-9]+$"))),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
