# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "auth_type"
  description   = "the Lambda function URL auth type\nThe AuthType parameter determines how Lambda authenticates or authorizes requests to your function URL. Must be either `AWS_IAM` or `NONE`."
  type          = "string"
  required      = false
  default_value = "NONE"
}

parameter {
  key           = "cors"
  description   = "CORS configuration for the function URL"
  type          = "category"
  required      = false
  default_value = "NONE"
}

parameter {
  key           = "cors.allow_credentials"
  description   = "Whether to allow cookies or other credentials in requests to your function URL."
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key           = "cors.allow_headers"
  description   = "The HTTP headers that origins can include in requests to your function URL. For example: Date, Keep-Alive, X-Custom-Header."
  type          = "list of string"
  required      = false
  default_value = "[]"
}

parameter {
  key           = "cors.allow_methods"
  description   = "The HTTP methods that are allowed when calling your function URL. For example: GET, POST, DELETE, or the wildcard character (*)."
  type          = "list of string"
  required      = false
  default_value = "[]"
}

parameter {
  key           = "cors.allow_origins"
  description   = "The origins that can access your function URL. You can list any number of specific origins, separated by a comma. You can grant access to all origins using the wildcard character (*)."
  type          = "list of string"
  required      = false
  default_value = "[]"
}

parameter {
  key           = "cors.expose_headers"
  description   = "The HTTP headers in your function response that you want to expose to origins that call your function URL. For example: Date, Keep-Alive, X-Custom-Header."
  type          = "list of string"
  required      = false
  default_value = "[]"
}

parameter {
  key           = "cors.max_age"
  description   = "The maximum amount of time, in seconds, that web browsers can cache results of a preflight request."
  type          = "int64"
  required      = false
  default_value = "0"
}

parameter {
  key           = "principal"
  description   = "the principal to use when auth_type is `AWS_IAM`\nThe Principal parameter specifies the principal that is allowed to invoke the function."
  type          = "string"
  required      = false
  default_value = "*"
}

