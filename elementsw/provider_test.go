package elementsw

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"elementsw": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ELEMENTSW_USERNAME"); v == "" {
		t.Fatal("ELEMENTSW_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_PASSWORD"); v == "" {
		t.Fatal("ELEMENTSW_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_SERVER"); v == "" {
		t.Fatal("ELEMENTSW_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("ELEMENTSW_API_VERSION"); v == "" {
		t.Fatal("ELEMENTSW_API_VERSION must be set for acceptance tests")
	}
}
