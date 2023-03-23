# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "AWS SSM"
  description = "The AWS SSM plugin reads configuration values from the AWS SSM Parameter Store."
  identifier  = "waypoint/aws-ssm"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "config-sourcer"
    name = "AWS SSM Config Sourcer"
    slug = "aws-ssm-config-sourcer"
  }
}
