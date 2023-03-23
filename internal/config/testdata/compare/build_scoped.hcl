# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        use "A" {}

        workspace "production" {
          use "B" {}
        }

        label "waypoint/workspace == staging" {
          use "C" {}
        }
    }
}
