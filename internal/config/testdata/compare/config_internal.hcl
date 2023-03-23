# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            static = "${config.internal.greeting}"
        }

        internal = {
            greeting = "hello"
        }
    }
}
