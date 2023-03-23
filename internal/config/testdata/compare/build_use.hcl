# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    build {
        use "test" {
            foo = path.app
        }
    }
}
