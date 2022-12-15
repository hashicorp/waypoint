parameter {
  key         = "capacity"
  description = <<EOT
The capacity details for the container.
EOT
  type        = "category"
  required    = true
}

parameter {
  key         = "capacity.cpu_count"
  description = <<EOT
Number of CPUs to allocate the container, min 1, max based on resource availability of the region.
EOT
  type        = "int"
  required    = false
}

parameter {
  key         = "capacity.memory"
  description = <<EOT
Memory to allocate the container specified in MB, min 1024, max based on resource availability of the region.
EOT
  type        = "int"
  required    = false
}

parameter {
  key         = "resource_group"
  description = <<EOT
The resource group to deploy the container to.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "volume"
  description = <<EOT
The volume details for a container.
EOT
  type        = "category"
  required    = true
}

parameter {
  key         = "volume.azure_file_share"
  description = <<EOT
The details for the Azure file share volume.
EOT
  type        = "aci.AzureFileShareVolume"
  required    = true
}

parameter {
  key         = "volume.git_repo"
  description = <<EOT
The details for GitHub repo to mount as a volume.
EOT
  type        = "aci.GitRepoVolume"
  required    = true
}

parameter {
  key         = "volume.name"
  description = <<EOT
The name of the volume to mount into the container.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "volume.path"
  description = <<EOT
The path to mount the volume to in the container.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "volume.read_only"
  description = <<EOT
Specify if the volume is read only.
EOT
  type        = "bool"
  required    = true
}

parameter {
  key         = "location"
  description = <<EOT
The resource location to deploy the container instance to.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "managed_identity"
  description = <<EOT
The managed identity assigned to the container group.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "ports"
  description = <<EOT
The ports the container is listening on, the first port in this list will be used by the entrypoint binary to direct traffic to your application.
EOT
  type        = "list of int"
  required    = false
}

parameter {
  key         = "static_environment"
  description = <<EOT
Environment variables to control broad modes of the application.

Environment variables that are meant to configure the application in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.
EOT
  type        = "map of stirng to string"
  required    = false
}

parameter {
  key         = "subscription_id"
  description = <<EOT
The Azure subscription id.

If not set uses the environment variable AZURE_SUBSCRIPTION_ID.
EOT
  type        = "string"
  required    = false
}

