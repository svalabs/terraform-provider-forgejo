package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"terraform-provider-forgejo/internal/schemavalidator"
)

// RequiresTrueIfConfigured checks that any Bool values in the paths described by the
// path.Expression are true if the current attribute value is configured.
// If you require the value described by the path.Expression to be set,
// combine this validator with the "AlsoRequires" validator.
func RequiresTrueIfConfigured(expressions ...path.Expression) validator.Bool {
	return &schemavalidator.RequiresTrueIfConfiguredValidator{
		Expressions: expressions,
	}
}
