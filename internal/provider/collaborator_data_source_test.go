package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCollaboratorDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
resource "forgejo_repository" "test" {
	name  = "test_repo"
}
resource "forgejo_user" "test" {
	login    = "tftest"
	email    = "tftest@localhost.localdomain"
	password = "passw0rd"
}
resource "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
	permission    = "read"
}
data "forgejo_collaborator" "test" {
	repository_id = forgejo_repository.test.id
	user          = forgejo_user.test.login
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_collaborator.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_collaborator.test", tfjsonpath.New("user"), knownvalue.StringExact("tftest")),
					statecheck.ExpectKnownValue("data.forgejo_collaborator.test", tfjsonpath.New("permission"), knownvalue.StringRegexp(regexp.MustCompile("^read|write|admin$"))),
				},
			},
		},
	})
}
