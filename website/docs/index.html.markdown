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
# Configure the ElementSW Provider
provider "elementsw" {
  username         = "${var.elementsw_username}"
  password         = "${var.elementsw_password}"
  elementsw_server = "${var.elementsw_server}"
  api_version      = "${var.elementsw_api_version}"
}

# Create an account
resource "elementsw_account" "main-account" {
  username = "main"
}

# Create a volume tied to an account
resource "elementsw_volume" "volume1" {
  name       = "main-volume"
  accountID  = "${elementsw_account.main-account.id}"
  totalSize  = 10000000000
  enable512e = true
  minIOPS    = 50
  maxIOPS    = 10000
  burstIOPS  = 10000
}

# Create a volume access group for the volume
resource "elementsw_volume_access_group" "main-group" {
  name    = "main-volume-access-group"
  volumes = ["${elementsw_volumes.volume1.id}"]
}

# Create an initiator for the volume access group
resource "elementsw_initiator" "main-initiator" {
  name = "qn.1998-01.com.vmware:test-terraform-00000000"
  alias = "Main Initiator"
  volumeAccessGroupID = "${elementsw_volume_access_group.main-group.id}"
}
```

## Argument Reference

The following arguments are used to configure the ElementSW Provider:

* `username` - (Required) This is the username for ElementSW API operations.
* `password` - (Required) This is the password for ElementSW API operations.
* `elementsw_server` - (Required) This is the ElementSW cluster name for ElementSW 
  API operations.
* `api_version` - (Required) This is the ElementSW cluster version for ElementSW
  API operations.

## Required Privileges

These settings were tested with NetApp ElementSW Element OS 11.1
For additional information on roles and permissions, please refer to official
ElementSW documentation.

