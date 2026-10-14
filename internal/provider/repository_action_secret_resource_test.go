package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryActionSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non-existent repository)
			{
				Config: providerConfig + `
resource "forgejo_repository_action_secret" "test" {
	repository_id = -1
	name          = "my_secret"
	data          = "my_secret_value"
}`,
				ExpectError: regexp.MustCompile("Repository with ID -1 not found"),
			},
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
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
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
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
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
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Recreate and Read testing (long name)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name        = "test_repo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_new_secret_with_a_very_long_name_that_is_over_30_characters_long"
	data          = "my_new_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_new_secret_with_a_very_long_name_that_is_over_30_characters_long")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryActionSecretValidationDataConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check data and data_wo are mutually exclusive
			{
				Config: providerConfig + `
resource "forgejo_repository_action_secret" "test" {
	repository_id = 1
	name          = "my_wo_secret"
	data          = "my_secret_value"
	data_wo       = "my_wo_secret_value"
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of`),
			},
			// Check one of data or data_wo is required
			{
				Config: providerConfig + `
resource "forgejo_repository_action_secret" "test" {
	repository_id = 1
	name          = "my_wo_secret"
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of`),
			},
			// Check data_wo_version requires data_wo
			{
				Config: providerConfig + `
resource "forgejo_repository_action_secret" "test" {
	repository_id   = 1
	name            = "my_wo_secret"
	data            = "my_secret_value"
	data_wo_version = 1
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute "data_wo" must be specified when`),
			},
		},
	})
}

func TestAccRepositoryActionSecretResourceWriteOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (write-only data without version)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_wo_secret"
	data_wo       = "my_wo_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (adding data_wo_version)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id   = forgejo_repository.test.id
	name            = "my_wo_secret"
	data_wo         = "my_wo_value"
	data_wo_version = 1
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Changing the write-only value alone is not detected (by design)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id   = forgejo_repository.test.id
	name            = "my_wo_secret"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 1
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (bumping data_wo_version)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id   = forgejo_repository.test.id
	name            = "my_wo_secret"
	data_wo         = "my_changed_wo_value"
	data_wo_version = 2
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (switching to data)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_wo_secret"
	data          = "my_secret_value"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectSensitiveValue("forgejo_repository_action_secret.test", tfjsonpath.New("data")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Update and Read testing (switching back to write-only data)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}
resource "forgejo_repository_action_secret" "test" {
	repository_id   = forgejo_repository.test.id
	name            = "my_wo_secret"
	data_wo         = "my_wo_value"
	data_wo_version = 3
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_repository_action_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("forgejo_repository_action_secret.test", tfjsonpath.New("repository_id"), "forgejo_repository.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("name"), knownvalue.StringExact("my_wo_secret")),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue("forgejo_repository_action_secret.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
