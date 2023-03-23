# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "p"

config {
  internal = {
    value = "V"
  }

  env = {
    "direct"       = config.internal.value
    "interpolated" = "value: ${config.internal.value}"
  }
}
