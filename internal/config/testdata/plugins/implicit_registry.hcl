# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "hello"

app "tubes" {
    build {
        use "docker" {}

        registry {
            use "aws-ecr" {}
        }
    }

    deploy {
        use "nomad" {}
    }
}
