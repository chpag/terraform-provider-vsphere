// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vmware/terraform-provider-vsphere/vsphere/internal/helper/testhelper"
)

func TestAccDataSourceVSphereVirtualMachines_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereVirtualMachinesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.vsphere_virtual_machines.vms", "virtual_machines.#",
					),
				),
			},
		},
	})
}

func TestAccDataSourceVSphereVirtualMachines_nameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereVirtualMachinesNameRegexConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vsphere_virtual_machines.vms", "virtual_machines.0.name", "TestVM",
					),
					resource.TestCheckResourceAttrSet(
						"data.vsphere_virtual_machines.vms", "virtual_machines.0.id",
					),
				),
			},
		},
	})
}

func testAccDataSourceVSphereVirtualMachinesConfig() string {
	return fmt.Sprintf(`
%s

data "vsphere_virtual_machines" "vms" {
  datacenter_id = data.vsphere_datacenter.rootdc1.id
}
`, testhelper.ConfigDataRootDC1())
}

func testAccDataSourceVSphereVirtualMachinesNameRegexConfig() string {
	return fmt.Sprintf(`
%s

data "vsphere_virtual_machines" "vms" {
  datacenter_id = data.vsphere_datacenter.rootdc1.id
  name_regex    = "^TestVM$"
}
`, testhelper.ConfigDataRootDC1())
}
