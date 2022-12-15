parameter {
  key           = "cpu"
  description   = <<EOT
Cpu resource request to be added to the task container.
EOT 
  type          = "k8s.ResourceConfig"
  required      = true

}

parameter {
  key           = "ephemeral_storage"
  description   = <<EOT
Ephemeral_storage resource request to be added to the task container.
EOT 
  type          = "k8s.ResourceConfig"
  required      = true

}

parameter {
  key           = "memory"
  description   = <<EOT
Memory resource request to be added to the task container.
EOT 
  type          = "k8s.ResourceConfig"
  required      = true

}

parameter {
  key           = "context"
  description   = <<EOT
The kubectl context to use, as defined in the kubeconfig file.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "image_pull_policy"
  description   = <<EOT
Pull policy to use for the task container image.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "image_secret"
  description   = <<EOT
Name of the Kubernetes secret to use for the image.

This references an existing secret; Waypoint does not create this secret.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "kubeconfig"
  description   = <<EOT
Path to the kubeconfig file to use.

By default uses from current user's home directory.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "namespace"
  description   = <<EOT
Namespace in which to launch task.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "service_account"
  description   = <<EOT
Service account name to be added to the application pod.

Service account is the name of the Kubernetes service account to add to the pod. This is useful to apply Kubernetes RBAC to the application.
EOT 
  type          = "string"
  required      = false

}

parameter {
  key           = "watchtask_startup_timeout_seconds"
  description   = <<EOT
This option configures how long the WatchTask should wait for a task pod to start-up before attempting to stream its logs. If the pod does not start up within the given timeout, WatchTask will exit.
EOT 
  type          = "int"
  required      = false
  default_value = "30"
}

