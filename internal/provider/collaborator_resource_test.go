package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCollaboratorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent repository)
			{
				Config: providerConfig + `
resource "forgejo_collaborator" "test" {
	repository_id = -1
	user          = "tftest"
	permission    = "read"
}`,
				ExpectError: regexp.MustCompile("Repository with id -1 not found"),
			},
			// Create and Read testing (non-existent user)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name  = "test_repo"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = "non_existent"
	permission    = "read"
}`,
				ExpectError: regexp.MustCompile("Input validation error: user does not exist"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name  = "test_repo"
}
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "read"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("user"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("permission"), knownvalue.StringExact("read")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name  = "test_repo"
}
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "write"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("user"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("permission"), knownvalue.StringExact("write")),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name  = "test_repo"
}
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "admin"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("user"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("forgejo_collaborator.test", tfjsonpath.New("permission"), knownvalue.StringExact("admin")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
