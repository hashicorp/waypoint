parameter {
  key         = "client_config"
  type        = "docker.ClientConfig"
  required    = true
}

parameter {
  key         = "binds"
  description = <<EOF
A 'source:destination' list of folders to mount onto the container from the host.

A list of folders to mount onto the container from the host. The expected format for each string entry in the list is `source:destination`. So for example: `binds: ["host_folder/scripts:/scripts"]`.
EOF
  type        = "list of string"
  required    = false
}

parameter {
  key         = "debug_containers"
  type        = "bool"
  required    = false
}

parameter {
  key         = "force_pull"
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
  description = "The resources that the tasks should use."
  type        = "category"
  required    = false
}

parameter {
  key         = "resources.cpu"
  description = "The cpu shares that the tasks should use."
  type        = "int64"
  required    = false
}

parameter {
  key         = "resources.memory"
  description = "The amount of memory to use. Defined as '512MB', '44kB', etc."
  type        = "string"
  required    = false
}

parameter {
  key         = "static_environment"
  description = <<EOF
Environment variables to expose to the application.

These variables are used to control all of a container's modes, such as configuring it to start a web app vs a background worker. These environment variables should not be common configuration variables normally set in `waypoint config`.
EOF
  type        = "map of string to string"
  required    = false
}
