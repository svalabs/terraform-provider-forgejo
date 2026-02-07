package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGPGKeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"gpg": {
				Source: "terraform-provider-gpg/gpg",
			},
		},
		Steps: []resource.TestStep{
			// Read testing (non-existent resource)
			{
				Config: providerConfig + `
data "forgejo_gpg_key" "test" {
	key_id = "non_existent"
}`,
				ExpectError: regexp.MustCompile(`GPG key with key_id "non_existent" not found`),
			},
			// Read testing (valid user, non-existent resource)
			{
				Config: providerConfig + `
data "forgejo_gpg_key" "test" {
	user   = "tfadmin"
	key_id = "non_existent"
}`,
				ExpectError: regexp.MustCompile(`GPG key with user "tfadmin" and key_id "non_existent" not found`),
			},
			// Read testing (non-existent user)
			{
				Config: providerConfig + `
data "forgejo_gpg_key" "test" {
	user   = "invalid"
	key_id = "non_existent"
}`,
				ExpectError: regexp.MustCompile(`GPG keys for user "invalid" not found`),
			},
			// Read testing (current user)
			{
				Config: providerConfig + fmt.Sprintf(`
resource "gpg_key_pair" "test" {
	identities = [{
		name  = "TF Admin"
		email = "%s"
	}]
	passphrase = "supersecret"
}
resource "forgejo_gpg_key" "test" {
	armored_public_key = gpg_key_pair.test.public_key
}
data "forgejo_gpg_key" "test" {
	key_id = gpg_key_pair.test.id

	depends_on = [forgejo_gpg_key.test]
}`, forgejoEmail),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("primary_key_id"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("public_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("can_sign"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("can_encrypt_comms"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("can_encrypt_storage"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("can_certify"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("expires_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("emails").AtSliceIndex(0).AtMapKey("email"), knownvalue.StringExact(forgejoEmail)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("emails").AtSliceIndex(0).AtMapKey("verified"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("primary_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("public_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("can_sign"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("can_encrypt_comms"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("can_encrypt_storage"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("can_certify"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_gpg_key.test", tfjsonpath.New("subkeys").AtSliceIndex(0).AtMapKey("expires_at"), knownvalue.NotNull()),
				},
			},
			// Read testing (explicit user)
			{
				Config: providerConfig + fmt.Sprintf(`
resource "gpg_key_pair" "test" {
	identities = [{
		name  = "TF Admin"
		email = "%s"
	}]
	passphrase = "supersecret"
}
resource "forgejo_gpg_key" "test" {
	armored_public_key = gpg_key_pair.test.public_key
}
data "forgejo_gpg_key" "test" {
	user   = "tfadmin"
	key_id = gpg_key_pair.test.id

	depends_on = [forgejo_gpg_key.test]
}`, forgejoEmail),
				ConfigStateChecks: []statecheck.StateCheck{
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
		},
	})
}
