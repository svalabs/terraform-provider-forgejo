package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ validator.Bool = &RequiresTrueIfConfiguredValidator{}
	_ validator.List = &RequiresTrueIfConfiguredValidator{}
	_ validator.Set  = &RequiresTrueIfConfiguredValidator{}
)

// RequiresTrueIfConfiguredValidator is the underlying type implementing RequiresTrueIfConfigured.
type RequiresTrueIfConfiguredValidator struct {
	Expressions path.Expressions
}

// Description returns a plaintext string describing the validator.
func (v RequiresTrueIfConfiguredValidator) Description(_ context.Context) string {
	return "If this attribute is configured, the attribute(s) at the given path(s) must be set to 'true'."
}

// MarkdownDescription returns a Markdown-formatted string describing the validator.
func (v RequiresTrueIfConfiguredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateList performs the validation logic for the validator if the attribute type is a List.
func (v RequiresTrueIfConfiguredValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}
	if req.ConfigValue.Length(opts) == 0 {
		return
	}

	validateReq := requiresTrueIfConfiguredValidatorRequest{
		Config:         req.Config,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}

	v.validate(
		ctx,
		validateReq,
		&resp.Diagnostics,
		fmt.Sprintf("If %s is not empty, %%s must be 'true'", req.Path),
	)
}

// ValidateSet performs the validation logic for the validator if the attribute type is a Set.
func (v RequiresTrueIfConfiguredValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	opts := basetypes.CollectionLengthOptions{
		UnhandledNullAsZero:    true,
		UnhandledUnknownAsZero: true,
	}
	if req.ConfigValue.Length(opts) == 0 {
		return
	}

	validateReq := requiresTrueIfConfiguredValidatorRequest{
		Config:         req.Config,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}

	v.validate(
		ctx,
		validateReq,
		&resp.Diagnostics,
		fmt.Sprintf("If %s is not empty, %%s must be 'true'", req.Path),
	)
}

// ValidateBool performs the validation logic for the validator if the attribute type is a Bool.
func (v RequiresTrueIfConfiguredValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !req.ConfigValue.ValueBool() {
		return
	}

	validateReq := requiresTrueIfConfiguredValidatorRequest{
		Config:         req.Config,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}

	v.validate(
		ctx,
		validateReq,
		&resp.Diagnostics,
		fmt.Sprintf("If %s is '%t', %%s must be 'true'", req.Path, req.ConfigValue.ValueBool()),
	)
}

// validate is the shared validation function for the validator.
func (v RequiresTrueIfConfiguredValidator) validate(ctx context.Context, req requiresTrueIfConfiguredValidatorRequest, diags *diag.Diagnostics, messageFmt string) {
	for _, expression := range req.PathExpression.MergeExpressions(v.Expressions...) {
		matchedPaths, d := req.Config.PathMatches(ctx, expression)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}

		for _, matchedPath := range matchedPaths {
			var matchedPathValue attr.Value

			diags.Append(req.Config.GetAttribute(
				ctx,
				matchedPath,
				&matchedPathValue,
			)...)
			if diags.HasError() {
				continue
			}

			if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
				continue
			}

			var matchedPathConfig types.Bool
			diags.Append(tfsdk.ValueAs(
				ctx,
				matchedPathValue,
				&matchedPathConfig,
			)...)
			if diags.HasError() {
				continue
			}

			if !matchedPathConfig.ValueBool() {
				diags.AddAttributeError(
					matchedPath,
					"Invalid Attribute Combination",
					fmt.Sprintf(messageFmt, matchedPath.String()),
				)
			}
		}
	}
}

type requiresTrueIfConfiguredValidatorRequest struct {
	Config         tfsdk.Config
	Path           path.Path
	PathExpression path.Expression
}
