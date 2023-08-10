# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Nomad"
  description = "The Nomad plugin deploys a Docker container to a Nomad cluster. It also launches on-demand runners to do operations remotely."
  identifier  = "waypoint/nomad"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Nomad Platform"
    slug = "nomad-platform"
  }
  component {
    type = "task"
    name = "Nomad Task"
    slug = "nomad-task"
  }
}
