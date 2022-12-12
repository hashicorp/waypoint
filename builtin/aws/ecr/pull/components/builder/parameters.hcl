parameter {
  key         = "repository"
  description = <<EOT
The AWS ECR repository name.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "tag"
  description = <<EOT
The tag of the image to pull.
EOT
  type        = "string"
  required    = true
}

parameter {
  key           = "force_architecture"
  description   = <<EOT
**Note**: This is a temporary field that enables overriding the `architecture` output attribute. Valid values are: `"x86_64"`, `"arm64"`.
EOT
  type          = "string"
  required      = false
  default_value = ""
}

parameter {
  key         = "region"
  description = <<EOT
The AWS region the ECR repository is in.

If not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT.
EOT
  type        = "string"
  required    = false
}

parameter {
  key           = "region"
  description   = <<EOT
The AWS region the ECR repository is in.

If not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT.
EOT
  type          = "string"
  required      = false
  default_value = ""
}
