# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "certificate"
  description = "ARN for the certificate to install on the ALB listener\nwhen this is set, the port automatically changes to 443 unless overriden in this configuration"
  type        = "string"
  required    = false
}

parameter {
  key         = "domain_name"
  description = "Fully qualified domain name to set for the ALB\nset along with zone_id to have DNS automatically setup for the ALB. this value should include the full hostname and domain name, for instance app.example.com"
  type        = "string"
  required    = false
}

parameter {
  key         = "listener_arn"
  description = "the ARN on an existing ALB to configure\nwhen this is set, no ALB or Listener is created. Instead the application is configured by manipulating this existing Listener. This allows users to configure their ALB outside waypoint but still have waypoint hook the application to that ALB"
  type        = "string"
  required    = false
}

parameter {
  key           = "name"
  description   = "the name to assign the ALB\nnames have to be unique per region"
  type          = "string"
  required      = false
  default_value = "derived from application name"
}

parameter {
  key           = "port"
  description   = "the TCP port to configure the ALB to listen on"
  type          = "int"
  required      = false
  default_value = "80 for HTTP, 443 for HTTPS"
}

parameter {
  key         = "security_group_ids"
  description = "the existing security groups to add to the ALB\na set of existing security groups to add to the ALB"
  type        = "list of string"
  required    = false
}

parameter {
  key           = "subnets"
  description   = "the subnet ids to allow the ALB to run in"
  type          = "list of string"
  required      = false
  default_value = "public subnets in the account default VPC"
}

parameter {
  key         = "zone_id"
  description = "Route53 ZoneID to create a DNS record into\nset along with domain_name to have DNS automatically setup for the ALB"
  type        = "string"
  required    = false
}

