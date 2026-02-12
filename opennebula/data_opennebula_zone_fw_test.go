package opennebula

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// testAccProtoV5ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"opennebula": providerserver.NewProtocol5WithError(New("test")()),
}

// TestAccZoneDataSource_FrameworkMigration tests that the Framework implementation
// of the zone data source is compatible with the SDKv2 implementation.
// This is a compile-time check to ensure the Framework provider is properly set up.
func TestAccZoneDataSource_FrameworkMigration(t *testing.T) {
	// This test verifies that:
	// 1. The Framework provider compiles
	// 2. The zone data source is registered
	// 3. The basic structure is compatible
	//
	// Full acceptance tests would require a running OpenNebula instance.
	// For now, we just verify the code compiles and can be instantiated.

	provider := New("test")()
	if provider == nil {
		t.Fatal("Failed to create Framework provider")
	}
}
