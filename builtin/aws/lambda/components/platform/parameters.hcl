parameter {
  key         = "region"
  description = <<EOT
The AWS region for the ECS cluster.
EOT
  type        = "string"
  required    = true
}

parameter {
  key           = "architecture"
  description   = <<EOT
The instruction set architecture that the function supports. Valid values are: "x86_64", "arm64".
EOT
  type          = "string"
  required      = false
  default_value = "x86_64"
}

parameter {
  key           = "iam_role"
  description   = <<EOT
An IAM Role specified by ARN that will be used by the Lambda at execution time.
EOT
  type          = "string"
  required      = false
  default_value = "cerated automatically"
}

parameter {
  key           = "memory"
  description   = <<EOT
The amount of memory, in megabytes, to assign the function.
EOT
  type          = "int"
  required      = false
  default_value = "256"
}

parameter {
  key         = "static_environment"
  description = <<EOT
Environment variables to expose to the lambda function.

Environment variables that are meant to configure the application in a static way. This might be to control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.
EOT
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "storagemb"
  description   = <<EOT
The storage size (in MB) of the Lambda function's `/tmp` directory. Must be a value between 512 and 10240.
EOT
  type          = "int"
  required      = false
  default_value = "512"
}

parameter {
  key           = "timeout"
  description   = <<EOT
The number of seconds a function has to return a result.
EOT
  type          = "int"
  required      = false
  default_value = "60"
}

