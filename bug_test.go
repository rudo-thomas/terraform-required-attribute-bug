package bug

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const template = `
terraform {
	required_version = ">= 0.15.4"
}
resource "test_test" "r1" {
	l = [
		[
			{
				// Attribute "s" is optional here, but if not present,
				// terraform 0.15.4 (and older) fails and says it is required.
				%s s = "foo"
			},
		],
	]
}
`

func runTest(t *testing.T, steps ...resource.TestStep) {
	t.Helper()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: steps,
	})

}

func TestOk(t *testing.T) {
	runTest(t, resource.TestStep{
		Config: fmt.Sprintf(template, ""),
		Check:  resource.TestCheckResourceAttr("test_test.r1", "l.0.0.s", "foo"),
	})
}

func TestBug(t *testing.T) {
	runTest(t, resource.TestStep{
		Config: fmt.Sprintf(template, "//"),
		Check:  resource.TestCheckResourceAttrSet("test_test.r1", "id"),
	})
}
