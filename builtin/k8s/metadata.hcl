integration {
  name        = "Kubernetes"
  description = "TODO"
  identifier  = "waypoint/kubernetes"
  components  = ["platform", "release-manager", "config-sourcer", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
