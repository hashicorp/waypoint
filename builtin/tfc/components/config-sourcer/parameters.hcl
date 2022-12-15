parameter {
  key         = "organization"
  description = <<EOT
The Terraform Cloud organization to query.
Within a single workspace, multiple dynamic values that use the same organization and workspace will only read the value once. This allows outputs to be extracted into multiple values. The example above shows this functionality by reading the username and password into separate values.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "output"
  description = <<EOT
The name of the output to read the value of.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "workspace"
  description = <<EOT
The Terraform Cloud workspace associated with the given organization to read the outputs of.
The outputs associtaed with the most recent state version for the given workspace are the ones that are used. These values are refreshed according to refreshInternal, a source field.

EOT
  type        = "string"
  required    = true

}
