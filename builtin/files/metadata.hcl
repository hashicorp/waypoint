integration {
  name        = "Files"
  description = "The Files plugin generates a value representing a path on disk, and can copy them to a specific directory."
  identifier  = "waypoint/files"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "Files Builder"
    slug = "files-builder"
  }
  component {
    type = "registry"
    name = "Files Registry"
    slug = "files-registry"
  }
}
