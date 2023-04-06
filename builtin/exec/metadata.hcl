# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "Exec"
  description = "The Exec plugin executes any command to perform a deploy. This enables the use of pre-existing deployment tools."
  identifier  = "waypoint/exec"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Exec Platform"
    slug = "exec-platform"
  }
}
