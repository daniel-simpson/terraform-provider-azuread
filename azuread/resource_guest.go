package azuread

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGuest() *schema.Resource {
	return &schema.Resource{
		Create: resourceGuestCreate,
		Read:   resourceGuestRead,
		Delete: resourceGuestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"mail": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGuestCreate(d *schema.ResourceData, meta interface{}) error {
	// TODO: implement
	// [Add user to tenant (via invitation)](https://docs.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-1.0)

	return nil
}

func resourceGuestRead(d *schema.ResourceData, meta interface{}) error {
	// TODO: implement
	//[Get User in Tenant](https://docs.microsoft.com/en-us/graph/api/user-get?view=graph-rest-1.0)

	return nil
}

func resourceGuestDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO: implement
	//[Delete user from tenant](https://docs.microsoft.com/en-us/graph/api/user-delete?view=graph-rest-1.0)

	return nil
}
