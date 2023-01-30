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
  key         = "binds"
  description = "A 'source:destination' list of folders to mount onto the container from the host.\nA list of folders to mount onto the container from the host. The expected format for each string entry in the list is `source:destination`. So for example: `binds: [\"host_folder/scripts:/scripts\"]"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "client_config"
  description = "client config for remote Docker engine\nthis config block can be used to configure a remote Docker engine. By default Waypoint will attempt to discover this configuration using the environment variables: `DOCKER_HOST` to set the url to the docker server. `DOCKER_API_VERSION` to set the version of the API to reach, leave empty for latest. `DOCKER_CERT_PATH` to load the TLS certificates from. `DOCKER_TLS_VERIFY` to enable or disable TLS verification, off by default."
  type        = "docker.ClientConfig"
  required    = false
}

parameter {
  key           = "command"
  description   = "the command to run to start the application in the container"
  type          = "list of string"
  required      = false
  default_value = "the image entrypoint"
}

parameter {
  key         = "extra_ports"
  description = "additional TCP ports the application is listening on to expose on the container\nUsed to define and expose multiple ports that the application is listening on for the container in use. These ports will get merged with service_port when creating the container if defined."
  type        = "list of uint"
  required    = false
}

parameter {
  key           = "force_pull"
  description   = "always pull the docker container from the registry"
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key         = "labels"
  description = "A map of key/value pairs to label the docker container with.\nA map of key/value pair(s), stored in docker as a string. Each key/value pair must be unique. Validiation occurs at the docker layer, not in Waypoint. Label keys are alphanumeric strings which may contain periods (.) and hyphens (-)."
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "networks"
  description   = "An list of strings with network names to connect the container to.\nA list of networks to connect the container to. By default the container will always connect to the `waypoint` network."
  type          = "list of string"
  required      = false
  default_value = "waypoint"
}

parameter {
  key         = "resources"
  description = "A map of resources to configure the container with, such as memory or cpu limits.\nthese options are used to configure the container used when deploying with docker. Currently, the supported resources are 'memory' and 'cpu' limits. The field 'memory' is expected to be defined as \"512MB\", \"44kB\", etc."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "scratch_path"
  description = "a path within the container to store temporary data\ndocker will mount a tmpfs at this path"
  type        = "string"
  required    = false
}

parameter {
  key           = "service_port"
  description   = "port that your service is running on in the container"
  type          = "uint"
  required      = false
  default_value = "3000"
}

parameter {
  key         = "static_environment"
  description = "environment variables to expose to the application\nthese environment variables should not be run of the mill configuration variables, use waypoint config for that. These variables are used to control over all container modes, such as configuring it to start a web app vs a background worker"
  type        = "map of string to string"
  required    = false
}

