# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "teeth" {
  default = configdynamic("static", {
    value = "hello"
  })
  type = string
}
