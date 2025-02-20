package veeam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_USERNAME", nil),
				Description: "The username to authenticate to the Veeam server.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_PASSWORD", nil),
				Description: "The password to authenticate to the Veeam server.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_ENDPOINT", nil),
				Description: "The endpoint of the Veeam server.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"veeam_job": resourceJob(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"veeam_job": dataSourceJob(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Endpoint: d.Get("endpoint").(string),
	}
	return config.Client()
}
