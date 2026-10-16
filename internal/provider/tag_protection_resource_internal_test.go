// This test lives in package provider (not provider_test) because it needs
// access to the unexported tagProtectionResource and its model in order to
// drive Read directly against an httptest server. Acceptance tests cannot
// distinguish a 403 from a 404 response, since both require a Forgejo
// instance in a very specific state.
package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// TestTagProtectionResourceRead verifies that Read removes the resource from
// state when the tag protection no longer exists (404), but keeps the state
// and reports an error when access is forbidden (403).
func TestTagProtectionResourceRead(t *testing.T) {
	testCases := []struct {
		name         string
		repoStatus   int
		ruleStatus   int
		expectRemove bool
		expectError  bool
	}{
		{name: "present", repoStatus: http.StatusOK, ruleStatus: http.StatusOK},
		{name: "repository removed", repoStatus: http.StatusNotFound, ruleStatus: http.StatusOK, expectError: true},
		{name: "repository forbidden", repoStatus: http.StatusForbidden, ruleStatus: http.StatusOK, expectError: true},
		{name: "tag protection removed", repoStatus: http.StatusOK, ruleStatus: http.StatusNotFound, expectRemove: true},
		{name: "tag protection forbidden", repoStatus: http.StatusOK, ruleStatus: http.StatusForbidden, expectError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				var body string
				switch r.URL.Path {
				case "/api/v1/repositories/1":
					w.WriteHeader(tc.repoStatus)
					body = `{"id":1,"name":"repo","owner":{"login":"owner"}}`
				case "/api/v1/repos/owner/repo/tag_protections/2":
					w.WriteHeader(tc.ruleStatus)
					body = `{"id":2,"name_pattern":"v*","whitelist_usernames":[],"whitelist_teams":[]}`
				default:
					t.Errorf("unexpected request path: %s", r.URL.Path)
					http.NotFound(w, r)

					return
				}

				if _, err := w.Write([]byte(body)); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()

			client, err := forgejo.NewClient(server.URL, forgejo.SetForgejoVersion("17.0.0"))
			if err != nil {
				t.Fatal(err)
			}

			r := tagProtectionResource{client: client}

			var schemaResp resource.SchemaResponse
			r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)

			state := tfsdk.State{Schema: schemaResp.Schema}
			diags := state.Set(t.Context(), tagProtectionResourceModel{
				ID:                 types.Int64Value(2),
				RepositoryID:       types.Int64Value(1),
				NamePattern:        types.StringValue("v*"),
				WhitelistUsernames: types.SetValueMust(types.StringType, nil),
				WhitelistTeams:     types.SetValueMust(types.StringType, nil),
				CreatedAt:          types.StringValue("2025-01-01T00:00:00Z"),
				UpdatedAt:          types.StringValue("2025-01-01T00:00:00Z"),
			})
			if diags.HasError() {
				t.Fatal(diags)
			}

			resp := resource.ReadResponse{State: state}
			r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)

			if resp.Diagnostics.HasError() != tc.expectError {
				t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if resp.State.Raw.IsNull() != tc.expectRemove {
				t.Errorf("state removed = %v, want %v", resp.State.Raw.IsNull(), tc.expectRemove)
			}
		})
	}
}
