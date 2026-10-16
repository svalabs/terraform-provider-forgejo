// This is an internal httptest-based test rather than an acceptance test: the
// invitation flow needs `[service] ADD_MEMBERS_BY_INVITATIONS`, which the
// acceptance environment does not enable, and it exercises the unexported
// setTeamMember helper, so it must live in package provider.
package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

func TestSetTeamMemberVerifiesMembership(t *testing.T) {
	for _, tc := range []struct {
		name       string
		addStatus  int
		readStatus int
		wantError  string
		wantReads  int32
	}{
		{"accepted", 204, 200, "", 1},
		{"invitation sent", 204, 404, "must accept the team invitation", 1},
		{"already invited", 409, 404, "team invitation for user 'member' to team 42 is already pending", 0},
		{"verification failure", 204, 500, "could not be verified", 1},
		{"missing user or team", 404, 404, "not found", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var reads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/teams/42/members/member" {
					t.Errorf("unexpected path: %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodPut {
					w.WriteHeader(tc.addStatus)
					return
				}
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", r.Method)
				}
				reads.Add(1)
				w.WriteHeader(tc.readStatus)
				if tc.readStatus == http.StatusOK {
					if err := json.NewEncoder(w).Encode(map[string]string{"login": "member"}); err != nil {
						t.Error(err)
					}
				}
			}))
			defer server.Close()
			client, err := forgejo.NewClient(server.URL, forgejo.SetForgejoVersion("17.0.0"))
			if err != nil {
				t.Fatal(err)
			}
			diags := setTeamMember(t.Context(), client, 42, "member")
			if tc.wantError == "" {
				if diags.HasError() {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
			} else if !diags.HasError() || !strings.Contains(diags[0].Summary()+" "+diags[0].Detail(), tc.wantError) {
				t.Fatalf("expected diagnostic containing %q, got %v", tc.wantError, diags)
			}
			if reads.Load() != tc.wantReads {
				t.Errorf("membership reads: got %d, want %d", reads.Load(), tc.wantReads)
			}
		})
	}
}
