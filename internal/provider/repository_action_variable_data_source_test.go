package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRepositoryActionVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing (non-existent repository)
			{
				Config: providerConfig + `
data "forgejo_repository_action_variable" "test" {
	repository_id = -1
	name          = "my_variable"
}`,
				ExpectError: regexp.MustCompile("Repository with ID -1 not found"),
			},
			// Read testing (non-existent variable)
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo"
}
data "forgejo_repository_action_variable" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_variable"
}`,
				ExpectError: regexp.MustCompile(`Action variable with owner "tfadmin", repo "test_repo" and name "my_variable"
not found`),
			},
			// Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name = "test_repo"
}
resource "forgejo_repository_action_variable" "test" {
	repository_id = forgejo_repository.test.id
	name          = "my_variable"
	data          = "my_variable_value"
}
data "forgejo_repository_action_variable" "test" {
	repository_id = forgejo_repository.test.id
	name          = forgejo_repository_action_variable.test.name
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_repository_action_variable.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_repository_action_variable.test", tfjsonpath.New("name"), knownvalue.StringExact("my_variable")),
					statecheck.ExpectKnownValue("data.forgejo_repository_action_variable.test", tfjsonpath.New("data"), knownvalue.StringExact("my_variable_value")),
				},
			},
		},
	})
}
