integration {
  name        = "CloudNative Buildpacks"
  description = "The Pack plugin creates a Docker image using CloudNative Buildpacks."
  identifier  = "waypoint/pack"
  components  = ["builder"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
