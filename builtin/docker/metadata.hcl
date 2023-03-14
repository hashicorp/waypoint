integration {
  name        = "Docker"
  description = "The Docker plugin can build a Docker image of an application, push a Docker image to a remote registry, and/or deploy the Docker image to a Docker daemon. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/docker"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "Docker Builder"
    slug = "docker-builder"
  }
  component {
    type = "platform"
    name = "Docker Platform"
    slug = "docker-platform"
  }
  component {
    type = "registry"
    name = "Docker Registry"
    slug = "docker-registry"
  }
  component {
    type = "task"
    name = "Docker Task"
    slug = "docker-task"
  }
}
