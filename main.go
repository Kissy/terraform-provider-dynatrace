package main

import (
	"github.com/Kissy/terraform-provider-dynatrace/dynatrace"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return dynatrace.Provider()
		},
	})
}