integration {
  name        = "Docker Pull"
  description = "The Docker Pull plugin pulls a Docker image from an existing Docker repository, and wraps the existing image entrypoint with the Waypoint entrypoint."
  identifier  = "waypoint/docker-pull"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "Docker Pull Builder"
    slug = "docker-pull-builder"
  }
}
