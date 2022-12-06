integration {
  name        = "Docker"
  description = "Use Waypoint on a Docker instance."
  identifier  = "waypoint/docker"
  components  = ["builder", "platform", "registry", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
