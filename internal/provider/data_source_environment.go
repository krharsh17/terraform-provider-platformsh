package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rubaiat-hossain/terraform-provider-platformsh/internal/platformsh"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	client *platformsh.Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	Projects []ProjectModel `tfsdk:"projects"`
}

type ProjectModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of projects available in Platform.sh",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "List of projects",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the project",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "Title of the project",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the project",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*platformsh.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *platformsh.Client",
		)
		return
	}

	d.client = client
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	// Fetch the projects
	projects, err := d.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Unable to read projects, got error: "+err.Error(),
		)
		return
	}

	// Map the projects to the Terraform data model
	for _, project := range projects {
		data.Projects = append(data.Projects, ProjectModel{
			ID:          types.StringValue(project.ID),
			Title:       types.StringValue(project.Title),
			Description: types.StringValue(project.Description),
		})
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
