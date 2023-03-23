# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "Consul"
  description = "The Consul plugin reads configuration values from the Consul KV store."
  identifier  = "waypoint/consul"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "Consul Config Sourcer"
    slug = "consul-config-sourcer"
  }
}
