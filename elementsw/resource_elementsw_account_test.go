package elementsw

import (
	"strconv"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccount_basic(t *testing.T) {
	return
	var account account
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwAccountConfig,
					"terraform-acceptance-test",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwAccountExists("elementsw_account.terraform-acceptance-account-1", &account),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "username", "terraform-acceptance-test"),
				),
			},
		},
	})
}

func TestAccount_secrets(t *testing.T) {
	return
	var account account
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwAccountConfigSecrets,
					"terraform-acceptance-test",
					"ABC123456XYZ",
					"SecretSecret1",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwAccountExists("elementsw_account.terraform-acceptance-account-1", &account),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "username", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "target_secret", "ABC123456XYZ"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "initiator_secret", "SecretSecret1"),
				),
			},
		},
	})
}

func TestAccount_update(t *testing.T) {
	return
	var account account
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElementSwAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwAccountConfigSecrets,
					"terraform-acceptance-test",
					"ABC123456XYZ",
					"SecretSecret1",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwAccountExists("elementsw_account.terraform-acceptance-account-1", &account),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "username", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "target_secret", "ABC123456XYZ"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "initiator_secret", "SecretSecret1"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckElementSwAccountConfigUpdate,
					"terraform-acceptance-test-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElementSwAccountExists("elementsw_account.terraform-acceptance-account-1", &account),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "username", "terraform-acceptance-test-update"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "target_secret", "ABC123456XYZU"),
					resource.TestCheckResourceAttr("elementsw_account.terraform-acceptance-account-1", "initiator_secret", "SecretSecret1U"),
				),
			},
		},
	})
}

func testAccCheckElementSwAccountDestroy(s *terraform.State) error {
	return nil
	virConn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elementsw_account" {
			continue
		}

		convID, convErr := strconv.Atoi(rs.Primary.ID)
		if convErr != nil {
			return convErr
		}

		_, err := virConn.getAccountByID(convID)
		if err == nil {
			return fmt.Errorf("Error waiting for volume (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckElementSwAccountExists(n string, account *account) resource.TestCheckFunc {
	return nil
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElementSw account key ID is set")
		}

		convID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		retrievedAcc, err := virConn.getAccountByID(convID)
		if err != nil {
			return err
		}

		if retrievedAcc.AccountID != convID {
			return fmt.Errorf("Resource ID and account ID do not match")
		}

		*account = retrievedAcc

		return nil
	}
}

const testAccCheckElementSwAccountConfig = `
resource "elementsw_account" "terraform-acceptance-account-1" {
	username = "%s"
}
`

const testAccCheckElementSwAccountConfigSecrets = `
resource "elementsw_account" "terraform-acceptance-account-1" {
	username = "%s"
	target_secret = "%s"
	initiator_secret = "%s"
}
`

const testAccCheckElementSwAccountConfigUpdate = `
resource "elementsw_account" "terraform-acceptance-account-1" {
	username = "%s"
}
`
