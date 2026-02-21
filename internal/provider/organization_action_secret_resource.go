package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &organizationActionSecretResource{}
	_ resource.ResourceWithConfigure = &organizationActionSecretResource{}
)

// organizationActionSecretResource is the resource implementation.
type organizationActionSecretResource struct {
	client *forgejo.Client
}

// organizationActionSecretResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#CreateSecretOption
type organizationActionSecretResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Name         types.String `tfsdk:"name"`
	Data         types.String `tfsdk:"data"`
	CreatedAt    types.String `tfsdk:"created_at"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *organizationActionSecretResourceModel) from(s *forgejo.Secret) {
	if s == nil {
		return
	}

	m.CreatedAt = types.StringValue(s.Created.Format(time.RFC3339))
}

// to is a helper function to save Terraform data model into an API struct.
func (m *organizationActionSecretResourceModel) to(o *forgejo.CreateSecretOption) {
	if o == nil {
		o = new(forgejo.CreateSecretOption)
	}

	o.Name = m.Name.ValueString()
	o.Data = m.Data.ValueString()
}

// Metadata returns the resource type name.
func (r *organizationActionSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_action_secret"
}

// Schema defines the schema for the resource.
func (r *organizationActionSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo organization action secret resource.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "Name of the organization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the secret.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 30),
				},
			},
			"data": schema.StringAttribute{
				Description: "Data of the secret.",
				Required:    true,
				Sensitive:   true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the secret was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *organizationActionSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *organizationActionSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create organization action secret resource"))

	var data organizationActionSecretResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create organization action secret", map[string]any{
		"org":  data.Organization.ValueString(),
		"name": data.Name.ValueString(),
		"data": data.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateSecretOption{}
	data.to(&opts)

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new organization action secret
	res, err := r.client.CreateOrgActionSecret(
		data.Organization.ValueString(),
		opts,
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 400:
				msg = fmt.Sprintf("Generic error: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Organization with name %s not found: %s",
					data.Organization.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create organization action secret", msg)

		return
	}

	// Use Forgejo client to get organization action secret
	secret, diags := r.getSecret(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.from(secret)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *organizationActionSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read organization action secret resource"))

	var data organizationActionSecretResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get organization action secret
	secret, diags := r.getSecret(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.from(secret)

	/*
	 * The secret exists, so we re-save the state from the prior state data.
	 * This is to signal to Terraform that the resource still exists without
	 * overriding the user's configuration casing for the name.
	 */
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *organizationActionSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update organization action secret resource"))

	var data organizationActionSecretResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update organization action secret", map[string]any{
		"org":  data.Organization.ValueString(),
		"name": data.Name.ValueString(),
		"data": data.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateSecretOption{}
	data.to(&opts)

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to update organization action secret
	res, err := r.client.CreateOrgActionSecret(
		data.Organization.ValueString(),
		opts,
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 400:
				msg = fmt.Sprintf("Generic error: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Organization with name %s not found: %s",
					data.Organization.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update organization action secret", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *organizationActionSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete organization action secret resource"))

	var data organizationActionSecretResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete organization action secret", map[string]any{
		"org":  data.Organization.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.AddWarning(
		"Resource cannot be deleted from Forgejo",
		fmt.Sprintf(
			"The Forgejo SDK does not currently support deleting organization action secrets. "+
				"Secret with org %s and name %s will be removed from Terraform state, but will remain in Forgejo.",
			data.Organization.String(),
			data.Name.String(),
		),
	)
}

// NewOrganizationActionSecretResource is a helper function to simplify the provider implementation.
func NewOrganizationActionSecretResource() resource.Resource {
	return &organizationActionSecretResource{}
}

// getSecret returns the secret with the given name from the organization.
func (r *organizationActionSecretResource) getSecret(ctx context.Context, data *organizationActionSecretResourceModel) (*forgejo.Secret, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "List organization action secrets", map[string]any{
		"org":  data.Organization.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to list organization action secrets
	secrets, res, err := r.client.ListOrgActionSecret(
		data.Organization.ValueString(),
		forgejo.ListOrgActionSecretOption{},
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 404:
				msg = fmt.Sprintf(
					"Organization action secrets with org %s not found: %s",
					data.Organization.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		diags.AddError("Unable to list organization action secrets", msg)

		return nil, diags
	}

	// Search for organization action secrets with given name
	idx := slices.IndexFunc(secrets, func(s *forgejo.Secret) bool {
		return strings.EqualFold(s.Name, data.Name.ValueString())
	})
	if idx == -1 {
		diags.AddError(
			"Unable to find organization action secret by name",
			fmt.Sprintf(
				"Organization action secret with org %s and name %s not found",
				data.Organization.String(),
				data.Name.String(),
			),
		)

		return nil, diags
	}

	return secrets[idx], diags
}
