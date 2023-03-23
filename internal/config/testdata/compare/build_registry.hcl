# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        registry {
            use "docker" {}
        }
    }
}
