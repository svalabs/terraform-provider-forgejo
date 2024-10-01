package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "forgejo_organization" "test" { name = "test1" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify organization to ensure all attributes are set
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "avatar_url", "http://localhost:3000/avatars/5a105e8b9d40e1329780d62ea2265d8a"),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "description", ""),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "full_name", ""),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "id", "2"),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "location", ""),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "name", "test1"),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "visibility", "public"),
					resource.TestCheckResourceAttr("data.forgejo_organization.test", "website", ""),
				),
			},
		},
	})
}
