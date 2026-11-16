package provider_test

import (
	"fmt"
	"os"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const importProviderConfig = providerConfig + `
terraform {
	required_providers {
		forgejo = {
			source = "svalabs/forgejo"
		}
	}
}
`

// Fixtures are created outside Terraform so import can persist without an already managed resource.
func testAccImportOrganization(t *testing.T, opts forgejo.CreateOrgOption) (*forgejo.Client, *forgejo.Organization) {
	t.Helper()
	testAccPreCheck(t)
	client, err := forgejo.NewClient(forgejoTestHost, forgejo.SetToken(os.Getenv("FORGEJO_API_TOKEN")))
	if err != nil {
		t.Fatal(err)
	}
	org, _, err := client.CreateOrg(opts)
	if err != nil {
		t.Fatalf("creating import fixture: %v", err)
	}
	t.Cleanup(func() {
		res, err := client.DeleteOrg(org.UserName)
		if err != nil && (res == nil || res.StatusCode != 404) {
			t.Errorf("cleaning import fixture: %v", err)
		}
	})

	return client, org
}

func TestAccOrganizationResourceImport(t *testing.T) {
	t.Setenv(resource.EnvTfAccProviderNamespace, "svalabs")
	for _, setting := range []string{"true", "false"} {
		t.Run(setting, func(t *testing.T) {
			name := "organization_import_" + setting
			config := func(description string) string {
				return importProviderConfig + `
resource "forgejo_organization" "test" {
	name        = "` + name + `"
	full_name   = "Existing organization"
	description = "` + description + `"
	website     = "https://example.org"
	location    = "Earth"
	visibility  = "private"
}`
			}
			var org *forgejo.Organization
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					_, org = testAccImportOrganization(t, forgejo.CreateOrgOption{
						Name: name, FullName: "Existing organization", Description: "Existing profile",
						Website: "https://example.org", Location: "Earth", Visibility: forgejo.VisibleType("private"),
						RepoAdminChangeTeamAccess: setting == "true",
					})
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: importProviderConfig},
					{
						Config:             config("Existing profile"),
						ResourceName:       "forgejo_organization.test",
						ImportState:        true,
						ImportStateId:      name,
						ImportStatePersist: true,
						ImportStateCheck: func(states []*terraform.InstanceState) error {
							if len(states) != 1 || states[0].Attributes["id"] != fmt.Sprint(org.ID) {
								return fmt.Errorf("organization identity changed during import")
							}

							return nil
						},
					},
					{
						Config: config("Existing profile"),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("repo_admin_change_team_access"), knownvalue.Null()),
						},
					},
					{
						Config: config("Changed profile"),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("forgejo_organization.test", plancheck.ResourceActionUpdate),
							},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("forgejo_organization.test", tfjsonpath.New("repo_admin_change_team_access"), knownvalue.Null()),
						},
					},
				},
			})
		})
	}
}

func TestAccTeamMemberResourceImport(t *testing.T) {
	t.Setenv(resource.EnvTfAccProviderNamespace, "svalabs")
	fixtures := importProviderConfig + `
data "forgejo_team" "test" {
	name         = "developers"
	organization = "membership_import_org"
}`
	config := fixtures + `
resource "forgejo_team_member" "test" {
	team_id = data.forgejo_team.test.id
	user    = "` + forgejoTestUser + `"
}`
	var client *forgejo.Client
	var team *forgejo.Team
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			client, _ = testAccImportOrganization(t, forgejo.CreateOrgOption{Name: "membership_import_org"})
			var err error
			team, _, err = client.CreateTeam("membership_import_org", forgejo.CreateTeamOption{
				Name: "developers", Permission: forgejo.AccessMode("read"), UnitsMap: map[string]string{"repo.code": "read"},
			})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.AddTeamMember(team.ID, forgejoTestUser); err != nil {
				t.Fatal(err)
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: importProviderConfig},
			{
				Config:       config,
				ResourceName: "forgejo_team_member.test",
				ImportState:  true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%d/%s", team.ID, forgejoTestUser), nil
				},
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 || states[0].Attributes["id"] != fmt.Sprintf("%d/%s", team.ID, forgejoTestUser) || states[0].Attributes["team_id"] != fmt.Sprint(team.ID) || states[0].Attributes["user"] != forgejoTestUser {
						return fmt.Errorf("membership identity changed during import")
					}

					return nil
				},
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: fixtures,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("forgejo_team_member.test", plancheck.ResourceActionDestroy),
					},
				},
				Check: func(*terraform.State) error {
					if _, _, err := client.GetUserInfo(forgejoTestUser); err != nil {
						return fmt.Errorf("membership deletion removed the user: %w", err)
					}
					if _, _, err := client.GetTeam(team.ID); err != nil {
						return fmt.Errorf("membership deletion removed the team: %w", err)
					}
					_, res, err := client.GetTeamMember(team.ID, forgejoTestUser)
					if err == nil || res == nil || res.StatusCode != 404 {
						return fmt.Errorf("membership was not removed")
					}

					return nil
				},
			},
		},
	})
}

func TestAccTeamRepositoryResourceImport(t *testing.T) {
	t.Setenv(resource.EnvTfAccProviderNamespace, "svalabs")
	fixtures := importProviderConfig + `
data "forgejo_team" "test" {
	name         = "developers"
	organization = "assignment_existing_org"
}
data "forgejo_repository" "test" {
	owner = "assignment_existing_org"
	name  = "api"
}`
	config := fixtures + `
resource "forgejo_team_repository" "test" {
	team_id       = data.forgejo_team.test.id
	repository_id = data.forgejo_repository.test.id
}`
	var client *forgejo.Client
	var team *forgejo.Team
	var repo *forgejo.Repository
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			client, _ = testAccImportOrganization(t, forgejo.CreateOrgOption{Name: "assignment_existing_org"})
			var err error
			team, _, err = client.CreateTeam("assignment_existing_org", forgejo.CreateTeamOption{
				Name: "developers", Permission: forgejo.AccessMode("read"), UnitsMap: map[string]string{"repo.code": "read"},
			})
			if err != nil {
				t.Fatal(err)
			}
			repo, _, err = client.CreateOrgRepo("assignment_existing_org", forgejo.CreateRepoOption{Name: "api", Private: true})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				res, err := client.DeleteRepo("assignment_existing_org", repo.Name)
				if err != nil && (res == nil || res.StatusCode != 404) {
					t.Errorf("cleaning repository fixture: %v", err)
				}
			})
			if _, err := client.AddTeamRepository(team.ID, "assignment_existing_org", "api"); err != nil {
				t.Fatal(err)
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: importProviderConfig},
			{
				Config:             config,
				ResourceName:       "forgejo_team_repository.test",
				ImportState:        true,
				ImportStateId:      "assignment_existing_org/developers/api",
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 || states[0].Attributes["id"] != fmt.Sprintf("%d/%d", team.ID, repo.ID) {
						return fmt.Errorf("assignment identity changed during import")
					}

					return nil
				},
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: fixtures,
				Check: func(*terraform.State) error {
					if _, _, err := client.GetTeam(team.ID); err != nil {
						return fmt.Errorf("assignment deletion removed the team: %w", err)
					}
					if _, _, err := client.GetRepoByID(repo.ID); err != nil {
						return fmt.Errorf("assignment deletion removed the repository: %w", err)
					}
					repos, _, err := client.ListTeamRepositories(team.ID, forgejo.ListTeamRepositoriesOptions{})
					if err != nil {
						return err
					}
					if len(repos) != 0 {
						return fmt.Errorf("assignment was not removed")
					}

					return nil
				},
			},
		},
	})
}
