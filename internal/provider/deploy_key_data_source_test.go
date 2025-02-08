package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDeployKeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "forgejo_repository" "test" {
  owner = {
    login = "achim"
  }
  name = "user_test_repo_1"
}
data "forgejo_deploy_key" "test" {
  repository_id = data.forgejo_repository.test.id
  title         = "test_1"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("fingerprint"), knownvalue.StringRegexp(regexp.MustCompile("^SHA256:.{43}$"))),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("read_only"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("title"), knownvalue.StringExact("test_1")),
					statecheck.ExpectKnownValue("data.forgejo_deploy_key.test", tfjsonpath.New("url"), knownvalue.StringRegexp(regexp.MustCompile("^http://localhost:3000/api/v1/repos/achim/user_test_repo_1/keys/[0-9]+$"))),
				},
			},
		},
	})
}
