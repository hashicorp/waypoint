# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "context"
  description = "the kubectl context to use, as defined in the kubeconfig file"
  type        = "string"
  required    = false
}

parameter {
  key         = "cpu"
  description = "cpu resource request to be added to the task container"
  type        = "k8s.ResourceConfig"
  required    = true
}

parameter {
  key         = "ephemeral_storage"
  description = "ephemeral_storage resource request to be added to the task container"
  type        = "k8s.ResourceConfig"
  required    = true
}

parameter {
  key         = "image_pull_policy"
  description = "pull policy to use for the task container image"
  type        = "string"
  required    = false
}

parameter {
  key         = "image_secret"
  description = "name of the Kubernetes secret to use for the image\nthis references an existing secret; Waypoint does not create this secret"
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
  key         = "memory"
  description = "memory resource request to be added to the task container"
  type        = "k8s.ResourceConfig"
  required    = true
}

parameter {
  key         = "namespace"
  description = "namespace in which to launch task"
  type        = "string"
  required    = false
}

parameter {
  key         = "security_context"
  description = ""
  type        = "k8s.PodSecurityContext"
  required    = true
}

parameter {
  key         = "service_account"
  description = "service account name to be added to the application pod\nservice account is the name of the Kubernetes service account to add to the pod. This is useful to apply Kubernetes RBAC to the application."
  type        = "string"
  required    = false
}

parameter {
  key           = "watchtask_startup_timeout_seconds"
  description   = "This option configures how long the WatchTask should wait for a task pod to start-up before attempting to stream its logs. If the pod does not start up within the given timeout, WatchTask will exit."
  type          = "int"
  required      = false
  default_value = "30"
}

