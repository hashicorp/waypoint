# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "architecture"
  description   = "The instruction set architecture that the function supports. Valid values are: \"x86_64\", \"arm64\""
  type          = "string"
  required      = false
  default_value = "x86_64"
}

parameter {
  key           = "iam_role"
  description   = "an IAM Role specified by ARN that will be used by the Lambda at execution time"
  type          = "string"
  required      = false
  default_value = "created automatically"
}

parameter {
  key           = "memory"
  description   = "the amount of memory, in megabytes, to assign the function"
  type          = "int"
  required      = false
  default_value = "256"
}

parameter {
  key         = "region"
  description = "the AWS region for the ECS cluster"
  type        = "string"
  required    = true
}

parameter {
  key         = "static_environment"
  description = "environment variables to expose to the lambda function\nenvironment variables that are meant to configure the application in a static way. This might be to control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands."
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "storagemb"
  description   = "The storage size (in MB) of the Lambda function's `/tmp` directory. Must be a value between 512 and 10240."
  type          = "int"
  required      = false
  default_value = "512"
}

parameter {
  key           = "timeout"
  description   = "the number of seconds a function has to return a result"
  type          = "int"
  required      = false
  default_value = "60"
}

