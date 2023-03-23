# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "client_id"
  description = "The OAuth2 Client ID for HCP API operations."
  type        = "string"
  required    = false
}

parameter {
  key         = "client_secret"
  description = "The OAuth2 Client Secret for HCP API operations."
  type        = "string"
  required    = false
}

parameter {
  key         = "organization_id"
  description = "The HCP organization ID."
  type        = "string"
  required    = true
}

parameter {
  key         = "project_id"
  description = "The HCP Project ID."
  type        = "string"
  required    = true
}

