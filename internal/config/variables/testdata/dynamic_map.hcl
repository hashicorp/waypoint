# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "teeth" {
  default = configdynamic("static", {
    json = <<-EOF
      {"k1":"v1", "k2":"v2"}
EOF
  })
  type = map(string)
}
