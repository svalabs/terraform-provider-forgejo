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
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("branch_filter"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("authorization_header"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("type"), knownvalue.StringExact("forgejo")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("events").AtSliceIndex(0), knownvalue.StringExact("push")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("content_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("url"), knownvalue.StringExact("http://example.com/abc12345")),
				},
			},
			// Import testing (invalid identifier)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile(`Expected import identifier with format: 'owner/name/webhookID', got:\s+'invalid'`),
			},
			// Import testing (non-existent repo)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "tfadmin/non-existent/1",
				ExpectError:   regexp.MustCompile("Repository with owner 'tfadmin' and name 'non-existent' not found"),
			},
			// Import testing (non-existent resource)
			{
				ResourceName:  "forgejo_repository_webhook.test",
				ImportState:   true,
				ImportStateId: "tfadmin/test_repo/0",
				ExpectError:   regexp.MustCompile("Repository webhook 'tfadmin/test_repo/0' not found"),
			},
			// Import testing
			{
				ResourceName:                         "forgejo_repository_webhook.test",
				ImportState:                          true,
				ImportStateIdFunc:                    testAccRepositoryWebhookConfigImportStateIdFunc("forgejo_repository_webhook.test"),
				ImportStateIdPrefix:                  "tfadmin/test_repo/",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "webhook_id",
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
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("content_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("url"), knownvalue.StringExact("http://example.com/abc12345")),
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
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("content_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("forgejo_repository_webhook.test", tfjsonpath.New("config").AtMapKey("url"), knownvalue.StringExact("http://example.com/abc12345")),
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
