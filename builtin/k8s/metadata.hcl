# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Kubernetes"
  description = "The Kubernetes plugin can deploy a Docker image of an application to Kubernetes, expose the Deployment with a Kubernetes Service, and source configuration from a Kubernetes Secret or ConfigMap. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/kubernetes"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Kubernetes Platform"
    slug = "kubernetes-platform"
  }
  component {
    type = "release-manager"
    name = "Kubernetes Release Manager"
    slug = "kubernetes-release-manager"
  }
  component {
    type = "config-sourcer"
    name = "Kubernetes Config Sourcer"
    slug = "kubernetes-config-sourcer"
  }
  component {
    type = "task"
    name = "Kubernetes Task"
    slug = "kubernetes-task"
  }
}
