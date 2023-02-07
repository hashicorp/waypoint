# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "annotations"
  description = "annotations to be added to the application pod\nannotations are added to the pod spec of the deployed application. This is useful when using mutating webhook admission controllers to further process pod events."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "autoscale"
  description = "sets up a horizontal pod autoscaler to scale deployments automatically\nThis configuration will automatically set up and associate the current deployment with a horizontal pod autoscaler in Kuberentes. Note that for this to work, you must also define resource limits and requests for a deployment otherwise the metrics-server will not be able to properly determine a deployments target CPU utilization"
  type        = "category"
  required    = false
}

parameter {
  key         = "autoscale.cpu_percent"
  description = "the target CPU percent utilization before the horizontal pod autoscaler scales up a deployments replicas"
  type        = "int32"
  required    = false
}

parameter {
  key         = "autoscale.max_replicas"
  description = "the maximum amount of pods to scale to for a deployment"
  type        = "int32"
  required    = false
}

parameter {
  key         = "autoscale.min_replicas"
  description = "the minimum amount of pods to have for a deployment"
  type        = "int32"
  required    = false
}

parameter {
  key         = "context"
  description = "the kubectl context to use, as defined in the kubeconfig file"
  type        = "string"
  required    = false
}

parameter {
  key         = "cpu"
  description = "cpu resource configuration\nCPU lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "cpu.limit"
  description = "maximum amount of cpu to give the container. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "cpu.request"
  description = "how much cpu to give the container in cpu cores. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "image_secret"
  description = "name of the Kubernetes secrete to use for the image\nthis references an existing secret, waypoint does not create this secret"
  type        = "string"
  required    = false
}

parameter {
  key         = "kubeconfig"
  description = "path to the kubeconfig file to use\nby default uses from current user's home directory"
  type        = "string"
  required    = false
}

parameter {
  key         = "labels"
  description = "a map of key value labels to apply to the deployment pod"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "memory"
  description = "memory resource configuration\nMemory lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "memory.limit"
  description = "maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "memory.request"
  description = "how much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "namespace"
  description = "namespace to target deployment into\nnamespace is the name of the Kubernetes namespace to apply the deployment in. This is useful to create deployments in non-default namespaces without creating kubeconfig contexts for each"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod"
  description = "the configuration for a pod\nPod describes the configuration for a pod when deploying"
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.container"
  description = "container describes the commands and arguments for a container config"
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.container.args"
  description = "an array of string arguments to pass through to the container"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "pod.container.command"
  description = "an array of strings to run for the container"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "pod.container.cpu"
  description = "cpu resource configuration\nCPU lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "pod.container.cpu.limit"
  description = "maximum amount of cpu to give the container. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.cpu.request"
  description = "how much cpu to give the container in cpu cores. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.memory"
  description = "memory resource configuration\nMemory lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "pod.container.memory.limit"
  description = "maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.memory.request"
  description = "how much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.name"
  description = "name of the container"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.port"
  description = "a port and options that the application is listening on\nused to define and expose multiple ports that the application or process is listening on for the container in use. Can be specified multiple times for many ports."
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.container.port.host_ip"
  description = "what host IP to bind the external port to"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.port.host_port"
  description = "the corresponding worker node port\nNumber of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this."
  type        = "uint"
  required    = false
}

parameter {
  key         = "pod.container.port.name"
  description = "name of the port\nIf specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services."
  type        = "string"
  required    = true
}

parameter {
  key         = "pod.container.port.port"
  description = "the port number\nNumber of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536."
  type        = "uint"
  required    = true
}

parameter {
  key           = "pod.container.port.protocol"
  description   = "protocol for port. Must be UDP, TCP, or SCTP"
  type          = "string"
  required      = false
  default_value = "TCP"
}

parameter {
  key         = "pod.container.probe"
  description = "configuration to control liveness and readiness probes\nProbe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic."
  type        = "category"
  required    = false
}

parameter {
  key           = "pod.container.probe.failure_threshold"
  description   = "number of times a liveness probe can fail before the container is killed\nfailureThreshold * TimeoutSeconds should be long enough to cover your worst case startup times"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "pod.container.probe.initial_delay"
  description   = "time in seconds to wait before performing the initial liveness and readiness probes"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "pod.container.probe.timeout"
  description   = "time in seconds before the probe fails"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key         = "pod.container.probe_path"
  description = "the HTTP path to request to test that the application is running\nwithout this, the test will simply be that the application has bound to the port"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.container.resources"
  description = "a map of resource limits and requests to apply to a container on deploy\nresource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "pod.container.static_environment"
  description = "environment variables to control broad modes of the application\nenvironment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "pod.security_context"
  description = "holds pod-level security attributes and container settings"
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.security_context.fs_group"
  description = "a special supplemental group that applies to all containers in a pod"
  type        = "int64"
  required    = true
}

parameter {
  key         = "pod.security_context.run_as_group"
  description = ""
  type        = "int64"
  required    = true
}

parameter {
  key         = "pod.security_context.run_as_non_root"
  description = "indicates that the container must run as a non-root user"
  type        = "bool"
  required    = true
}

parameter {
  key         = "pod.security_context.run_as_user"
  description = "the UID to run the entrypoint of the container process"
  type        = "int64"
  required    = true
}

parameter {
  key         = "pod.sidecar"
  description = "a sidecar container within the same pod\nAnother container to run alongside the app container in the kubernetes pod. Can be specified multiple times for multiple sidecars."
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.sidecar.container"
  description = "container describes the commands and arguments for a container config"
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.sidecar.container.args"
  description = "an array of string arguments to pass through to the container"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.command"
  description = "an array of strings to run for the container"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.cpu"
  description = "cpu resource configuration\nCPU lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.cpu.limit"
  description = "maximum amount of cpu to give the container. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.cpu.request"
  description = "how much cpu to give the container in cpu cores. Supports m to indicate milli-cores"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.memory"
  description = "memory resource configuration\nMemory lets you define resource limits and requests for a container in a deployment."
  type        = "category"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.memory.limit"
  description = "maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.memory.request"
  description = "how much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.name"
  description = "name of the container"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.port"
  description = "a port and options that the application is listening on\nused to define and expose multiple ports that the application or process is listening on for the container in use. Can be specified multiple times for many ports."
  type        = "category"
  required    = true
}

parameter {
  key         = "pod.sidecar.container.port.host_ip"
  description = "what host IP to bind the external port to"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.port.host_port"
  description = "the corresponding worker node port\nNumber of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this."
  type        = "uint"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.port.name"
  description = "name of the port\nIf specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services."
  type        = "string"
  required    = true
}

parameter {
  key         = "pod.sidecar.container.port.port"
  description = "the port number\nNumber of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536."
  type        = "uint"
  required    = true
}

parameter {
  key           = "pod.sidecar.container.port.protocol"
  description   = "protocol for port. Must be UDP, TCP, or SCTP"
  type          = "string"
  required      = false
  default_value = "TCP"
}

parameter {
  key         = "pod.sidecar.container.probe"
  description = "configuration to control liveness and readiness probes\nProbe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic."
  type        = "category"
  required    = false
}

parameter {
  key           = "pod.sidecar.container.probe.failure_threshold"
  description   = "number of times a liveness probe can fail before the container is killed\nfailureThreshold * TimeoutSeconds should be long enough to cover your worst case startup times"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "pod.sidecar.container.probe.initial_delay"
  description   = "time in seconds to wait before performing the initial liveness and readiness probes"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "pod.sidecar.container.probe.timeout"
  description   = "time in seconds before the probe fails"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key         = "pod.sidecar.container.probe_path"
  description = "the HTTP path to request to test that the application is running\nwithout this, the test will simply be that the application has bound to the port"
  type        = "string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.resources"
  description = "a map of resource limits and requests to apply to a container on deploy\nresource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "pod.sidecar.container.static_environment"
  description = "environment variables to control broad modes of the application\nenvironment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "pod.sidecar.image"
  description = "image of the sidecar container"
  type        = "string"
  required    = true
}

parameter {
  key         = "probe"
  description = "configuration to control liveness and readiness probes\nProbe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic."
  type        = "category"
  required    = false
}

parameter {
  key           = "probe.failure_threshold"
  description   = "number of times a liveness probe can fail before the container is killed\nfailureThreshold * TimeoutSeconds should be long enough to cover your worst case startup times"
  type          = "uint"
  required      = false
  default_value = "30"
}

parameter {
  key           = "probe.initial_delay"
  description   = "time in seconds to wait before performing the initial liveness and readiness probes"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "probe.timeout"
  description   = "time in seconds before the probe fails"
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key         = "probe_path"
  description = "the HTTP path to request to test that the application is running\nwithout this, the test will simply be that the application has bound to the port"
  type        = "string"
  required    = false
}

parameter {
  key         = "replicas"
  description = "the number of replicas to maintain\nif the replica count is maintained outside waypoint, for instance by a pod autoscaler, do not set this variable"
  type        = "int32"
  required    = false
}

parameter {
  key         = "resources"
  description = "a map of resource limits and requests to apply to a container on deploy\nresource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "scratch_path"
  description = "a path for the service to store temporary data\na path to a directory that will be created for the service to store temporary data using EmptyDir."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "service_account"
  description = "service account name to be added to the application pod\nservice account is the name of the Kubernetes service account to add to the pod. This is useful to apply Kubernetes RBAC to the application."
  type        = "string"
  required    = false
}

parameter {
  key           = "service_port"
  description   = "the TCP port that the application is listening on\nby default, this config variable is used for exposing a single port for the container in use. For multi-port configuration, use 'ports' instead."
  type          = "uint"
  required      = false
  default_value = "3000"
}

parameter {
  key         = "static_environment"
  description = "environment variables to control broad modes of the application\nenvironment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands"
  type        = "map of string to string"
  required    = false
}

