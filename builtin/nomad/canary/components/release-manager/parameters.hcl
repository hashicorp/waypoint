parameter {
  key         = "fail_deployment"
  description = <<EOT
If true, marks the deployment as failed.

EOT
  type        = "bool"
  required    = false

}

parameter {
  key         = "groups"
  description = <<EOT
List of task group names which are to be promoted.

EOT
  type        = "list of string"
  required    = false

}

