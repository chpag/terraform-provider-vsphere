---
subcategory: "Virtual Machine"
page_title: "VMware vSphere: vsphere_virtual_machines"
sidebar_current: "docs-vsphere-data-source-virtual-machines"
description: |-
  Provides a VMware vSphere data source to search for virtual machines.
---

# vsphere_virtual_machines

The `vsphere_virtual_machines` data source can be used to search for virtual
machines in a datacenter or folder and return a list of matching results. You
can optionally filter results by name using a regular expression.

## Example Usage

### List all virtual machines in a datacenter

```hcl
data "vsphere_datacenter" "dc" {
  name = "dc-01"
}

data "vsphere_virtual_machines" "all" {
  datacenter_id = data.vsphere_datacenter.dc.id
}
```

### Filter by name regex

```hcl
data "vsphere_datacenter" "dc" {
  name = "dc-01"
}

data "vsphere_virtual_machines" "web" {
  datacenter_id = data.vsphere_datacenter.dc.id
  name_regex    = "^web-"
}

output "web_vm_ids" {
  value = [for vm in data.vsphere_virtual_machines.web.virtual_machines : vm.id]
}
```

### Restrict search to a specific folder

```hcl
data "vsphere_datacenter" "dc" {
  name = "dc-01"
}

data "vsphere_virtual_machines" "prod" {
  datacenter_id = data.vsphere_datacenter.dc.id
  folder        = "prod/workloads"
  name_regex    = "^app-"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter_id` - (Optional) The managed object ID of the datacenter to
  search in. If not set, the search starts at the root of the inventory.
* `folder` - (Optional) Relative folder path within the datacenter (e.g.
  `"prod/workloads"`) to restrict the search to. If not set, the entire
  datacenter is searched.
* `name_regex` - (Optional) A regular expression used to filter virtual
  machines by name. If not set, all virtual machines found under the search
  root are returned.

## Attribute Reference

* `id` - A deterministic identifier derived from the search parameters.
* `virtual_machines` - A list of virtual machines matching the search criteria.
  Each element contains:
  * `id` - The managed object ID (MOID) of the virtual machine.
  * `name` - The name of the virtual machine.
