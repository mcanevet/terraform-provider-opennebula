package opennebula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &frameworkProvider{}
)

// frameworkProvider is the provider implementation for the Plugin Framework.
// This will gradually replace the SDKv2 provider as resources are migrated.
type frameworkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance testing.
	version string
}

// frameworkProviderModel maps provider schema data to a Go type.
type frameworkProviderModel struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	FlowEndpoint types.String `tfsdk:"flow_endpoint"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

// New is a helper function to simplify provider server implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &frameworkProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *frameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opennebula"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *frameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "OpenNebula Provider (Plugin Framework)",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The URL to your public or private OpenNebula",
				Required:    true,
			},
			"flow_endpoint": schema.StringAttribute{
				Description: "The URL to your public or private OpenNebula Flow server",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The ID of the user to identify as",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user",
				Required:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Disable TLS validation",
				Optional:    true,
			},
		},
	}
}

// Configure prepares an OpenNebula API client for data sources and resources.
func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config frameworkProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Initialize OpenNebula client configuration
	// For now, we'll leave this empty until we migrate actual resources
	// The SDKv2 provider will continue handling all resources

	// Example configuration would look like:
	// client, err := goca.NewClient(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Unable to Create OpenNebula API Client", err.Error())
	//     return
	// }
	// resp.DataSourceData = client
	// resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *frameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Data sources will be added here as they are migrated from SDKv2
		// Example: NewTemplateDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *frameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Resources will be added here as they are migrated from SDKv2
		// Example: NewTemplateResource,
	}
}
