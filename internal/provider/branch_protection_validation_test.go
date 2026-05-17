package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccBranchProtectionValidationPushConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check push_enabled and push_whitelist_enabled attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_push           = false
	enable_push_whitelist = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Cannot enable push whitelist without enabling push protection"),
			},
			// Check push_whitelist_usernames attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name              = "main"
	repository_id            = forgejo_repository.test.id
	enable_push              = false
	enable_push_whitelist    = false
	push_whitelist_usernames = ["user1", "user2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Push Whitelist configuration is not valid when 'enable_push_whitelist' is false"),
			},
			// Check push_whitelist_teams attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_push           = false
	enable_push_whitelist = false
	push_whitelist_teams  = ["team1", "team2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Push Whitelist configuration is not valid when 'enable_push_whitelist' is false"),
			},
			// Check push_whitelist_deploy_keys attribute
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = false
	enable_push_whitelist      = false
	push_whitelist_deploy_keys = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Push Whitelist configuration is not valid when 'enable_push_whitelist' is false"),
			},
			// enabling push, but not push_whitelist while configuring whitelist attributes ought to also error
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = true
	enable_push_whitelist      = false
	push_whitelist_deploy_keys = true
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Push Whitelist configuration is not valid when 'enable_push_whitelist' is false"),
			},
			// successful push configuration with push_whitelist
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                = "main"
	repository_id              = forgejo_repository.test.id
	enable_push                = true
	enable_push_whitelist      = true
	push_whitelist_deploy_keys = true
	push_whitelist_usernames   = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_push_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_usernames"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_teams"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("push_whitelist_deploy_keys"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func TestAccBranchProtectionValidationStatusCheckConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check status_check_enabled and status_check_contexts attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = false
	status_check_contexts = ["context1", "context2"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Cannot specify status check contexts without enabling status check"),
			},
			// Valid. status_check_enabled is false and contexts are empty
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = false
	status_check_contexts = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(0)),
				},
			},
			// Valid. status_check_enabled is enabled and contexts are set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name           = "main"
	repository_id         = forgejo_repository.test.id
	enable_status_check   = true
	status_check_contexts = ["*"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_status_check"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("status_check_contexts"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func TestAccBranchProtectionValidationMergeConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check merge_settings attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = false
	merge_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Merge configuration is not valid when 'enable_merge_whitelist' is false"),
			},
			// Valid. merge_whitelist_usernames is set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = true
	merge_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(1)),
				},
			},
			// Valid. merge configuration is disabled
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name               = "main"
	repository_id             = forgejo_repository.test.id
	enable_merge_whitelist    = false
	merge_whitelist_usernames = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_merge_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("merge_whitelist_usernames"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}

func TestAccBranchProtectionValidationApprovalsConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check approvals attributes
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = false
	approvals_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Approvals configuration is not valid when 'enable_approvals_whitelist' is false"),
			},
			// Valid. approvals_whitelist_usernames is set
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = true
	approvals_whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(1)),
				},
			},
			// Valid. merge configuration is disabled
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_branch_protection"
}
resource "forgejo_branch_protection" "test" {
	branch_name                   = "main"
	repository_id                 = forgejo_repository.test.id
	enable_approvals_whitelist    = false
	approvals_whitelist_usernames = []
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_branch_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("branch_name"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("enable_approvals_whitelist"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_branch_protection.test", tfjsonpath.New("approvals_whitelist_usernames"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}
