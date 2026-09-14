package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"terraform-provider-forgejo/internal/schemavalidator"
)

func TestRequiresTrueIfConfiguredValidatorValidateBool(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target": schema.BoolAttribute{Optional: true},
			"self":   schema.BoolAttribute{Optional: true},
		},
	}

	testCases := map[string]struct {
		self      tftypes.Value
		target    tftypes.Value
		configVal types.Bool
		expErrors int
	}{
		"self-null": {
			self:      tftypes.NewValue(tftypes.Bool, nil),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.BoolNull(),
		},
		"self-unknown": {
			self:      tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.BoolUnknown(),
		},
		"self-false-target-false": {
			self:      tftypes.NewValue(tftypes.Bool, false),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.BoolValue(false),
			expErrors: 1,
		},
		"self-true-target-false": {
			self:      tftypes.NewValue(tftypes.Bool, true),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.BoolValue(true),
			expErrors: 1,
		},
		"self-true-target-true": {
			self:      tftypes.NewValue(tftypes.Bool, true),
			target:    tftypes.NewValue(tftypes.Bool, true),
			configVal: types.BoolValue(true),
		},
		"self-false-target-null": {
			self:      tftypes.NewValue(tftypes.Bool, false),
			target:    tftypes.NewValue(tftypes.Bool, nil),
			configVal: types.BoolValue(false),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.BoolRequest{
				Path:           path.Root("self"),
				PathExpression: path.MatchRoot("self"),
				ConfigValue:    tc.configVal,
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"target": tftypes.Bool,
							"self":   tftypes.Bool,
						},
					}, map[string]tftypes.Value{
						"target": tc.target,
						"self":   tc.self,
					}),
				},
			}
			resp := &validator.BoolResponse{}

			v := schemavalidator.RequiresTrueIfConfiguredValidator{
				Expressions: path.Expressions{path.MatchRoot("target")},
			}
			v.ValidateBool(context.Background(), req, resp)

			if len(resp.Diagnostics) != tc.expErrors {
				t.Fatalf("expected %d errors, got %d: %v", tc.expErrors, len(resp.Diagnostics), resp.Diagnostics)
			}
		})
	}
}

func TestRequiresTrueIfConfiguredValidatorValidateString(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target": schema.BoolAttribute{Optional: true},
			"self":   schema.StringAttribute{Optional: true},
		},
	}

	testCases := map[string]struct {
		self      tftypes.Value
		target    tftypes.Value
		configVal types.String
		expErrors int
	}{
		"self-null": {
			self:      tftypes.NewValue(tftypes.String, nil),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.StringNull(),
		},
		"self-unknown": {
			self:      tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.StringUnknown(),
		},
		"self-empty-target-false": {
			self:      tftypes.NewValue(tftypes.String, ""),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.StringValue(""),
		},
		"self-nonempty-target-false": {
			self:      tftypes.NewValue(tftypes.String, "value"),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.StringValue("value"),
			expErrors: 1,
		},
		"self-nonempty-target-true": {
			self:      tftypes.NewValue(tftypes.String, "value"),
			target:    tftypes.NewValue(tftypes.Bool, true),
			configVal: types.StringValue("value"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:           path.Root("self"),
				PathExpression: path.MatchRoot("self"),
				ConfigValue:    tc.configVal,
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"target": tftypes.Bool,
							"self":   tftypes.String,
						},
					}, map[string]tftypes.Value{
						"target": tc.target,
						"self":   tc.self,
					}),
				},
			}
			resp := &validator.StringResponse{}

			v := schemavalidator.RequiresTrueIfConfiguredValidator{
				Expressions: path.Expressions{path.MatchRoot("target")},
			}
			v.ValidateString(context.Background(), req, resp)

			if len(resp.Diagnostics) != tc.expErrors {
				t.Fatalf("expected %d errors, got %d: %v", tc.expErrors, len(resp.Diagnostics), resp.Diagnostics)
			}
		})
	}
}

func TestRequiresTrueIfConfiguredValidatorValidateObject(t *testing.T) {
	t.Parallel()

	objType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"inner": tftypes.Bool,
		},
	}
	objAttrTypes := map[string]attr.Type{
		"inner": types.BoolType,
	}

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target": schema.BoolAttribute{Optional: true},
			"self": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"inner": schema.BoolAttribute{Optional: true},
				},
			},
		},
	}

	testCases := map[string]struct {
		self      tftypes.Value
		target    tftypes.Value
		configVal types.Object
		expErrors int
	}{
		"self-null": {
			self:      tftypes.NewValue(objType, nil),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ObjectNull(objAttrTypes),
		},
		"self-unknown": {
			self:      tftypes.NewValue(objType, tftypes.UnknownValue),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ObjectUnknown(objAttrTypes),
		},
		"self-empty-object-target-false": {
			self: tftypes.NewValue(objType, map[string]tftypes.Value{
				"inner": tftypes.NewValue(tftypes.Bool, nil),
			}),
			target: tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"inner": types.BoolNull(),
			}),
			expErrors: 1,
		},
		"self-empty-object-target-true": {
			self: tftypes.NewValue(objType, map[string]tftypes.Value{
				"inner": tftypes.NewValue(tftypes.Bool, nil),
			}),
			target: tftypes.NewValue(tftypes.Bool, true),
			configVal: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"inner": types.BoolNull(),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.ObjectRequest{
				Path:           path.Root("self"),
				PathExpression: path.MatchRoot("self"),
				ConfigValue:    tc.configVal,
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"target": tftypes.Bool,
							"self":   objType,
						},
					}, map[string]tftypes.Value{
						"target": tc.target,
						"self":   tc.self,
					}),
				},
			}
			resp := &validator.ObjectResponse{}

			v := schemavalidator.RequiresTrueIfConfiguredValidator{
				Expressions: path.Expressions{path.MatchRoot("target")},
			}
			v.ValidateObject(context.Background(), req, resp)

			if len(resp.Diagnostics) != tc.expErrors {
				t.Fatalf("expected %d errors, got %d: %v", tc.expErrors, len(resp.Diagnostics), resp.Diagnostics)
			}
		})
	}
}

func TestRequiresTrueIfConfiguredValidatorValidateList(t *testing.T) {
	t.Parallel()

	listType := tftypes.List{ElementType: tftypes.String}

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target": schema.BoolAttribute{Optional: true},
			"self": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}

	testCases := map[string]struct {
		self      tftypes.Value
		target    tftypes.Value
		configVal types.List
		expErrors int
	}{
		"self-null": {
			self:      tftypes.NewValue(listType, nil),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ListNull(types.StringType),
		},
		"self-unknown": {
			self:      tftypes.NewValue(listType, tftypes.UnknownValue),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ListUnknown(types.StringType),
		},
		"self-empty-target-false": {
			self:      tftypes.NewValue(listType, []tftypes.Value{}),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ListValueMust(types.StringType, []attr.Value{}),
		},
		"self-nonempty-target-false": {
			self: tftypes.NewValue(listType, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "a"),
			}),
			target: tftypes.NewValue(tftypes.Bool, false),
			configVal: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("a"),
			}),
			expErrors: 1,
		},
		"self-nonempty-target-true": {
			self: tftypes.NewValue(listType, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "a"),
			}),
			target: tftypes.NewValue(tftypes.Bool, true),
			configVal: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("a"),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.ListRequest{
				Path:           path.Root("self"),
				PathExpression: path.MatchRoot("self"),
				ConfigValue:    tc.configVal,
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"target": tftypes.Bool,
							"self":   listType,
						},
					}, map[string]tftypes.Value{
						"target": tc.target,
						"self":   tc.self,
					}),
				},
			}
			resp := &validator.ListResponse{}

			v := schemavalidator.RequiresTrueIfConfiguredValidator{
				Expressions: path.Expressions{path.MatchRoot("target")},
			}
			v.ValidateList(context.Background(), req, resp)

			if len(resp.Diagnostics) != tc.expErrors {
				t.Fatalf("expected %d errors, got %d: %v", tc.expErrors, len(resp.Diagnostics), resp.Diagnostics)
			}
		})
	}
}

func TestRequiresTrueIfConfiguredValidatorValidateSet(t *testing.T) {
	t.Parallel()

	setType := tftypes.Set{ElementType: tftypes.String}

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target": schema.BoolAttribute{Optional: true},
			"self": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}

	testCases := map[string]struct {
		self      tftypes.Value
		target    tftypes.Value
		configVal types.Set
		expErrors int
	}{
		"self-null": {
			self:      tftypes.NewValue(setType, nil),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.SetNull(types.StringType),
		},
		"self-unknown": {
			self:      tftypes.NewValue(setType, tftypes.UnknownValue),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.SetUnknown(types.StringType),
		},
		"self-empty-target-false": {
			self:      tftypes.NewValue(setType, []tftypes.Value{}),
			target:    tftypes.NewValue(tftypes.Bool, false),
			configVal: types.SetValueMust(types.StringType, []attr.Value{}),
		},
		"self-nonempty-target-false": {
			self: tftypes.NewValue(setType, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "a"),
			}),
			target: tftypes.NewValue(tftypes.Bool, false),
			configVal: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("a"),
			}),
			expErrors: 1,
		},
		"self-nonempty-target-true": {
			self: tftypes.NewValue(setType, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "a"),
			}),
			target: tftypes.NewValue(tftypes.Bool, true),
			configVal: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("a"),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.SetRequest{
				Path:           path.Root("self"),
				PathExpression: path.MatchRoot("self"),
				ConfigValue:    tc.configVal,
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"target": tftypes.Bool,
							"self":   setType,
						},
					}, map[string]tftypes.Value{
						"target": tc.target,
						"self":   tc.self,
					}),
				},
			}
			resp := &validator.SetResponse{}

			v := schemavalidator.RequiresTrueIfConfiguredValidator{
				Expressions: path.Expressions{path.MatchRoot("target")},
			}
			v.ValidateSet(context.Background(), req, resp)

			if len(resp.Diagnostics) != tc.expErrors {
				t.Fatalf("expected %d errors, got %d: %v", tc.expErrors, len(resp.Diagnostics), resp.Diagnostics)
			}
		})
	}
}
