package elementsw

import (
	"strconv"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestVolume_basic(t *testing.T) {
	var volume volume
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwVolumeConfig,
					"terraform-acceptance-test",
					"1000000000",
					"true",
					"500",
					"10000",
					"10000",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwVolumeExists("elementsw_volume.terraform-acceptance-test-1", &volume),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "total_size", "1000000000"),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "enable512e", "true"),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "min_iops", "500"),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "max_iops", "10000"),
					resource.TestCheckResourceAttr("elementsw_volume.terraform-acceptance-test-1", "burst_iops", "10000"),
				),
			},
		},
	})
}

func testAccCheckElementSwVolumeDestroy(s *terraform.State) error {
	virConn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elementsw_volume" {
			continue
		}
		id, _ := strconv.Atoi(rs.Primary.ID)
		i := make([]int, 1)
		i[0] = id
		volumeInput := listVolumesRequest{}
		volumeInput.Volumes = i
		// should return an empty list
		retrievedVol, err := virConn.listVolumes(volumeInput)
		if err != nil {
			return err
		}
		if len(retrievedVol.Volumes) != 0 {
			return fmt.Errorf("Error waiting for volume (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckElementSwVolumeExists(n string, volume *volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElementSw volume key ID is set")
		}

		id, _ := strconv.Atoi(rs.Primary.ID)
		i := make([]int, 1)
		i[0] = id
		volumeInput := listVolumesRequest{}
		volumeInput.Volumes = i
		retrievedVol, err := virConn.listVolumes(volumeInput)
		if err != nil {
			return err
		}

		if retrievedVol.Volumes[0].VolumeID != volumeInput.Volumes[0] {
			return fmt.Errorf("Resource ID and volume ID do not match")
		}

		*volume = retrievedVol.Volumes[0]

		return nil
	}
}

const testAccCheckElementSwVolumeConfig = `
resource "elementsw_volume" "terraform-acceptance-test-1" {
	name = "%s"
	account = "${elementsw_account.terraform-acceptance-test-1.id}"
	total_size = "%s"
	enable512e = "%s"
	min_iops = "%s"
	max_iops = "%s"
	burst_iops = "%s"
}
resource "elementsw_account" "terraform-acceptance-test-1" {
	username = "terraform-acceptance-test-volume"
}
`
