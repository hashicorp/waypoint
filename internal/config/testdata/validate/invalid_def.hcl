# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}

variable "bees" {
  default = "buzz"
  type = bool
}