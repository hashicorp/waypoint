integration {
  name        = "Kubernetes"
  description = "The Kubernetes plugin can deploy a Docker image of an application to Kubernetes, expose the Deployment with a Kubernetes Service, and source configuration from a Kubernetes Secret or ConfigMap. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/kubernetes"
  components  = ["platform", "release-manager", "config-sourcer", "task"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
