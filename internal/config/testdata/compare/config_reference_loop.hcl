# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            v1 = "${config.env.v2}"
            v2 = "${config.env.v3}"
            v3 = "${config.env.v1}"
        }
    }
}
