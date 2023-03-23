# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

app "api" {
  config {
    env = { "foo" = "bar" }

    workspace "dev" {
      env = { "bar" = "baz" }
    }
  }
}
