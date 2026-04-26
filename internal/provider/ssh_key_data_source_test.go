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

func TestAccSSHKeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Read testing (non-existent user)
			{
				Config: providerConfig + `
data "forgejo_ssh_key" "test" {
	user  = "non_existent"
	title = "tftest"
}`,
				ExpectError: regexp.MustCompile("SSH keys for user \"non_existent\" not found"),
			},
			// Read testing (non-existent resource)
			{
				Config: providerConfig + `
data "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	title = "non_existent"
}`,
				ExpectError: regexp.MustCompile("SSH key with user \"" + forgejoTestUser + "\" and title \"non_existent\" not found"),
			},
			// Read testing
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
data "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	title = forgejo_ssh_key.test.title
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("data.forgejo_ssh_key.test", plancheck.ResourceActionRead),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("read_only"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^"+forgejoTestHost+"/api/v1/user/keys/[0-9]+$"))),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact(forgejoTestUser)),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
		},
	})
}
