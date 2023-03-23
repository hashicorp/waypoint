# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "AWS Application Load Balancer"
  description = "The AWS ALB plugin releases applications deployed to AWS by attaching target groups to an ALB."
  identifier  = "waypoint/aws-alb"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "release-manager"
    name = "AWS ALB Release Manager"
    slug = "aws-alb-release-manager"
  }
}
