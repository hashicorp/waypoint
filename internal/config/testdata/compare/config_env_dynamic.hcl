# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "test" {
    config {
        env = {
            DATABASE_URL = configdynamic("vault", {
                path = "foo/"
            })
        }
    }
}
