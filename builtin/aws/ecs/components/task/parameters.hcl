parameter {
  key         = "odr_image"
  description = <<EOT
Docker image for the Waypoint On-Demand Runners.

Docker image for the Waypoint On-Demand Runners. This will
default to the server image with the name (not label) suffixed with '-odr'.".
EOT
  type        = "string"
  required    = true

}
parameter {
  key         = "security_group_id"
  description = <<EOT
Security Group ID to place the On-Demand Runner task in.

Security Group ID to place the On-Demand Runner task in. This defaults to the security group used for the Waypoint server.
EOT
  type        = "string"
  required    = true

}
parameter {
  key         = "subnets"
  description = <<EOT
List of subnets to place the On-Demand Runner task in.

List of subnets to place the On-Demand Runner task in. This defaults to the list of subnets configured for the Waypoint server and must be either identical or a subset of the subnets used by the Waypoint server.
EOT
  type        = "list of string"
  required    = true

}
parameter {
  key         = "cluster"
  description = <<EOT
Cluster name to place On-Demand runner tasks in.

ECS Cluster to place On-Demand runners in. This defaults to the cluster used by the Waypoint server.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "execution_role_name"
  description = <<EOT
The name of the AWS IAM role to apply to the task's Execution Role.

ExecutionRoleName is the name of the AWS IAM role to apply to the task's Execution Role. At this time we reuse the same Role as the Waypoint server Execution Role.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "log_group"
  description = <<EOT
Cloud Watch Log Group to use for On-Demand Runners.

Cloud Watch Log Group to use for On-Demand Runners. Defaults to the log group used for runners (waypoint-runner).
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "odr_cpu"
  description = <<EOT
CPU to use for the On-Demand runners.

Configure the CPU for the On-Demand runners. The default is 512. See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html for valid values.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "odr_memory"
  description = <<EOT
Memory to use for the On-Demand runners.

Configure the memory for the On-Demand runners. The default is 1024. See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html for valid values.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "region"
  description = <<EOT
AWS Region to use.

AWS region to use. Defaults to the region used for the Waypoint Server.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "task_role_name"
  description = <<EOT
The name of the AWS IAM role to apply to the task's Task Role.

TaskRoleName is the name of the AWS IAM role to apply to the task. This role determines the privileges the ODR builder. If no role name is given, an IAM role will be created with the required policies.
EOT
  type        = "string"
  required    = false

}
