---
layout: "elementsw"
page_title: "ElementSW: elementsw_volume"
sidebar_current: "docs-elementsw-resource-volume"
description: |-
  Provides an ElementSW cluster volume resource. This can be used to create a new (empty) volume on the cluster. As soon
  as the volume creation is complete, the volume is available for connection via iSCSI.
---

# elementsw\_volume

Provides an ElementSW cluster volume resource. This can be used to create a new (empty) volume on the cluster. As soon
as the volume creation is complete, the volume is available for connection via iSCSI.

## Example Usages

**Create ElementSW cluster volume:**

```
resource "elementsw_volume" "main-volume" {
  name        = "main-volume"
  account     = "1"
  total_size  = 10000000000
  enable512e  = true
  min_iops    = 50
  max_iops    = 10000
  burst_iops  = 10000
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ElementSW volume.
* `account` - (Required) The unique identifier or name of the ElementSW account owner.
* `total_size` - (Required) The total size of the volume, in bytes. Size is rounded up to the nearest 1MB size.
* `enable512e` - (Required) Whether to enable 512-byte sector emulation. The setting needs to
  be enabled if using VMWare.
* `min_iops` - (Optional) The minimum initial quality of service.
* `max_iops` - (Optional) The maximum initial quality of service.
* `burst_iops` - (Optional) The burst initial quality of service.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume.

## Unique id versus name

With ElementSW, every resource has a unique id, but names are not necessarily unique.

In order to make it easier to automate, we are assuming that volume names are unique within an
account (identified by a unique account id, or a customer enforced unique account name).
