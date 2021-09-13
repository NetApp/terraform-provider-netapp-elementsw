---
layout: "elementsw"
page_title: "Provider: ElementSW"
sidebar_current: "docs-elementsw-index"
description: |-
  The ElementSW provider is used to interact with the resources supported by
  ElementSW. The provider needs to be configured with the proper credentials
  before it can be used.
---

# ElementSW Provider

The ElementSW provider is used to interact with the resources supported by
ElementSW.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The ElementSW Provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it.

## Example Usage

```
# Configure ElementSW provider
provider "netapp-elementsw" {
  username         = var.elementsw_username
  password         = var.elementsw_password
  elementsw_server = var.elementsw_cluster
  api_version      = var.elementsw_api_version
}

# Specify ElementSW resources
resource "elementsw_account" test-account {
  provider = netapp-elementsw
  username = "testAccount"
}

resource "elementsw_volume" test-volume {
  # Create N instances
  count      = length(var.volume_size_list)
  provider   = netapp-elementsw
  name       = "${var.volume_name}-${count.index}"
  account    = elementsw_account.test-account.id
  total_size = var.volume_size_list[count.index]
  enable512e = true
  min_iops   = 100
  max_iops   = 500
  burst_iops = 1000
}

resource "elementsw_volume_access_group" test-group {
  provider = netapp-elementsw
  name     = "testGroup"
  volumes  = elementsw_volume.test-volume.*.id
}

resource "elementsw_initiator" test-initiator {
  provider               = netapp-elementsw
  name                   = "iqn.1998-01.com.vmware:test-terraform-000000"
  alias                  = "testIQN"
  volume_access_group_id = elementsw_volume_access_group.test-group.id
  iqns                   = elementsw_volume.test-volume.*.iqn
}
```

## Argument Reference

The following arguments are used to configure the ElementSW Provider:

* `elementsw_username` - (Required) This is the username for ElementSW API operations.
* `elementsw_password` - (Required) This is the password for ElementSW API operations.
* `elementsw_cluster` - (Required) This is the ElementSW cluster MVIP for ElementSW
  API operations.
* `elementsw_api_version` - (Required) This is the ElementSW API version for ElementSW
  API operations.

## Required Privileges

These settings were tested with NetApp ElementSW (Element OS, SolidFire) 11.7
For additional information on roles and permissions, please refer to official
ElementSW documentation.
