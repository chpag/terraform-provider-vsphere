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

func TestAccDataSourceVSphereContentLibraryItems_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereContentLibraryItemsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vsphere_content_library_items.items", "items.#", "1",
					),
					resource.TestCheckResourceAttr(
						"data.vsphere_content_library_items.items", "items.0.name", "TinyVM",
					),
					resource.TestCheckResourceAttr(
						"data.vsphere_content_library_items.items", "items.0.type", "ovf",
					),
				),
			},
		},
	})
}

func TestAccDataSourceVSphereContentLibraryItems_allItems(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereContentLibraryItemsAllConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.vsphere_content_library_items.all", "items.#",
					),
				),
			},
		},
	})
}

func testAccDataSourceVSphereContentLibraryItemsConfig() string {
	return fmt.Sprintf(`
%s

variable "file" {
  type    = string
  default = "%s"
}

resource "vsphere_content_library" "library" {
  name            = "ContentLibrary_items_test"
  storage_backing = [data.vsphere_datastore.rootds1.id]
  description     = "Library Description"
}

resource "vsphere_content_library_item" "item" {
  name       = "TinyVM"
  library_id = vsphere_content_library.library.id
  type       = "ova"
  file_url   = var.file
}

data "vsphere_content_library_items" "items" {
  library_id = vsphere_content_library.library.id
  name       = vsphere_content_library_item.item.name
  type       = "ovf"
}
`, testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootDS1()),
		testhelper.TestOva,
	)
}

func testAccDataSourceVSphereContentLibraryItemsAllConfig() string {
	return fmt.Sprintf(`
%s

variable "file" {
  type    = string
  default = "%s"
}

resource "vsphere_content_library" "library" {
  name            = "ContentLibrary_items_all_test"
  storage_backing = [data.vsphere_datastore.rootds1.id]
  description     = "Library Description"
}

resource "vsphere_content_library_item" "item" {
  name       = "TinyVM"
  library_id = vsphere_content_library.library.id
  type       = "ova"
  file_url   = var.file
}

data "vsphere_content_library_items" "all" {
  library_id = vsphere_content_library.library.id

  depends_on = [vsphere_content_library_item.item]
}
`, testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootDS1()),
		testhelper.TestOva,
	)
}
