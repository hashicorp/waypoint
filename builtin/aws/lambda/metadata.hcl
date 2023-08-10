# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

integration {
  name        = "AWS Lambda"
  description = "The AWS Lambda plugin deploys OCI images as functions to AWS Lambda."
  identifier  = "waypoint/aws-lambda"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "AWS Lambda Platform"
    slug = "aws-lambda-platform"
  }
}
