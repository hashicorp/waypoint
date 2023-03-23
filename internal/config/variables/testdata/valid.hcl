# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "art" {
  default = null
  type = string
}

variable "is_good" {
  default = false
  type = bool
}

variable "whatdoesittaketobenumber" {
  default = 1
  type = number
  sensitive = true
}

variable "envs" {
  default = 1
  type = number
  env = ["foo", "bar"]
}

variable "dynamic" {
  type = string
  default = configdynamic("vault", {})
}
