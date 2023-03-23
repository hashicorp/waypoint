# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            static = file("${path.project}/config_escape_data.hcl")
            more = templatestring(file("${path.project}/config_escape_data.hcl"), {
              pass = config.internal.pass,
            })
        }

        internal = {
          pass = configdynamic("tfc", {})
        }
    }
}
