# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "Null"
  description = "The Null plugin is used for testing and experimentation with the different plugin components."
  identifier  = "waypoint/null"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "Null Config Sourcer"
    slug = "null-config-sourcer"
  }
  component {
    type = "builder"
    name = "Null Builder"
    slug = "null-builder"
  }
  component {
    type = "platform"
    name = "Null Platform"
    slug = "null-platform"
  }
  component {
    type = "release-manager"
    name = "Null Release Manager"
    slug = "null-release-manager"
  }
}
