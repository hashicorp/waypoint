# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "Vault"
  description = "The Vault plugin reads configuration values from Vault."
  identifier  = "waypoint/vault"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "Vault Config Sourcer"
    slug = "vault-config-sourcer"
  }
}
