# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "consul_token"
  description = "The Consul ACL token used to register services with the Nomad job.\nUses the runner config environment variable CONSUL_HTTP_TOKEN."
  type        = ""
  required    = true
}

parameter {
  key           = "hcl1"
  description   = "Parses jobspec as HCL1 instead of HCL2."
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key         = "jobspec"
  description = "Path to a Nomad job specification file."
  type        = "string"
  required    = true
}

parameter {
  key         = "vault_token"
  description = "The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.\nUses the runner config environment variable VAULT_TOKEN."
  type        = ""
  required    = true
}

