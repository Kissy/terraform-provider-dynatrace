package dynatrace

import (
	"fmt"
	apiclient "github.com/Kissy/go-dynatrace/dynatrace/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},

			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/api/config/v1",
			},

			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYNATRACE_TOKEN", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_alerting_profile":   resourceAlertingProfile(),
			"dynatrace_maintenance_window": resourceMaintenanceWindow(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	transport := httptransport.New(d.Get("host").(string), d.Get("base_path").(string), nil)
	apiClient := apiclient.New(transport, strfmt.Default)

	authorizationHeader := fmt.Sprintf("Api-Token %s", d.Get("token").(string))
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", authorizationHeader)
	return apiClient, nil
}
