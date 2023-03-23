# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

parameter {
  key         = "key"
  description = <<EOT
The key in the ConfigMap or Secret to read the value from.

ConfigMaps and Secrets store data in key/value format. This specifies the key to read from the resource. If you want multiple values you must specify multiple dynamic values.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "name"
  description = <<EOT
The name of the ConfigMap of Secret.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "namespace"
  description = <<EOT
The namespace to load the ConfigMap or Secret from.

By default this will use the namespace of the running pod. If this config source is used outside of a pod, this will use the namespace from the kubeconfig.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "secret"
  description = <<EOT
This must be set to true to read from a Secret. If it is false we read from a ConfigMap.
EOT
  type        = "string"
  required    = false
}

