# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "base_url"
  description   = "The scheme, host, and port to calculate the URL to fetch using\nThis is provided to allow users to query values from Terraform Enterprise installations"
  type          = "string"
  required      = false
  default_value = "https://api.terraform.io"
}

parameter {
  key           = "refresh_interval"
  description   = "How often the outputs should be fetch.\nThe format of this value is the Go time duration format. Specifically a whole number followed by: s for seconds, m for minutes, h for hours. The minimum value for this setting is 60 seconds, with no specified maximum."
  type          = "string"
  required      = false
  default_value = "10m0s"
}

parameter {
  key         = "skip_verify"
  description = "Do not validate the TLS cert presented by Terraform Cloud.\nThis is not recommended unless absolutely necessary."
  type        = "bool"
  required    = false
}

parameter {
  key         = "token"
  description = "The Terraform Cloud API token\nThe token is used to authenticate access to the specific organization and workspace. Tokens are managed at https://app.terraform.io/app/settings/tokens."
  type        = "string"
  required    = true
}

