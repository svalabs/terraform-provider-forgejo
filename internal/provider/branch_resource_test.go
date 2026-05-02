package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccBranchResource_OrgRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with organization
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_branch"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo_branch"
}
resource "forgejo_branch" "test" {
	name       = "test"
	owner	   = forgejo_organization.test.name
	repository = forgejo_repository.test.name
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("repository"), knownvalue.StringExact("test_repo_branch")),
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org_branch")),
				},
			},
		},
	})
}

func TestAccBranchResource_UserRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with user
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user_branch"
	email    = "test_user_branch@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.test.login
	name  = "test_repo_branch"
}
resource "forgejo_branch" "test" {
	name       = "test"
	owner	   = forgejo_user.test.login
	repository = forgejo_repository.test.name
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("repository"), knownvalue.StringExact("test_repo_branch")),
					statecheck.ExpectKnownValue("forgejo_branch.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user_branch")),
				},
			},
		},
	})
}

func TestAccBranchResource_InvalidRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "forgejo_branch_protection" "test" {
	branch_name   = "main"
	repository_id = 123
}`,
				ExpectError: regexp.MustCompile("Error: Unable to read repository"),
			},
		},
	})
}
