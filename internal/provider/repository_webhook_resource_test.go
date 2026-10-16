package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent repository)
			{
				Config: providerConfig + `
resource "forgejo_repository_webhook" "test" {
	repository_id = -1
	type          = "forgejo"
    events        = ["push"]
    config        = {
        "content_type" = "json"
        "url"          = "http://example.com/abc12345"
    }
}`,
				ExpectError: regexp.MustCompile("Repository with ID -1 not found"),
			},
			// Create and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo"
}
resource "forgejo_repository_webhook" "test" {
	repository_id = forgejo_repository.test.id
	type          = "forgejo"
    events        = ["push"]
    config        = {
        "content_type" = "json"
        "url"          = "http://example.com/abc12345"
        "secret"       = "supersecret"
    }
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("webhook_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("active"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("branch_filter"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("authorization_header"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("type"), knownvalue.StringExact("forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("events").AtSliceIndex(0), knownvalue.StringExact("push")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config"), knownvalue.MapExact(map[string]knownvalue.Check{
						"content_type": knownvalue.StringExact("json"),
						"url":          knownvalue.StringExact("http://example.com/abc12345"),
						"secret":       knownvalue.StringExact("supersecret"),
					})),
				},
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError: regexp.MustCompile(`Expected import identifier with format: 'owner/name/webhookID', got:
'invalid'`),
			},
			// Import testing (non-existent repo)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "tfadmin/non-existent/1",
				ExpectError: regexp.MustCompile(`Repository webhook with owner 'tfadmin', repo 'non-existent' and ID 1 not
found`),
			},
			// Import testing (invalid ID)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "tfadmin/test_repo/invalid",
				ExpectError:   regexp.MustCompile("Failed to parse webhook ID"),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "tfadmin/test_repo/0",
				ExpectError:   regexp.MustCompile("Repository webhook with owner 'tfadmin', repo 'test_repo' and ID 0 not found"),
			},
			// Import testing
			{
				ResourceName:                         "forgejo_repository_webhook.test",
				ImportState:                          true,
				ImportStateIdFunc:                    testAccRepositoryWebhookConfigImportStateIdFunc("forgejo_repository_webhook.test"),
				ImportStateIdPrefix:                  "tfadmin/test_repo/",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "webhook_id",
				ImportStateVerifyIgnore:              []string{"config.secret", "config.%"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo"
}
resource "forgejo_repository_webhook" "test" {
	repository_id = forgejo_repository.test.id
	type          = "forgejo"
    events        = ["push"]
    config = {
        "content_type" = "json"
        "url"          = "http://example.com/abc12345"
        "secret"       = "rotatedsecret"
    }
    active               = true
    authorization_header = "Bearer token123456"
    branch_filter        = "*"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("webhook_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("active"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("branch_filter"), knownvalue.StringExact("*")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("authorization_header"), knownvalue.StringExact("Bearer token123456")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("type"), knownvalue.StringExact("forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("events").AtSliceIndex(0), knownvalue.StringExact("push")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config"), knownvalue.MapExact(map[string]knownvalue.Check{
						"content_type": knownvalue.StringExact("json"),
						"url":          knownvalue.StringExact("http://example.com/abc12345"),
						"secret":       knownvalue.StringExact("rotatedsecret"),
					})),
				},
			},
			// Recreate and Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo"
}
resource "forgejo_repository_webhook" "test" {
	repository_id = forgejo_repository.test.id
	type          = "gitea"
    events        = ["push"]
    config        = {
        "content_type" = "json"
        "url"          = "http://example.com/abc12345"
    }
    active               = true
    authorization_header = "Bearer token123456"
    branch_filter        = "*"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("webhook_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("active"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("branch_filter"), knownvalue.StringExact("*")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("authorization_header"), knownvalue.StringExact("Bearer token123456")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("type"), knownvalue.StringExact("gitea")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("events").AtSliceIndex(0), knownvalue.StringExact("push")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config"), knownvalue.MapExact(map[string]knownvalue.Check{
						"content_type": knownvalue.StringExact("json"),
						"url":          knownvalue.StringExact("http://example.com/abc12345"),
					})),
				},
			},
		},
	})
}

func TestAccRepositoryWebhookResourceEventGroups(t *testing.T) {
	config := func(events string) string {
		return providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_event_groups"
}
resource "forgejo_repository_webhook" "test" {
	repository_id = forgejo_repository.test.id
	type          = "forgejo"
	events        = ` + events + `
	config        = {
		"content_type" = "json"
		"url"          = "http://example.com/event-groups"
	}
}`
	}
	expectEvents := func(events ...string) statecheck.StateCheck {
		checks := make([]knownvalue.Check, 0, len(events))
		for _, event := range events {
			checks = append(checks, knownvalue.StringExact(event))
		}

		return statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("events"), knownvalue.SetExact(checks))
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (issues group)
			{
				Config:            config(`["issues"]`),
				ConfigStateChecks: []statecheck.StateCheck{expectEvents("issues")},
			},
			// Update and Read testing (pull_request group)
			{
				Config:            config(`["pull_request", "push"]`),
				ConfigStateChecks: []statecheck.StateCheck{expectEvents("pull_request", "push")},
			},
			// Update and Read testing (pull_request_comment implies pull_request_review_comment)
			{
				Config:            config(`["pull_request_comment"]`),
				ConfigStateChecks: []statecheck.StateCheck{expectEvents("pull_request_comment")},
			},
			// Update and Read testing (group and member listed together)
			{
				Config:            config(`["issues", "issue_assign"]`),
				ConfigStateChecks: []statecheck.StateCheck{expectEvents("issues", "issue_assign")},
			},
			// Update and Read testing (members without their group)
			{
				Config:            config(`["issue_assign", "pull_request_sync"]`),
				ConfigStateChecks: []statecheck.StateCheck{expectEvents("issue_assign", "pull_request_sync")},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// retrieves the dynamically generated repository webhook ID from state for testing import.
func testAccRepositoryWebhookConfigImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["webhook_id"], nil
	}
}
