integration {
  name        = "Null"
  description = "The Null plugin is used for testing and experimentation with the different plugin components."
  identifier  = "waypoint/null"
  components  = ["config-sourcer", "builder", "platform", "release-manager"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
