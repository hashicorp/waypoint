# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            static = "${config.internal.greeting} ${config.internal.suffix}"
            extra = "extra: ${config.env.static}"
        }

        internal = {
            suffix = "ok?"
            greeting = configdynamic("foo", {})
        }
    }
}
