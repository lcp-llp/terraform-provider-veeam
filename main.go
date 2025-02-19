package main

import (
	"github.com/lcp-llp/terraform-provider-veeam/veeam"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return veeam.Provider()
		},
	})
}
