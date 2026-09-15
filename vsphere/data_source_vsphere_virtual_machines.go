// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"crypto/sha256"
	"fmt"
	"log"
	"path"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
)

func dataSourceVSphereVirtualMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereVirtualMachinesRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The managed object ID of the datacenter to search in. If not set, the search is performed on the root of the inventory.",
			},
			"folder": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The absolute path of the folder to restrict the search to (e.g. 'folder/subfolder'). If not set, the entire datacenter is searched.",
			},
			"name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regular expression used to filter virtual machines by name. If not set, all virtual machines are returned.",
			},
			"virtual_machines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of virtual machines matching the search criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The managed object ID (MOID) of the virtual machine.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the virtual machine.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVSphereVirtualMachinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient

	// Resolve the search root path
	searchRoot := "/*"
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		dc, err := datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
		log.Printf("[DEBUG] dataSourceVSphereVirtualMachinesRead: Searching in datacenter %s", dc.InventoryPath)
		searchRoot = dc.InventoryPath
		if folderName, ok := d.GetOk("folder"); ok {
			searchRoot = path.Join(searchRoot, "vm", folderName.(string))
		}
	} else if folderName, ok := d.GetOk("folder"); ok {
		searchRoot = "/" + folderName.(string)
	}

	log.Printf("[DEBUG] dataSourceVSphereVirtualMachinesRead: Listing VMs under %s", searchRoot)
	vms, err := virtualmachine.ListFromPath(client, searchRoot)
	if err != nil {
		return fmt.Errorf("error listing virtual machines: %s", err)
	}

	// Compile optional name filter
	var re *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		re, err = regexp.Compile(nameRegex.(string))
		if err != nil {
			return fmt.Errorf("invalid name_regex: %s", err)
		}
	}

	flat := make([]map[string]interface{}, 0, len(vms))
	for _, vm := range vms {
		name := path.Base(vm.InventoryPath)
		if re != nil && !re.MatchString(name) {
			continue
		}
		flat = append(flat, map[string]interface{}{
			"id":   vm.Reference().Value,
			"name": name,
		})
	}

	if err := d.Set("virtual_machines", flat); err != nil {
		return err
	}

	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%s|%s|%s", d.Get("datacenter_id"), d.Get("folder"), d.Get("name_regex"))
	d.SetId(fmt.Sprintf("%x", h.Sum(nil)))

	log.Printf("[DEBUG] dataSourceVSphereVirtualMachinesRead: Found %d VM(s)", len(flat))
	return nil
}
