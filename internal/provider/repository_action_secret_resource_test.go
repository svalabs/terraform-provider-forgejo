package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryActionSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name        = "test_repo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_secret"
	data          = "my_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name        = "test_repo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_new_secret"
	data          = "my_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name        = "test_repo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_new_secret"
	data          = "my_new_secret_value"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
