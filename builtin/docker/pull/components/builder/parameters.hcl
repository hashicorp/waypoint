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

parameter {
  key         = "auth"
  description = <<EOT
The authentication information to log into the docker repository.
EOT
  type        = "category"
  required    = false
}

parameter {
  key         = "auth.auth"
  description = ""
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.email"
  description = ""
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.hostname"
  description = <<EOT
Hostname of Docker registry.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.identityToken"
  description = <<EOT
Token used to authenticate user.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.password"
  description = <<EOT
Password of Docker registry account.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.registryToken"
  description = <<EOT
Bearer tokens to be sent to Docker registry.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.serverAddress"
  description = <<EOT
Address of Docker registry.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.username"
  description = <<EOT
Username of Docker registry account.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "disable_entrypoint"
  description = <<EOT
If set, the entrypoint binary won't be injected into the image.

The entrypoint binary is what provides extended functionality such as logs and exec. If it is not injected at build time the expectation is that the image already contains it.
EOT
  type        = "bool"
  required    = false
}

parameter {
  key         = "encoded_auth"
  description = <<EOT
The authentication information to log into the docker repository.

WARNING: be very careful to not leak the authentication information by hardcoding it here. Use a helper function like `file()` to read the information from a file not stored in VCS.
EOT
  type        = "string"
  required    = false
}

