parameter {
  key         = "tag"
  description = <<EOT
The docker tag to assign to the new image.
EOT
  type        = "string"
  required    = true

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
  key         = "repository"
  description = <<EOT
The ECR repository to store the image into.

This defaults to waypoint- then the application name. The repository will be automatically created if needed.
EOT
  type        = "string"
  required    = false

}
