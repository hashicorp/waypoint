# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "rate" {
  default = configdynamic("vault", {})
  type = number
}
