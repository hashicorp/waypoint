# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "annotations"
  description = "Annotations to be applied to the kube service"
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "context"
  description = "the kubectl context to use, as defined in the kubeconfig file"
  type        = "string"
  required    = false
}

parameter {
  key         = "ingress"
  description = "Configuration to set up an ingress resource to route traffic to the given application from an ingress controller\nAn ingress resource can be created on release that will route traffic to the Kubernetes service. Note that before this happens, the Kubernetes cluster must already be configured with an Ingress controller. Otherwise there won't be a way for inbound traffic to be routed to the ingress resource."
  type        = "category"
  required    = true
}

parameter {
  key         = "ingress.annotations"
  description = "Annotations to be applied to the ingress resource"
  type        = "map of string to string"
  required    = false
}

parameter {
  key           = "ingress.default"
  description   = "sets the ingress resource to be the default backend for any traffic that doesn't match existing ingress rule paths"
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key         = "ingress.host"
  description = "If set, will configure the ingress resource to have the ingress controller route traffic for any inbound requests that match this host. IP addresses are not allowed, nor are ':' delimiters. Wildcards are allowed to a certain extent. For more details check out the Kubernetes documentation"
  type        = "string"
  required    = false
}

parameter {
  key           = "ingress.path"
  description   = "The route rule that should be used to route requests to this ingress resource. A path must begin with a '/'."
  type          = "string"
  required      = false
  default_value = "/"
}

parameter {
  key           = "ingress.path_type"
  description   = "defines the kind of rule the path will be for the ingress controller. Valid path types are 'Exact', 'Prefix', and 'ImplementationSpecific'."
  type          = "string"
  required      = false
  default_value = "Prefix"
}

parameter {
  key         = "ingress.tls"
  description = "A stanza of TLS configuration options for traffic to the ingress resource"
  type        = "category"
  required    = true
}

parameter {
  key         = "ingress.tls.hosts"
  description = "A list of hosts included in the TLS certificate"
  type        = ""
  required    = true
}

parameter {
  key         = "ingress.tls.secret_name"
  description = "The Kubernetes secret name that should be used to look up or store TLS configs"
  type        = ""
  required    = true
}

parameter {
  key         = "kubeconfig"
  description = "path to the kubeconfig file to use\nby default uses from current user's home directory"
  type        = "string"
  required    = false
}

parameter {
  key         = "load_balancer"
  description = "indicates if the Kubernetes Service should LoadBalancer type\nif the Kubernetes Service is not a LoadBalancer and node_port is not set, then the Service uses ClusterIP"
  type        = "bool"
  required    = false
}

parameter {
  key         = "namespace"
  description = "namespace to create Service in\nnamespace is the name of the Kubernetes namespace to create the deployment in This is useful to create Services in non-default namespaces without creating kubeconfig contexts for each"
  type        = "string"
  required    = false
}

parameter {
  key         = "node_port"
  description = "the TCP port that the Service should consume as a NodePort\nif this is set but load_balancer is not, the service will be NodePort type, but if load_balancer is also set, it will be LoadBalancer"
  type        = "uint"
  required    = false
}

parameter {
  key           = "port"
  description   = "the TCP port that the application is listening on"
  type          = "uint"
  required      = false
  default_value = "80"
}

parameter {
  key         = "ports"
  description = "a map of ports and options that the application is listening on\nused to define and configure multiple ports that the application is listening on. Available keys are 'port', 'node_port', 'name', and 'target_port'. If 'node_port' is set but 'load_balancer' is not, the service will be NodePort type. If 'load_balancer' is also set, it will be LoadBalancer. Ports defined will be TCP protocol. Note that 'name' is required if defining more than one port."
  type        = "list of map of string to string"
  required    = false
}

