package provider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGPGKeyResource(t *testing.T) {
	key, armoredPubKey := createGPGKey(t)
	newKey, newArmoredPubKey := createGPGKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "forgejo_gpg_key" "test" {
	armored_public_key = <<EOT
%s
EOT
}`, armoredPubKey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("armored_public_key"), knownvalue.StringFunc(func(value string) error {
						if strings.TrimSpace(value) != strings.TrimSpace(armoredPubKey) {
							return fmt.Errorf(`expected "%s" to equal "%s"`, value, armoredPubKey)
						}
						return nil
					})),
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
			// Recreate and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "forgejo_gpg_key" "test" {
	armored_public_key = <<EOT
%s
EOT
}`, newArmoredPubKey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("armored_public_key"), knownvalue.StringFunc(func(value string) error {
						if strings.TrimSpace(value) != strings.TrimSpace(newArmoredPubKey) {
							return fmt.Errorf(`expected "%s" to equal "%s"`, value, newArmoredPubKey)
						}
						return nil
					})),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_gpg_key.test", tfjsonpath.New("key_id"), knownvalue.StringExact(strings.ToUpper(newKey.GetHexKeyID()))),
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

func createGPGKey(t *testing.T) (*crypto.Key, string) {
	t.Helper()

	pgp := crypto.PGP()
	key, err := pgp.KeyGeneration().
		AddUserId("Test User", forgejoEmail).
		New().
		GenerateKey()
	if err != nil {
		t.Fatalf("error generating gpg private key: %s", err)
	}

	armoredPubkey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("Error generating armored public key: %s", err)
	}

	return key, armoredPubkey
}
