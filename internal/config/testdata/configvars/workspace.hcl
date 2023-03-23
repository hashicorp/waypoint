# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

config {
  env = { "foo" = "bar" }

  workspace "dev" {
    env = { "bar" = "baz" }
  }
}
