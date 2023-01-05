# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "auth_type"
  description   = "the Lambda function URL auth type\nThe AuthType parameter determines how Lambda authenticates or authorizes requests to your function URL. Must be either `AWS_IAM` or `NONE`."
  type          = "string"
  required      = false
  default_value = "NONE"
}

parameter {
  key           = "principal"
  description   = "the principal to use when auth_type is `AWS_IAM`\nThe Principal parameter specifies the principal that is allowed to invoke the function."
  type          = "string"
  required      = false
  default_value = "*"
}

