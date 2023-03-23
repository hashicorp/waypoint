# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "chart"
  description = "The name or path of the chart to install.\nIf you're installing a local chart, this is the path to the chart. If you're installing a chart from a repository (have the `repository` configuration set), then this is the name of the chart in the repository."
  type        = "string"
  required    = true
}

parameter {
  key         = "context"
  description = "The kubectl context to use, as defined in the kubeconfig file."
  type        = "string"
  required    = false
}

parameter {
  key         = "create_namespace"
  description = "Create Namespace if it doesn't exist.\nThis option will instruct Helm to create a namespace if it doesn't exist."
  type        = "bool"
  required    = false
}

parameter {
  key           = "devel"
  description   = "True to considered non-released chart versions for installation.\nThis is equivalent to the `--devel` flag to `helm install`."
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key           = "driver"
  description   = "The Helm storage driver to use.\nThis can be one of `configmap`, `secret`, `memory`, or `sql`. The SQL connection string can not be set currently so this must be set on the runners."
  type          = "string"
  required      = false
  default_value = "secret"
}

parameter {
  key         = "kubeconfig"
  description = "Path to the kubeconfig file to use.\nIf this isn't set, the default lookup used by `kubectl` will be used."
  type        = "string"
  required    = false
}

parameter {
  key         = "name"
  description = "Name of the Helm release.\nThis must be globally unique within the context of your Helm installation."
  type        = "string"
  required    = true
}

parameter {
  key         = "namespace"
  description = "Namespace to deploy the Helm chart.\nThis will be created if it does not exist (see create_namespace). This defaults to the current namespace of the auth settings."
  type        = "string"
  required    = false
}

parameter {
  key         = "repository"
  description = "URL of the Helm repository that contains the chart.\nThis only needs to be set if you're NOT using a local chart."
  type        = "string"
  required    = false
}

parameter {
  key         = "set"
  description = "A single value to set. This can be repeated multiple times.\nThis sets a single value. Separate nested values with a `.`. This is the same as the `--set` flag on `helm install`."
  type        = "list of struct { Name string \"hcl:\\\"name,attr\\\"\"; Value string \"hcl:\\\"value,attr\\\"\"; Type string \"hcl:\\\"type,optional\\\"\" }"
  required    = true
}

parameter {
  key         = "skip_crds"
  description = "Do not create CRDs\nThis option will tell Helm to skip the creation of CRDs."
  type        = "bool"
  required    = false
}

parameter {
  key         = "values"
  description = "Values in raw YAML to configure the Helm chart.\nThese values are usually loaded from files using HCL functions such as `file` or templating with `templatefile`. Multiple values will be merged using the same logic as the `-f` flag with Helm."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "version"
  description = "The version of the chart to install."
  type        = "string"
  required    = false
}

