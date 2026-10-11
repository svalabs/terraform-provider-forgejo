package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// actionSecretDataAttributes returns the schema attributes holding the value
// of an action secret: either "data" (persisted in state) or the write-only
// "data_wo" (never persisted) together with "data_wo_version".
func actionSecretDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"data": schema.StringAttribute{
			MarkdownDescription: "Data of the secret. The value is stored in the Terraform state. **Note**: Exactly one of `data` or `data_wo` must be specified.",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.Expressions{
					path.MatchRoot("data_wo"),
				}...),
			},
		},
		"data_wo": schema.StringAttribute{
			MarkdownDescription: "Write-only data of the secret. The value is never stored in the Terraform state or plan; change `data_wo_version` to update it. Requires Terraform >= 1.11 or OpenTofu >= 1.11. **Note**: Exactly one of `data` or `data_wo` must be specified.",
			Optional:            true,
			Sensitive:           true,
			WriteOnly:           true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.Expressions{
					path.MatchRoot("data"),
				}...),
			},
		},
		"data_wo_version": schema.Int64Attribute{
			MarkdownDescription: "Version of `data_wo`. Terraform cannot detect changes of write-only values, so the secret is only updated when this value changes.",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AlsoRequires(path.Expressions{
					path.MatchRoot("data_wo"),
				}...),
			},
		},
	}
}

// actionSecretData returns the secret value to send to the API: "data_wo" from
// the configuration (write-only values are only available there) or "data"
// from the plan.
func actionSecretData(ctx context.Context, config tfsdk.Config, data types.String) (string, diag.Diagnostics) {
	var dataWO types.String

	diags := config.GetAttribute(ctx, path.Root("data_wo"), &dataWO)
	if diags.HasError() {
		return "", diags
	}

	if !dataWO.IsNull() && !dataWO.IsUnknown() {
		return dataWO.ValueString(), diags
	}

	return data.ValueString(), diags
}

// mergeAttributes merges schema attribute maps; later maps win on conflicts.
func mergeAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	merged := map[string]schema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}
