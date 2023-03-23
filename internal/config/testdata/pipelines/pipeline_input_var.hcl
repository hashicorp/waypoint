# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

pipeline "foo" {
  step "test" {
    image_url = var.image_url

    use "test" {
      foo = "bar"
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

variable "image_url" {
  default = "example.com/test"
  type    = string
}
