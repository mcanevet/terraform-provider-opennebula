package opennebula

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVirtualNetworkAddressRange_UpdateIP reproduces GitHub issue #505:
// changing ip4 of an AR that has active leases (e.g. held IPs or VM NIC leases)
// was failing with "Address Range has leases in use" because the update path was
// trying to remove and recreate the AR instead of updating it in-place.
//
// The test simulates the lease-in-use scenario by keeping hold_ips unchanged
// while changing ip4, so the held lease is still present when the AR update runs.
func TestAccVirtualNetworkAddressRange_UpdateIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVNetARIPUpdateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "ip4", "172.16.200.10"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "size", "20"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "hold_ips.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "hold_ips.0", "172.16.200.15"),
				),
			},
			{
				// Change ip4 while hold_ips stays the same.
				// The held IP creates an active lease in the AR.
				// With the old remove+add approach this fails with:
				//   "Address Range has leases in use"
				// With the fix (in-place UpdateAR) this should succeed.
				Config: testAccVNetARIPUpdateConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "ip4", "172.16.200.11"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "size", "20"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "hold_ips.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_network_address_range.test", "hold_ips.0", "172.16.200.15"),
				),
			},
		},
	})
}

var testAccVNetARIPUpdateConfigBasic = `
resource "opennebula_virtual_network" "test" {
  name   = "test-issue-505"
  type   = "dummy"
  bridge = "onebr"
}

resource "opennebula_virtual_network_address_range" "test" {
  virtual_network_id = opennebula_virtual_network.test.id
  ar_type            = "IP4"
  size               = 20
  ip4                = "172.16.200.10"
  hold_ips           = ["172.16.200.15"]
}
`

var testAccVNetARIPUpdateConfigUpdated = `
resource "opennebula_virtual_network" "test" {
  name   = "test-issue-505"
  type   = "dummy"
  bridge = "onebr"
}

resource "opennebula_virtual_network_address_range" "test" {
  virtual_network_id = opennebula_virtual_network.test.id
  ar_type            = "IP4"
  size               = 20
  ip4                = "172.16.200.11"
  hold_ips           = ["172.16.200.15"]
}
`
