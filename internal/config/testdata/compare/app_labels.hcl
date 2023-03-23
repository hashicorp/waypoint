# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "bar" {
    path = "./bar"

    labels = {
        "pwd": path.pwd,
        "project": path.project,
        "app": path.app,
    }
}
