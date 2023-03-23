# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "AWS ECS"
  description = "The AWS ECS plugin deploys an application image to an AWS ECS cluster. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/aws-ecs"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "AWS ECS Platform"
    slug = "aws-ecs-platform"
  }
  component {
    type = "task"
    name = "AWS ECS Task"
    slug = "aws-ecs-task"
  }
}
