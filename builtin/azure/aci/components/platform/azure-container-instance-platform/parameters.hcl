# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "capacity"
  description = "the capacity details for the container"
  type        = "category"
  required    = true
}

parameter {
  key           = "capacity.cpu"
  description   = "number of CPUs to allocate the container, min 1, max based on resource availability of the region."
  type          = ""
  required      = true
  default_value = "1"
}

parameter {
  key         = "capacity.cpu_count"
  description = ""
  type        = "int"
  required    = true
}

parameter {
  key           = "capacity.memory"
  description   = "memory to allocate the container specified in MB, min 1024, max based on resource availability of the region."
  type          = "int"
  required      = true
  default_value = "1024"
}

parameter {
  key         = "location"
  description = "the resource location to deploy the container instance to"
  type        = "string"
  required    = false
}

parameter {
  key         = "managed_identity"
  description = "the managed identity assigned to the container group"
  type        = "string"
  required    = false
}

parameter {
  key         = "ports"
  description = "the ports the container is listening on, the first port in this list will be used by the entrypoint binary to direct traffic to your application"
  type        = "list of int"
  required    = false
}

parameter {
  key         = "resource_group"
  description = "the resource group to deploy the container to"
  type        = "string"
  required    = true
}

parameter {
  key         = "static_environment"
  description = "environment variables to control broad modes of the application\nenvironment variables that are meant to configure the application in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "subscription_id"
  description = "the Azure subscription id\nif not set uses the environment variable AZURE_SUBSCRIPTION_ID"
  type        = "string"
  required    = false
}

parameter {
  key         = "volume"
  description = "the volume details for a container"
  type        = "category"
  required    = true
}

parameter {
  key         = "volume.azure_file_share"
  description = "the details for the Azure file share volume"
  type        = "aci.AzureFileShareVolume"
  required    = true
}

parameter {
  key         = "volume.git_repo"
  description = "the details for GitHub repo to mount as a volume"
  type        = "aci.GitRepoVolume"
  required    = true
}

parameter {
  key         = "volume.name"
  description = "the name of the volume to mount into the container"
  type        = "string"
  required    = true
}

parameter {
  key         = "volume.path"
  description = "the path to mount the volume to in the container"
  type        = "string"
  required    = true
}

parameter {
  key         = "volume.read_only"
  description = "specify if the volume is read only"
  type        = "bool"
  required    = true
}

