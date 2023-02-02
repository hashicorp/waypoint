# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "force_architecture"
  description   = "**Note**: This is a temporary field that enables overriding the `architecture` output attribute. Valid values are: `\"x86_64\"`, `\"arm64\"`"
  type          = "string"
  required      = false
  default_value = "`\"\"`"
}

parameter {
  key         = "region"
  description = "the AWS region the ECR repository is in\nif not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT."
  type        = "string"
  required    = false
}

parameter {
  key         = "repository"
  description = "the AWS ECR repository name"
  type        = "string"
  required    = true
}

parameter {
  key         = "tag"
  description = "the tag of the image to pull"
  type        = "string"
  required    = true
}

