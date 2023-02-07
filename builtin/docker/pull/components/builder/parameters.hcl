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
  key         = "disable_entrypoint"
  description = "if set, the entrypoint binary won't be injected into the image\nThe entrypoint binary is what provides extended functionality such as logs and exec. If it is not injected at build time the expectation is that the image already contains it"
  type        = "bool"
  required    = false
}

parameter {
  key         = "encoded_auth"
  description = "the authentication information to log into the docker repository\nWARNING: be very careful to not leak the authentication information by hardcoding it here. Use a helper function like `file()` to read the information from a file not stored in VCS"
  type        = "string"
  required    = false
}

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

