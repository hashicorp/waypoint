# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "fail_deployment"
  description = "If true, marks the deployment as failed."
  type        = "bool"
  required    = false
}

parameter {
  key         = "groups"
  description = "List of task group names which are to be promoted."
  type        = "list of string"
  required    = false
}

