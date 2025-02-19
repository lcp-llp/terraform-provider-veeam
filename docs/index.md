# Veeam Provider

This provider can used to interact with a Veeam Backup and Replication Server.

## Example Usage


```terraform
# This example fetches a known order ID
terraform {
  required_providers {
    bytes = {
      version = "~> 0.2.5"
      source  = "lcp-llp/veeam"
    }
  }
}

provider "veeam" {
  username = "example"
  password = "example"
  endpoint = "server.example.com
}

## Schema

### Required

- `password` (String) Password used for authentication to API Endpoints
- `username` (String) Username used for authentication to API Endpoints
- `identity_api_url` (String) The Endpoint of the Veeam server