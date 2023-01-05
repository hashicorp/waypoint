# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "auth"
  description = "The credentials for docker registry."
  type        = "nomad.AuthConfig"
  required    = true
}

parameter {
  key         = "consul_token"
  description = "The Consul ACL token used to register services with the Nomad job.\nUses the runner config environment variable CONSUL_HTTP_TOKEN."
  type        = ""
  required    = true
}

parameter {
  key           = "datacenter"
  description   = "The Nomad datacenter to deploy the job to."
  type          = "string"
  required      = false
  default_value = "dc1"
}

parameter {
  key         = "namespace"
  description = "The Nomad namespace to deploy the job to."
  type        = "string"
  required    = false
}

parameter {
  key           = "region"
  description   = "The Nomad region to deploy the job to."
  type          = "string"
  required      = false
  default_value = "global"
}

parameter {
  key           = "replicas"
  description   = "The replica count for the job."
  type          = "int"
  required      = false
  default_value = "1"
}

parameter {
  key         = "resources"
  description = "The amount of resources to allocate to the deployed allocation."
  type        = "category"
  required    = true
}

parameter {
  key           = "resources.cpu"
  description   = "Amount of CPU in MHz to allocate to this task"
  type          = "int"
  required      = false
  default_value = "100"
}

parameter {
  key           = "resources.memorymb"
  description   = "Amount of memory in MB to allocate to this task."
  type          = "int"
  required      = false
  default_value = "300"
}

parameter {
  key         = "service_port"
  description = "TCP port the job is listening on."
  type        = "uint"
  required    = false
}

parameter {
  key           = "service_provider"
  description   = "Specifies the service registration provider to use for registering a service for the job"
  type          = "string"
  required      = false
  default_value = "consul"
}

parameter {
  key         = "static_environment"
  description = "Environment variables to add to the job."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "vault_token"
  description = "The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.\nUses the runner config environment variable VAULT_TOKEN."
  type        = ""
  required    = true
}

