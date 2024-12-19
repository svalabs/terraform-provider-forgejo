package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &organizationResource{}
	_ resource.ResourceWithConfigure = &organizationResource{}
)

// organizationResource is the resource implementation.
type organizationResource struct {
	client *forgejo.Client
}

// organizationResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Organization
type organizationResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	FullName    types.String `tfsdk:"full_name"`
	AvatarURL   types.String `tfsdk:"avatar_url"`
	Description types.String `tfsdk:"description"`
	Website     types.String `tfsdk:"website"`
	Location    types.String `tfsdk:"location"`
	Visibility  types.String `tfsdk:"visibility"`
}

func (m *organizationResourceModel) from(o *forgejo.Organization) {
	m.ID = types.Int64Value(o.ID)
	m.Name = types.StringValue(o.UserName)
	m.FullName = types.StringValue(o.FullName)
	m.AvatarURL = types.StringValue(o.AvatarURL)
	m.Description = types.StringValue(o.Description)
	m.Website = types.StringValue(o.Website)
	m.Location = types.StringValue(o.Location)
	m.Visibility = types.StringValue(o.Visibility)
}

// Metadata returns the resource type name.
func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo organization resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the organization.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the organization.",
				Required:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the organization.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"website": schema.StringAttribute{
				Description: "Website of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"location": schema.StringAttribute{
				Description: "Location of the organization.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"visibility": schema.StringAttribute{
				Description: "Visibility of the organization. Possible values are 'public' (default), 'limited', or 'private'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("public"),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create organization resource"))

	var data organizationResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create organization", map[string]any{
		"name":        data.Name.ValueString(),
		"full_name":   data.FullName.ValueString(),
		"description": data.Description.ValueString(),
		"website":     data.Website.ValueString(),
		"location":    data.Location.ValueString(),
		"visibility":  data.Visibility.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateOrgOption{
		Name:        data.Name.ValueString(),
		FullName:    data.FullName.ValueString(),
		Description: data.Description.ValueString(),
		Website:     data.Website.ValueString(),
		Location:    data.Location.ValueString(),
		Visibility:  forgejo.VisibleType(data.Visibility.ValueString()),
	}

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new organization
	org, res, err := r.client.CreateOrg(opts)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf("Organization with name %s forbidden: %s", data.Name.String(), err)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create organization", msg)

		return
	}

	// Map response body to model
	data.from(org)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read organization resource"))

	var data organizationResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get organization by name", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to get organization by name
	org, res, err := r.client.GetOrg(data.Name.ValueString())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf("Organization with name %s not found: %s", data.Name.String(), err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get organization by name", msg)

		return
	}

	// Map response body to model
	data.from(org)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update organization resource"))

	var data organizationResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update organization", map[string]any{
		"name":        data.Name.ValueString(),
		"full_name":   data.FullName.ValueString(),
		"description": data.Description.ValueString(),
		"website":     data.Website.ValueString(),
		"location":    data.Location.ValueString(),
		"visibility":  data.Visibility.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.EditOrgOption{
		FullName:    data.FullName.ValueString(),
		Description: data.Description.ValueString(),
		Website:     data.Website.ValueString(),
		Location:    data.Location.ValueString(),
		Visibility:  forgejo.VisibleType(data.Visibility.ValueString()),
	}

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to update existing organization
	res, err := r.client.EditOrg(
		data.Name.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf("Organization with name %s not found: %s", data.Name.String(), err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update organization", msg)

		return
	}

	// Use Forgejo client to fetch updated organization
	org, res, err := r.client.GetOrg(data.Name.ValueString())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf("Organization with name %s not found: %s", data.Name.String(), err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get organization by name", msg)

		return
	}

	// Map response body to model
	data.from(org)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete organization resource"))

	var data organizationResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete organization", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing organization
	res, err := r.client.DeleteOrg(data.Name.ValueString())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf("Organization with name %s not found: %s", data.Name.String(), err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete organization", msg)

		return
	}
}

// NewOrganizationResource is a helper function to simplify the provider implementation.
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}
