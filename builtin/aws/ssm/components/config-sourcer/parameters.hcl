# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "access_key"
  description = "This is the AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable, or via a shared credentials file if `profile` is specified"
  type        = "string"
  required    = false
}

parameter {
  key         = "assume_role_arn"
  description = "Amazon Resource Name (ARN) of the IAM Role to assume."
  type        = "string"
  required    = false
}

parameter {
  key         = "assume_role_duration_seconds"
  description = "Number of seconds to restrict the assume role session duration."
  type        = "int"
  required    = false
}

parameter {
  key         = "assume_role_external_id"
  description = "External identifier to use when assuming the role."
  type        = "string"
  required    = false
}

parameter {
  key         = "assume_role_policy"
  description = "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed."
  type        = "string"
  required    = false
}

parameter {
  key         = "assume_role_session_name"
  description = "Session name to use when assuming the role."
  type        = "string"
  required    = false
}

parameter {
  key         = "iam_endpoint"
  description = "Custom endpoint address for the IAM service."
  type        = "string"
  required    = false
}

parameter {
  key           = "insecure"
  description   = "Explicitly allow the provider to perform \"insecure\" SSL requests."
  type          = "bool"
  required      = false
  default_value = "false"
}

parameter {
  key           = "max_retries"
  description   = "This is the maximum number of times an API call is retried, in the case where requests are being throttled or experiencing transient failures. The delay between the subsequent API calls increases exponentially."
  type          = "int"
  required      = false
  default_value = "25"
}

parameter {
  key         = "profile"
  description = "This is the AWS profile name as set in the shared credentials file."
  type        = "string"
  required    = false
}

parameter {
  key         = "region"
  description = "This is the AWS region. It must be provided, but it can also be sourced from the `AWS_DEFAULT_REGION` environment variables, or via a shared credentials file if profile is specified."
  type        = "string"
  required    = false
}

parameter {
  key         = "secret_key"
  description = "This is the AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable, or via a shared credentials file if `profile` is specified."
  type        = "string"
  required    = false
}

parameter {
  key         = "shared_credentials_file"
  description = "This is the path to the shared credentials file. If this is not set and a profile is specified, `~/.aws/credentials` will be used."
  type        = "string"
  required    = false
}

parameter {
  key         = "skip_credentials_validation"
  description = "Skip the credentials validation via the STS API. Useful for AWS API implementations that do not have STS available or implemented."
  type        = "bool"
  required    = false
}

parameter {
  key         = "skip_metadata_api_check"
  description = "Skip the AWS Metadata API check. Useful for AWS API implementations that do not have a metadata API endpoint. Setting to true prevents Terraform from authenticating via the Metadata API. You may need to use other authentication methods like static credentials, configuration variables, or environment variables."
  type        = "bool"
  required    = false
}

parameter {
  key         = "skip_requesting_account_id"
  description = "Skip requesting the account ID. Useful for AWS API implementations that do not have the IAM, STS API, or metadata API."
  type        = "bool"
  required    = false
}

parameter {
  key         = "sts_endpoint"
  description = "Custom endpoint for the STS service."
  type        = "string"
  required    = false
}

parameter {
  key         = "token"
  description = ""
  type        = "string"
  required    = false
}

