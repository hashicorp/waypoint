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
  key         = "build_args"
  description = <<EOF
Build args to pass to docker for the build step.

A map of key/value pairs passed as build-args to docker for the build step.
EOF
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "buildkit"
  description = "If set, use the buildkit builder from Docker."
  type        = "bool"
  required    = false
}

parameter {
  key         = "context"
  description = "Build context path."
  type        = "string"
  required    = false
}

parameter {
  key         = "disable_entrypoint"
  description = <<EOF
If set, the entrypoint binary won't be injected into the image.

The entrypoint binary is what provides extended functionality such as logs and exec. If it is not injected at build time the expectation is that the image already contains it."
EOF
  type        = "bool"
  required    = false
}

parameter {
  key         = "dockerfile"
  description = <<EOF
The path to the Dockerfile.

Set this when the Dockerfile is not APP-PATH/Dockerfile.
EOF
  type        = "string"
  required    = false
}

parameter {
  key         = "no_cache"
  description = <<EOF
Do not use cache when building the image.

Ensures a clean image build.
EOF
  type        = "bool"
  required    = false
}

parameter {
  key         = "platform"
  description = <<EOF
Set target platform to build container if server is multi-platform capable.

Must enable Docker buildkit to use the 'platform' flag.
EOF
  type        = "string"
  required    = false
}


parameter {
  key         = "target"
  description = <<EOF
The target build stage in a multi-stage Dockerfile.

If buildkit is enabled unused stages will be skipped.
EOF
  type        = "string"
  required    = false
}
