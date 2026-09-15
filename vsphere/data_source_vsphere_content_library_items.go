// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"github.com/vmware/terraform-provider-vsphere/vsphere/internal/helper/provider"
)

func dataSourceVSphereContentLibraryItems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereContentLibraryItemsRead,
		Schema: map[string]*schema.Schema{
			"library_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the content library to search.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Filter by name of the content library item. If omitted, all items matching the type filter are returned.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Filter by type of the content library item (e.g. ovf, iso, vm-template). If omitted, all items matching the name filter are returned.",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of content library items matching the search criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the content library item.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the content library item.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the content library item.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVSphereContentLibraryItemsRead(d *schema.ResourceData, meta interface{}) error {
	rc := meta.(*Client).restClient
	libraryID := d.Get("library_id").(string)
	name := d.Get("name").(string)
	itemType := d.Get("type").(string)

	items, err := contentlibrary.ItemsFromCriteria(rc, libraryID, name, itemType)
	if err != nil {
		return provider.Error(libraryID, "dataSourceVSphereContentLibraryItemsRead", err)
	}

	flat := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		flat = append(flat, map[string]interface{}{
			"id":   item.ID,
			"name": item.Name,
			"type": item.Type,
		})
	}
	if err := d.Set("items", flat); err != nil {
		return err
	}

	// Derive a stable ID from the search parameters
	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%s|%s|%s", libraryID, name, itemType)
	d.SetId(fmt.Sprintf("%x", h.Sum(nil)))
	return nil
}
