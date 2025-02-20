# veeam_job

Provides details about a specific Veeam job.

## Example Usage

```hcl
data "veeam_job" "example" {
  id = "job-id"
}

output "job_name" {
  value = data.veeam_job.example.name
}

output "job_type" {
  value = data.veeam_job.example.type
}

output "job_description" {
  value = data.veeam_job.example.description
}

output "job_is_high_priority" {
  value = data.veeam_job.example.is_high_priority
}
```

## Argument Reference

* `id` - (Required) The ID of the job.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the job.
* `type` - The type of the job.
* `is_disabled` - Whether the job is disabled.
* `description` - The description of the job.
* `is_high_priority` - Whether the job is high priority.
