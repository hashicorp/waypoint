# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
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
