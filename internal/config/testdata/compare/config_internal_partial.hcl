# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            static = lower(config.internal.greeting, config.internal.extra)
        }

        internal = {
            greeting = configdynamic("tfc", {})
            extra = "FOO"
        }
    }
}
