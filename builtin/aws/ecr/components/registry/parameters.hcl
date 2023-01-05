# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "region"
  description = "the AWS region the ECR repository is in\nif not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT"
  type        = "string"
  required    = false
}

parameter {
  key         = "repository"
  description = "the ECR repository to store the image into\nThis defaults to waypoint- then the application name. The repository will be automatically created if needed"
  type        = "string"
  required    = false
}

parameter {
  key         = "tag"
  description = "the docker tag to assign to the new image"
  type        = "string"
  required    = true
}

