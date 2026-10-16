package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

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

func TestAccRepositoryWebhookResourceWorkflowEvents(t *testing.T) {
	const allEvents = `["workflow_job_blocked", "workflow_job_cancelled", "workflow_job_failure", "workflow_job_running", "workflow_job_skipped", "workflow_job_success", "workflow_job_waiting", "workflow_run_blocked", "workflow_run_cancelled", "workflow_run_failure", "workflow_run_running", "workflow_run_skipped", "workflow_run_success", "workflow_run_waiting"]`
	config := func(events string) string {
		return providerConfig + fmt.Sprintf(`
resource "forgejo_repository" "workflow_events" {
	name = "test_workflow_events"
}
resource "forgejo_repository_webhook" "workflow_events" {
	repository_id = forgejo_repository.workflow_events.id
	type          = "forgejo"
	events        = %s
	config = {
		content_type = "json"
		url          = "http://example.com/workflow-events"
	}
}
`, events)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			client, err := forgejo.NewClient(forgejoTestHost, forgejo.SetToken(os.Getenv("FORGEJO_API_TOKEN")))
			if err != nil {
				t.Fatal(err)
			}
			if err := client.CheckServerVersionConstraint(">= 17.0.0"); err != nil {
				t.Skipf("Workflow webhook events require Forgejo 17: %s", err)
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{Config: config(`["push"]`)},
			// Recreate and Read testing (first workflow events added)
			{
				Config: config(allEvents),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.workflow_events", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.workflow_events", tfjsonpath.New("events"), knownvalue.SetSizeExact(14)),
				},
			},
			// ImportState testing
			{
				ResourceName:                         "forgejo_repository_webhook.workflow_events",
				ImportState:                          true,
				ImportStateIdFunc:                    testAccRepositoryWebhookConfigImportStateIdFunc("forgejo_repository_webhook.workflow_events"),
				ImportStateIdPrefix:                  "tfadmin/test_workflow_events/",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "webhook_id",
			},
			// Recreate and Read testing (workflow events removed)
			{
				Config: config(`["push", "workflow_job_failure", "workflow_run_failure"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.workflow_events", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.workflow_events", tfjsonpath.New("events"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("push"),
						knownvalue.StringExact("workflow_job_failure"),
						knownvalue.StringExact("workflow_run_failure"),
					})),
				},
			},
			// Update and Read testing (workflow events unchanged, push replaced by create)
			{
				Config: config(`["create", "workflow_job_failure", "workflow_run_failure"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.workflow_events", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.workflow_events", tfjsonpath.New("events"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("create"),
						knownvalue.StringExact("workflow_job_failure"),
						knownvalue.StringExact("workflow_run_failure"),
					})),
				},
			},
			// Recreate and Read testing (last workflow events removed)
			{
				Config: config(`["create"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_webhook.workflow_events", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("forgejo_repository_webhook.workflow_events", tfjsonpath.New("events"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("create")})),
				},
			},
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
