package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &userResource{}
	_ resource.ResourceWithConfigure = &userResource{}
)

// userResource is the resource implementation.
type userResource struct {
	client *forgejo.Client
}

// userResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#User
type userResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"login"`
	LoginName               types.String `tfsdk:"login_name"`
	SourceID                types.Int64  `tfsdk:"source_id"`
	FullName                types.String `tfsdk:"full_name"`
	Email                   types.String `tfsdk:"email"`
	HTMLURL                 types.String `tfsdk:"html_url"`
	AvatarURL               types.String `tfsdk:"avatar_url"`
	Language                types.String `tfsdk:"language"`
	Admin                   types.Bool   `tfsdk:"admin"`
	LastLogin               types.String `tfsdk:"last_login"`
	Created                 types.String `tfsdk:"created_at"`
	Restricted              types.Bool   `tfsdk:"restricted"`
	Active                  types.Bool   `tfsdk:"active"`
	ProhibitLogin           types.Bool   `tfsdk:"prohibit_login"`
	Location                types.String `tfsdk:"location"`
	Website                 types.String `tfsdk:"website"`
	Description             types.String `tfsdk:"description"`
	Visibility              types.String `tfsdk:"visibility"`
	FollowerCount           types.Int64  `tfsdk:"followers_count"`
	FollowingCount          types.Int64  `tfsdk:"following_count"`
	StarredRepoCount        types.Int64  `tfsdk:"starred_repos_count"`
	Password                types.String `tfsdk:"password"`
	MustChangePassword      types.Bool   `tfsdk:"must_change_password"`
	SendNotify              types.Bool   `tfsdk:"send_notify"`
	AllowGitHook            types.Bool   `tfsdk:"allow_git_hook"`
	AllowImportLocal        types.Bool   `tfsdk:"allow_import_local"`
	AllowCreateOrganization types.Bool   `tfsdk:"allow_create_organization"`
	MaxRepoCreation         types.Int64  `tfsdk:"max_repo_creation"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *userResourceModel) from(u *forgejo.User) {
	if u == nil {
		return
	}

	m.ID = types.Int64Value(u.ID)
	m.Name = types.StringValue(u.UserName)
	m.LoginName = types.StringValue(u.LoginName)
	m.SourceID = types.Int64Value(u.SourceID)
	m.FullName = types.StringValue(u.FullName)
	m.Email = types.StringValue(u.Email)
	m.HTMLURL = types.StringValue(u.HTMLURL)
	m.AvatarURL = types.StringValue(u.AvatarURL)
	m.Language = types.StringValue(u.Language)
	m.Admin = types.BoolValue(u.IsAdmin)
	m.LastLogin = types.StringValue(u.LastLogin.Format(time.RFC3339))
	m.Created = types.StringValue(u.Created.Format(time.RFC3339))
	m.Restricted = types.BoolValue(u.Restricted)
	m.Active = types.BoolValue(u.IsActive)
	m.ProhibitLogin = types.BoolValue(u.ProhibitLogin)
	m.Location = types.StringValue(u.Location)
	m.Website = types.StringValue(u.Website)
	m.Description = types.StringValue(u.Description)
	m.Visibility = types.StringValue(string(u.Visibility))
	m.FollowerCount = types.Int64Value(int64(u.FollowerCount))
	m.FollowingCount = types.Int64Value(int64(u.FollowingCount))
	m.StarredRepoCount = types.Int64Value(int64(u.StarredRepoCount))
}

// to is a helper function to save Terraform data model into an API struct.
func (m *userResourceModel) to(s *userResourceModel, o *forgejo.EditUserOption) {
	if o == nil {
		o = new(forgejo.EditUserOption)
	}

	o.SourceID = m.SourceID.ValueInt64()
	o.LoginName = m.LoginName.ValueString()
	o.Email = m.Email.ValueStringPointer()
	o.FullName = m.FullName.ValueStringPointer()

	if s != nil && !m.Password.Equal(s.Password) {
		// only update password if it has changed
		o.Password = m.Password.ValueString()
	}

	o.Description = m.Description.ValueStringPointer()
	o.MustChangePassword = m.MustChangePassword.ValueBoolPointer()
	o.Website = m.Website.ValueStringPointer()
	o.Location = m.Location.ValueStringPointer()
	o.Active = m.Active.ValueBoolPointer()
	o.Admin = m.Admin.ValueBoolPointer()
	o.AllowGitHook = m.AllowGitHook.ValueBoolPointer()
	o.AllowImportLocal = m.AllowImportLocal.ValueBoolPointer()
	o.ProhibitLogin = m.ProhibitLogin.ValueBoolPointer()
	o.AllowCreateOrganization = m.AllowCreateOrganization.ValueBoolPointer()
	o.Restricted = m.Restricted.ValueBoolPointer()

	vt := forgejo.VisibleType(m.Visibility.ValueString())
	o.Visibility = &vt

	mrc := int(m.MaxRepoCreation.ValueInt64())
	o.MaxRepoCreation = &mrc
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo user resource.

**Note**: Managing users requires administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the user.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"login": schema.StringAttribute{
				Description: "Name of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"login_name": schema.StringAttribute{
				Description: "Login name of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"source_id": schema.Int64Attribute{
				Description: "Numeric identifier of the user's authentication source.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				Description: "Email address of the user.",
				Required:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the user.",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "URL to the user's profile page.",
				Computed:    true,
			},
			"language": schema.StringAttribute{
				Description: "Locale of the user.",
				Computed:    true,
			},
			"admin": schema.BoolAttribute{
				Description: "Is the user an administrator?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"last_login": schema.StringAttribute{
				Description: "Time at which the user last logged in.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the user was created.",
				Computed:    true,
				// 6b66d9e: standardize on formatting temporal data in RFC3339 format
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"restricted": schema.BoolAttribute{
				Description: "Is the user restricted?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"active": schema.BoolAttribute{
				Description: "Is the user active?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"prohibit_login": schema.BoolAttribute{
				Description: "Are user logins prohibited?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"location": schema.StringAttribute{
				Description: "Location of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"website": schema.StringAttribute{
				Description: "Website of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "Description of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"visibility": schema.StringAttribute{
				Description: "Visibility of the user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"public",
						"limited",
						"private",
					),
				},
			},
			"followers_count": schema.Int64Attribute{
				Description: "Number of following users.",
				Computed:    true,
			},
			"following_count": schema.Int64Attribute{
				Description: "Number of users followed.",
				Computed:    true,
			},
			"starred_repos_count": schema.Int64Attribute{
				Description: "Number of starred repositories.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password of the user.",
				Required:    true,
				Sensitive:   true,
			},
			"must_change_password": schema.BoolAttribute{
				Description: "Require user to change password?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"send_notify": schema.BoolAttribute{
				Description: "Send notification to administrators?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_git_hook": schema.BoolAttribute{
				Description: "Allow user to create Git hooks?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_import_local": schema.BoolAttribute{
				Description: "Allow user to import local repositories?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_create_organization": schema.BoolAttribute{
				Description: "Allow user to create organizations?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"max_repo_creation": schema.Int64Attribute{
				Description: "Maximum number of repositories user can create. A value of -1 means no limit.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(-1),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create user resource"))

	var data userResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create user", map[string]any{
		"source_id":            data.SourceID.ValueInt64(),
		"login_name":           data.LoginName.ValueString(),
		"login":                data.Name.ValueString(),
		"full_name":            data.FullName.ValueString(),
		"email":                data.Email.ValueString(),
		"password":             strings.Repeat("*", len(data.Password.ValueString())),
		"must_change_password": data.MustChangePassword.ValueBool(),
		"send_notify":          data.SendNotify.ValueBool(),
		"visibility":           data.Visibility.ValueString(),
	})

	// Generate API request body from plan
	vt := forgejo.VisibleType(data.Visibility.ValueString())
	copts := forgejo.CreateUserOption{
		SourceID:           data.SourceID.ValueInt64(),
		LoginName:          data.LoginName.ValueString(),
		Username:           data.Name.ValueString(),
		FullName:           data.FullName.ValueString(),
		Email:              data.Email.ValueString(),
		Password:           data.Password.ValueString(),
		MustChangePassword: data.MustChangePassword.ValueBoolPointer(),
		SendNotify:         data.SendNotify.ValueBool(),
		Visibility:         &vt,
	}

	// Validate API request body
	err := copts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new user
	usr, res, err := r.client.AdminCreateUser(copts)
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
			case 403:
				msg = fmt.Sprintf(
					"User with name %s forbidden: %s",
					data.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create user", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(usr.ID)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update user", map[string]any{
		"source_id":                 data.SourceID.ValueInt64(),
		"login_name":                data.LoginName.ValueString(),
		"email":                     data.Email.ValueString(),
		"full_name":                 data.FullName.ValueString(),
		"password":                  strings.Repeat("*", len(data.Password.ValueString())),
		"description":               data.Description.ValueString(),
		"must_change_password":      data.MustChangePassword.ValueBool(),
		"website":                   data.Website.ValueString(),
		"location":                  data.Location.ValueString(),
		"active":                    data.Active.ValueBool(),
		"admin":                     data.Admin.ValueBool(),
		"allow_git_hook":            data.AllowGitHook.ValueBool(),
		"allow_import_local":        data.AllowImportLocal.ValueBool(),
		"max_repo_creation":         data.MaxRepoCreation.ValueInt64(),
		"prohibit_login":            data.ProhibitLogin.ValueBool(),
		"allow_create_organization": data.AllowCreateOrganization.ValueBool(),
		"restricted":                data.Restricted.ValueBool(),
		"visibility":                data.Visibility.ValueString(),
	})

	// Generate API request body from plan
	eopts := forgejo.EditUserOption{}
	data.to(nil, &eopts)

	// Validate API request body
	// err := eopts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// Use Forgejo client to update existing user
	res, err = r.client.AdminEditUser(
		data.Name.ValueString(),
		eopts,
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
			case 403:
				msg = fmt.Sprintf(
					"User with name %s forbidden: %s",
					data.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"User with name %s not found: %s",
					data.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update user", msg)

		return
	}

	tflog.Info(ctx, "Read user", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to fetch updated user
	usr, res, err = r.client.GetUserInfo(data.Name.ValueString())
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
					"User with name %s not found: %s",
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read user", msg)

		return
	}

	// Map response body to model
	data.from(usr)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read user resource"))

	var data userResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read user", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to get user by name
	usr, res, err := r.client.GetUserInfo(data.Name.ValueString())
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
					"User with name %s not found: %s",
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read user", msg)

		return
	}

	// Map response body to model
	data.from(usr)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update user resource"))

	var (
		plan  userResourceModel
		state userResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into model
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update user", map[string]any{
		"source_id":                 plan.SourceID.ValueInt64(),
		"login_name":                plan.LoginName.ValueString(),
		"email":                     plan.Email.ValueString(),
		"full_name":                 plan.FullName.ValueString(),
		"password":                  strings.Repeat("*", len(plan.Password.ValueString())),
		"description":               plan.Description.ValueString(),
		"must_change_password":      plan.MustChangePassword.ValueBool(),
		"website":                   plan.Website.ValueString(),
		"location":                  plan.Location.ValueString(),
		"active":                    plan.Active.ValueBool(),
		"admin":                     plan.Admin.ValueBool(),
		"allow_git_hook":            plan.AllowGitHook.ValueBool(),
		"allow_import_local":        plan.AllowImportLocal.ValueBool(),
		"max_repo_creation":         plan.MaxRepoCreation.ValueInt64(),
		"prohibit_login":            plan.ProhibitLogin.ValueBool(),
		"allow_create_organization": plan.AllowCreateOrganization.ValueBool(),
		"restricted":                plan.Restricted.ValueBool(),
		"visibility":                plan.Visibility.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.EditUserOption{}
	plan.to(&state, &opts)

	// Validate API request body
	// err := opts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// Use Forgejo client to update existing user
	res, err := r.client.AdminEditUser(
		plan.Name.ValueString(),
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
			case 403:
				msg = fmt.Sprintf(
					"User with name %s forbidden: %s",
					plan.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"User with name %s not found: %s",
					plan.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update user", msg)

		return
	}

	tflog.Info(ctx, "Read user", map[string]any{
		"name": plan.Name.ValueString(),
	})

	// Use Forgejo client to fetch updated user
	usr, res, err := r.client.GetUserInfo(plan.Name.ValueString())
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
					"User with name %s not found: %s",
					plan.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read user", msg)

		return
	}

	// Map response body to model
	plan.from(usr)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete user resource"))

	var data userResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete user", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing user
	res, err := r.client.AdminDeleteUser(data.Name.ValueString())
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
					"User with name %s forbidden: %s",
					data.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"User with name %s not found: %s",
					data.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete user", msg)

		return
	}
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import user resource"))

	var state userResourceModel

	tflog.Info(ctx, "Read user", map[string]any{
		"username": req.ID,
	})

	// Use Forgejo client to get user by name
	usr, res, err := r.client.GetUserInfo(req.ID)
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
					"User with name '%s' not found: %s",
					req.ID,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read user", msg)

		return
	}

	// Map response body to model
	state.from(usr)

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}
