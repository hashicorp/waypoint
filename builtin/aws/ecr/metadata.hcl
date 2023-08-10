# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "AWS ECR"
  description = "The AWS ECR plugin pushes a Docker image to an Elastic Container Registry on AWS."
  identifier  = "waypoint/aws-ecr"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "registry"
    name = "AWS ECR Registry"
    slug = "aws-ecr-registry"
  }
}
