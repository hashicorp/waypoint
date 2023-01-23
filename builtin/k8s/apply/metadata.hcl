integration {
  name        = "Kubernetes Apply"
  description = "The Kubernetes Apply plugin deploys Kubernetes resources directly from a single file or a directory of YAML or JSON files."
  identifier  = "waypoint/kubernetes-apply"
  components  = ["platform"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
