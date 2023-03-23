# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        internal = {
            greeting = "hello"
        }

        file = {
          "blah.yml" = templatestring(file("${path.project}/blah.yml"), {
            greeting = config.internal.greeting,
          })

          "foo.yml" = "foo: ${config.internal.greeting}"
        }
    }
}
