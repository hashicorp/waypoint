parameter {
  key           = "auth_type"
  description   = <<EOT
The Lambda function URL auth type.

The AuthType parameter determines how Lambda authenticates or authorizes requests to your function URL. Must be either `AWS_IAM` or `NONE`.
EOT
  type          = "string"
  required      = false
  default_value = "NONE"
}

parameter {
  key           = "principal"
  description   = <<EOT
The principal to use when auth_type is `AWS_IAM`.

The Principal parameter specifies the principal that is allowed to invoke the function.
EOT
  type          = "string"
  required      = false
  default_value = "*"
}

