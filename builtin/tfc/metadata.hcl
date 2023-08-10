# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Terraform Cloud"
  description = "The Terraform Cloud plugin reads Terraform state outputs from Terraform Cloud."
  identifier  = "waypoint/terraform-cloud"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "Terraform Cloud Config Sourcer"
    slug = "terraform-cloud-config-sourcer"
  }
}
