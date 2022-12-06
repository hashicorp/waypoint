parameter {
  key         = "client_config"
  description = <<EOF

Client config for remote Docker engine.

This config block can be used to configure a remote Docker engine. By default Waypoint will attempt to discover this configuration using the environment variables: `DOCKER_HOST` to set the url to the docker server. `DOCKER_API_VERSION` to set the version of the API to reach, leave empty for latest. `DOCKER_CERT_PATH` to load the TLS certificates from. `DOCKER_TLS_VERIFY` to enable or disable TLS verification, off by default.
EOF
  type        = "docker.ClientConfig"
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
  key         = "binds"
  description = <<EOF
A 'source:destination' list of folders to mount onto the container from the host.

A list of folders to mount onto the container from the host. The expected format for each string entry in the list is `source:destination`. So for example: `binds: ["host_folder/scripts:/scripts"].
EOF
  type        = "list of string"
  required    = false
}

parameter {
  key         = "command"
  description = "The command to run to start the application in the container."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "extra_ports"
  description = <<EOF
Additional TCP ports the application is listening on to expose on the container.

Used to define and expose multiple ports that the application is listening on for the container in use. These ports will get merged with service_port when creating the container if defined.
EOF
  type        = "list of uint"
  required    = false
}

parameter {
  key         = "force_pull"
  description = "Always pull the docker container from the registry."
  type        = "bool"
  required    = false
}

parameter {
  key         = "labels"
  description = <<EOF
A map of key/value pairs to label the docker container with.

A map of key/value pair(s), stored in docker as a string. Each key/value pair must be unique. Validiation occurs at the docker layer, not in Waypoint. Label keys are alphanumeric strings which may contain periods (.) and hyphens (-).
EOF
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "networks"
  description = <<EOF
A list of strings with network names to connect the container to.

A list of networks to connect the container to. By default the container will always connect to the `waypoint` network.
EOF
  type        = "list of string"
  required    = false
  default     = "waypoint"
}

parameter {
  key         = "resources"
  description = <<EOF
A map of resources to configure the container with, such as memory or cpu limits.

These options are used to configure the container used when deploying with docker. Currently, the supported resources are 'memory' and 'cpu' limits. The field 'memory' is expected to be defined as "512MB", "44kB", etc.
EOF
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "scratch_path"
  description = <<EOF
A path within the container to store temporary data.

Docker will mount a tmpfs at this path.
EOF
  type        = "string"
  required    = false
}

parameter {
  key         = "service_port"
  description = "Port that your service is running on in the container."
  type        = "uint"
  required    = false
  default     = "3000"
}

parameter {
  key         = "static_environment"
  description = <<EOF
Environment variables to expose to the application.

These environment variables should not be run of the mill configuration variables, use waypoint config for that. These variables are used to control over all container modes, such as configuring it to start a web app vs a background worker.
EOF
  type        = "map of string to string"
  required    = false
}
