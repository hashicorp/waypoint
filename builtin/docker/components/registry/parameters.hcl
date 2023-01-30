# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "auth"
  description = "the authentication information to log into the docker repository"
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
  description = "Hostname of Docker registry"
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.identityToken"
  description = "Token used to authenticate user"
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.password"
  description = "Password of Docker registry account"
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.registryToken"
  description = "Bearer tokens to be sent to Docker registry"
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.serverAddress"
  description = "Address of Docker registry"
  type        = "string"
  required    = false
}

parameter {
  key         = "auth.username"
  description = "Username of Docker registry account"
  type        = "string"
  required    = false
}

parameter {
  key         = "encoded_auth"
  description = "the authentication information to log into the docker repository\nThe format of this is base64-encoded JSON. The structure is the [`AuthConfig`](https://pkg.go.dev/github.com/docker/cli/cli/config/types#AuthConfig) structure used by Docker.\n  WARNING: be very careful to not leak the authentication information by hardcoding it here. Use a helper function like `file()` to read the information from a file not stored in VCS"
  type        = "string"
  required    = false
}

parameter {
  key         = "image"
  description = "the image to push the local image to, fully qualified\nthis value must be the fully qualified name to the image. for example: gcr.io/waypoint-demo/demo"
  type        = "string"
  required    = true
}

parameter {
  key         = "insecure"
  description = "access the registry via http rather than https\nThis indicates that the registry should be accessed via http rather than https. Not recommended for production usage."
  type        = "bool"
  required    = false
}

parameter {
  key         = "local"
  description = "if set, the image will only be tagged locally and not pushed to a remote repository"
  type        = "bool"
  required    = false
}

parameter {
  key         = "password"
  description = "password associated with username on the registry\nThis optional conflicts with encoded_auth and thusly only one can be used at a time. If both are used, encoded_auth takes precedence."
  type        = "string"
  required    = false
}

parameter {
  key         = "tag"
  description = "the tag for the new image\nthis is added to image to provide the full image reference"
  type        = "string"
  required    = true
}

parameter {
  key         = "username"
  description = "username to authenticate with the registry\nThis optional conflicts with encoded_auth and thusly only one can be used at a time. If both are used, encoded_auth takes precedence."
  type        = "string"
  required    = false
}

