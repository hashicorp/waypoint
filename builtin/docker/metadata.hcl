integration {
  name        = "Docker"
  description = "The Docker plugin can build a Docker image of an application, push a Docker image to a remote registry, and/or deploy the Docker image to a Docker daemon. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/docker"
  components  = ["builder", "platform", "registry", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
