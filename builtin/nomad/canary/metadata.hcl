# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Nomad Jobspec Canary"
  description = "The Nomad Jobspec Canary plugin promotes a Nomad canary deployment initiated by a Nomad jobspec deployment."
  identifier  = "waypoint/nomad-jobspec-canary"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "release-manager"
    name = "Nomad Jobspec Canary Release Manager"
    slug = "nomad-jobspec-canary-release-manager"
  }
}
