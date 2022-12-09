parameter {
  key         = "certificate"
  description = <<EOT
ARN for the certificate to install on the ALB listener.

When this is set, the port automatically changes to 443 unless overriden in this configuration.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "domain_name"
  description = <<EOT
Fully qualified domain name to set for the ALB.

Set along with zone_id to have DNS automatically setup for the ALB. this value should include the full hostname and domain name, for instance app.example.com.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "domain_name"
  description = <<EOT
The ARN on an existing ALB to configure.

When this is set, no ALB or Listener is created. Instead the application is configured by manipulating this existing Listener. This allows users to configure their ALB outside waypoint but still have waypoint hook the application to that ALB.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "listener_arn"
  description = <<EOT
The ARN on an existing ALB to configure.

When this is set, no ALB or Listener is created. Instead the application is configured by manipulating this existing Listener. This allows users to configure their ALB outside waypoint but still have waypoint hook the application to that ALB.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "name"
  description = <<EOT
The name to assign the ALB.

Names have to be unique per region.
EOT
  type        = "string"
  required    = false
  default_value = "derived from application name"
}

parameter {
  key         = "port"
  description = <<EOT
The TCP port to configure the ALB to listen on.
EOT
  type        = "int"
  required    = false
  default_value = "80 for HTTP, 443 for HTTPS"
}

parameter {
  key         = "security_group_ids"
  description = <<EOT
The existing security groups to add to the ALB.

A set of existing security groups to add to the ALB.
EOT
  type        = "list of string"
  required    = false
}

parameter {
  key         = "subnets"
  description = <<EOT
The subnet ids to allow the ALB to run in.
EOT
  type        = "list of string"
  required    = false
  default_value = "public subnets in the account default VPC"
}

parameter {
  key         = "zone_id"
  description = <<EOT
Route53 ZoneID to create a DNS record into.

Set along with domain_name to have DNS automatically setup for the ALB.
EOT
  type        = "string"
  required    = false
}