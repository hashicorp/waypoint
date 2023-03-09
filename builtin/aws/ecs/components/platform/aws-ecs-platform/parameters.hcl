# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "alb"
  description = "Provides additional configuration for using an ALB with ECS"
  type        = "category"
  required    = true
}

parameter {
  key         = "alb.certificate"
  description = "the ARN of an AWS Certificate Manager cert to associate with the ALB"
  type        = "string"
  required    = false
}

parameter {
  key         = "alb.domain_name"
  description = "Fully qualified domain name to set for the ALB\nset along with zone_id to have DNS automatically setup for the ALB. this value should include the full hostname and domain name, for instance app.example.com"
  type        = "string"
  required    = false
}

parameter {
  key         = "alb.ingress_port"
  description = "Internet-facing traffic port. Defaults to 80 if 'certificate' is unset, 443 if set.\nused to set the ALB listener port, and the ALB security group ingress port"
  type        = "int64"
  required    = false
}

parameter {
  key         = "alb.internal"
  description = "Whether or not the created ALB should be internal\nused when listener_arn is not set. If set, the created ALB will have a scheme of `internal`, otherwise by default it has a scheme of `internet-facing`."
  type        = "bool"
  required    = false
}

parameter {
  key         = "alb.load_balancer_arn"
  description = "the ARN on an existing ALB to configure\nwhen this is set, Waypoint will use this ALB instead of creating its own. A target group will still be created for each deployment, and will be added to a listener on the configured ALB port (Waypoint will the listener if it doesn't exist). This allows users to configure their ALB outside Waypoint but still have Waypoint hook the application to that ALB"
  type        = "string"
  required    = true
}

parameter {
  key         = "alb.security_group_ids"
  description = ""
  type        = "list of string"
  required    = false
}

parameter {
  key           = "alb.subnets"
  description   = "the VPC subnets to use for the ALB"
  type          = "list of string"
  required      = false
  default_value = "public subnets in the default VPC"
}

parameter {
  key         = "alb.zone_id"
  description = "Route53 ZoneID to create a DNS record into\nset along with alb.domain_name to have DNS automatically setup for the ALB"
  type        = "string"
  required    = false
}

parameter {
  key         = "architecture"
  description = "the instruction set CPU architecture that the Amazon ECS supports. Valid values are: \"x86_64\", \"arm64\""
  type        = "string"
  required    = false
}

parameter {
  key           = "assign_public_ip"
  description   = "assign a public ip address to tasks. Defaults to true. Ignored if using an ec2 cluster.\nIf this is set to false, deployments will fail unless tasks are able to egress to the container registry by some other means (i.e. a subnet default route to a NAT gateway)."
  type          = "bool"
  required      = false
  default_value = "true"
}

parameter {
  key         = "cluster"
  description = "the name of the ECS cluster to deploy into\nthe ECS cluster that will run the application as a Service. if there is no ECS cluster with this name, the ECS cluster will be created and configured to use Fargate to run containers."
  type        = "string"
  required    = false
}

parameter {
  key         = "count"
  description = "how many instances of the application should run"
  type        = "int"
  required    = false
}

parameter {
  key         = "cpu"
  description = "how many cpu shares the container running the application is allowed\non Fargate, possible values for this are configured by the amount of memory the container is using. Here is a complete listing of possible values: 512MB: 256\n1024MB: 256, 512\n2048MB: 256, 512, 1024\n3072MB: 512, 1024\n4096MB: 512, 1024\n5120MB: 1024\n6144MB: 1024\n7168MB: 1024\n8192MB: 1024"
  type        = "int"
  required    = false
}

parameter {
  key         = "disable_alb"
  description = "do not create a load balancer assigned to the service"
  type        = "bool"
  required    = false
}

parameter {
  key         = "ec2_cluster"
  description = "indicate if the ECS cluster should be EC2 type rather than Fargate\nthis controls if we should verify the ECS cluster in EC2 type. The cluster will not be created if it doesn't exist, only that there as existing cluster this is using EC2 and not Fargate"
  type        = "bool"
  required    = false
}

parameter {
  key           = "execution_role_name"
  description   = "the name of the IAM role to use for ECS execution"
  type          = "string"
  required      = false
  default_value = "create a new exeuction IAM role based on the application name"
}

parameter {
  key         = "health_check"
  description = "Health check settings for the app."
  type        = "category"
  required    = false
}

parameter {
  key         = "health_check.grpc_code"
  description = ""
  type        = "string"
  required    = false
}

parameter {
  key           = "health_check.healthy_threshold_count"
  description   = "The number of consecutive successful health checks required toconsider a target healthy."
  type          = "int64"
  required      = false
  default_value = "5"
}

parameter {
  key         = "health_check.http_code"
  description = ""
  type        = "string"
  required    = false
}

parameter {
  key         = "health_check.interval"
  description = "The amount of time, in seconds, between health checks."
  type        = "int64"
  required    = false
}

parameter {
  key         = "health_check.matcher"
  description = "The range of HTTP codes to use when checking for a successful response fromthe target."
  type        = ""
  required    = true
}

parameter {
  key         = "health_check.path"
  description = "The destination of the ping path for the target health check."
  type        = "string"
  required    = false
}

parameter {
  key           = "health_check.protocol"
  description   = "The protocol for the health check to use."
  type          = "string"
  required      = false
  default_value = "HTTP"
}

parameter {
  key         = "health_check.timeout"
  description = "The amount of time, in seconds, for which no target response means a failure."
  type        = "int64"
  required    = false
}

parameter {
  key           = "health_check.unhealthy_threshold_count"
  description   = "The number of consecutive failed health checks required to consider a target unhealthy."
  type          = "int64"
  required      = false
  default_value = "2"
}

parameter {
  key           = "log_group"
  description   = "the CloudWatchLogs log group to store container logs into"
  type          = "string"
  required      = false
  default_value = "derived from the application name"
}

parameter {
  key         = "logging"
  description = "Provides additional configuration for logging flags for ECS\nPart of the ecs task definition.  These configuration flags help control how the awslogs log driver is configured."
  type        = "category"
  required    = false
}

parameter {
  key         = "logging.create_group"
  description = "Enables creation of the aws logs group if not present"
  type        = "bool"
  required    = false
}

parameter {
  key         = "logging.datetime_format"
  description = "Defines the multiline start pattern in Python strftime format"
  type        = "string"
  required    = false
}

parameter {
  key         = "logging.max_buffer_size"
  description = "When using non-blocking logging mode, this is the buffer size for message storage"
  type        = "string"
  required    = false
}

parameter {
  key         = "logging.mode"
  description = "Delivery method for log messages, either 'blocking' or 'non-blocking'"
  type        = "string"
  required    = false
}

parameter {
  key         = "logging.multiline_pattern"
  description = "Defines the multiline start pattern using a regular expression"
  type        = "string"
  required    = false
}

parameter {
  key           = "logging.region"
  description   = "The region the logs are to be shipped to"
  type          = ""
  required      = true
  default_value = "The same region the task is to be running"
}

parameter {
  key           = "logging.stream_prefix"
  description   = "Prefix for application in cloudwatch logs path"
  type          = "string"
  required      = false
  default_value = "Generated based off timestamp"
}

parameter {
  key         = "memory"
  description = "how much memory to assign to the container running the application\nwhen running in Fargate, this must be one of a few values, specified in MB: 512, 1024, 2048, 3072, 4096, 5120, and up to 16384 in increments of 1024. The memory value also controls the possible values for cpu"
  type        = "int"
  required    = true
}

parameter {
  key         = "memory_reservation"
  description = ""
  type        = "int"
  required    = false
}

parameter {
  key         = "region"
  description = "the AWS region for the ECS cluster"
  type        = "string"
  required    = true
}

parameter {
  key         = "secrets"
  description = "secret key/values to pass to the ECS container"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "security_group_ids"
  description = "Security Group IDs of existing security groups to use for the ECS service's network access\nlist of existing group IDs to use for the ECS service's network access. If none are specified, waypoint will create one. If DisableALB is false (the default), waypoint will only allow ingress from the ALB's security group"
  type        = "list of string"
  required    = false
}

parameter {
  key           = "service_port"
  description   = "the TCP port that the application is listening on"
  type          = "int64"
  required      = false
  default_value = "3000"
}

parameter {
  key         = "sidecar"
  description = "Additional container to run as a sidecar.\nThis runs additional containers in addition to the main container that comes from the build phase."
  type        = "category"
  required    = true
}

parameter {
  key         = "sidecar.container_port"
  description = "The port number for the container"
  type        = "int"
  required    = false
}

parameter {
  key         = "sidecar.health_check"
  description = ""
  type        = "ecs.HealthCheckConfig"
  required    = true
}

parameter {
  key         = "sidecar.host_port"
  description = "The port number on the host to reserve for the container"
  type        = "int"
  required    = false
}

parameter {
  key         = "sidecar.image"
  description = "Image of the sidecar container"
  type        = "string"
  required    = true
}

parameter {
  key         = "sidecar.memory"
  description = "The amount (in MiB) of memory to present to the container"
  type        = "int"
  required    = false
}

parameter {
  key         = "sidecar.memory_reservation"
  description = "The soft limit (in MiB) of memory to reserve for the container"
  type        = "int"
  required    = false
}

parameter {
  key         = "sidecar.name"
  description = "Name of the container"
  type        = "string"
  required    = true
}

parameter {
  key         = "sidecar.protocol"
  description = "The protocol used for port mapping."
  type        = "string"
  required    = false
}

parameter {
  key         = "sidecar.secrets"
  description = "Secrets to expose to this container"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "sidecar.static_environment"
  description = "Environment variables to expose to this container"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "static_environment"
  description = "static environment variables to make available"
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "subnets"
  description   = "the VPC subnets to use for the service\nyou may set a list of private subnets here to prevent your tasks from being directly exposed publicly"
  type          = "list of string"
  required      = false
  default_value = "public subnets in the default VPC"
}

parameter {
  key         = "task_role_name"
  description = "the name of the task IAM role to assign.\nIf no role exists and a one or more task role policies are requested, a role with this name will be created."
  type        = "string"
  required    = false
}

parameter {
  key         = "task_role_policy_arns"
  description = "IAM Policy arns for attaching to the task role.\nIf no task role name is specified a task role with a default name will be created for this app, and these policies will be attached."
  type        = "list of string"
  required    = false
}

