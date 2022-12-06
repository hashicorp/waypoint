parameter {
  key         = "image"
  description = <<EOF
The image to push the local image to, fully qualified.

This value must be the fully qualified name to the image. for example: gcr.io/waypoint-demo/demo.
EOF
  type        = "string"
  required    = true
}

parameter {
  key         = "tag"
  description = <<EOF
The tag for the new image.

This is added to image to provide the full image reference.
EOF
  type        = "string"
  required    = true
}

parameter {
  key         = "auth"
  description = "The authentication information to log into the docker repository."
  type        = "category"
  required    = false
}

parameter {
  key      = "auth.auth"
  type     = "string"
  required = false
}

parameter {
  key      = "auth.email"
  type     = "string"
  required = false
}

parameter {
  key         = "auth.hostname"
  description = "Hostname of Docker registry."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.identityToken"
  description = "Token used to authenticate user."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.password"
  description = "Password of Docker registry account."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.registryToken"
  description = "Bearer tokens to be sent to Docker registry."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.serverAddress"
  description = "Address of Docker registry."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.username"
  description = "Username of Docker registry account."
  type        = "string"
  required    = false
}


parameter {
  key         = "encoded_auth"
  description = <<EOF
The authentication information to log into the docker repository.

The format of this is base64-encoded JSON. The structure is the [`AuthConfig`](https://pkg.go.dev/github.com/docker/cli/cli/config/types#AuthConfig) structure used by Docker. WARNING: be very careful to not leak the authentication information by hardcoding it here. Use a helper function like `file()` to read the information from a file not stored in VCS.
EOF
  type        = "string"
  required    = false
}

parameter {
  key         = "insecure"
  description = <<EOF
Access the registry via http rather than https.

This indicates that the registry should be accessed via http rather than https. Not recommended for production usage.
EOF
  type        = "bool"
  required    = false
}

parameter {
  key         = "local"
  description = "If set, the image will only be tagged locally and not pushed to a remote repository."
  type        = "bool"
  required    = false
}

parameter {
  key         = "password"
  description = <<EOF
Password associated with username on the registry.

This optional conflicts with encoded_auth and thusly only one can be used at a time. If both are used, encoded_auth takes precedence.
EOF
  type        = "string"
  required    = false
}


parameter {
  key         = "username"
  description = <<EOF
Username to authenticate with the registry.

This optional conflicts with encoded_auth and thusly only one can be used at a time. If both are used, encoded_auth takes precedence.
EOF
  type        = "string"
  required    = false
}
