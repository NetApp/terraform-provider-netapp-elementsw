---
layout: "elementsw"
page_title: "ElementSW: elementsw_account"
sidebar_current: "docs-elementsw-resource-account"
description: |-
  Provides an ElementSW cluster account resource. This can be used to ...
---

# elementsw\_account

Provides an ElementSW cluster account resource. This can be used to ...

## Example Usages

**Create ElementSW cluster account:**

```
resource "elementsw_account" "main-account" {
  username         = "main"
  initiator_secret = "s!39naDlLa9"
  target_secret    = "2Z>D0jf3Dpa"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The name of the ElementSW account.
* `initiator_secret` - (Optional) The initiator secret. If not specified, the ElementSW cluster will autogenerate
  an initiator secret.
* `target_secret` - (Optional) The target secret. If not specified, the ElementSW cluster will autogenerate
  an initiator secret.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the account.
