parameter {
  key         = "image"
  description = <<EOT
The image to pull.

This should NOT include the tag (the value following the ':' in a Docker image). Use `tag` to define the image tag.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "tag"
  description = <<EOT
The tag of the image to pull.
EOT
  type        = "string"
  required    = true
}

