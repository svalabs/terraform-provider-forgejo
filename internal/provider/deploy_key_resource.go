package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deployKeyResource{}
	_ resource.ResourceWithConfigure = &deployKeyResource{}
)

// deployKeyResource is the resource implementation.
type deployKeyResource struct {
	client *forgejo.Client
}

// deployKeyResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#DeployKey
type deployKeyResourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	KeyID        types.Int64  `tfsdk:"key_id"`
	Key          types.String `tfsdk:"key"`
	URL          types.String `tfsdk:"url"`
	Title        types.String `tfsdk:"title"`
	Fingerprint  types.String `tfsdk:"fingerprint"`
	Created      types.String `tfsdk:"created_at"`
	ReadOnly     types.Bool   `tfsdk:"read_only"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *deployKeyResourceModel) from(k *forgejo.DeployKey) {
	m.KeyID = types.Int64Value(k.ID)
	m.Key = types.StringValue(k.Key)
	m.URL = types.StringValue(k.URL)
	m.Title = types.StringValue(k.Title)
	m.Fingerprint = types.StringValue(k.Fingerprint)
	m.Created = types.StringValue(k.Created.String())
	m.ReadOnly = types.BoolValue(k.ReadOnly)
}

// to is a helper function to save Terraform data model into an API struct.
func (m *deployKeyResourceModel) to(o *forgejo.CreateKeyOption) {
	if o == nil {
		o = new(forgejo.CreateKeyOption)
	}

	o.Title = m.Title.ValueString()
	o.Key = m.Key.ValueString()
	o.ReadOnly = m.ReadOnly.ValueBool()
}

// Metadata returns the resource type name.
func (r *deployKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy_key"
}

// Schema defines the schema for the resource.
func (r *deployKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository deploy key resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"key_id": schema.Int64Attribute{
				Description: "Numeric identifier of the deploy key.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "Armored SSH key. Trailing newlines must be removed (e.g. using trimspace() function).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "URL of the deploy key.",
				Computed:    true,
			},
			"title": schema.StringAttribute{
				Description: "Title of the deploy key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fingerprint": schema.StringAttribute{
				Description: "Fingerprint of the deploy key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the deploy key was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"read_only": schema.BoolAttribute{
				Description: "Does the key have only read access?",
				Required:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *deployKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *deployKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create deploy key resource"))

	var (
		repo repositoryResourceModel
		data deployKeyResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Create deploy key", map[string]any{
		"user":      repo.Owner.ValueString(),
		"repo":      repo.Name.ValueString(),
		"title":     data.Title.ValueString(),
		"key":       data.Key.ValueString(),
		"read_only": data.ReadOnly.ValueBool(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateKeyOption{}
	data.to(&opts)

	// Validate API request body
	// err := opts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// Use Forgejo client to create new deploy key
	key, res, err := r.client.CreateDeployKey(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create deploy key", msg)

		return
	}

	// Map response body to model
	data.from(key)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *deployKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read deploy key resource"))

	var (
		repo repositoryResourceModel
		data deployKeyResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Get deploy key by id", map[string]any{
		"user":   repo.Owner.ValueString(),
		"repo":   repo.Name.ValueString(),
		"key_id": data.KeyID.ValueInt64(),
	})

	// Use Forgejo client to get deploy key
	key, res, err := r.client.GetDeployKey(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.KeyID.ValueInt64(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Deploy key with user %s repo %s and id %d not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.KeyID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get deploy key by id", msg)

		return
	}

	// Map response body to model
	data.from(key)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deployKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update deploy key resource"))

	/*
	 * Deploy keys can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deployKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete deploy key resource"))

	var (
		repo repositoryResourceModel
		data deployKeyResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Delete deploy key", map[string]any{
		"owner":  repo.Owner.ValueString(),
		"repo":   repo.Name.ValueString(),
		"key_id": data.KeyID.ValueInt64(),
	})

	// Use Forgejo client to delete existing deploy key
	res, err = r.client.DeleteDeployKey(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.KeyID.ValueInt64(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Deploy key with owner %s repo %s and id %d forbidden: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.KeyID.ValueInt64(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Deploy key with owner %s repo %s and id %d not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.KeyID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete deploy key", msg)

		return
	}
}

// NewDeployKeyResource is a helper function to simplify the provider implementation.
func NewDeployKeyResource() resource.Resource {
	return &deployKeyResource{}
}
