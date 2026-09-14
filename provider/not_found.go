package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

/*
 * The Terraform Plugin Framework requires Read to distinguish "this object is
 * gone" from "the read failed". A vanished object must be removed from state
 * with resp.State.RemoveResource so the next plan is a create; reporting it as
 * an error instead leaves the resource permanently unreconcilable, because it
 * can never be created without first being read.
 *
 * The lookup helpers in this package are shared between resources and data
 * sources, and a data source referencing a missing object *should* fail. So
 * rather than change what the helpers return, a not-found is tagged with a
 * distinguishable diagnostic type: data sources keep surfacing it as an error
 * verbatim, while resource Read implementations test for it with isNotFound
 * and drop the object from state.
 */

// notFoundDiagnostic is an error diagnostic raised because the remote object
// does not exist, as opposed to the read having failed.
type notFoundDiagnostic struct {
	summary string
	detail  string
}

func (d notFoundDiagnostic) Severity() diag.Severity { return diag.SeverityError }
func (d notFoundDiagnostic) Summary() string         { return d.summary }
func (d notFoundDiagnostic) Detail() string          { return d.detail }

// Equal compares by content rather than by type, so a notFoundDiagnostic and a
// plain error diagnostic carrying the same text stay interchangeable.
func (d notFoundDiagnostic) Equal(other diag.Diagnostic) bool {
	return other != nil &&
		other.Severity() == d.Severity() &&
		other.Summary() == d.Summary() &&
		other.Detail() == d.Detail()
}

// newNotFoundDiagnostic builds a diagnostic marking the remote object as absent.
func newNotFoundDiagnostic(summary, detail string) diag.Diagnostic {
	return notFoundDiagnostic{summary: summary, detail: detail}
}

// addReadError appends msg as an error diagnostic, tagging HTTP 404 responses
// as not-found.
func addReadError(diags *diag.Diagnostics, res *forgejo.Response, summary, msg string) {
	if res != nil && res.Response != nil && res.StatusCode == 404 {
		diags.Append(newNotFoundDiagnostic(summary, msg))

		return
	}

	diags.AddError(summary, msg)
}

// isNotFound reports whether diags say the remote object does not exist.
func isNotFound(diags diag.Diagnostics) bool {
	for _, d := range diags {
		if _, ok := d.(notFoundDiagnostic); ok {
			return true
		}
	}

	return false
}
