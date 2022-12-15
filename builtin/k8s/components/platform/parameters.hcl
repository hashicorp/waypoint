parameter {
  key         = "pod"
  description = <<EOT
The configuration for a pod.
Pod describes the configuration for a pod when deploying.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.container"
  description = <<EOT
Container describes the commands and arguments for a container config.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.container.args"
  description = <<EOT
An array of string arguments to pass through to the container.

EOT
  type        = "list of string"
  required    = true

}

parameter {
  key         = "pod.container.command"
  description = <<EOT
An array of strings to run for the container.

EOT
  type        = "list of string"
  required    = true

}

parameter {
  key         = "pod.container.cpu"
  description = <<EOT
Cpu resource configuration.
CPU lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.container.cpu.limit"
  description = <<EOT
Maximum amount of cpu to give the container. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.cpu.request"
  description = <<EOT
How much cpu to give the container in cpu cores. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.memory"
  description = <<EOT
Memory resource configuration.
Memory lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.container.memory.limit"
  description = <<EOT
Maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.memory.request"
  description = <<EOT
How much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.name"
  description = <<EOT
Name of the container.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.port"
  description = <<EOT
A port and options that the application is listening on.
Used to define and expose multiple ports that the application or process is listening on for the container in use. Can be specified multiple times for many ports.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.container.port.host_ip"
  description = <<EOT
What host IP to bind the external port to.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.port.host_port"
  description = <<EOT
The corresponding worker node port.
Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.

EOT
  type        = "uint"
  required    = true

}

parameter {
  key         = "pod.container.port.name"
  description = <<EOT
Name of the port.
If specified, this must be an IANA\_SVC\_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.port.port"
  description = <<EOT
The port number.
Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.

EOT
  type        = "uint"
  required    = true

}

parameter {
  key           = "pod.container.port.protocol"
  description   = <<EOT
Protocol for port. Must be UDP, TCP, or SCTP.

EOT
  type          = "string"
  required      = true
  default_value = "TCP"
}

parameter {
  key         = "pod.container.probe"
  description = <<EOT
Configuration to control liveness and readiness probes.
Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.

EOT
  type        = "category"
  required    = true

}

parameter {
  key           = "pod.container.probe.failure_threshold"
  description   = <<EOT
Number of times a liveness probe can fail before the container is killed.
FailureThreshold \* TimeoutSeconds should be long enough to cover your worst case startup times.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key           = "pod.container.probe.initial_delay"
  description   = <<EOT
Time in seconds to wait before performing the initial liveness and readiness probes.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key           = "pod.container.probe.timeout"
  description   = <<EOT
Time in seconds before the probe fails.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key         = "pod.container.probe_path"
  description = <<EOT
The HTTP path to request to test that the application is running.
Without this, the test will simply be that the application has bound to the port.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.container.resources"
  description = <<EOT
A map of resource limits and requests to apply to a container on deploy.
Resource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored.

EOT
  type        = "map of string to string"
  required    = true

}

parameter {
  key         = "pod.container.static_environment"
  description = <<EOT
Environment variables to control broad modes of the application.
Environment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.

EOT
  type        = "map of string to string"
  required    = true

}

parameter {
  key         = "pod.security_context"
  description = <<EOT
Holds pod-level security attributes and container settings.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.security_context.fs_group"
  description = <<EOT
A special supplemental group that applies to all containers in a pod.

EOT
  type        = "int64"
  required    = true

}

parameter {
  key         = "pod.security_context.run_as_non_root"
  description = <<EOT
Indicates that the container must run as a non-root user.

EOT
  type        = "bool"
  required    = true

}

parameter {
  key         = "pod.security_context.run_as_user"
  description = <<EOT
The UID to run the entrypoint of the container process.

EOT
  type        = "int64"
  required    = true

}

parameter {
  key         = "pod.sidecar"
  description = <<EOT
A sidecar container within the same pod.
Another container to run alongside the app container in the kubernetes pod. Can be specified multiple times for multiple sidecars.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.sidecar.container"
  description = <<EOT
Container describes the commands and arguments for a container config.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.args"
  description = <<EOT
An array of string arguments to pass through to the container.

EOT
  type        = "list of string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.command"
  description = <<EOT
An array of strings to run for the container.

EOT
  type        = "list of string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.cpu"
  description = <<EOT
Cpu resource configuration.
CPU lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.cpu.limit"
  description = <<EOT
Maximum amount of cpu to give the container. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.cpu.request"
  description = <<EOT
How much cpu to give the container in cpu cores. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.memory"
  description = <<EOT
Memory resource configuration.
Memory lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.memory.limit"
  description = <<EOT
Maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.memory.request"
  description = <<EOT
How much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.name"
  description = <<EOT
Name of the container.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.port"
  description = <<EOT
A port and options that the application is listening on.
Used to define and expose multiple ports that the application or process is listening on for the container in use. Can be specified multiple times for many ports.

EOT
  type        = "category"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.port.host_ip"
  description = <<EOT
What host IP to bind the external port to.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.port.host_port"
  description = <<EOT
The corresponding worker node port.
Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.

EOT
  type        = "uint"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.port.name"
  description = <<EOT
Name of the port.
If specified, this must be an IANA\_SVC\_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.port.port"
  description = <<EOT
The port number.
Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.

EOT
  type        = "uint"
  required    = true

}

parameter {
  key           = "pod.sidecar.container.port.protocol"
  description   = <<EOT
Protocol for port. Must be UDP, TCP, or SCTP.

EOT
  type          = "string"
  required      = true
  default_value = "TCP"
}

parameter {
  key         = "pod.sidecar.container.probe"
  description = <<EOT
Configuration to control liveness and readiness probes.
Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.

EOT
  type        = "category"
  required    = true

}

parameter {
  key           = "pod.sidecar.container.probe.failure_threshold"
  description   = <<EOT
Number of times a liveness probe can fail before the container is killed.
FailureThreshold \* TimeoutSeconds should be long enough to cover your worst case startup times.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key           = "pod.sidecar.container.probe.initial_delay"
  description   = <<EOT
Time in seconds to wait before performing the initial liveness and readiness probes.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key           = "pod.sidecar.container.probe.timeout"
  description   = <<EOT
Time in seconds before the probe fails.

EOT
  type          = "uint"
  required      = true
  default_value = "5"
}

parameter {
  key         = "pod.sidecar.container.probe_path"
  description = <<EOT
The HTTP path to request to test that the application is running.
Without this, the test will simply be that the application has bound to the port.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.resources"
  description = <<EOT
A map of resource limits and requests to apply to a container on deploy.
Resource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored.

EOT
  type        = "map of string to string"
  required    = true

}

parameter {
  key         = "pod.sidecar.container.static_environment"
  description = <<EOT
Environment variables to control broad modes of the application.
Environment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.

EOT
  type        = "map of string to string"
  required    = true

}

parameter {
  key         = "pod.sidecar.image"
  description = <<EOT
Image of the sidecar container.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "annotations"
  description = <<EOT
Annotations to be added to the application pod.
Annotations are added to the pod spec of the deployed application. This is useful when using mutating webhook admission controllers to further process pod events.

EOT
  type        = "map of string to string"
  required    = false

}

parameter {
  key         = "autoscale"
  description = <<EOT
Sets up a horizontal pod autoscaler to scale deployments automatically.
This configuration will automatically set up and associate the current deployment with a horizontal pod autoscaler in Kuberentes. Note that for this to work, you must also define resource limits and requests for a deployment otherwise the metrics-server will not be able to properly determine a deployments target CPU utilization.

EOT
  type        = "category"
  required    = false

}

parameter {
  key         = "autoscale.cpu_percent"
  description = <<EOT
The target CPU percent utilization before the horizontal pod autoscaler scales up a deployments replicas.

EOT
  type        = "int32"
  required    = false

}

parameter {
  key         = "autoscale.max_replicas"
  description = <<EOT
The maximum amount of pods to scale to for a deployment.

EOT
  type        = "int32"
  required    = false

}

parameter {
  key         = "autoscale.min_replicas"
  description = <<EOT
The minimum amount of pods to have for a deployment.

EOT
  type        = "int32"
  required    = false

}

parameter {
  key         = "context"
  description = <<EOT
The kubectl context to use, as defined in the kubeconfig file.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "cpu"
  description = <<EOT
Cpu resource configuration.
CPU lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = false

}

parameter {
  key         = "cpu.limit"
  description = <<EOT
Maximum amount of cpu to give the container. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "cpu.request"
  description = <<EOT
How much cpu to give the container in cpu cores. Supports m to indicate milli-cores.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "image_secret"
  description = <<EOT
Name of the Kubernetes secrete to use for the image.
This references an existing secret, waypoint does not create this secret.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "kubeconfig"
  description = <<EOT
Path to the kubeconfig file to use.
By default uses from current user's home directory.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "labels"
  description = <<EOT
A map of key value labels to apply to the deployment pod.

EOT
  type        = "map of string to string"
  required    = false

}

parameter {
  key         = "memory"
  description = <<EOT
Memory resource configuration.
Memory lets you define resource limits and requests for a container in a deployment.

EOT
  type        = "category"
  required    = false

}

parameter {
  key         = "memory.limit"
  description = <<EOT
Maximum amount of memory to give the container. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "memory.request"
  description = <<EOT
How much memory to give the container in bytes. Supports k for kilobytes, m for megabytes, and g for gigabytes.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "namespace"
  description = <<EOT
Namespace to target deployment into.
Namespace is the name of the Kubernetes namespace to apply the deployment in. This is useful to create deployments in non-default namespaces without creating kubeconfig contexts for each.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "probe"
  description = <<EOT
Configuration to control liveness and readiness probes.
Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.

EOT
  type        = "category"
  required    = false

}

parameter {
  key           = "probe.failure_threshold"
  description   = <<EOT
Number of times a liveness probe can fail before the container is killed.
FailureThreshold \* TimeoutSeconds should be long enough to cover your worst case startup times.

EOT
  type          = "uint"
  required      = false
  default_value = "30"
}

parameter {
  key           = "probe.initial_delay"
  description   = <<EOT
Time in seconds to wait before performing the initial liveness and readiness probes.

EOT
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key           = "probe.timeout"
  description   = <<EOT
Time in seconds before the probe fails.

EOT
  type          = "uint"
  required      = false
  default_value = "5"
}

parameter {
  key         = "probe_path"
  description = <<EOT
The HTTP path to request to test that the application is running.
Without this, the test will simply be that the application has bound to the port.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "replicas"
  description = <<EOT
The number of replicas to maintain.
If the replica count is maintained outside waypoint, for instance by a pod autoscaler, do not set this variable.

EOT
  type        = "int32"
  required    = false

}

parameter {
  key         = "resources"
  description = <<EOT
A map of resource limits and requests to apply to a container on deploy.
Resource limits and requests for a container. This exists to allow any possible resources. For cpu and memory, use those relevant settings instead. Keys must start with either `limits_` or `requests_`. Any other options will be ignored.

EOT
  type        = "map of string to string"
  required    = false

}

parameter {
  key         = "scratch_path"
  description = <<EOT
A path for the service to store temporary data.
A path to a directory that will be created for the service to store temporary data using EmptyDir.

EOT
  type        = "list of string"
  required    = false

}

parameter {
  key         = "service_account"
  description = <<EOT
Service account name to be added to the application pod.
Service account is the name of the Kubernetes service account to add to the pod. This is useful to apply Kubernetes RBAC to the application.

EOT
  type        = "string"
  required    = false

}

parameter {
  key           = "service_port"
  description   = <<EOT
The TCP port that the application is listening on.
By default, this config variable is used for exposing a single port for the container in use. For multi-port configuration, use 'ports' instead.

EOT
  type          = "uint"
  required      = false
  default_value = "3000"
}

parameter {
  key         = "static_environment"
  description = <<EOT
Environment variables to control broad modes of the application.
Environment variables that are meant to configure the container in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.

EOT
  type        = "map of string to string"
  required    = false

}

