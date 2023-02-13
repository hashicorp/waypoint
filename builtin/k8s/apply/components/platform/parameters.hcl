# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "context"
  description = "The kubectl context to use, as defined in the kubeconfig file."
  type        = "string"
  required    = false
}

parameter {
  key         = "kubeconfig"
  description = "Path to the kubeconfig file to use.\nIf this isn't set, the default lookup used by `kubectl` will be used."
  type        = "string"
  required    = false
}

parameter {
  key         = "path"
  description = "Path to a file or directory of YAML or  JSON files.\nThis will be used for `kubectl apply` to create a set of Kubernetes resources. Pair this with `templatefile` or `templatedir` [templating functions](/waypoint/docs/waypoint-hcl/functions/template) to inject dynamic elements into your Kubernetes resources. Subdirectories are included recursively."
  type        = "string"
  required    = true
}

parameter {
  key         = "prune_allowlist"
  description = ""
  type        = "list of string"
  required    = false
}

parameter {
  key         = "prune_label"
  description = "Label selector to prune resources that aren't present in the `path`.\nThis is a label selector that is used to search for any resources that are NOT present in the configured `path` and delete them."
  type        = "string"
  required    = true
}

