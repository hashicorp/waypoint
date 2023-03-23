# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "image"
  description = "The image to pull.\nThis should NOT include the tag (the value following the ':' in a Docker image). Use `tag` to define the image tag."
  type        = "string"
  required    = true
}

parameter {
  key         = "tag"
  description = "The tag of the image to pull."
  type        = "string"
  required    = true
}

