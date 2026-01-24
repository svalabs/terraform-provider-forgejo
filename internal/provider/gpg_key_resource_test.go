package provider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const forgejoEmail = "tfadmin@localhost"

func TestAccGPGKeyResource(t *testing.T) {
	validateIsArmoredGPGKey := func(value string) error {
		lines := strings.Split(value, "\n")
		if lines[0] != "-----BEGIN PGP PUBLIC KEY BLOCK-----" {
			return fmt.Errorf(`expected "%s" to start with "-----BEGIN PGP PUBLIC KEY BLOCK-----"`, value)
		}
		if lines[len(lines)-1] != "-----END PGP PUBLIC KEY BLOCK-----" {
			return fmt.Errorf(`expected "%s" to end with "-----END PGP PUBLIC KEY BLOCK-----"`, value)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"gpg": {
				Source: "terraform-provider-gpg/gpg",
			},
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "gpg_key" "test" {
	identities = [{
		name  = "TF Admin"
		email = "%s"
	}]
	passphrase = "supersecret"
}
resource "forgejo_gpg_key" "test" {
	armored_public_key = gpg_key.test.public_key
}`, forgejoEmail),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("armored_public_key"), knownvalue.StringFunc(validateIsArmoredGPGKey)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
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
			// Recreate and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "gpg_key" "test" {
	identities = [{
		name  = "TF Admin"
		email = "%s"
	}]
	passphrase = "supersecret"
}
resource "forgejo_gpg_key" "test" {
	armored_public_key = gpg_key.test.public_key
}`, forgejoEmail),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("armored_public_key"), knownvalue.StringFunc(validateIsArmoredGPGKey)),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
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
			// Delete testing automatically occurs in TestCase
		},
	})
}
