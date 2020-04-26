# Specify ElementSW resources
resource "elementsw_account" test-account {
    provider = netapp-elementsw
    username = "testAccount"
}

resource "elementsw_volume" test-volume {
    # Create N instances
    count = length(var.volume_size_list)
    provider = netapp-elementsw
    name = "${var.volume_name}-${count.index}"

    account = elementsw_account.test-account.id
    total_size = var.volume_size_list[count.index]
    enable512e = true
    min_iops = 100
    max_iops = 500
    burst_iops = 1000
}

resource "elementsw_volume_access_group" test-group {
    provider = netapp-elementsw
    name = "testGroup"
    volumes = elementsw_volume.test-volume.*.id
}

resource "elementsw_initiator" test-initiator {
    provider = netapp-elementsw
    name = "iqn.1998-01.com.vmware:test-es65-7f17a50c"
    alias = "testVAG"
    volume_access_group_id = elementsw_volume_access_group.test-group.id
    iqns = elementsw_volume.test-volume.*.iqn
}
