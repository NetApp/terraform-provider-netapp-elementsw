package elementsw

import (
	"strconv"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestInitiator_basic(t *testing.T) {
	var initiator initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwInitiatorExists("elementsw_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
		},
	})
}

func TestInitiator_update(t *testing.T) {
	var initiator initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwInitiatorExists("elementsw_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwInitiatorConfigUpdate,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwInitiatorExists("elementsw_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias-update"),
				),
			},
		},
	})
}

func TestInitiator_removeVolumeAccessGroup(t *testing.T) {
	var initiator initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwInitiatorExists("elementsw_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwInitiatorConfigRemoveVAG,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwInitiatorExists("elementsw_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias-update"),
				),
			},
		},
	})
}

func testAccCheckElementSwInitiatorDestroy(s *terraform.State) error {
	virConn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elementsw_initiator" {
			continue
		}

		_, err := virConn.getInitiatorByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Error waiting for initiator (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckElementSwInitiatorExists(n string, initiator *initiator) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElementSw initiator key ID is set")
		}

		retrievedInit, err := virConn.getInitiatorByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		convID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		if retrievedInit.InitiatorID != convID {
			return fmt.Errorf("Resource ID and initiator ID do not match")
		}

		*initiator = retrievedInit

		return nil
	}
}

const testAccCheckElementSwInitiatorConfig = `
resource "elementsw_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
	volume_access_group_id = "${elementsw_volume_access_group.terraform-acceptance-test-1.id}"
}

resource "elementsw_volume_access_group" "terraform-acceptance-test-1" {
	name = "terraform-acceptance-test-group"
}
`

const testAccCheckElementSwInitiatorConfigUpdate = `
resource "elementsw_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
	volume_access_group_id = "${elementsw_volume_access_group.terraform-acceptance-test-2.id}"
}

resource "elementsw_volume_access_group" "terraform-acceptance-test-2" {
	name = "terraform-acceptance-test-group-2"
}
`

const testAccCheckElementSwInitiatorConfigRemoveVAG = `
resource "elementsw_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
}
`
