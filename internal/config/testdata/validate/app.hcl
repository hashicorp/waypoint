# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

app "foo" {
    build {
      use "docker" {}
    }

    deploy {
      use "docker" {}
    }
}

app "relative_above_root" {
    path = "../nope"

    build {
      use "docker" {}
    }

    deploy {
      use "docker" {}
    }
}

app "system_label" {
    labels = {
        "waypoint/foo" = "bar"
    }

    build {
      use "docker" {}
    }

    deploy {
      use "docker" {}
    }
}
