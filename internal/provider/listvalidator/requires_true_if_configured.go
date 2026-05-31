package listvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"terraform-provider-forgejo/internal/schemavalidator"
)

// RequiresTrueIfConfigured checks that any Bool values in the paths described by the
// path.Expression are true if the current attribute value is configured to a non-empty
// list.
func RequiresTrueIfConfigured(expressions ...path.Expression) validator.List {
	return &schemavalidator.RequiresTrueIfConfiguredValidator{
		Expressions: expressions,
	}
}
