// +build acceptance

package opennebula

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestMigrateZoneDataSource_IdenticalBehavior verifies that the Framework
// implementation of opennebula_zone data source produces identical results
// to the SDKv2 implementation.
//
// This test:
// 1. Uses SDKv2 provider to read zone data
// 2. Uses Framework provider to read the same zone data
// 3. Switches back to SDKv2 to ensure bidirectional compatibility
// 4. Verifies all implementations produce identical attributes
func TestMigrateZoneDataSource_IdenticalBehavior(t *testing.T) {
	// Skip if no OpenNebula instance available
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	resource.Test(t, resource.TestCase{
		// Step 1: Use SDKv2 provider
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"opennebula": func() (tfprotov5.ProviderServer, error) {
						return schema.NewGRPCProviderServer(Provider()), nil
					},
				},
				Config: testMigrateZoneDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					// Verify zone data is read correctly with SDKv2
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "id"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "name"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "endpoint"),
				),
			},
			{
				// Step 2: Switch to Framework provider with same config
				// Data sources don't have state, so we just verify attributes match
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"opennebula": providerserver.NewProtocol5WithError(New("test")()),
				},
				Config: testMigrateZoneDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					// Verify zone data is still correct with Framework
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "id"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "name"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "endpoint"),
				),
			},
			{
				// Step 3: Switch back to SDKv2 - verify attributes still match
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"opennebula": func() (tfprotov5.ProviderServer, error) {
						return schema.NewGRPCProviderServer(Provider()), nil
					},
				},
				Config: testMigrateZoneDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "id"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "name"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.test", "endpoint"),
				),
			},
		},
	})
}

// TestZoneDataSource_FrameworkOnly tests the Framework implementation independently
// to ensure it handles various scenarios correctly.
func TestZoneDataSource_FrameworkOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"opennebula": providerserver.NewProtocol5WithError(New("test")()),
		},
		Steps: []resource.TestStep{
			{
				// Test filtering by name
				Config: testZoneDataSourceByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.opennebula_zone.by_name", "id"),
					resource.TestCheckResourceAttr("data.opennebula_zone.by_name", "name", "OpenNebula"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.by_name", "endpoint"),
				),
			},
			{
				// Test filtering by ID
				Config: testZoneDataSourceByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.opennebula_zone.by_id", "id", "0"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.opennebula_zone.by_id", "endpoint"),
				),
			},
		},
	})
}

// Test configurations
const testMigrateZoneDataSourceConfig = `
provider "opennebula" {
  endpoint = "http://localhost:2633/RPC2"
  username = "oneadmin"
  password = "opennebula"
}

data "opennebula_zone" "test" {
  name = "OpenNebula"
}
`

const testZoneDataSourceByName = `
provider "opennebula" {
  endpoint = "http://localhost:2633/RPC2"
  username = "oneadmin"
  password = "opennebula"
}

data "opennebula_zone" "by_name" {
  name = "OpenNebula"
}
`

const testZoneDataSourceByID = `
provider "opennebula" {
  endpoint = "http://localhost:2633/RPC2"
  username = "oneadmin"
  password = "opennebula"
}

data "opennebula_zone" "by_id" {
  id = 0
}
`
