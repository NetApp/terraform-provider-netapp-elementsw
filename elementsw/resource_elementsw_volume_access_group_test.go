package elementsw

import (
	"strconv"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestVolumeAccessGroup_basic(t *testing.T) {
	var volumeAccessGroup volumeAccessGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwVolumeAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwVolumeAccessGroupConfig,
					"terraform-acceptance-test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwVolumeAccessGroupExists("elementsw_volume_access_group.terraform-acceptance-test-1", &volumeAccessGroup),
					resource.TestCheckResourceAttr("elementsw_volume_access_group.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttrSet("elementsw_volume_access_group.terraform-acceptance-test-1", "id"),
				),
			},
		},
	})
}

func TestVolumeAccessGroup_update(t *testing.T) {
	var volumeAccessGroup volumeAccessGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwVolumeAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwVolumeAccessGroupConfig,
					"terraform-acceptance-test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwVolumeAccessGroupExists("elementsw_volume_access_group.terraform-acceptance-test-1", &volumeAccessGroup),
					resource.TestCheckResourceAttr("elementsw_volume_access_group.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwVolumeAccessGroupConfigUpdate,
					"terraform-acceptance-test-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwVolumeAccessGroupExists("elementsw_volume_access_group.terraform-acceptance-test-1", &volumeAccessGroup),
					resource.TestCheckResourceAttr("elementsw_volume_access_group.terraform-acceptance-test-1", "name", "terraform-acceptance-test-update"),
				),
			},
		},
	})
}

// func TestVolumeAccessGroup_removeVolumes(t *testing.T) {
// 	var volumeAccessGroup volumeAccessGroup
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckElementSwVolumeAccessGroupDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: fmt.Sprintf(
// 					testAccCheckElementSwVolumeAccessGroupConfig,
// 					"terraform-acceptance-test",
// 				),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckElementSwVolumeAccessGroupExists("elementsw_volume_access_group.terraform-acceptance-test-1", &volumeAccessGroup),
// 					resource.TestCheckResourceAttr("elementsw_volume_access_group.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
// 				),
// 			},
// 			{
// 				Config: fmt.Sprintf(
// 					testAccCheckElementSwVolumeAccessGroupConfigRemoveVolumes,
// 					"terraform-acceptance-test-update",
// 				),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckElementSwVolumeAccessGroupExists("elementsw_volume_access_group.terraform-acceptance-test-1", &volumeAccessGroup),
// 					resource.TestCheckResourceAttr("elementsw_volume_access_group.terraform-acceptance-test-1", "name", "terraform-acceptance-test-update"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckElementSwVolumeAccessGroupDestroy(s *terraform.State) error {
	virConn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elementsw_volume_access_group" {
			continue
		}

		_, err := virConn.getVolumeAccessGroupByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Error waiting for volume access group (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckElementSwVolumeAccessGroupExists(n string, volume *volumeAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElementSw volume access group key ID is set")
		}

		retrievedVAG, err := virConn.getVolumeAccessGroupByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		convID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		if retrievedVAG.VolumeAccessGroupID != convID {
			return fmt.Errorf("Resource ID and volume access group ID do not match")
		}

		*volume = retrievedVAG

		return nil
	}
}

const testAccCheckElementSwVolumeAccessGroupConfig = `
resource "elementsw_volume_access_group" "terraform-acceptance-test-1" {
	name = "%s"
	volumes = ["${elementsw_volume.terraform-acceptance-test-1.id}"]
}
resource "elementsw_volume" "terraform-acceptance-test-1" {
	name = "Terraform-Acceptance-Volume-1"
	account = "${elementsw_account.terraform-acceptance-test-1.id}"
	total_size = "1000000000"
	enable512e = "true"
	min_iops = "600"
	max_iops = "8000"
	burst_iops = "8000"
}
resource "elementsw_account" "terraform-acceptance-test-1" {
	username = "terraform-acceptance-test-vag"
}
`

const testAccCheckElementSwVolumeAccessGroupConfigUpdate = `
resource "elementsw_volume_access_group" "terraform-acceptance-test-1" {
	name = "%s"
	volumes = ["${elementsw_volume.terraform-acceptance-test-1.id}", "${elementsw_volume.terraform-acceptance-test-2.id}"]
}
resource "elementsw_volume" "terraform-acceptance-test-1" {
	name = "Terraform-Acceptance-Volume-1"
	account = "${elementsw_account.terraform-acceptance-test-1.id}"
	total_size = "1000000000"
	enable512e = "true"
	min_iops = "600"
	max_iops = "8000"
	burst_iops = "8000"
}
resource "elementsw_volume" "terraform-acceptance-test-2" {
	name = "Terraform-Acceptance-Volume-2"
	account = "${elementsw_account.terraform-acceptance-test-1.id}"
	total_size = "1000000000"
	enable512e = "true"
	min_iops = "600"
	max_iops = "8000"
	burst_iops = "8000"
}
resource "elementsw_account" "terraform-acceptance-test-1" {
	username = "terraform-acceptance-test-vag"
}
`

const testAccCheckElementSwVolumeAccessGroupConfigRemoveVolumes = `
resource "elementsw_volume_access_group" "terraform-acceptance-test-1" {
	name = "%s"
	volumes = ["%s"]
}
`
