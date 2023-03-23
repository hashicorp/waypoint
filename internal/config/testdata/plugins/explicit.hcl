# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "hello"

plugin "go1" {
    type {
        mapper = true
    }
}

plugin "go2" {
    type {
        registry = true
    }
}
