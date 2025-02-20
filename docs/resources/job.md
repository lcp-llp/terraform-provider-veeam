# veeam_job

Provision a Veeam job resource.

## Example Usage

```hcl
resource "veeam_job" "example" {
  type            = "Backup"
  description     = "This is an example job"
  is_high_priority = true

  virtual_machines {
    includes = [
      {
        // Define fields for InventoryObjectModel
      }
    ]
    excludes = {
      // Define fields for BackupJobExclusionsSpec
    }
  }

  storage {
    backup_repository_id = "repository-id"
    backup_proxies = {
      // Define fields for BackupProxiesSettingsModel
    }
    retention_policy = {
      // Define fields for BackupJobRetentionPolicySettingsModel
    }
    gfs_policy = {
      // Define fields for GFSPolicySettingsModel
    }
    advanced_settings = {
      // Define fields for BackupJobAdvancedSettingsModel
    }
  }

  guest_processing {
    app_aware_processing = {
      // Define fields for BackupApplicationAwareProcessingModel
    }
    guest_fs_indexing = {
      // Define fields for GuestFileSystemIndexingModel
    }
    guest_interaction_proxies = {
      // Define fields for GuestInteractionProxiesSettingsModel
    }
    guest_credentials = {
      // Define fields for GuestOsCredentialsModel
    }
  }

  schedule {
    run_automatically = true
    daily = {
      // Define fields for ScheduleDailyModel
    }
    monthly = {
      // Define fields for ScheduleMonthlyModel
    }
    periodically = {
      // Define fields for SchedulePeriodicallyModel
    }
    continuously = {
      // Define fields for ScheduleBackupWindowModel
    }
    after_this_job = {
      // Define fields for ScheduleAfterThisJobModel
    }
    retry = {
      // Define fields for ScheduleRetryModel
    }
    backup_window = {
      // Define fields for ScheduleBackupWindowModel
    }
  }
}
```

## Argument Reference

* `name` - (Required) Name of the job.
* `type` - (Required) Type of the job. Valid values are `Backup`, `VSphereReplica`, `CloudDirectorBackup`, `EntraIDTenantBackup`, `EntraIDAuditLogBackup`, `FileBackupCopy`.
* `description` - (Required) Description of the job.
* `is_high_priority` - (Optional) If true, the resource scheduler prioritizes this job higher than other similar jobs and allocates resources to it in the first place. Defaults to `false`.

### Nested Blocks

#### `virtual_machines`

* `includes` - (Required) Array of VMs and VM containers processed by the job.
* `excludes` - (Optional) Objects excluded from the job.

#### `storage`

* `backup_repository_id` - (Required) ID of the backup repository.
* `backup_proxies` - (Required) Backup proxy settings.
* `retention_policy` - (Required) Retention policy settings.
* `gfs_policy` - (Optional) GFS retention policy settings.
* `advanced_settings` - (Optional) Advanced settings of the backup job.

#### `guest_processing`

* `app_aware_processing` - (Required) Application-aware processing settings.
* `guest_fs_indexing` - (Required) VM guest OS file indexing.
* `guest_interaction_proxies` - (Optional) Guest interaction proxy used to deploy the runtime process on the VM guest OS.
* `guest_credentials` - (Optional) VM custom credentials.

#### `schedule`

* `run_automatically` - (Required) If true, job scheduling is enabled. Defaults to `false`.
* `daily` - (Optional) Daily scheduling options.
* `monthly` - (Optional) Monthly scheduling options.
* `periodically` - (Optional) Periodic scheduling options.
* `continuously` - (Optional) Backup window settings.
* `after_this_job` - (Optional) Job chaining options.
* `retry` - (Optional) Retry options.
* `backup_window` - (Optional) Backup window settings.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the job.
