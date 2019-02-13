package azuread

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataSourceGuest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGuestRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"mail": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func dataSourceGuestRead(d *schema.ResourceData, meta interface{}) error {
	// TODO: implement
	return nil
}
