package veeam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_high_priority": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{}
	id := d.Get("id").(string)
	url := fmt.Sprintf("%s/api/v1/jobs/%s", m.(string), id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("x-api-version", "1.2-rev0")
	req.Header.Set("Authorization", "Bearer "+m.(string))

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to get job: received status code %d", resp.StatusCode)
	}

	var job map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("name", job["name"])
	d.Set("type", job["type"])
	d.Set("is_disabled", job["isDisabled"])
	d.Set("description", job["description"])
	d.Set("is_high_priority", job["isHighPriority"])
	// ...set other fields as needed...

	return diags
}
