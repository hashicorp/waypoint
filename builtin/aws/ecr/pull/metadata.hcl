integration {
  name        = "AWS ECR Pull"
  description = "TODO"
  identifier  = "waypoint/aws-ecr-pull"
  components  = ["builder"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
