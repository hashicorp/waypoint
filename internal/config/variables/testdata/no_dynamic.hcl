# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "art" {
  default = null
  type = string
}

variable "mug" {
  default = "clay"
  type = string
}

variable "is_good" {
  default = false
  type = bool
}

variable "whatdoesittaketobenumber" {
  default = 1
  type = number
}

variable "envs" {
  default = 1
  type = number
  env = ["foo", "bar"]
}
