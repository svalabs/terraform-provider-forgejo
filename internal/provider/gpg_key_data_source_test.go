package provider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGPGKeyDataSource(t *testing.T) {
	key, armoredPubKey := createGPGKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent resource)
			{
				Config: providerConfig + `
data "forgejo_gpg_key" "test" {
	key_id= "non_existent"
}`,
				ExpectError: regexp.MustCompile("GPG key with key_id \"non_existent\" not found"),
			},
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "forgejo_gpg_key" "test" {
	armored_public_key = <<EOT
%s
EOT
}
data "forgejo_gpg_key" "test" {
	key_id = forgejo_gpg_key.test.key_id
}`, armoredPubKey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("key_id"), knownvalue.StringExact(strings.ToUpper(key.GetHexKeyID()))),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("primary_key_id"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("public_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("can_sign"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("can_encrypt_comms"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("can_encrypt_storage"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("can_certify"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("expires_at"), knownvalue.NotNull()),
				},
			},
		},
	})
}
