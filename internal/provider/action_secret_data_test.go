package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// actionSecretDataSteps returns acceptance test steps covering the "data",
// "data_wo" and "data_wo_version" attributes of an action secret resource.
// parent is the configuration of the owning resource, owner the attribute
// assignment that links the secret to it.
func actionSecretDataSteps(resourceType, parent, owner string) []resource.TestStep {
	address := resourceType + ".test"
	config := func(attrs string) string {
		return providerConfig + parent + fmt.Sprintf(`
resource %q "test" {
	%s
	name = "my_wo_secret"
	%s
}`, resourceType, owner, attrs)
	}

	writeOnlyState := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(address, tfjsonpath.New("data"), knownvalue.Null()),
		statecheck.ExpectKnownValue(address, tfjsonpath.New("data_wo"), knownvalue.Null()),
		statecheck.ExpectKnownValue(address, tfjsonpath.New("created_at"), knownvalue.NotNull()),
	}

	return []resource.TestStep{
		// Validation: data and data_wo are mutually exclusive
		{
			Config:      config(`data = "a"` + "\n\tdata_wo = \"b\""),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		// Validation: one of data or data_wo is required
		{
			Config:      config(""),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		// Validation: data_wo_version requires data_wo
		{
			Config:      config(`data = "a"` + "\n\tdata_wo_version = 1"),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		// Create with write-only data
		{
			Config: config(`data_wo = "my_wo_value"` + "\n\tdata_wo_version = 1"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionCreate),
				},
			},
			ConfigStateChecks: append(writeOnlyState,
				statecheck.ExpectKnownValue(address, tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(1)),
			),
		},
		// Changing the write-only value alone is not detected (by design)
		{
			Config: config(`data_wo = "my_changed_wo_value"` + "\n\tdata_wo_version = 1"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
		},
		// Bumping the version updates the secret in place
		{
			Config: config(`data_wo = "my_changed_wo_value"` + "\n\tdata_wo_version = 2"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate),
				},
			},
			ConfigStateChecks: append(writeOnlyState,
				statecheck.ExpectKnownValue(address, tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(2)),
			),
		},
		// Switch to data persisted in state
		{
			Config: config(`data = "my_value"`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate),
				},
			},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectSensitiveValue(address, tfjsonpath.New("data")),
				statecheck.ExpectKnownValue(address, tfjsonpath.New("data_wo"), knownvalue.Null()),
				statecheck.ExpectKnownValue(address, tfjsonpath.New("data_wo_version"), knownvalue.Null()),
			},
		},
		// Switch back to write-only data: the value is removed from state
		{
			Config: config(`data_wo = "my_wo_value"` + "\n\tdata_wo_version = 3"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate),
				},
			},
			ConfigStateChecks: writeOnlyState,
		},
		// Delete testing automatically occurs in TestCase
	}
}

func TestAccRepositoryActionSecretResourceWriteOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: actionSecretDataSteps(
			"forgejo_repository_action_secret",
			`
resource "forgejo_repository" "test" {
	name = "test_repo_wo"
}`,
			"repository_id = forgejo_repository.test.id",
		),
	})
}

func TestAccOrganizationActionSecretResourceWriteOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: actionSecretDataSteps(
			"forgejo_organization_action_secret",
			`
resource "forgejo_organization" "test" {
	name = "test_org_wo"
}`,
			"organization = forgejo_organization.test.name",
		),
	})
}
