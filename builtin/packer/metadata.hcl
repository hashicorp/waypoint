# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Packer"
  description = "The Packer plugin retrieves the image ID of an image whose metadata is pushed to an HCP Packer registry. The image ID is that of the HCP Packer bucket iteration assigned to the configured channel, with a matching cloud provider and region."
  identifier  = "waypoint/packer"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "Packer Config Sourcer"
    slug = "packer-config-sourcer"
  }
}
