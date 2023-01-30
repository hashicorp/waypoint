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
  key         = "build_args"
  description = "build args to pass to docker for the build step\nA map of key/value pairs passed as build-args to docker for the build step."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "buildkit"
  description = "if set, use the buildkit builder from Docker"
  type        = "bool"
  required    = false
}

parameter {
  key         = "context"
  description = "Build context path"
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
  key         = "dockerfile"
  description = "The path to the Dockerfile.\nSet this when the Dockerfile is not APP-PATH/Dockerfile"
  type        = "string"
  required    = false
}

parameter {
  key         = "no_cache"
  description = "Do not use cache when building the image\nEnsures a clean image build."
  type        = "bool"
  required    = false
}

parameter {
  key         = "platform"
  description = "set target platform to build container if server is multi-platform capable\nMust enable Docker buildkit to use the 'platform' flag."
  type        = "string"
  required    = false
}

parameter {
  key         = "target"
  description = "the target build stage in a multi-stage Dockerfile\nIf buildkit is enabled unused stages will be skipped"
  type        = "string"
  required    = false
}

