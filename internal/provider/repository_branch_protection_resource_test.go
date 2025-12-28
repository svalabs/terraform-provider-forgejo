package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryBranchProtectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name   = "main"
	repository_id = forgejo_repository.test.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
				},
			},
			// ImportState testing
			{
				ResourceName:                         "forgejo_repository_branch_protection.test",
				ImportStateId:                        "tfadmin/test_repo_branch_protection/main",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "branch_name",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name   = "develop"
	repository_id = forgejo_repository.test.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("develop")),
				},
			},
		},
	})
}

func TestAccRepositoryBranchProtectionResource_OrgRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with organization
			{
				Config: providerConfig + `
resource "forgejo_organization" "test" {
	name = "test_org_branch_protection"
}
resource "forgejo_repository" "test" {
	owner = forgejo_organization.test.name
	name  = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name    = "main"
	repository_id  = forgejo_repository.test.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
				},
			},
		},
	})
}

func TestAccRepositoryBranchProtectionResource_UserRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with user
			{
				Config: providerConfig + `
resource "forgejo_user" "test" {
	login    = "test_user_branch_protection"
	email    = "test_user_branch_protection@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_repository" "test" {
	owner = forgejo_user.test.login
	name  = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name		= "main"
	repository_id	= forgejo_repository.test.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
				},
			},
		},
	})
}

func TestAccRepositoryBranchProtectionResource_WithOptionalAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name                       = "main"
	repository_id      				  = forgejo_repository.test.id
	enable_push                       = true
	enable_push_whitelist             = true
	push_whitelist_usernames          = ["tfadmin"]
	require_signed_commits            = true
	required_approvals                = 2
	block_on_rejected_reviews         = true
	block_on_official_review_requests = true
	dismiss_stale_approvals           = true
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("require_signed_commits"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("required_approvals"), knownvalue.Int64Exact(2)),
				},
			},
		},
	})
}

func TestAccRepositoryBranchProtectionResource_BranchPattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name 	     = "release/*"
	repository_id        = forgejo_repository.test.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("release/*")),
				},
			},
		},
	})
}

func TestAccRepositoryBranchProtectionResource_InvalidRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "forgejo_repository_branch_protection" "test" {
	branch_name   = "main"
	repository_id = 123
}
`,
				ExpectError: regexp.MustCompile("Error: Unable to fetch repository details"),
			},
		},
	})
}
