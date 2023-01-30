# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "binds"
  description = "A 'source:destination' list of folders to mount onto the container from the host.\nA list of folders to mount onto the container from the host. The expected format for each string entry in the list is `source:destination`. So for example: `binds: [\"host_folder/scripts:/scripts\"]`"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "client_config"
  description = ""
  type        = "docker.ClientConfig"
  required    = false
}

parameter {
  key         = "debug_containers"
  description = ""
  type        = "bool"
  required    = false
}

parameter {
  key         = "force_pull"
  description = ""
  type        = "bool"
  required    = false
}

parameter {
  key         = "labels"
  description = "A map of key/value pairs to label the docker container with.\nA map of key/value pair(s), stored in docker as a string. Each key/value pair must be unique. Validiation occurs at the docker layer, not in Waypoint. Label keys are alphanumeric strings which may contain periods (.) and hyphens (-)."
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "networks"
  description   = "A list of strings with network names to connect the container to.\nA list of networks to connect the container to. By default the container will always connect to the `waypoint` network."
  type          = "list of string"
  required      = false
  default_value = "waypoint"
}

parameter {
  key         = "resources"
  description = "The resources that the tasks should use."
  type        = "category"
  required    = false
}

parameter {
  key         = "resources.cpu"
  description = "The cpu shares that the tasks should use"
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
  description = "environment variables to expose to the application\nThese variables are used to control all of a container's modes, such as configuring it to start a web app vs a background worker. These environment variables should not be common configuration variables normally set in `waypoint config`."
  type        = "map of string to string"
  required    = false
}

