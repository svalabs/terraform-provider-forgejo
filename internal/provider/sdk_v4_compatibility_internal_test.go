package provider

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPersonalAccessTokenDeleteSDKV4Response(t *testing.T) {
	for _, status := range []int{204, 403, 404, 500, 0} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			requests := 0
			client := importTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				requests++
				if req.Method != http.MethodDelete || req.URL.Path != "/api/v1/users/alice/tokens/42" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				if status != http.StatusNoContent {
					fmt.Fprint(w, `{"message":"test deletion failure"}`)
				}
			})
			if status == 0 {
				client.SetHTTPClient(&http.Client{Transport: importFailingTransport{}})
			}
			r := &personalAccessTokenResource{client: client}
			data := personalAccessTokenResourceModel{ID: types.Int64Value(42), User: types.StringValue("alice"), Name: types.StringValue("test"), Scopes: types.SetNull(types.StringType)}
			state := importTestState(t, r, data)
			resp := resource.DeleteResponse{State: state}
			r.Delete(t.Context(), resource.DeleteRequest{State: state}, &resp)
			if resp.Diagnostics.HasError() != (status != http.StatusNoContent) {
				t.Fatal(resp.Diagnostics)
			}
			if !resp.State.Raw.Equal(state.Raw) {
				t.Fatal("token deletion discarded state")
			}
			if (requests == 1) != (status != 0) {
				t.Fatalf("unexpected request count: %d", requests)
			}
		})
	}
}
