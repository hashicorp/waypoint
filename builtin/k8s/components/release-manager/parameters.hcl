parameter {
  key         = "ingress"
  description = <<EOT
Configuration to set up an ingress resource to route traffic to the given application from an ingress controller.

An ingress resource can be created on release that will route traffic to the Kubernetes service. Note that before this happens, the Kubernetes cluster must already be configured with an Ingress controller. Otherwise there won't be a way for inbound traffic to be routed to the ingress resource.
EOT 
  type        = "category"
  required    = true

}

parameter {
  key         = "ingress.annotations"
  description = <<EOT
Annotations to be applied to the ingress resource.
EOT 
  type        = "map of string to string"
  required    = true

}

parameter {
  key           = "ingress.default"
  description   = <<EOT
Sets the ingress resource to be the default backend for any traffic that doesn't match existing ingress rule paths.
EOT 
  type          = "bool"
  required      = true
  default_value = "false"
}

parameter {
  key         = "ingress.host"
  description = <<EOT
If set, will configure the ingress resource to have the ingress controller route traffic for any inbound requests that match this host. IP addresses are not allowed, nor are ':' delimiters. Wildcards are allowed to a certain extent. For more details check out the Kubernetes documentation.
EOT 
  type        = "string"
  required    = true

}

parameter {
  key           = "ingress.path"
  description   = <<EOT
The route rule that should be used to route requests to this ingress resource. A path must begin with a '/'.
EOT 
  type          = "string"
  required      = true
  default_value = "/"
}

parameter {
  key           = "ingress.path_type"
  description   = <<EOT
Defines the kind of rule the path will be for the ingress controller. Valid path types are 'Exact', 'Prefix', and 'ImplementationSpecific'.
EOT 
  type          = "string"
  required      = true
  default_value = "Prefix"
}

parameter {
  key         = "ingress.tls"
  description = <<EOT
A stanza of TLS configuration options for traffic to the ingress resource.
EOT 
  type        = "category"
  required    = true

}

parameter {
  key         = "ingress.tls.hosts"
  description = <<EOT
A list of hosts included in the TLS certificate.
EOT 
  type        = "list of string" # WARNING: no type was documented. This will be a best effort choice.
  required    = true

}

parameter {
  key         = "ingress.tls.secret_name"
  description = <<EOT
The Kubernetes secret name that should be used to look up or store TLS configs.
EOT 
  type        = "string" # WARNING: no type was documented. This will be a best effort choice.
  required    = true

}

parameter {
  key         = "annotations"
  description = <<EOT
Annotations to be applied to the kube service.
EOT 
  type        = "map of string to string"
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
  key         = "kubeconfig"
  description = <<EOT
Path to the kubeconfig file to use.

By default uses from current user's home directory.
EOT 
  type        = "string"
  required    = false

}

parameter {
  key         = "load_balancer"
  description = <<EOT
Indicates if the Kubernetes Service should LoadBalancer type.

If the Kubernetes Service is not a LoadBalancer and node_port is not set, then the Service uses ClusterIP.
EOT 
  type        = "bool"
  required    = false

}

parameter {
  key         = "namespace"
  description = <<EOT
Namespace to create Service in.

Namespace is the name of the Kubernetes namespace to create the deployment in This is useful to create Services in non-default namespaces without creating kubeconfig contexts for each.
EOT 
  type        = "string"
  required    = false

}

parameter {
  key         = "node_port"
  description = <<EOT
The TCP port that the Service should consume as a NodePort.

If this is set but load_balancer is not, the service will be NodePort type, but if load_balancer is also set, it will be LoadBalancer.
EOT 
  type        = "uint"
  required    = false

}

parameter {
  key           = "port"
  description   = <<EOT
The TCP port that the application is listening on.
EOT 
  type          = "uint"
  required      = false
  default_value = "80"
}

parameter {
  key         = "ports"
  description = <<EOT
A map of ports and options that the application is listening on.

Used to define and configure multiple ports that the application is listening on. Available keys are 'port', 'node_port', 'name', and 'target_port'. If 'node_port' is set but 'load_balancer' is not, the service will be NodePort type. If 'load_balancer' is also set, it will be LoadBalancer. Ports defined will be TCP protocol. Note that 'name' is required if defining more than one port.
EOT 
  type        = "list of map of string to string"
  required    = false

}

