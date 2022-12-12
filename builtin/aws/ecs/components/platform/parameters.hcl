parameter {
  key         = "logging"
  description = <<EOT
Provides additional configuration for logging flags for ECS.

Part of the ecs task definition. These configuration flags help control how the awslogs log driver is configured.
EOT
  type        = "category"
  required    = true

}
parameter {
  key         = "logging.create_group"
  description = <<EOT
Enables creation of the aws logs group if not present.
EOT
  type        = "bool"
  required    = false

}
parameter {
  key         = "logging.datetime_format"
  description = <<EOT
Defines the multiline start pattern in Python strftime format.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "logging.max_buffer_size"
  description = <<EOT
When using non-blocking logging mode, this is the buffer size for message storage.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "logging.mode"
  description = <<EOT
Delivery method for log messages, either 'blocking' or 'non-blocking'.
EOT
  type        = "string"
  required    = false
}
parameter {
  key         = "logging.multiline_pattern"
  description = <<EOT
Defines the multiline start pattern using a regular expression.
EOT
  type        = "string"
  required    = false
}
parameter {
  key         = "logging.region"
  description = <<EOT
The region the logs are to be shipped to.
EOT
  type        = "string"
  required    = false

}
parameter {
  key           = "logging.stream_prefix"
  description   = <<EOT
Prefix for application in cloudwatch logs path.
EOT
  type          = "string"
  required      = false
  default_value = "Generated based off timestamp"

}
parameter {
  key         = "memory"
  description = <<EOT
How much memory to assign to the container running the application.

When running in Fargate, this must be one of a few values, specified in MB: 512, 1024, 2048, 3072, 4096, 5120, and up to 16384 in increments of 1024. The memory value also controls the possible values for cpu.
EOT
  type        = "int"
  required    = true

}
parameter {
  key         = "region"
  description = <<EOT
The AWS region for the ECS cluster.
EOT
  type        = "string"
  required    = true

}
parameter {
  key         = "sidecar"
  description = <<EOT
Additional container to run as a sidecar.

This runs additional containers in addition to the main container that comes from the build phase.
EOT
  type        = "category"
  required    = true
}
parameter {
  key         = "sidecar.container_port"
  description = <<EOT
The port number for the container.
EOT
  type        = "int"
  required    = false

}
parameter {
  key         = "sidecar.health_check"
  description = <<EOT
The port number on the host to reserve for the container.
EOT
  type        = "ecs.HealthCheckConfig"
  required    = false

}
parameter {
  key         = "sidecar.host_port"
  description = <<EOT
The port number on the host to reserve for the container.
EOT
  type        = "int"
  required    = false

}
parameter {
  key         = "sidecar.image"
  description = <<EOT
Image of the sidecar container.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "sidecar.memory"
  description = <<EOT
The amount (in MiB) of memory to present to the container.
EOT
  type        = "int"
  required    = false

}
parameter {
  key         = "sidecar.memory_reservation"
  description = <<EOT
The soft limit (in MiB) of memory to reserve for the container.
EOT
  type        = "int"
  required    = false

}
parameter {
  key         = "sidecar.name"
  description = <<EOT
Name of the container.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "sidecar.protocol"
  description = <<EOT
The protocol used for port mapping.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "sidecar.secrets"
  description = <<EOT
Secrets to expose to this container.
EOT
  type        = "map of string to string"
  required    = false

}
parameter {
  key         = "sidecar.static_environment"
  description = <<EOT
Environment variables to expose to this container.
EOT
  type        = "map of string to string"
  required    = false

}
parameter {
  key         = "alb"
  description = <<EOT
Provides additional configuration for using an ALB with ECS.
EOT
  type        = "category"
  required    = false

}
parameter {
  key         = "alb.certificate"
  description = <<EOT
The ARN of an AWS Certificate Manager cert to associate with the ALB.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "alb.domain_name"
  description = <<EOT
Fully qualified domain name to set for the ALB.

Set along with zone_id to have DNS automatically setup for the ALB. this value should include the full hostname and domain name, for instance app.example.com.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "alb.ingress_port"
  description = <<EOT
Internet-facing traffic port. Defaults to 80 if 'certificate' is unset, 443 if set.

Used to set the ALB listener port, and the ALB security group ingress port.
EOT
  type        = "int64"
  required    = false

}
parameter {
  key         = "alb.internal"
  description = <<EOT
Whether or not the created ALB should be internal.

Used when listener_arn is not set. If set, the created ALB will have a scheme of `internal`, otherwise by default it has a scheme of `internet-facing`.
EOT
  type        = "bool"
  required    = false

}
parameter {
  key         = "alb.listener_arn"
  description = <<EOT
The ARN on an existing ALB to configure.

When this is set, no ALB or Listener is created. Instead the application is configured by manipulating this existing Listener. This allows users to configure their ALB outside waypoint but still have waypoint hook the application to that ALB.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "alb.security_group_ids"
  description = ""
  type        = "list of string"
  required    = false
}
parameter {
  key           = "alb.subnets"
  description   = <<EOT
The VPC subnets to use for the ALB.
EOT
  type          = "list of string"
  required      = false
  default_value = "public subnets in the default VPC"
}
parameter {
  key         = "alb.zone_id"
  description = <<EOT
Route53 ZoneID to create a DNS record into.

Set along with alb.domain_name to have DNS automatically setup for the ALB.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "architecture"
  description = <<EOT
The instruction set CPU architecture that the Amazon ECS supports. Valid values are: "x86_64", "arm64".
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "assign_public_ip"
  description = <<EOT
Assign a public ip address to tasks. Defaults to true. Ignored if using an ec2 cluster.

If this is set to false, deployments will fail unless tasks are able to egress to the container registry by some other means (i.e. a subnet default route to a NAT gateway).
EOT
  type        = "bool"
  required    = false

}
parameter {
  key         = "cluster"
  description = <<EOT
The name of the ECS cluster to deploy into.

The ECS cluster that will run the application as a Service. if there is no ECS cluster with this name, the ECS cluster will be created and configured to use Fargate to run containers.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "count"
  description = <<EOT
How many instances of the application should run.
EOT
  type        = "int"
  required    = false

}
parameter {
  key         = "cpu"
  description = <<EOT
How many cpu shares the container running the application is allowed.

On Fargate, possible values for this are configured by the amount of memory the container is using. Here is a complete listing of possible values:
512MB: 256
1024MB: 256, 512
2048MB: 256, 512, 1024
3072MB: 512, 1024
4096MB: 512, 1024
5120MB: 1024
6144MB: 1024
7168MB: 1024
8192MB: 1024.
EOT
  type        = "int"
  required    = false
}
parameter {
  key         = "disable_alb"
  description = <<EOT
Do not create a load balancer assigned to the service.
EOT
  type        = "bool"
  required    = false
}
parameter {
  key         = "ec2_cluster"
  description = <<EOT
Indicate if the ECS cluster should be EC2 type rather than Fargate.

This controls if we should verify the ECS cluster in EC2 type. The cluster will not be created if it doesn't exist, only that there as existing cluster this is using EC2 and not Fargate.
EOT
  type        = "bool"
  required    = false

}
parameter {
  key           = "execution_role_name"
  description   = <<EOT
The name of the IAM role to use for ECS execution.
EOT
  type          = "string"
  required      = false
  default_value = "create a new execution IAM role based on the application name"

}
parameter {
  key           = "log_group"
  description   = <<EOT
The CloudWatchLogs log group to store container logs into.
EOT
  type          = "string"
  required      = false
  default_value = "derived from the application name"

}
parameter {
  key         = "memory_reservation"
  description = ""
  type        = "int"
  required    = false

}
parameter {
  key         = "secrets"
  description = <<EOT
Secret key/values to pass to the ECS container.
EOT
  type        = "map of string to string"
  required    = false

}
parameter {
  key         = "security_group_ids"
  description = <<EOT
Security Group IDs of existing security groups to use for the ECS service's network access.

List of existing group IDs to use for the ECS service's network access. If none are specified, waypoint will create one. If DisableALB is false (the default), waypoint will only allow ingress from the ALB's security group.
EOT
  type        = "list of string"
  required    = false

}
parameter {
  key           = "service_port"
  description   = <<EOT
The TCP port that the application is listening on.
EOT
  type          = "int64"
  required      = false
  default_value = "3000"
}
parameter {
  key         = "static_environment"
  description = <<EOT
Static environment variables to make available.
EOT
  type        = "map of string to string"
  required    = false

}
parameter {
  key           = "subnets"
  description   = <<EOT
The VPC subnets to use for the service.

You may set a list of private subnets here to prevent your tasks from being directly exposed publicly.
EOT
  type          = "list of string"
  required      = false
  default_value = "public subnets in the default VPC"
}
parameter {
  key         = "task_role_name"
  description = <<EOT
The name of the task IAM role to assign.

If no role exists and a one or more task role policies are requested, a role with this name will be created.
EOT
  type        = "string"
  required    = false

}
parameter {
  key         = "task_role_policy_arns"
  description = <<EOT

IAM Policy arns for attaching to the task role.

If no task role name is specified a task role with a default name will be created for this app, and these policies will be attached.
EOT
  type        = "list of string"
  required    = false

}
