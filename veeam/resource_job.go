package veeam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobCreate,
		ReadContext:   resourceJobRead,
		UpdateContext: resourceJobUpdate,
		DeleteContext: resourceJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Backup",
					"VSphereReplica",
					"CloudDirectorBackup",
					"EntraIDTenantBackup",
					"EntraIDAuditLogBackup",
					"FileBackupCopy",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_high_priority": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"virtual_machines": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"includes": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Define fields for InventoryObjectModel
								},
							},
						},
						"excludes": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for BackupJobExclusionsSpec
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_repository_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"backup_proxies": {
							Type:     schema.TypeMap,
							Required: true,
							// Define fields for BackupProxiesSettingsModel
						},
						"retention_policy": {
							Type:     schema.TypeMap,
							Required: true,
							// Define fields for BackupJobRetentionPolicySettingsModel
						},
						"gfs_policy": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for GFSPolicySettingsModel
						},
						"advanced_settings": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for BackupJobAdvancedSettingsModel
						},
					},
				},
			},
			"guest_processing": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_aware_processing": {
							Type:     schema.TypeMap,
							Required: true,
							// Define fields for BackupApplicationAwareProcessingModel
						},
						"guest_fs_indexing": {
							Type:     schema.TypeMap,
							Required: true,
							// Define fields for GuestFileSystemIndexingModel
						},
						"guest_interaction_proxies": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for GuestInteractionProxiesSettingsModel
						},
						"guest_credentials": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for GuestOsCredentialsModel
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_automatically": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the job should run automatically.",
						},
						"daily": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleDailyModel
						},
						"monthly": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleMonthlyModel
						},
						"periodically": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for SchedulePeriodicallyModel
						},
						"continuously": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleBackupWindowModel
						},
						"after_this_job": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleAfterThisJobModel
						},
						"retry": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleRetryModel
						},
						"backup_window": {
							Type:     schema.TypeMap,
							Optional: true,
							// Define fields for ScheduleBackupWindowModel
						},
					},
				},
			},
		},
	}
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*http.Client)
	baseURL := m.(string)
	apiVersion := "1.2-rev0"
	url := fmt.Sprintf("%s/api/v1/jobs", baseURL)

	job := map[string]interface{}{
		"name":            d.Get("name").(string),
		"type":            d.Get("type").(string),
		"description":     d.Get("description").(string),
		"isHighPriority":  d.Get("is_high_priority").(bool),
		"virtualMachines": d.Get("virtual_machines").([]interface{}),
		"storage":         d.Get("storage").([]interface{}),
		"guestProcessing": d.Get("guest_processing").([]interface{}),
		"schedule":        d.Get("schedule").([]interface{}),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshaling job: %s", err))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating request: %s", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-version", apiVersion)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", baseURL))

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making request: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.FromErr(fmt.Errorf("error creating job: %s", resp.Status))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding response: %s", err))
	}

	d.SetId(result["id"].(string))
	return resourceJobRead(ctx, d, m)
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*http.Client)
	baseURL := m.(string)
	apiVersion := "1.2-rev0"
	jobID := d.Id()
	url := fmt.Sprintf("%s/api/v1/jobs/%s", baseURL, jobID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating request: %s", err))
	}

	req.Header.Set("x-api-version", apiVersion)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", baseURL))

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making request: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("error reading job: %s", resp.Status))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding response: %s", err))
	}

	d.Set("name", result["name"])
	d.Set("type", result["type"])
	d.Set("description", result["description"])
	d.Set("is_high_priority", result["isHighPriority"])
	d.Set("virtual_machines", result["virtualMachines"])
	d.Set("storage", result["storage"])
	d.Set("guest_processing", result["guestProcessing"])
	d.Set("schedule", result["schedule"])

	return nil
}

func resourceJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*http.Client)
	baseURL := m.(string)
	apiVersion := "1.2-rev0"
	jobID := d.Id()
	url := fmt.Sprintf("%s/api/v1/jobs/%s", baseURL, jobID)

	job := map[string]interface{}{
		"id":              jobID,
		"name":            d.Get("name").(string),
		"type":            d.Get("type").(string),
		"description":     d.Get("description").(string),
		"isHighPriority":  d.Get("is_high_priority").(bool),
		"virtualMachines": d.Get("virtual_machines").([]interface{}),
		"storage":         d.Get("storage").([]interface{}),
		"guestProcessing": d.Get("guest_processing").([]interface{}),
		"schedule":        d.Get("schedule").([]interface{}),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshaling job: %s", err))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating request: %s", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-version", apiVersion)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", baseURL))

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making request: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("error updating job: %s", resp.Status))
	}

	return resourceJobRead(ctx, d, m)
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*http.Client)
	baseURL := m.(string)
	apiVersion := "1.2-rev0"
	jobID := d.Id()
	url := fmt.Sprintf("%s/api/v1/jobs/%s", baseURL, jobID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating request: %s", err))
	}

	req.Header.Set("x-api-version", apiVersion)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", baseURL))

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making request: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.FromErr(fmt.Errorf("error deleting job: %s", resp.Status))
	}

	d.SetId("")
	return nil
}
