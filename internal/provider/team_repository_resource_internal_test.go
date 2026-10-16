// These tests live in package provider (not provider_test) because they
// exercise the unexported Read implementation against a fake Forgejo API.
// The behaviour they cover -- pagination without a 'Link' header and the
// distinction between a missing (404) and a forbidden (403) team -- can not
// be provoked reliably in acceptance tests against a real server.
package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// newTeamRepositoryTestClient returns a Forgejo client pointed at the given
// fake API server.
func newTeamRepositoryTestClient(t *testing.T, server *httptest.Server) *forgejo.Client {
	t.Helper()

	client, err := forgejo.NewClient(server.URL, forgejo.SetForgejoVersion("17.0.0"))
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// readTeamRepository runs the resource Read method against the given state
// and returns the resulting response.
func readTeamRepository(t *testing.T, client *forgejo.Client, teamID, repositoryID int64) resource.ReadResponse {
	t.Helper()

	r := teamRepositoryResource{client: client}

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)

	state := tfsdk.State{Schema: schemaResp.Schema}
	diags := state.Set(t.Context(), teamRepositoryResourceModel{
		ID:           types.StringValue(teamRepositoryID(teamID, repositoryID)),
		TeamID:       types.Int64Value(teamID),
		RepositoryID: types.Int64Value(repositoryID),
	})
	if diags.HasError() {
		t.Fatal(diags)
	}

	resp := resource.ReadResponse{State: state}
	r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)

	return resp
}

func TestTeamRepositoryReadFollowsPaginationWithoutLinkHeader(t *testing.T) {
	// Three pages of repositories; the server caps the page size at 2 and,
	// like Forgejo, sets only 'X-Total-Count' and no 'Link' header.
	pages := map[string]string{
		"1": `[{"id":1},{"id":2}]`,
		"2": `[{"id":3},{"id":4}]`,
		"3": `[{"id":5}]`,
		"4": `[]`,
	}
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/api/v1/teams/42/repos" {
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)

			return
		}
		body, ok := pages[r.URL.Query().Get("page")]
		if !ok {
			t.Errorf("unexpected page %q", r.URL.Query().Get("page"))
			http.NotFound(w, r)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "5")
		fmt.Fprint(w, body)
	}))
	defer server.Close()
	client := newTeamRepositoryTestClient(t, server)

	for _, tc := range []struct {
		repositoryID int64
		found        bool
		requests     int
	}{
		// Found on the first page: no further pages are requested
		{1, true, 1},
		// Found on the last non-empty page
		{5, true, 3},
		// Not found: all pages up to the empty one are requested
		{6, false, 4},
	} {
		t.Run(strconv.FormatInt(tc.repositoryID, 10), func(t *testing.T) {
			requests = 0
			resp := readTeamRepository(t, client, 42, tc.repositoryID)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if removed := resp.State.Raw.IsNull(); removed == tc.found {
				t.Errorf("state removed=%v; want %v", removed, !tc.found)
			}
			if requests != tc.requests {
				t.Errorf("requests=%d; want %d", requests, tc.requests)
			}
		})
	}
}

func TestTeamRepositoryReadDistinguishesMissingAndForbidden(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		removed bool
		failure bool
	}{
		// The association exists: state is kept
		{"present", http.StatusOK, `[{"id":99}]`, false, false},
		// The association is gone: state is removed
		{"association missing", http.StatusOK, `[]`, true, false},
		// The team is gone: state is removed
		{"team missing", http.StatusNotFound, `{"message":"not found"}`, true, false},
		// Access is denied: state is kept and an error is reported
		{"team forbidden", http.StatusForbidden, `{"message":"forbidden"}`, false, true},
		// Any other failure: state is kept and an error is reported
		{"server error", http.StatusInternalServerError, `{"message":"error"}`, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/teams/42/repos" {
					t.Errorf("unexpected path %s", r.URL.Path)
					http.NotFound(w, r)

					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}))
			defer server.Close()
			client := newTeamRepositoryTestClient(t, server)

			resp := readTeamRepository(t, client, 42, 99)
			if resp.Diagnostics.HasError() != tc.failure {
				t.Errorf("diagnostics error=%v; want %v: %v", resp.Diagnostics.HasError(), tc.failure, resp.Diagnostics)
			}
			if removed := resp.State.Raw.IsNull(); removed != tc.removed {
				t.Errorf("state removed=%v; want %v", removed, tc.removed)
			}
		})
	}
}
