parameter {
  key         = "chart"
  description = <<EOT
If you're installing a local chart, this is the path to the chart. If you're installing a chart from a repository (have the `repository` configuration set), then this is the name of the chart in the repository.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "name"
  description = <<EOT
Name of the Helm release.

This must be globally unique within the context of your Helm installation.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "set"
  description = <<EOT
A single value to set. This can be repeated multiple times.

This sets a single value. Separate nested values with a `.`. This is the same as the `--set` flag on `helm install`.
EOT
  type        = "list of struct { Name string \"hcl:\"name,attr\"\"; Value string \"hcl:\"value,attr\"\"; Type string \"hcl:\"type,optional\"\" }"
  required    = true
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
  key         = "create_namespace"
  description = <<EOT
Create Namespace if it doesn't exist.

This option will instruct Helm to create a namespace if it doesn't exist.
EOT
  type        = "string"
  required    = false
}

parameter {
  key           = "devel"
  description   = <<EOT
True to considered non-released chart versions for installation.

This is equivalent to the `--devel` flag to `helm install`.
EOT
  type          = "string"
  required      = false
  default_value = "false"
}

parameter {
  key           = "driver"
  description   = <<EOT
The Helm storage driver to use.

This can be one of `configmap`, `secret`, `memory`, or `sql`. The SQL connection string can not be set currently so this must be set on the runners.
EOT
  type          = "string"
  required      = false
  default_value = "secret"
}

parameter {
  key         = "kubeconfig"
  description = <<EOT
Path to the kubeconfig file to use.

If this isn't set, the default lookup used by `kubectl` will be used.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "namespace"
  description = <<EOT
Namespace to deploy the Helm chart.

This will be created if it does not exist (see create_namespace). This defaults to the current namespace of the auth settings.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "repository"
  description = <<EOT
URL of the Helm repository that contains the chart.

This only needs to be set if you're NOT using a local chart.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "skip_crds"
  description = <<EOT
Do not create CRDs.

This option will tell Helm to skip the creation of CRDs.
EOT
  type        = "bool"
  required    = false
}

parameter {
  key         = "values"
  description = <<EOT
Values in raw YAML to configure the Helm chart.

These values are usually loaded from files using HCL functions such as `file` or templating with `templatefile`. Multiple values will be merged using the same logic as the `-f` flag with Helm.
EOT
  type        = "list of string"
  required    = false
}

parameter {
  key         = "version"
  description = <<EOT
The version of the chart to install.
EOT
  type        = "string"
  required    = false
}

