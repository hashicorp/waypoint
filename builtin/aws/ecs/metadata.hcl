integration {
  name        = "AWS ECS"
  description = "The AWS ECS plugin deploys an application image to an AWS ECS cluster. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/aws-ecs"
  components  = ["platform", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
