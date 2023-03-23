# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

pipeline "foo" {
  step "bad" {
    pipeline "bad" {
    }
    use "bad" {
    }
  }
}

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}
