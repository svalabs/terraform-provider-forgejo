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
	branch_name = "main"
	repo        = forgejo_repository.test.name
	owner       = forgejo_repository.test.owner
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("repo"), knownvalue.StringExact("test_repo_branch_protection")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("owner"), knownvalue.StringExact("tfadmin")),
				},
			},
			// ImportState testing
			{
				ResourceName:      "forgejo_repository_branch_protection.test",
				ImportStateId:     "tfadmin/test_repo_branch_protection/main",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_repository_branch_protection" "test" {
	branch_name = "develop"
	repo        = forgejo_repository.test.name
	owner       = forgejo_repository.test.owner
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("develop")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("repo"), knownvalue.StringExact("test_repo_branch_protection")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("owner"), knownvalue.NotNull()),
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
	branch_name = "main"
	repo        = forgejo_repository.test.name
	owner       = forgejo_repository.test.owner
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("repo"), knownvalue.StringExact("test_repo_branch_protection")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_org_branch_protection")),
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
	branch_name = "main"
	repo        = forgejo_repository.test.name
	owner       = forgejo_repository.test.owner
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("repo"), knownvalue.StringExact("test_repo_branch_protection")),
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("owner"), knownvalue.StringExact("test_user_branch_protection")),
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
	branch_name                      = "main"
	repo                              = forgejo_repository.test.name
	owner                             = forgejo_repository.test.owner
	enable_push                       = false
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
					statecheck.ExpectKnownValue("forgejo_repository_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(false)),
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
	branch_name = "release/*"
	repo        = forgejo_repository.test.name
	owner       = forgejo_repository.test.owner
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
	branch_name = "main"
	repo        = "nonexistent_repo"
	owner       = "tfadmin"
}
`,
				ExpectError: regexp.MustCompile("Repository not found|404"),
			},
		},
	})
}
