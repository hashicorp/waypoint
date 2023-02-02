integration {
  name        = "Null"
  description = "TODO"
  identifier  = "waypoint/null"
  components  = ["config-sourcer", "builder", "platform", "release-manager"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
