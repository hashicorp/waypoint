parameter {
  key         = "auth"
  description = <<EOT
The credentials for docker registry.

EOT
  type        = "nomad.AuthConfig"
  required    = true

}

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
  key         = "vault_token"
  description = <<EOT
The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.
Uses the runner config environment variable VAULT\_TOKEN.

EOT
  type        = "string" # WARNING: no type was documented. This will be a best effort choice.
  required    = true

}

parameter {
  key           = "datacenter"
  description   = <<EOT
The Nomad datacenter to deploy the job to.

EOT
  type          = "string"
  required      = false
  default_value = "dc1"
}

parameter {
  key         = "namespace"
  description = <<EOT
The Nomad namespace to deploy the job to.

EOT
  type        = "string"
  required    = false

}

parameter {
  key           = "region"
  description   = <<EOT
The Nomad region to deploy the job to.

EOT
  type          = "string"
  required      = false
  default_value = "global"
}

parameter {
  key           = "replicas"
  description   = <<EOT
The replica count for the job.

EOT
  type          = "int"
  required      = false
  default_value = "1"
}

parameter {
  key         = "resources"
  description = <<EOT
The amount of resources to allocate to the deployed allocation.

EOT
  type        = "category"
  required    = false

}

parameter {
  key           = "resources.cpu"
  description   = <<EOT
Amount of CPU in MHz to allocate to this task.

EOT
  type          = "int"
  required      = false
  default_value = "100"
}

parameter {
  key           = "resources.memorymb"
  description   = <<EOT
Amount of memory in MB to allocate to this task.

EOT
  type          = "int"
  required      = false
  default_value = "300"
}

parameter {
  key         = "service_port"
  description = <<EOT
TCP port the job is listening on.

EOT
  type        = "uint"
  required    = false

}

parameter {
  key           = "service_provider"
  description   = <<EOT
Specifies the service registration provider to use for registering a service for the job.

EOT
  type          = "string"
  required      = false
  default_value = "consul"
}

parameter {
  key         = "static_environment"
  description = <<EOT
Environment variables to add to the job.

EOT
  type        = "map of string to string"
  required    = false

}

