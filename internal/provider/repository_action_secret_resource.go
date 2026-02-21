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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repositoryActionSecretResource{}
	_ resource.ResourceWithConfigure = &repositoryActionSecretResource{}
)

// repositoryActionSecretResource is the resource implementation.
type repositoryActionSecretResource struct {
	client *forgejo.Client
}

// repositoryActionSecretResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#CreateSecretOption
type repositoryActionSecretResourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	Name         types.String `tfsdk:"name"`
	Data         types.String `tfsdk:"data"`
	CreatedAt    types.String `tfsdk:"created_at"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *repositoryActionSecretResourceModel) from(s *forgejo.Secret) {
	if s == nil {
		return
	}

	m.CreatedAt = types.StringValue(s.Created.Format(time.RFC3339))
}

// to is a helper function to save Terraform data model into an API struct.
func (m *repositoryActionSecretResourceModel) to(o *forgejo.CreateSecretOption) {
	if o == nil {
		o = new(forgejo.CreateSecretOption)
	}

	o.Name = m.Name.ValueString()
	o.Data = m.Data.ValueString()
}

// Metadata returns the resource type name.
func (r *repositoryActionSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_action_secret"
}

// Schema defines the schema for the resource.
func (r *repositoryActionSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository action secret resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
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
func (r *repositoryActionSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryActionSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Create repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
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

	// Use Forgejo client to create new repository action secret
	res, err := r.client.CreateRepoActionSecret(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Repository with owner %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create repository action secret", msg)

		return
	}

	// Use Forgejo client to get repository action secret
	secret, diags := r.getSecret(
		ctx,
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		&data,
	)
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
func (r *repositoryActionSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	// Use Forgejo client to get repository action secret
	secret, diags := r.getSecret(
		ctx,
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		&data,
	)
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
func (r *repositoryActionSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Update repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
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

	// Use Forgejo client to update repository action secret
	res, err := r.client.CreateRepoActionSecret(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Repository with owner %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update repository action secret", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryActionSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Delete repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.AddWarning(
		"Resource cannot be deleted from Forgejo",
		fmt.Sprintf(
			"The Forgejo SDK does not currently support deleting repository action secrets. "+
				"Secret with owner %s repo %s and name %s will be removed from Terraform state, but will remain in Forgejo.",
			repo.Owner.String(),
			repo.Name.String(),
			data.Name.String(),
		),
	)
}

// NewRepositoryActionSecretResource is a helper function to simplify the provider implementation.
func NewRepositoryActionSecretResource() resource.Resource {
	return &repositoryActionSecretResource{}
}

// getSecret returns the secret with the given name from the repository.
func (r *repositoryActionSecretResource) getSecret(ctx context.Context, owner, repoName string, data *repositoryActionSecretResourceModel) (*forgejo.Secret, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "List repository action secrets", map[string]any{
		"user": owner,
		"repo": repoName,
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to list repository action secrets
	secrets, res, err := r.client.ListRepoActionSecret(
		owner,
		repoName,
		forgejo.ListRepoActionSecretOption{},
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
					"Repository action secrets with user '%s' and repo '%s' not found: %s",
					owner,
					repoName,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		diags.AddError("Unable to list repository action secrets", msg)

		return nil, diags
	}

	// Search for repository action secrets with given name
	idx := slices.IndexFunc(secrets, func(s *forgejo.Secret) bool {
		return strings.EqualFold(s.Name, data.Name.ValueString())
	})
	if idx == -1 {
		diags.AddError(
			"Unable to find repository action secret by name",
			fmt.Sprintf(
				"Repository action secret with user '%s' repo '%s' and name %s not found",
				owner,
				repoName,
				data.Name.String(),
			),
		)

		return nil, diags
	}

	return secrets[idx], diags
}
