package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repositoryResource{}
	_ resource.ResourceWithConfigure = &repositoryResource{}
)

// repositoryResource is the resource implementation.
type repositoryResource struct {
	client *forgejo.Client
}

// repositoryResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Repository
type repositoryResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Owner                     types.Object `tfsdk:"owner"`
	Name                      types.String `tfsdk:"name"`
	FullName                  types.String `tfsdk:"full_name"`
	Description               types.String `tfsdk:"description"`
	Empty                     types.Bool   `tfsdk:"empty"`
	Private                   types.Bool   `tfsdk:"private"`
	Fork                      types.Bool   `tfsdk:"fork"`
	Template                  types.Bool   `tfsdk:"template"`
	ParentID                  types.Int64  `tfsdk:"parent_id"`
	Mirror                    types.Bool   `tfsdk:"mirror"`
	Size                      types.Int64  `tfsdk:"size"`
	HTMLURL                   types.String `tfsdk:"html_url"`
	SSHURL                    types.String `tfsdk:"ssh_url"`
	CloneURL                  types.String `tfsdk:"clone_url"`
	OriginalURL               types.String `tfsdk:"original_url"`
	Website                   types.String `tfsdk:"website"`
	Stars                     types.Int64  `tfsdk:"stars_count"`
	Forks                     types.Int64  `tfsdk:"forks_count"`
	Watchers                  types.Int64  `tfsdk:"watchers_count"`
	OpenIssues                types.Int64  `tfsdk:"open_issues_count"`
	OpenPulls                 types.Int64  `tfsdk:"open_pr_counter"`
	Releases                  types.Int64  `tfsdk:"release_counter"`
	DefaultBranch             types.String `tfsdk:"default_branch"`
	Archived                  types.Bool   `tfsdk:"archived"`
	Created                   types.String `tfsdk:"created_at"`
	Updated                   types.String `tfsdk:"updated_at"`
	Permissions               types.Object `tfsdk:"permissions"`
	HasIssues                 types.Bool   `tfsdk:"has_issues"`
	InternalTracker           types.Object `tfsdk:"internal_tracker"`
	ExternalTracker           types.Object `tfsdk:"external_tracker"`
	HasWiki                   types.Bool   `tfsdk:"has_wiki"`
	ExternalWiki              types.Object `tfsdk:"external_wiki"`
	HasPullRequests           types.Bool   `tfsdk:"has_pull_requests"`
	HasProjects               types.Bool   `tfsdk:"has_projects"`
	HasReleases               types.Bool   `tfsdk:"has_releases"`
	HasPackages               types.Bool   `tfsdk:"has_packages"`
	HasActions                types.Bool   `tfsdk:"has_actions"`
	IgnoreWhitespaceConflicts types.Bool   `tfsdk:"ignore_whitespace_conflicts"`
	AllowMerge                types.Bool   `tfsdk:"allow_merge_commits"`
	AllowRebase               types.Bool   `tfsdk:"allow_rebase"`
	AllowRebaseMerge          types.Bool   `tfsdk:"allow_rebase_explicit"`
	AllowSquash               types.Bool   `tfsdk:"allow_squash_merge"`
	AvatarURL                 types.String `tfsdk:"avatar_url"`
	Internal                  types.Bool   `tfsdk:"internal"`
	MirrorInterval            types.String `tfsdk:"mirror_interval"`
	MirrorUpdated             types.String `tfsdk:"mirror_updated"`
	DefaultMergeStyle         types.String `tfsdk:"default_merge_style"`
	AutoInit                  types.Bool   `tfsdk:"auto_init"`
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#User
type repositoryResourceUser struct {
	ID        types.Int64  `tfsdk:"id"`
	UserName  types.String `tfsdk:"login"`
	LoginName types.String `tfsdk:"login_name"`
	FullName  types.String `tfsdk:"full_name"`
	Email     types.String `tfsdk:"email"`
}

func (m repositoryResourceUser) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.Int64Type,
		"login":      types.StringType,
		"login_name": types.StringType,
		"full_name":  types.StringType,
		"email":      types.StringType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Permission
type repositoryResourcePermissions struct {
	Admin types.Bool `tfsdk:"admin"`
	Push  types.Bool `tfsdk:"push"`
	Pull  types.Bool `tfsdk:"pull"`
}

func (m repositoryResourcePermissions) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"admin": types.BoolType,
		"push":  types.BoolType,
		"pull":  types.BoolType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#InternalTracker
type repositoryResourceInternalTracker struct {
	EnableTimeTracker                types.Bool `tfsdk:"enable_time_tracker"`
	AllowOnlyContributorsToTrackTime types.Bool `tfsdk:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          types.Bool `tfsdk:"enable_issue_dependencies"`
}

func (m repositoryResourceInternalTracker) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enable_time_tracker":                   types.BoolType,
		"allow_only_contributors_to_track_time": types.BoolType,
		"enable_issue_dependencies":             types.BoolType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalTracker
type repositoryResourceExternalTracker struct {
	ExternalTrackerURL    types.String `tfsdk:"external_tracker_url"`
	ExternalTrackerFormat types.String `tfsdk:"external_tracker_format"`
	ExternalTrackerStyle  types.String `tfsdk:"external_tracker_style"`
}

func (m repositoryResourceExternalTracker) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_tracker_url":    types.StringType,
		"external_tracker_format": types.StringType,
		"external_tracker_style":  types.StringType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalWiki
type repositoryResourceExternalWiki struct {
	ExternalWikiURL types.String `tfsdk:"external_wiki_url"`
}

func (m repositoryResourceExternalWiki) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_wiki_url": types.StringType,
	}
}

// Metadata returns the resource type name.
func (r *repositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

// Schema defines the schema for the resource.
func (r *repositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Computed:    true,
			},
			"owner": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Numeric identifier of the user.",
						Computed:    true,
					},
					"login": schema.StringAttribute{
						Description: "Name of the user.",
						Required:    true,
					},
					"login_name": schema.StringAttribute{
						Description: "Login name of the user.",
						Computed:    true,
					},
					"full_name": schema.StringAttribute{
						Description: "Full name of the user.",
						Computed:    true,
					},
					"email": schema.StringAttribute{
						Description: "Email address of the user.",
						Computed:    true,
					},
				},
				Description: "Owner of the repository.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the repository.",
				Required:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the repository.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the repository.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"empty": schema.BoolAttribute{
				Description: "Is the repository empty?",
				Computed:    true,
			},
			"private": schema.BoolAttribute{
				Description: "Is the repository private?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"fork": schema.BoolAttribute{
				Description: "Is the repository a fork?",
				Computed:    true,
			},
			"template": schema.BoolAttribute{
				Description: "Is the repository a template?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"parent_id": schema.Int64Attribute{
				Description: "Numeric identifier of the parent repository.",
				Computed:    true,
			},
			"mirror": schema.BoolAttribute{
				Description: "Is the repository a mirror?",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Size of the repository in KiB.",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "HTML URL of the repository.",
				Computed:    true,
			},
			"ssh_url": schema.StringAttribute{
				Description: "SSH URL of the repository.",
				Computed:    true,
			},
			"clone_url": schema.StringAttribute{
				Description: "Clone URL of the repository.",
				Computed:    true,
			},
			"original_url": schema.StringAttribute{
				Description: "Original URL of the repository.",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website of the repository.",
				Computed:    true,
			},
			"stars_count": schema.Int64Attribute{
				Description: "Number of stars of the repository.",
				Computed:    true,
			},
			"forks_count": schema.Int64Attribute{
				Description: "Number of forks of the repository.",
				Computed:    true,
			},
			"watchers_count": schema.Int64Attribute{
				Description: "Number of watchers of the repository.",
				Computed:    true,
			},
			"open_issues_count": schema.Int64Attribute{
				Description: "Number of open issues of the repository.",
				Computed:    true,
			},
			"open_pr_counter": schema.Int64Attribute{
				Description: "Number of open pull requests of the repository.",
				Computed:    true,
			},
			"release_counter": schema.Int64Attribute{
				Description: "Number of releases of the repository.",
				Computed:    true,
			},
			"default_branch": schema.StringAttribute{
				Description: "Default branch of the repository.",
				Computed:    true,
			},
			"archived": schema.BoolAttribute{
				Description: "Is the repository archived?",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the repository was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Time at which the repository was updated.",
				Computed:    true,
			},
			"permissions": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"admin": schema.BoolAttribute{
						Description: "Allowed to administer?",
						Computed:    true,
					},
					"push": schema.BoolAttribute{
						Description: "Allowed to push?",
						Computed:    true,
					},
					"pull": schema.BoolAttribute{
						Description: "Allowed to pull?",
						Computed:    true,
					},
				},
				Description: "Permissions of the repository.",
				Computed:    true,
			},
			"has_issues": schema.BoolAttribute{
				Description: "Is the repository issue tracker enabled?",
				Computed:    true,
			},
			"internal_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_time_tracker": schema.BoolAttribute{
						Description: "Enable time tracking.",
						Computed:    true,
					},
					"allow_only_contributors_to_track_time": schema.BoolAttribute{
						Description: "Let only contributors track time.",
						Computed:    true,
					},
					"enable_issue_dependencies": schema.BoolAttribute{
						Description: "Enable dependencies for issues and pull requests.",
						Computed:    true,
					},
				},
				Description: "Settings for built-in issue tracker.",
				Computed:    true,
			},
			"external_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_tracker_url": schema.StringAttribute{
						Description: "URL of external issue tracker.",
						Computed:    true,
					},
					"external_tracker_format": schema.StringAttribute{
						Description: "External issue tracker URL format.",
						Computed:    true,
					},
					"external_tracker_style": schema.StringAttribute{
						Description: "External issue tracker number format.",
						Computed:    true,
					},
				},
				Description: "Settings for external issue tracker.",
				Computed:    true,
			},
			"has_wiki": schema.BoolAttribute{
				Description: "Is the repository wiki enabled?",
				Computed:    true,
			},
			"external_wiki": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_wiki_url": schema.StringAttribute{
						Description: "URL of external wiki.",
						Computed:    true,
					},
				},
				Description: "Settings for external wiki.",
				Computed:    true,
			},
			"has_pull_requests": schema.BoolAttribute{
				Description: "Are repository pull requests enabled?",
				Computed:    true,
			},
			"has_projects": schema.BoolAttribute{
				Description: "Are repository projects enabled?",
				Computed:    true,
			},
			"has_releases": schema.BoolAttribute{
				Description: "Are repository releases enabled?",
				Computed:    true,
			},
			"has_packages": schema.BoolAttribute{
				Description: "Is the repository package registry enabled?",
				Computed:    true,
			},
			"has_actions": schema.BoolAttribute{
				Description: "Are integrated CI/CD pipelines enabled?",
				Computed:    true,
			},
			"ignore_whitespace_conflicts": schema.BoolAttribute{
				Description: "Are whitespace conflicts ignored?",
				Computed:    true,
			},
			"allow_merge_commits": schema.BoolAttribute{
				Description: "Allowed to create merge commit?",
				Computed:    true,
			},
			"allow_rebase": schema.BoolAttribute{
				Description: "Allowed to rebase then fast-forward?",
				Computed:    true,
			},
			"allow_rebase_explicit": schema.BoolAttribute{
				Description: "Allowed to rebase then create merge commit?",
				Computed:    true,
			},
			"allow_squash_merge": schema.BoolAttribute{
				Description: "Allowed to create squash commit?",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the repository.",
				Computed:    true,
			},
			"internal": schema.BoolAttribute{
				Description: "Is the repository internal?",
				Computed:    true,
			},
			"mirror_interval": schema.StringAttribute{
				Description: "Mirror interval of the repository.",
				Computed:    true,
			},
			"mirror_updated": schema.StringAttribute{
				Description: "Time at which the repository mirror was updated.",
				Computed:    true,
			},
			"default_merge_style": schema.StringAttribute{
				Description: "Default merge style of the repository.",
				Computed:    true,
			},
			"auto_init": schema.BoolAttribute{
				Description: "Whether the repository should be auto-intialized?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *repositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository resource"))

	var (
		data  repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create repository", map[string]any{
		"name":        data.Name.ValueString(),
		"description": data.Description.ValueString(),
		"private":     data.Private.ValueBool(),
		//"issue_labels": ,
		"auto_init": data.AutoInit.ValueBool(),
		"template":  data.Template.ValueBool(),
		//"gitignores": ,
		//"license": ,
		//"readme": ,
		//"default_branch" ,
		//"trust_model": forgejo.TrustModel(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateRepoOption{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Private:     data.Private.ValueBool(),
		//IssueLabels: ,
		AutoInit: data.AutoInit.ValueBool(),
		Template: data.Template.ValueBool(),
		//Gitignores: ,
		//License: ,
		//Readme: ,
		//DefaultBranch ,
		//TrustModel: forgejo.TrustModel(),
	}

	// Validate API request body
	err := opts.Validate(r.client)
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Read repository owner into model
	diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		repo     *forgejo.Repository
		response *forgejo.Response
		error    error
	)

	tflog.Info(ctx, "Get organization by name", map[string]any{
		"name": owner.UserName.ValueString(),
	})

	// Use Forgejo client to check if owner is org or user
	_, _, err = r.client.GetOrg(owner.UserName.ValueString())
	if err == nil {
		// Owner is org
		// -> use Forgejo client to create new org repository
		repo, response, error = r.client.CreateOrgRepo(
			owner.UserName.ValueString(),
			opts,
		)
	} else {
		// Assume owner is user
		// -> use Forgejo client to create new user repository
		repo, response, error = r.client.CreateRepo(opts)
	}
	if error != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": response.Status,
		})

		var msg string
		switch response.StatusCode {
		case 403:
			msg = fmt.Sprintf("Repository with name %s forbidden: %s", data.Name.String(), error)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", error)
		default:
			msg = fmt.Sprintf("Unknown error: %s", error)
		}
		resp.Diagnostics.AddError("Unable to create repository", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(repo.ID)
	data.FullName = types.StringValue(repo.FullName)
	data.Description = types.StringValue(repo.Description)
	data.Empty = types.BoolValue(repo.Empty)
	data.Private = types.BoolValue(repo.Private)
	data.Fork = types.BoolValue(repo.Fork)
	data.Template = types.BoolValue(repo.Template)
	if repo.Parent != nil {
		data.ParentID = types.Int64Value(repo.Parent.ID)
	} else {
		data.ParentID = types.Int64Null()
	}
	data.Mirror = types.BoolValue(repo.Mirror)
	data.Size = types.Int64Value(int64(repo.Size))
	data.HTMLURL = types.StringValue(repo.HTMLURL)
	data.SSHURL = types.StringValue(repo.SSHURL)
	data.CloneURL = types.StringValue(repo.CloneURL)
	data.OriginalURL = types.StringValue(repo.OriginalURL)
	data.Website = types.StringValue(repo.Website)
	data.Stars = types.Int64Value(int64(repo.Stars))
	data.Forks = types.Int64Value(int64(repo.Forks))
	data.Watchers = types.Int64Value(int64(repo.Watchers))
	data.OpenIssues = types.Int64Value(int64(repo.OpenIssues))
	data.OpenPulls = types.Int64Value(int64(repo.OpenPulls))
	data.Releases = types.Int64Value(int64(repo.Releases))
	data.DefaultBranch = types.StringValue(repo.DefaultBranch)
	data.Archived = types.BoolValue(repo.Archived)
	data.Created = types.StringValue(repo.Created.String())
	data.Updated = types.StringValue(repo.Updated.String())
	data.HasIssues = types.BoolValue(repo.HasIssues)
	data.HasWiki = types.BoolValue(repo.HasWiki)
	data.HasPullRequests = types.BoolValue(repo.HasPullRequests)
	data.HasProjects = types.BoolValue(repo.HasProjects)
	data.HasReleases = types.BoolValue(repo.HasReleases)
	data.HasPackages = types.BoolValue(repo.HasPackages)
	data.HasActions = types.BoolValue(repo.HasActions)
	data.IgnoreWhitespaceConflicts = types.BoolValue(repo.IgnoreWhitespaceConflicts)
	data.AllowMerge = types.BoolValue(repo.AllowMerge)
	data.AllowRebase = types.BoolValue(repo.AllowRebase)
	data.AllowRebaseMerge = types.BoolValue(repo.AllowRebaseMerge)
	data.AllowSquash = types.BoolValue(repo.AllowSquash)
	data.AvatarURL = types.StringValue(repo.AvatarURL)
	data.Internal = types.BoolValue(repo.Internal)
	data.MirrorInterval = types.StringValue(repo.MirrorInterval)
	data.MirrorUpdated = types.StringValue(repo.MirrorUpdated.String())
	data.DefaultMergeStyle = types.StringValue(string(repo.DefaultMergeStyle))

	// Repository owner
	if repo.Owner != nil {
		ownerElement := repositoryResourceUser{
			ID:        types.Int64Value(repo.Owner.ID),
			UserName:  types.StringValue(repo.Owner.UserName),
			LoginName: types.StringValue(repo.Owner.LoginName),
			FullName:  types.StringValue(repo.Owner.FullName),
			Email:     types.StringValue(repo.Owner.Email),
		}
		ownerValue, diags := types.ObjectValueFrom(
			ctx,
			ownerElement.AttributeTypes(),
			ownerElement,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Owner = ownerValue
	} else {
		data.Owner = types.ObjectNull(
			repositoryResourceUser{}.AttributeTypes(),
		)
	}

	// Repository permissions
	if repo.Permissions != nil {
		perms := repositoryResourcePermissions{
			Admin: types.BoolValue(repo.Permissions.Admin),
			Push:  types.BoolValue(repo.Permissions.Push),
			Pull:  types.BoolValue(repo.Permissions.Pull),
		}
		permsValue, diags := types.ObjectValueFrom(
			ctx,
			perms.AttributeTypes(),
			perms,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Permissions = permsValue
	} else {
		data.Permissions = types.ObjectNull(
			repositoryResourcePermissions{}.AttributeTypes(),
		)
	}

	// Internal issue tracker
	if repo.InternalTracker != nil {
		intTracker := repositoryResourceInternalTracker{
			EnableTimeTracker:                types.BoolValue(repo.InternalTracker.EnableTimeTracker),
			AllowOnlyContributorsToTrackTime: types.BoolValue(repo.InternalTracker.AllowOnlyContributorsToTrackTime),
			EnableIssueDependencies:          types.BoolValue(repo.InternalTracker.EnableIssueDependencies),
		}
		intTrackerValue, diags := types.ObjectValueFrom(
			ctx,
			intTracker.AttributeTypes(),
			intTracker,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.InternalTracker = intTrackerValue
	} else {
		data.InternalTracker = types.ObjectNull(
			repositoryResourceInternalTracker{}.AttributeTypes(),
		)
	}

	// External issue tracker
	if repo.ExternalTracker != nil {
		extTracker := repositoryResourceExternalTracker{
			ExternalTrackerURL:    types.StringValue(repo.ExternalTracker.ExternalTrackerURL),
			ExternalTrackerFormat: types.StringValue(repo.ExternalTracker.ExternalTrackerFormat),
			ExternalTrackerStyle:  types.StringValue(repo.ExternalTracker.ExternalTrackerStyle),
		}
		extTrackerValue, diags := types.ObjectValueFrom(
			ctx,
			extTracker.AttributeTypes(),
			extTracker,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ExternalTracker = extTrackerValue
	} else {
		data.ExternalTracker = types.ObjectNull(
			repositoryResourceExternalTracker{}.AttributeTypes(),
		)
	}

	// External wiki
	if repo.ExternalWiki != nil {
		wiki := repositoryResourceExternalWiki{
			ExternalWikiURL: types.StringValue(repo.ExternalWiki.ExternalWikiURL),
		}
		wikiValue, diags := types.ObjectValueFrom(
			ctx,
			wiki.AttributeTypes(),
			wiki,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ExternalWiki = wikiValue
	} else {
		data.ExternalWiki = types.ObjectNull(
			repositoryResourceExternalWiki{}.AttributeTypes(),
		)
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository resource"))

	// var data repositoryResourceModel

	// // Read Terraform configuration data into model
	// diags := req.State.Get(ctx, &data)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// tflog.Info(ctx, "Get repository by name", map[string]any{
	// 	"name": data.Name.ValueString(),
	// })

	// // Use Forgejo client to get repository by name
	// o, re, err := r.client.GetOrg(data.Name.ValueString())
	// if err != nil {
	// 	tflog.Error(ctx, "Error", map[string]any{
	// 		"status": re.Status,
	// 	})

	// 	var msg string
	// 	switch re.StatusCode {
	// 	case 404:
	// 		msg = fmt.Sprintf("Repository with name %s not found: %s", data.Name.String(), err)
	// 	default:
	// 		msg = fmt.Sprintf("Unknown error: %s", err)
	// 	}
	// 	resp.Diagnostics.AddError("Unable to get repository by name", msg)

	// 	return
	// }

	// // Map response body to model
	// data.ID = types.Int64Value(o.ID)
	// data.Name = types.StringValue(o.UserName)
	// data.FullName = types.StringValue(o.FullName)
	// data.AvatarURL = types.StringValue(o.AvatarURL)
	// data.Description = types.StringValue(o.Description)
	// data.Website = types.StringValue(o.Website)
	// data.Location = types.StringValue(o.Location)
	// data.Visibility = types.StringValue(o.Visibility)

	// // Save data into Terraform state
	// diags = resp.State.Set(ctx, &data)
	// resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository resource"))

	// var data repositoryResourceModel

	// // Read Terraform plan data into model
	// diags := req.Plan.Get(ctx, &data)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// tflog.Info(ctx, "Update repository", map[string]any{
	// 	"name":        data.Name.ValueString(),
	// 	"full_name":   data.FullName.ValueString(),
	// 	"description": data.Description.ValueString(),
	// 	"website":     data.Website.ValueString(),
	// 	"location":    data.Location.ValueString(),
	// 	"visibility":  data.Visibility.ValueString(),
	// })

	// // Generate API request body from plan
	// opts := forgejo.EditOrgOption{
	// 	FullName:    data.FullName.ValueString(),
	// 	Description: data.Description.ValueString(),
	// 	Website:     data.Website.ValueString(),
	// 	Location:    data.Location.ValueString(),
	// 	Visibility:  forgejo.VisibleType(data.Visibility.ValueString()),
	// }

	// // Validate API request body
	// err := opts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// // Use Forgejo client to update existing repository
	// re, err := r.client.EditOrg(data.Name.ValueString(), opts)
	// if err != nil {
	// 	tflog.Error(ctx, "Error", map[string]any{
	// 		"status": re.Status,
	// 	})

	// 	var msg string
	// 	switch re.StatusCode {
	// 	case 404:
	// 		msg = fmt.Sprintf("Repository with name %s not found: %s", data.Name.String(), err)
	// 	default:
	// 		msg = fmt.Sprintf("Unknown error: %s", err)
	// 	}
	// 	resp.Diagnostics.AddError("Unable to update repository", msg)

	// 	return
	// }

	// // Use Forgejo client to fetch updated repository
	// o, re, err := r.client.GetOrg(data.Name.ValueString())
	// if err != nil {
	// 	tflog.Error(ctx, "Error", map[string]any{
	// 		"status": re.Status,
	// 	})

	// 	var msg string
	// 	switch re.StatusCode {
	// 	case 404:
	// 		msg = fmt.Sprintf("Repository with name %s not found: %s", data.Name.String(), err)
	// 	default:
	// 		msg = fmt.Sprintf("Unknown error: %s", err)
	// 	}
	// 	resp.Diagnostics.AddError("Unable to get repository by name", msg)

	// 	return
	// }

	// // Map response body to model
	// data.ID = types.Int64Value(o.ID)
	// data.Name = types.StringValue(o.UserName)
	// data.FullName = types.StringValue(o.FullName)
	// data.AvatarURL = types.StringValue(o.AvatarURL)
	// data.Description = types.StringValue(o.Description)
	// data.Website = types.StringValue(o.Website)
	// data.Location = types.StringValue(o.Location)
	// data.Visibility = types.StringValue(o.Visibility)

	// // Save data into Terraform state
	// diags = resp.State.Set(ctx, &data)
	// resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository resource"))

	var (
		data  repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read repository owner into model
	diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete repository", map[string]any{
		"owner": owner.UserName.ValueString(),
		"name":  data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing repository
	res, err := r.client.DeleteRepo(
		owner.UserName.ValueString(),
		data.Name.ValueString(),
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
				owner.UserName.String(),
				data.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete repository", msg)

		return
	}
}

// NewRepositoryResource is a helper function to simplify the provider implementation.
func NewRepositoryResource() resource.Resource {
	return &repositoryResource{}
}
