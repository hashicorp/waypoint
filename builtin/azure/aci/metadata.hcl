# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Azure Container Instance"
  description = "The Azure ACI plugin deploys a container to Azure Container Instances."
  identifier  = "waypoint/azure-container-instance"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Azure Container Instance Platform"
    slug = "azure-container-instance-platform"
  }
}
