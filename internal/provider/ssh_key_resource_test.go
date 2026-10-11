package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
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

func TestAccSSHKeyResourceForAuthedUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Create and Read testing, without a user set
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest_userless"
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
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("title"), knownvalue.StringExact("tftest_userless")),
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
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest_self"
}
resource "forgejo_ssh_key" "duplicate" {
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest_self"
}`,
				ExpectError: regexp.MustCompile("Input validation error: Key content has been used as non-deploy key"),
			},

			// Create and recreate testing in a migration from key with user to key without
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest_owned"
}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_ssh_key.test", plancheck.ResourceActionCreate),
					},
				},
			},
			{
				Config: providerConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest_owned"
}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_ssh_key.test", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// forgejoTestNonAdmin* defines a non-admin user that is created out of band
// (see testAccSSHKeyNonAdminUser) to test the non-admin route of
// forgejo_ssh_key, i.e. a provider authenticated as a regular user.
const (
	forgejoTestNonAdminUser     = "tftest_nonadmin"
	forgejoTestNonAdminPassword = "passw0rd"

	// providerNonAdminConfig authenticates as the non-admin user. It must be the
	// only "forgejo" provider configuration in a step: all configurations share
	// a single in-process provider server during acceptance testing, so mixing
	// admin and non-admin credentials in the same step is unreliable.
	// api_token is set to "" to override FORGEJO_API_TOKEN from the environment.
	providerNonAdminConfig = `provider "forgejo" {
host      = "` + forgejoTestHost + `"
username  = "` + forgejoTestNonAdminUser + `"
password  = "` + forgejoTestNonAdminPassword + `"
api_token = ""
}`
)

// testAccAdminClient returns a Forgejo client authenticated as the admin user
// used for acceptance testing.
func testAccAdminClient(t *testing.T) *forgejo.Client {
	t.Helper()

	client, err := forgejo.NewClient(
		forgejoTestHost,
		forgejo.SetToken(os.Getenv("FORGEJO_API_TOKEN")),
	)
	if err != nil {
		t.Fatalf("Unable to create Forgejo admin client: %s", err)
	}

	return client
}

// testAccSSHKeyNonAdminUser (re-)creates the non-admin test user and removes
// it after the test. The user is managed outside of Terraform, because the
// provider configuration authenticating as this user requires it to exist
// before the first plan. Deleting the user also deletes all of its SSH keys,
// so nothing is left behind even if the test fails.
func testAccSSHKeyNonAdminUser(t *testing.T) {
	t.Helper()

	client := testAccAdminClient(t)

	// Remove leftovers from previous runs (the dev instance is not reset)
	if res, err := client.AdminDeleteUser(forgejoTestNonAdminUser); err != nil && (res == nil || res.StatusCode != 404) {
		t.Fatalf("Unable to delete user '%s': %s", forgejoTestNonAdminUser, err)
	}

	_, _, err := client.AdminCreateUser(forgejo.CreateUserOption{
		Username:           forgejoTestNonAdminUser,
		Email:              forgejoTestNonAdminUser + "@localhost.localdomain",
		Password:           forgejoTestNonAdminPassword,
		MustChangePassword: forgejo.OptionalBool(false),
	})
	if err != nil {
		t.Fatalf("Unable to create user '%s': %s", forgejoTestNonAdminUser, err)
	}

	t.Cleanup(func() {
		if _, err := client.AdminDeleteUser(forgejoTestNonAdminUser); err != nil {
			t.Errorf("Unable to delete user '%s': %s", forgejoTestNonAdminUser, err)
		}
	})
}

// testAccCheckSSHKeyNonAdminDestroy verifies that the non-admin user has no
// SSH keys left after destroy, i.e. that Delete really removed them.
func testAccCheckSSHKeyNonAdminDestroy(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		keys, _, err := testAccAdminClient(t).ListPublicKeys(
			forgejoTestNonAdminUser,
			forgejo.ListPublicKeysOptions{ListOptions: forgejo.ListOptions{Page: -1}},
		)
		if err != nil {
			return fmt.Errorf("unable to list SSH keys of user '%s': %w", forgejoTestNonAdminUser, err)
		}
		if len(keys) > 0 {
			return fmt.Errorf("user '%s' still has %d SSH key(s) after destroy", forgejoTestNonAdminUser, len(keys))
		}

		return nil
	}
}

// TestAccSSHKeyResourceNonAdmin tests the use case from issue #170: managing
// the SSH keys of the authenticated user without administrative privileges.
func TestAccSSHKeyResourceNonAdmin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccSSHKeyNonAdminUser(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		CheckDestroy: testAccCheckSSHKeyNonAdminDestroy(t),
		Steps: []resource.TestStep{
			// Create and Read testing (other user requires admin privileges)
			{
				Config: providerNonAdminConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
	user  = "` + forgejoTestUser + `"
	key   = trimspace(tls_private_key.test.public_key_openssh)
	title = "tftest"
}`,
				ExpectError: regexp.MustCompile("SSH key with user \"" + forgejoTestUser + "\" forbidden"),
			},
			// Create and Read testing (authenticated user)
			{
				Config: providerNonAdminConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
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
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact(forgejoTestNonAdminUser)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
			// Recreate and Read testing (deletes the previous key as non-admin)
			{
				Config: providerNonAdminConfig + `
resource "tls_private_key" "test" {
	algorithm = "ED25519"
}
resource "forgejo_ssh_key" "test" {
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
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("user"), knownvalue.StringExact(forgejoTestNonAdminUser)),
					statecheck.ExpectKnownValue("forgejo_ssh_key.test", tfjsonpath.New("key_type"), knownvalue.StringExact("user")),
				},
			},
			// Delete testing automatically occurs in TestCase (as non-admin)
		},
	})
}
