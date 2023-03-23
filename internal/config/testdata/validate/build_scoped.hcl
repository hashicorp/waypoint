# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "foo" {
    build {
        workspace "foo" {
            use "docker" {}
        }

        label "bar" {}
    }

    deploy {}
}
