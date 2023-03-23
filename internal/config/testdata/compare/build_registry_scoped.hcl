# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        use "docker" {}

        registry {
          use "A" {}

          workspace "production" {
            use "B" {}
          }
        }
    }
}
