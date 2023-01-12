integration {
  name        = "Docker Ref"
  description = "The Docker Ref plugin refers to an existing Docker image, passing its image information - the image name and tag - to the Waypoint lifecycle."
  identifier  = "waypoint/docker-ref"
  components  = ["builder"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
