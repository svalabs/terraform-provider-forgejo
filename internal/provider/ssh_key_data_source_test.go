package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
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
			// Read testing
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "tfadmin"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}
data "forgejo_ssh_key" "test" {
	user  = "tfadmin"
	title = forgejo_ssh_key.test.title
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("read_only"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/api/v1/user/keys/[0-9]+$"))),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact("tfadmin")),
					statecheck.ExpectKnownValue("data.forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
		},
	})
}
