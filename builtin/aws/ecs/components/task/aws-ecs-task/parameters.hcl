# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "cluster"
  description = "Cluster name to place On-Demand runner tasks in\nECS Cluster to place On-Demand runners in. This defaults to the cluster used by the Waypoint server"
  type        = "string"
  required    = false
}

parameter {
  key         = "execution_role_name"
  description = "The name of the AWS IAM role to apply to the task's Execution Role\nExecutionRoleName is the name of the AWS IAM role to apply to the task's Execution Role. At this time we reuse the same Role as the Waypoint server Execution Role."
  type        = "string"
  required    = false
}

parameter {
  key         = "log_group"
  description = "Cloud Watch Log Group to use for On-Demand Runners\nCloud Watch Log Group to use for On-Demand Runners. Defaults to the log group used for runners (waypoint-runner)."
  type        = "string"
  required    = false
}

parameter {
  key         = "odr_cpu"
  description = "CPU to use for the On-Demand runners.\nConfigure the CPU for the On-Demand runners. The default is 512. See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html for valid values"
  type        = "string"
  required    = false
}

parameter {
  key         = "odr_image"
  description = "Docker image for the Waypoint On-Demand Runners\nDocker image for the Waypoint On-Demand Runners. This will\ndefault to the server image with the name (not label) suffixed with '-odr'.\""
  type        = ""
  required    = true
}

parameter {
  key         = "odr_memory"
  description = "Memory to use for the On-Demand runners.\nConfigure the memory for the On-Demand runners. The default is 1024. See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html for valid values"
  type        = "string"
  required    = false
}

parameter {
  key         = "region"
  description = "AWS Region to use\nAWS region to use. Defaults to the region used for the Waypoint Server."
  type        = "string"
  required    = false
}

parameter {
  key         = "security_group_id"
  description = "Security Group ID to place the On-Demand Runner task in\nSecurity Group ID to place the On-Demand Runner task in. This defaults to the security group used for the Waypoint server"
  type        = "string"
  required    = true
}

parameter {
  key         = "subnets"
  description = "List of subnets to place the On-Demand Runner task in.\nList of subnets to place the On-Demand Runner task in. This defaults to the list of subnets configured for the Waypoint server and must be either identical or a subset of the subnets used by the Waypoint server"
  type        = "string"
  required    = true
}

parameter {
  key         = "task_role_name"
  description = "The name of the AWS IAM role to apply to the task's Task Role\nTaskRoleName is the name of the AWS IAM role to apply to the task. This role determines the privileges the ODR builder. If no role name is given, an IAM role will be created with the required policies"
  type        = "string"
  required    = false
}

