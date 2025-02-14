package provider

import (
	"context"
	"os"
	"terraform-provider-kuma/tools"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &kumaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kumaProvider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type kumaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type kumaProviderModel struct {
	HOST     types.String `tfsdk:"host"`
	USERNAME types.String `tfsdk:"username"`
	PASSWORD types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *kumaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kuma"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *kumaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *kumaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config kumaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.HOST.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown HOST",
			"",
		)
	}

	if config.USERNAME.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown username",
			"",
		)
	}

	if config.PASSWORD.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown password",
			"",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("KUMA_HOST")
	username := os.Getenv("KUMA_USERNAME")
	password := os.Getenv("KUMA_PASSWORD")

	if !config.HOST.IsNull() {
		host = config.HOST.ValueString()
	}

	if !config.USERNAME.IsNull() {
		username = config.USERNAME.ValueString()
	}

	if !config.PASSWORD.IsNull() {
		password = config.PASSWORD.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing host",
			"",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing username",
			"",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing password",
			"",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	kumaClient := tools.KumaClient{
		Host:     host,
		Username: username,
		Password: password,
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.

	resp.DataSourceData = &kumaClient
	resp.ResourceData = &kumaClient
}

// DataSources defines the data sources implemented in the provider.
func (p *kumaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewmonitorDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *kumaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
	}
}
