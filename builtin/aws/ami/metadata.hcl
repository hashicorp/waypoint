# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "AWS AMI"
  description = "The AWS AMI plugin searches for and returns an existing AMI, to be deployed as an EC2."
  identifier  = "waypoint/aws-ami"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "AWS AMI Builder"
    slug = "aws-ami-builder"
  }
}
