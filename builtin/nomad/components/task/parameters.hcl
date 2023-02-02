# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "datacenter"
  description   = "The Nomad datacenter to deploy the on-demand runner task to."
  type          = "string"
  required      = false
  default_value = "dc1"
}

parameter {
  key           = "namespace"
  description   = "The Nomad namespace to deploy the on-demand runner task to."
  type          = "string"
  required      = false
  default_value = "default"
}

parameter {
  key           = "nomad_host"
  description   = "Hostname of the Nomad server to use for launching on-demand tasks."
  type          = "string"
  required      = false
  default_value = "http://localhost:4646"
}

parameter {
  key           = "region"
  description   = "The Nomad region to deploy the on-demand runner task to."
  type          = "string"
  required      = false
  default_value = "global"
}

parameter {
  key           = "resources_cpu"
  description   = "Amount of CPU in MHz to allocate to this task. This can be overriden with the '-nomad-runner-cpu' flag on server install."
  type          = "int"
  required      = false
  default_value = "200"
}

parameter {
  key           = "resources_memory"
  description   = "Amount of memory in MB to allocate to this task. This can be overriden with the '-nomad-runner-memory' flag on server install."
  type          = "int"
  required      = false
  default_value = "2000"
}

