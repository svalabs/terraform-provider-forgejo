package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// Metadata returns the resource type name.
func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "",
				Computed:    true,
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
			fmt.Sprintf("Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Create organization resource - begin")

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

	// Use Forgejo client to create new organization
	o, re, err := r.client.CreateOrg(opts)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{"status": re.Status})

		var msg string
		switch re.StatusCode {
		case 403:
			msg = fmt.Sprintf("Permission error '%s'.", err)
		case 422:
			msg = fmt.Sprintf("Input validation error '%s'.", err)
		default:
			msg = fmt.Sprintf("Unknown error '%s'.", err)
		}
		resp.Diagnostics.AddError("Unable to create organization", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(o.ID)
	data.Name = types.StringValue(o.UserName)
	data.FullName = types.StringValue(o.FullName)
	data.AvatarURL = types.StringValue(o.AvatarURL)
	data.Description = types.StringValue(o.Description)
	data.Website = types.StringValue(o.Website)
	data.Location = types.StringValue(o.Location)
	data.Visibility = types.StringValue(o.Visibility)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "Create organization resource - end", map[string]any{"success": true})
}

// Read refreshes the Terraform state with the latest data.
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// NewOrganizationResource is a helper function to simplify the provider implementation.
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}
