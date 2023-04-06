# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "CloudNative Buildpacks"
  description = "The Pack plugin creates a Docker image using CloudNative Buildpacks."
  identifier  = "waypoint/pack"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "CloudNative Buildpacks Builder"
    slug = "cloudnative-buildpacks-builder"
  }
}
