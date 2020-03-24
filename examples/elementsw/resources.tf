# Specify ElementSW resources
resource "elementsw_account" "test-account" {
    provider = "netapp-elementsw"
    username = "test Account"
}

resource "elementsw_volume" "test-volume" {
    provider = "netapp-elementsw"
    name = "test-Volume-${count.index}"
    account = "${elementsw_account.test-account.id}"
    total_size = "${var.total_size[count.index]}"
    enable512e = true
    min_iops = 50
    max_iops = 10000
    burst_iops = 10000

    # Create N instances
    count = "${length(var.total_size)}"
}

resource "elementsw_volume_access_group" "test-group" {
    provider = "netapp-elementsw"
    name = "test-Group"
    volumes = "${elementsw_volume.test-volume.*.id}"
}

resource "elementsw_initiator" "test-initiator" {
    provider = "netapp-elementsw"
    name = "iqn.1998-01.com.vmware:test-es65-7f17a50c"
    alias = "test EUI Cluster"
    volume_access_group_id = "${elementsw_volume_access_group.test-group.id}"
    iqns = "${elementsw_volume.test-volume.*.iqn}"
}