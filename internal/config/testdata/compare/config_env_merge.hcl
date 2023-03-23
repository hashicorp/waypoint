# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

config {
  env = {
    parent = "1"
  }
}

app "test" {
    config {
        env = {
            child = "2"
        }
    }
}
