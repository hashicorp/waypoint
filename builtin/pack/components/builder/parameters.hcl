parameter {
  key           = "builder"
  description   = <<EOT
The buildpack builder image to use.

EOT
  type          = "string"
  required      = false
  default_value = "heroku/buildpacks:20"
}

parameter {
  key         = "buildpacks"
  description = <<EOT
The exact buildpacks to use.
If set, the builder will run these buildpacks in the specified order. They can be listed using several [URI formats](https://buildpacks.io/docs/app-developer-guide/specific-buildpacks).

EOT
  type        = "list of string"
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
  key         = "ignore"
  description = <<EOT
File patterns to match files which will not be included in the build.
Each pattern follows the semantics of .gitignore. This is a summarized version:

EOT
  type        = "list of string"
  required    = false

}

parameter {
  key         = "process_type"
  description = <<EOT
The process type to use from your Procfile. if not set, defaults to `web`.
The process type is used to control over all container modes, such as configuring it to start a web app vs a background worker.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "static_environment"
  description = <<EOT
Environment variables to expose to the buildpack.
These environment variables should not be run of the mill configuration variables, use waypoint config for that. These variables are used to control over all container modes, such as configuring it to start a web app vs a background worker.

EOT
  type        = "map of string to string"
  required    = false

}

