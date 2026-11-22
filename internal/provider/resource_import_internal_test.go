package provider

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v4"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func importTestClient(t *testing.T, handler http.HandlerFunc) *forgejo.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := forgejo.NewClient(server.URL, forgejo.SetForgejoVersion("17.0.0"))
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func importTestState(t *testing.T, r resource.Resource, model any) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	if model != nil {
		if diags := state.Set(t.Context(), model); diags.HasError() {
			t.Fatal(diags)
		}
	}

	return state
}

func importTestResource(t *testing.T, r resource.ResourceWithImportState, id string) resource.ImportStateResponse {
	t.Helper()
	resp := resource.ImportStateResponse{State: importTestState(t, r, nil)}
	r.ImportState(t.Context(), resource.ImportStateRequest{ID: id}, &resp)

	return resp
}

func TestImportRejectsMalformedIdentifiersBeforeHTTP(t *testing.T) {
	client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected HTTP request: %s %s", r.Method, r.URL)
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	})
	for _, tc := range []struct {
		name string
		r    resource.ResourceWithImportState
		ids  []string
	}{
		{"organization", &organizationResource{client: client}, []string{"", " ", "org/team", "/org"}},
		{"membership", &teamMemberResource{client: client}, []string{"", "alice", "/alice", "42/", "42/ ", "42/alice/extra", "zero/alice", "0/alice", "-1/alice", "9223372036854775808/alice"}},
		{"assignment", &teamRepositoryResource{client: client}, []string{"", "org/team", "/team/repo", "org//repo", "org/team/", "org/ /repo", "org/team/repo/extra"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, id := range tc.ids {
				t.Run(id, func(t *testing.T) {
					resp := importTestResource(t, tc.r, id)
					if !resp.Diagnostics.HasError() {
						t.Fatal("expected malformed identifier error")
					}
				})
			}
		})
	}
}

func TestOrganizationAndMembershipImportAreReadOnly(t *testing.T) {
	client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("import mutated the server: %s %s", r.Method, r.URL)
			http.Error(w, "mutation forbidden", http.StatusForbidden)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/orgs/acme":
			fmt.Fprint(w, `{"id":7,"username":"acme","full_name":"Acme","description":"Existing profile","website":"https://example.org","location":"Earth","visibility":"private"}`)
		case "/api/v1/teams/42/members/alice":
			fmt.Fprint(w, `{"id":8,"login":"alice"}`)
		default:
			t.Errorf("unexpected lookup (no admin enumeration allowed): %s", r.URL)
			http.NotFound(w, r)
		}
	})

	org := importTestResource(t, &organizationResource{client: client}, "acme")
	if org.Diagnostics.HasError() {
		t.Fatal(org.Diagnostics)
	}
	var profile organizationResourceModel
	if diags := org.State.Get(t.Context(), &profile); diags.HasError() {
		t.Fatal(diags)
	}
	if profile.ID.ValueInt64() != 7 || profile.Description.ValueString() != "Existing profile" || profile.Visibility.ValueString() != "private" {
		t.Fatalf("unexpected imported profile: %+v", profile)
	}
	if !profile.RepoAdminChangeTeamAccess.IsNull() {
		t.Fatal("import invented the unreadable organization setting")
	}

	member := importTestResource(t, &teamMemberResource{client: client}, "42/alice")
	if member.Diagnostics.HasError() {
		t.Fatal(member.Diagnostics)
	}
	var membership teamMemberResourceModel
	if diags := member.State.Get(t.Context(), &membership); diags.HasError() {
		t.Fatal(diags)
	}
	if membership.ID.ValueString() != "42/alice" || membership.TeamID.ValueInt64() != 42 || membership.User.ValueString() != "alice" {
		t.Fatalf("unexpected imported membership: %+v", membership)
	}
}

func TestResourceReadDistinguishesAbsenceFromFailure(t *testing.T) {
	for _, status := range []int{200, 401, 403, 404, 429, 500} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			for _, name := range []string{"organization", "membership", "assignment"} {
				t.Run(name, func(t *testing.T) {
					client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodGet {
							t.Errorf("Read mutated the server: %s", r.Method)
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(status)
						if status != http.StatusOK {
							fmt.Fprint(w, `{"message":"request failed"}`)
						} else {
							switch name {
							case "organization":
								fmt.Fprint(w, `{"id":7,"username":"acme","visibility":"private"}`)
							case "membership":
								fmt.Fprint(w, `{"login":"alice"}`)
							case "assignment":
								fmt.Fprint(w, `[{"id":99}]`)
							}
						}
					})
					var r resource.Resource
					var model any
					switch name {
					case "organization":
						r = &organizationResource{client: client}
						model = organizationResourceModel{ID: types.Int64Value(7), Name: types.StringValue("acme")}
					case "membership":
						r = &teamMemberResource{client: client}
						model = teamMemberResourceModel{TeamID: types.Int64Value(42), User: types.StringValue("alice")}
					case "assignment":
						r = &teamRepositoryResource{client: client}
						model = teamRepositoryResourceModel{ID: types.StringValue("42/99"), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)}
					}
					state := importTestState(t, r, model)
					resp := resource.ReadResponse{State: state}
					r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)
					wantError := status != 200 && status != 404
					if resp.Diagnostics.HasError() != wantError {
						t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
					}
					if removed := resp.State.Raw.IsNull(); removed != (status == 404) {
						t.Fatalf("removed=%v for HTTP %d", removed, status)
					}
					if wantError && !resp.State.Raw.Equal(state.Raw) {
						t.Fatal("state changed after failed read")
					}
					if status == 200 && name == "membership" {
						var member teamMemberResourceModel
						if diags := resp.State.Get(t.Context(), &member); diags.HasError() {
							t.Fatal(diags)
						}
						if member.ID.ValueString() != "42/alice" {
							t.Fatal("Read did not populate the ID for an existing membership state")
						}
					}
				})
			}
		})
	}
}

func TestImportFailsForMissingAndForbiddenRelationships(t *testing.T) {
	for _, status := range []int{200, 403, 404, 500} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected mutation: %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/api/v1/orgs/acme/teams":
					fmt.Fprint(w, `[{"id":42,"name":"developers"}]`)
				case "/api/v1/repos/acme/api":
					fmt.Fprint(w, `{"id":99,"name":"api"}`)
				default:
					w.WriteHeader(status)
					if status == 200 {
						fmt.Fprint(w, `[]`)
					} else {
						fmt.Fprint(w, `{"message":"request failed"}`)
					}
				}
			})
			assignment := importTestResource(t, &teamRepositoryResource{client: client}, "acme/developers/api")
			if !assignment.Diagnostics.HasError() {
				t.Fatal("import adopted an absent or unreadable assignment")
			}
			if status != 200 {
				membership := importTestResource(t, &teamMemberResource{client: client}, "42/alice")
				if !membership.Diagnostics.HasError() {
					t.Fatal("import adopted an absent or unreadable membership")
				}
				org := importTestResource(t, &organizationResource{client: client}, "acme")
				if !org.Diagnostics.HasError() {
					t.Fatal("import adopted an absent or unreadable organization")
				}
			}
		})
	}
}

func TestTeamRepositoryImportPaginatesTeamsAndRepositories(t *testing.T) {
	teamPages, repoPages := 0, 0
	client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("import mutated the server: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/orgs/acme/teams":
			teamPages++
			if r.URL.Query().Get("page") == "1" {
				fmt.Fprint(w, `[{"id":1,"name":"other"}]`)
			} else if r.URL.Query().Get("page") == "2" {
				fmt.Fprint(w, `[{"id":42,"name":"developers"}]`)
			} else {
				t.Errorf("unexpected team page: %s", r.URL)
				fmt.Fprint(w, `[]`)
			}
		case "/api/v1/repos/acme/api":
			fmt.Fprint(w, `{"id":99,"name":"api"}`)
		case "/api/v1/teams/42/repos":
			repoPages++
			if r.URL.Query().Get("page") == "1" {
				fmt.Fprint(w, `[{"id":1}]`)
			} else if r.URL.Query().Get("page") == "2" {
				fmt.Fprint(w, `[{"id":99}]`)
			} else {
				t.Errorf("unexpected repository page: %s", r.URL)
				fmt.Fprint(w, `[]`)
			}
		default:
			t.Errorf("unexpected path: %s", r.URL)
			http.NotFound(w, r)
		}
	})
	resp := importTestResource(t, &teamRepositoryResource{client: client}, "acme/developers/api")
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	var data teamRepositoryResourceModel
	if diags := resp.State.Get(t.Context(), &data); diags.HasError() {
		t.Fatal(diags)
	}
	if data.ID.ValueString() != "42/99" || teamPages != 2 || repoPages != 2 {
		t.Fatalf("unexpected state/pages: %+v, %d, %d", data, teamPages, repoPages)
	}
}

func TestPaginationExhaustionAndLaterPageFailures(t *testing.T) {
	for _, endpoint := range []string{"teams", "repositories"} {
		for _, status := range []int{200, 403, 500} {
			t.Run(endpoint+"/"+strconv.Itoa(status), func(t *testing.T) {
				pages := 0
				client := importTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					pages++
					if r.Method != http.MethodGet {
						t.Errorf("unexpected mutation: %s", r.Method)
					}
					w.Header().Set("Content-Type", "application/json")
					switch r.URL.Query().Get("page") {
					case "1":
						fmt.Fprint(w, `[{"id":1,"name":"other"}]`)
					case "2":
						w.WriteHeader(status)
						if status == 200 {
							fmt.Fprint(w, `[{"id":2,"name":"also-other"}]`)
						} else {
							fmt.Fprint(w, `{"message":"request failed"}`)
						}
					case "3":
						fmt.Fprint(w, `[]`)
					default:
						t.Errorf("unexpected page: %s", r.URL)
						http.NotFound(w, r)
					}
				})
				if endpoint == "teams" {
					team, diags := getOrgTeamByName(t.Context(), client, "acme", "developers")
					if team != nil || !diags.HasError() {
						t.Fatal("expected team lookup error")
					}
				} else {
					r := &teamRepositoryResource{client: client}
					state := importTestState(t, r, teamRepositoryResourceModel{ID: types.StringValue("42/99"), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)})
					resp := resource.ReadResponse{State: state}
					r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)
					if resp.Diagnostics.HasError() != (status != 200) {
						t.Fatal(resp.Diagnostics)
					}
					if resp.State.Raw.IsNull() != (status == 200) {
						t.Fatal("state removed after a failed page, or retained after exhaustion")
					}
				}
				wantPages := 2
				if status == 200 {
					wantPages = 3
				}
				if pages != wantPages {
					t.Fatalf("pages=%d; want %d", pages, wantPages)
				}
			})
		}
	}
}

func TestRelationshipDeleteOnlyRemovesTheRelationship(t *testing.T) {
	for _, name := range []string{"membership", "assignment"} {
		for _, status := range []int{204, 403, 404, 500} {
			t.Run(name+"/"+strconv.Itoa(status), func(t *testing.T) {
				deletes := 0
				client := importTestClient(t, func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					switch {
					case req.Method == http.MethodGet && req.URL.Path == "/api/v1/teams/42/repos":
						fmt.Fprint(w, `[{"id":99}]`)
					case req.Method == http.MethodGet && req.URL.Path == "/api/v1/repositories/99":
						fmt.Fprint(w, `{"id":99,"name":"api","owner":{"id":7,"login":"acme"}}`)
					case req.Method == http.MethodDelete && ((name == "membership" && req.URL.Path == "/api/v1/teams/42/members/alice") || (name == "assignment" && req.URL.Path == "/api/v1/teams/42/repos/acme/api")):
						deletes++
						w.WriteHeader(status)
						if status != 204 {
							fmt.Fprint(w, `{"message":"request failed"}`)
						}
					default:
						t.Errorf("unexpected request (parent deletion forbidden): %s %s", req.Method, req.URL)
						http.Error(w, "unexpected request", http.StatusForbidden)
					}
				})
				var r resource.Resource
				var model any
				if name == "membership" {
					r = &teamMemberResource{client: client}
					model = teamMemberResourceModel{TeamID: types.Int64Value(42), User: types.StringValue("alice")}
				} else {
					r = &teamRepositoryResource{client: client}
					model = teamRepositoryResourceModel{ID: types.StringValue("42/99"), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)}
				}
				state := importTestState(t, r, model)
				resp := resource.DeleteResponse{State: state}
				r.Delete(t.Context(), resource.DeleteRequest{State: state}, &resp)
				if resp.Diagnostics.HasError() != (status == 403 || status == 500) {
					t.Fatal(resp.Diagnostics)
				}
				if deletes != 1 {
					t.Fatalf("DELETE requests=%d; want 1", deletes)
				}
			})
		}
	}
}

type importFailingTransport struct{}

func (importFailingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("test connection failure")
}

func TestReadKeepsStateOnTransportFailure(t *testing.T) {
	client, err := forgejo.NewClient("http://unused.invalid",
		forgejo.SetForgejoVersion("17.0.0"),
		forgejo.SetHTTPClient(&http.Client{Transport: importFailingTransport{}}),
	)
	if err != nil {
		t.Fatal(err)
	}
	r := &teamRepositoryResource{client: client}
	state := importTestState(t, r, teamRepositoryResourceModel{ID: types.StringValue("42/99"), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)})
	resp := resource.ReadResponse{State: state}
	r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)
	if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(state.Raw) {
		t.Fatalf("transport failure did not preserve state: %v", resp.Diagnostics)
	}
}

func TestTeamRepositoryCreateValidatesParentsBeforeMutation(t *testing.T) {
	const team = `{"id":42,"name":"developers","organization":{"id":7,"username":"acme"}}`
	const repo = `{"id":99,"name":"api","owner":{"id":7,"login":"acme"}}`
	for _, tc := range []struct {
		name      string
		team      string
		repo      string
		putStatus int
		wantError bool
	}{
		{"valid", team, repo, 204, false},
		{"forbidden assignment", team, repo, 403, true},
		{"server error", team, repo, 500, true},
		{"different organizations", team, `{"id":99,"owner":{"id":8,"login":"elsewhere"}}`, 0, true},
		{"missing team organization", `{"id":42}`, repo, 0, true},
		{"missing repository owner", team, `{"id":99}`, 0, true},
		{"all repositories", `{"id":42,"includes_all_repositories":true}`, repo, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			puts := 0
			client := importTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case req.Method == http.MethodGet && req.URL.Path == "/api/v1/teams/42":
					fmt.Fprint(w, tc.team)
				case req.Method == http.MethodGet && req.URL.Path == "/api/v1/repositories/99":
					fmt.Fprint(w, tc.repo)
				case req.Method == http.MethodPut && req.URL.Path == "/api/v1/teams/42/repos/acme/api":
					puts++
					if tc.putStatus == 0 {
						t.Error("invalid parents reached the mutation endpoint")
						http.Error(w, "invalid parents", 500)

						return
					}
					w.WriteHeader(tc.putStatus)
					if tc.putStatus != 204 {
						fmt.Fprint(w, `{"message":"request failed"}`)
					}
				default:
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					http.NotFound(w, req)
				}
			})
			r := &teamRepositoryResource{client: client}
			planned := importTestState(t, r, teamRepositoryResourceModel{ID: types.StringUnknown(), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)})
			resp := resource.CreateResponse{State: importTestState(t, r, nil)}
			r.Create(t.Context(), resource.CreateRequest{Plan: tfsdk.Plan(planned)}, &resp)
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Fatal(resp.Diagnostics)
			}
			if (puts == 1) != (tc.putStatus != 0) {
				t.Fatalf("unexpected PUT count: %d", puts)
			}
			if !tc.wantError {
				var data teamRepositoryResourceModel
				if diags := resp.State.Get(t.Context(), &data); diags.HasError() {
					t.Fatal(diags)
				}
				if data.ID.ValueString() != "42/99" {
					t.Fatal("incorrect composite assignment ID")
				}
			}
		})
	}
}

func TestTeamRepositoryImportRejectsImplicitAccess(t *testing.T) {
	client := importTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/api/v1/orgs/acme/teams" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			http.NotFound(w, req)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":42,"name":"developers","includes_all_repositories":true}]`)
	})
	resp := importTestResource(t, &teamRepositoryResource{client: client}, "acme/developers/api")
	if !resp.Diagnostics.HasError() {
		t.Fatal("implicit all-repositories access was adopted as an individual assignment")
	}
}

func TestTeamRepositoryDeleteAlreadyMissing(t *testing.T) {
	for _, status := range []int{200, 404} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			client := importTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/api/v1/teams/42/repos" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					http.NotFound(w, req)

					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				if status == 200 {
					fmt.Fprint(w, `[]`)
				} else {
					fmt.Fprint(w, `{"message":"team not found"}`)
				}
			})
			r := &teamRepositoryResource{client: client}
			state := importTestState(t, r, teamRepositoryResourceModel{ID: types.StringValue("42/99"), TeamID: types.Int64Value(42), RepositoryID: types.Int64Value(99)})
			resp := resource.DeleteResponse{State: state}
			r.Delete(t.Context(), resource.DeleteRequest{State: state}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
		})
	}
}
