parameter {
  key         = "path"
  description = <<EOT
Path to a file or directory of YAML or JSON files.

This will be used for `kubectl apply` to create a set of Kubernetes resources. Pair this with `templatefile` or `templatedir` [templating functions](/docs/waypoint-hcl/functions/template) to inject dynamic elements into your Kubernetes resources. Subdirectories are included recursively.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "prune_label"
  description = <<EOT
Label selector to prune resources that aren't present in the `path`.

This is a label selector that is used to search for any resources that are NOT present in the configured `path` and delete them.
EOT
  type        = "string"
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
  key         = "kubeconfig"
  description = <<EOT
Path to the kubeconfig file to use.

If this isn't set, the default lookup used by `kubectl` will be used.
EOT
  type        = "string"
  required    = false
}

