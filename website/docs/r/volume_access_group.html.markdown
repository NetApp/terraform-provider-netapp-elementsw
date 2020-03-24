---
layout: "elementsw"
page_title: "ElementSW: elementsw_volume_access_group"
sidebar_current: "docs-elementsw-resource-volume-access-group"
description: |-
  Provides an ElementSW cluster volume access group resource. This can be used to create a new volume access group.
  Any initiator IQN that you add to the volume access group is able to access any volume in the group without CHAP
  authentication.
---

# elementsw\_volume\_access\_group

Provides an ElementSW cluster volume access group resource. This can be used to create a new volume access group.
Any initiator IQN that you add to the volume access group is able to access any volume in the group without CHAP
authentication.

## Example Usages

**Create ElementSW cluster volume access group:**

```
resource "elementsw_volume_access_group" "main-group" {
  name = "terraform-main-group"
  volumes = ["12345", "67890"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ElementSW volume access group.
* `volumes` - (Optional) The IDs of the ElementSW volumes to add to the
  ElementSW volume access group.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume access group.
* `initiators` - Any initiators tied to the volume access group.
