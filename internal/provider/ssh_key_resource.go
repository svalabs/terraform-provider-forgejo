package provider

import (
	"context"
	"fmt"
	"time"

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
	_ resource.Resource              = &sshKeyResource{}
	_ resource.ResourceWithConfigure = &sshKeyResource{}
)

// sshKeyResource is the resource implementation.
type sshKeyResource struct {
	client *forgejo.Client
}

// sshKeyResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#PublicKey
type sshKeyResourceModel struct {
	User        types.String `tfsdk:"user"`
	KeyID       types.Int64  `tfsdk:"key_id"`
	Key         types.String `tfsdk:"key"`
	URL         types.String `tfsdk:"url"`
	Title       types.String `tfsdk:"title"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created_at"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	KeyType     types.String `tfsdk:"key_type"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *sshKeyResourceModel) from(k *forgejo.PublicKey) {
	if k == nil {
		return
	}

	m.KeyID = types.Int64Value(k.ID)
	m.Key = types.StringValue(k.Key)
	m.URL = types.StringValue(k.URL)
	m.Title = types.StringValue(k.Title)
	m.Fingerprint = types.StringValue(k.Fingerprint)
	m.Created = types.StringValue(k.Created.Format(time.RFC3339))
	m.ReadOnly = types.BoolValue(k.ReadOnly)
	m.KeyType = types.StringValue(k.KeyType)
}

// to is a helper function to save Terraform data model into an API struct.
func (m *sshKeyResourceModel) to(o *forgejo.CreateKeyOption) {
	if o == nil {
		o = new(forgejo.CreateKeyOption)
	}

	o.Title = m.Title.ValueString()
	o.Key = m.Key.ValueString()
}

// Metadata returns the resource type name.
func (r *sshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

// Schema defines the schema for the resource.
func (r *sshKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Forgejo user SSH key resource.
Note: Managing user SSH keys requires administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description: "Name of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_id": schema.Int64Attribute{
				Description: "Numeric identifier of the SSH key.",
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
				Description: "URL of the SSH key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Title of the SSH key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fingerprint": schema.StringAttribute{
				Description: "Fingerprint of the SSH key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the SSH key was created.",
				Computed:    true,
				// 6b66d9e: standardize on formatting temporal data in RFC3339 format
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"read_only": schema.BoolAttribute{
				Description: "Does the key have only read access?",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "Type of the SSH key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *sshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *sshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create SSH key resource"))

	var data sshKeyResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create SSH key", map[string]any{
		"user":  data.User.ValueString(),
		"title": data.Title.ValueString(),
		"key":   data.Key.ValueString(),
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

	// Use Forgejo client to create new SSH key
	key, res, err := r.client.AdminCreateUserPublicKey(
		data.User.ValueString(),
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
			case 403:
				msg = fmt.Sprintf(
					"SSH key with user %s forbidden: %s",
					data.User.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"SSH key with user %s not found: %s",
					data.User.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create SSH key", msg)

		return
	}

	// Map response body to model
	data.from(key)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read SSH key resource"))

	var data sshKeyResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read SSH key", map[string]any{
		"user":   data.User.ValueString(),
		"key_id": data.KeyID.ValueInt64(),
	})

	// Use Forgejo client to get SSH key
	key, res, err := r.client.GetPublicKey(data.KeyID.ValueInt64())
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf(
					"SSH key with user %s and id %d forbidden: %s",
					data.User.String(),
					data.KeyID.ValueInt64(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"SSH key with user %s and id %d not found: %s",
					data.User.String(),
					data.KeyID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read SSH key", msg)

		return
	}

	// Map response body to model
	data.from(key)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update SSH key resource"))

	/*
	 * SSH keys can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete SSH key resource"))

	var data sshKeyResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete SSH key", map[string]any{
		"user":   data.User.ValueString(),
		"key_id": data.KeyID.ValueInt64(),
	})

	// Use Forgejo client to delete existing SSH key
	res, err := r.client.AdminDeleteUserPublicKey(
		data.User.ValueString(),
		int(data.KeyID.ValueInt64()),
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
			case 403:
				msg = fmt.Sprintf(
					"SSH key with user %s and id %d forbidden: %s",
					data.User.String(),
					data.KeyID.ValueInt64(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"SSH key with user %s and id %d not found: %s",
					data.User.String(),
					data.KeyID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete SSH key", msg)

		return
	}
}

// NewSSHKeyResource is a helper function to simplify the provider implementation.
func NewSSHKeyResource() resource.Resource {
	return &sshKeyResource{}
}
