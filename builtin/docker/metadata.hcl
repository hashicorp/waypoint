integration {
  name        = "Docker"
  description = "TODO"
  identifier  = "waypoint/docker"
  components  = ["builder", "platform", "registry", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
