# Specify ElementSW resources
resource "elementsw_account" test-account {
  provider = netapp-elementsw
  username = var.elementsw_tenant_name 
}

resource "elementsw_volume" test-volume {
  # Create N instances
  count      = length(var.volume_size_list)
  provider   = netapp-elementsw
  name       = "${var.volume_name}-${count.index}"
  account    = elementsw_account.test-account.id
  total_size = var.volume_size_list[count.index]
  enable512e = var.sectorsize_512e
  min_iops   = var.qos.min
  max_iops   = var.qos.max
  burst_iops = var.qos.burst
}

resource "elementsw_volume_access_group" test-group {
  provider = netapp-elementsw
  name     = var.volume_group_name
  volumes  = elementsw_volume.test-volume.*.id
}

resource "elementsw_initiator" test-initiator {
  provider               = netapp-elementsw
  name                   = var.elementsw_initiator.name
  alias                  = var.elementsw_initiator.alias
  volume_access_group_id = elementsw_volume_access_group.test-group.id
  iqns                   = elementsw_volume.test-volume.*.iqn
}
