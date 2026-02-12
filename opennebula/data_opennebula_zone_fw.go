package opennebula

import (
	"context"
	"fmt"

	zoneSc "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/zone"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &zoneDataSource{}
	_ datasource.DataSourceWithConfigure = &zoneDataSource{}
)

// NewZoneDataSource is a helper function to simplify the provider implementation.
func NewZoneDataSource() datasource.DataSource {
	return &zoneDataSource{}
}

// zoneDataSource is the data source implementation.
type zoneDataSource struct {
	config *Configuration
}

// zoneDataSourceModel maps the data source schema data.
type zoneDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Endpoint types.String `tfsdk:"endpoint"`
}

// Metadata returns the data source type name.
func (d *zoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

// Schema defines the schema for the data source.
func (d *zoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about an OpenNebula zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID of the zone. Defaults to -1 (any zone).",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the zone.",
				Optional:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint of the zone.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *zoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Configuration)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Configuration, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

// Read refreshes the Terraform state with the data source information.
func (d *zoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zoneDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get zone information from OpenNebula
	zone, err := d.filterZone(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Zone",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.Int64Value(int64(zone.ID))
	state.Name = types.StringValue(zone.Name)
	state.Endpoint = types.StringValue(zone.Template.Endpoint)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// filterZone finds a zone matching the configured criteria.
func (d *zoneDataSource) filterZone(model zoneDataSourceModel) (*zoneSc.Zone, error) {
	controller := d.config.Controller

	zones, err := controller.Zones().Info()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve zones: %w", err)
	}

	// Filter zones with user-defined criteria
	var id int64 = -1
	if !model.ID.IsNull() {
		id = model.ID.ValueInt64()
	}

	name := model.Name.ValueString()
	hasName := !model.Name.IsNull()

	match := make([]*zoneSc.Zone, 0, 1)
	for i, zone := range zones.Zones {
		// Filter by ID if specified and not -1
		if id != -1 && int64(zone.ID) != id {
			continue
		}

		// Filter by name if specified
		if hasName && zone.Name != name {
			continue
		}

		match = append(match, &zones.Zones[i])
	}

	// Check filtering results
	if len(match) == 0 {
		return nil, fmt.Errorf("no zone matches the constraints")
	} else if len(match) > 1 {
		return nil, fmt.Errorf("several zones match the constraints")
	}

	return match[0], nil
}
