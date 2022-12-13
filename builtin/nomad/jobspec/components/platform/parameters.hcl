parameter {
  key         = "consul_token"
  description = <<EOT
The Consul ACL token used to register services with the Nomad job.
Uses the runner config environment variable CONSUL\_HTTP\_TOKEN.

EOT
  type        = "string" # WARNING: no type was documented. This will be a best effort choice.
  required    = true

}

parameter {
  key         = "jobspec"
  description = <<EOT
Path to a Nomad job specification file.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "vault_token"
  description = <<EOT
The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.
Uses the runner config environment variable VAULT\_TOKEN.

EOT
  type        = "string" # WARNING: no type was documented. This will be a best effort choice.
  required    = true

}

parameter {
  key           = "hcl1"
  description   = <<EOT
Parses jobspec as HCL1 instead of HCL2.

EOT
  type          = "bool"
  required      = false
  default_value = "false"
}

