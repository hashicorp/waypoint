# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Nomad Jobspec"
  description = "The Nomad Jobspec plugin deploys to a Nomad cluster from a pre-existing Nomad job specification file."
  identifier  = "waypoint/nomad-jobspec"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Nomad Jobspec Platform"
    slug = "nomad-jobspec-platform"
  }
}
