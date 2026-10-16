package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTagProtectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent repository)
			{
				Config: providerConfig + `
resource "forgejo_tag_protection" "test" {
	repository_id       = -1
	name_pattern        = "v*"
	whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ExpectError: regexp.MustCompile("Repository with ID -1 not found"),
			},
			// Create and Read testing (empty whitelist)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_tag_protection" "test" {
	repository_id = forgejo_repository.test.id
	name_pattern  = "v*"
}`,
				ExpectError: regexp.MustCompile(`requires\s+at\s+least\s+one\s+whitelisted\s+user\s+or\s+team`),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_tag_protection" "test" {
	repository_id       = forgejo_repository.test.id
	name_pattern        = "v*"
	whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_tag_protection.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_tag_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("name_pattern"), knownvalue.StringExact("v*")),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_usernames"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(forgejoTestUser),
					})),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_teams"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_tag_protection.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile(`Expected\s+import\s+identifier\s+with\s+format:\s+'owner/repo/protectionID'`),
			},
			// Import testing (invalid ID)
			{
				ResourceName:  "forgejo_tag_protection.test",
				ImportState:   true,
				ImportStateId: forgejoTestUser + "/test_repo_tag_protection/invalid",
				ExpectError:   regexp.MustCompile("Failed to parse protection ID"),
			},
			// Import testing (non-existent repo)
			{
				ResourceName:  "forgejo_tag_protection.test",
				ImportState:   true,
				ImportStateId: forgejoTestUser + "/non-existent/1",
				ExpectError:   regexp.MustCompile(`Tag\s+protection\s+with\s+owner\s+'` + forgejoTestUser + `',\s+repo\s+'non-existent'\s+and\s+ID\s+1\s+not\s+found`),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_tag_protection.test",
				ImportState:   true,
				ImportStateId: forgejoTestUser + "/test_repo_tag_protection/0",
				ExpectError:   regexp.MustCompile(`Tag\s+protection\s+with\s+owner\s+'` + forgejoTestUser + `',\s+repo\s+'test_repo_tag_protection'`),
			},
			// Import testing
			{
				ResourceName:                         "forgejo_tag_protection.test",
				ImportState:                          true,
				ImportStateIdFunc:                    testAccTagProtectionImportStateIdFunc("forgejo_tag_protection.test"),
				ImportStateIdPrefix:                  forgejoTestUser + "/test_repo_tag_protection/",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
			},
			// Update and Read testing (name pattern)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_tag_protection" "test" {
	repository_id       = forgejo_repository.test.id
	name_pattern        = "release-*"
	whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_tag_protection.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_tag_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("name_pattern"), knownvalue.StringExact("release-*")),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_usernames"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(forgejoTestUser),
					})),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_teams"), knownvalue.SetSizeExact(0)),
				},
			},
			// Update and Read testing (whitelist with collaborator)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_user" "test" {
	login    = "test_user_tag_protection"
	email    = "test_user_tag_protection@localhost"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "write"
}
resource "forgejo_tag_protection" "test" {
	depends_on = [forgejo_collaborator.test]

	repository_id = forgejo_repository.test.id
	name_pattern  = "release-*"
	whitelist_usernames = [
		"` + forgejoTestUser + `",
		forgejo_user.test.login,
	]
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_tag_protection.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_tag_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("name_pattern"), knownvalue.StringExact("release-*")),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_usernames"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(forgejoTestUser),
						knownvalue.StringExact("test_user_tag_protection"),
					})),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_teams"), knownvalue.SetSizeExact(0)),
				},
			},
			// Update and Read testing (whitelisted user without repository access)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_user" "test" {
	login    = "test_user_tag_protection"
	email    = "test_user_tag_protection@localhost"
	password = "passw0rd"
}
resource "forgejo_user" "no_access" {
	login    = "test_user_tag_protection_no_access"
	email    = "test_user_tag_protection_no_access@localhost"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "write"
}
resource "forgejo_tag_protection" "test" {
	depends_on = [forgejo_collaborator.test]

	repository_id = forgejo_repository.test.id
	name_pattern  = "release-*"
	whitelist_usernames = [
		"` + forgejoTestUser + `",
		forgejo_user.test.login,
		forgejo_user.no_access.login,
	]
}`,
				ExpectError: regexp.MustCompile(`have\s+no\s+access\s+to\s+the\s+repository:\s+test_user_tag_protection_no_access`),
			},
			// Update and Read testing (empty whitelists)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_tag_protection"
}
resource "forgejo_tag_protection" "test" {
	repository_id = forgejo_repository.test.id
	name_pattern  = "release-*"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_tag_protection.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_tag_protection.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("name_pattern"), knownvalue.StringExact("release-*")),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_usernames"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue("forgejo_tag_protection.test", tfjsonpath.New("whitelist_teams"), knownvalue.SetSizeExact(0)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagProtectionValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check name_pattern attribute (empty)
			{
				Config: providerConfig + `
resource "forgejo_tag_protection" "test" {
	repository_id       = 1
	name_pattern        = ""
	whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ExpectError: regexp.MustCompile(`must\s+not\s+be\s+empty\s+or\s+have\s+leading\s+or\s+trailing\s+whitespace`),
			},
			// Check name_pattern attribute (whitespace)
			{
				Config: providerConfig + `
resource "forgejo_tag_protection" "test" {
	repository_id       = 1
	name_pattern        = " v* "
	whitelist_usernames = ["` + forgejoTestUser + `"]
}`,
				ExpectError: regexp.MustCompile(`must\s+not\s+be\s+empty\s+or\s+have\s+leading\s+or\s+trailing\s+whitespace`),
			},
		},
	})
}

// retrieves the dynamically generated tag protection ID from state for testing import.
func testAccTagProtectionImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["id"], nil
	}
}
