# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "addr"
  description = "The address to the Vault server.\nIf this is not set, the VAULT_ADDR environment variable will be read."
  type        = "string"
  required    = false
}

parameter {
  key         = "agent_addr"
  description = "The address to the Vault agent.\nIf this is not set, Vault agent will not be used. This should only be set if you're deploying to an environment with a Vault agent."
  type        = "string"
  required    = false
}

parameter {
  key         = "approle_role_id"
  description = "The role ID of the approle auth method to use for Vault.\nThis is required for the `approle` auth method."
  type        = "string"
  required    = false
}

parameter {
  key         = "approle_secret_id"
  description = "The secret ID of the approle auth method to use for Vault.\nThis is required for the `approle` auth method."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth_method"
  description = "The authentication method to use for Vault.\nThis can be one of: `aws`, `approle`, `kubernetes`, `gcp`. When this is set, configuration fields prefixed with the auth method type should be set, if required. Configuration fields prefixed with non-matching auth method types will be ignored (except for type validation).  If no auth method is set, Vault assumes proper environment variables are set for Vault to find and connect to the Vault server. When this is set, `auth_method_mount_path` is required."
  type        = "string"
  required    = false
}

parameter {
  key         = "auth_method_mount_path"
  description = "The path where the configured auth method is mounted in Vault.\nThis is required when `auth_method` is set."
  type        = "string"
  required    = false
}

parameter {
  key         = "aws_access_key"
  description = "The access key to use for authentication to the IAM service, if needed.\nThis usually isn't needed since IAM instance profiles are used."
  type        = "string"
  required    = false
}

parameter {
  key           = "aws_credential_poll_interval"
  description   = "The interval in seconds to wait before checking for new credentials."
  type          = "int"
  required      = false
  default_value = "60"
}

parameter {
  key         = "aws_header_value"
  description = "The value to match the [`iam_server_id_header_value`](/vault/api-docs/auth/aws#iam_server_id_header_value) if set."
  type        = "string"
  required    = false
}

parameter {
  key           = "aws_region"
  description   = "The region for the STS endpoint when using that method of auth."
  type          = "string"
  required      = false
  default_value = "us-east-1"
}

parameter {
  key         = "aws_role"
  description = "The role to use for AWS authentication.\nThis is required for the `aws` auth method. This depends on how you configured the Vault [AWS Auth Method](/vault/docs/auth/aws)."
  type        = "string"
  required    = false
}

parameter {
  key         = "aws_secret_key"
  description = "The secret key to use for authentication to the IAM service, if needed.\nThis usually isn't needed since IAM instance profiles are used."
  type        = "string"
  required    = false
}

parameter {
  key         = "aws_type"
  description = "The type of authentication to use for AWS: either `iam` or `ec2`.\nThis is required for the `aws` auth method. This depends on how you configured the Vault [AWS Auth Method](/vault/docs/auth/aws)."
  type        = "string"
  required    = false
}

parameter {
  key         = "ca_cert"
  description = "The path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate."
  type        = "string"
  required    = false
}

parameter {
  key         = "ca_path"
  description = "The path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate."
  type        = "string"
  required    = false
}

parameter {
  key         = "client_cert"
  description = "The path to a PEM-encoded certificate to present as a client certificate.\nThis only needs to be set if Vault is configured to expect a client cert."
  type        = "string"
  required    = false
}

parameter {
  key         = "client_key"
  description = "The path to a private key for the client cert.\nThis only needs to be set if Vault is configured to expect a client cert."
  type        = "string"
  required    = false
}

parameter {
  key         = "gcp_credentials"
  description = "When using static credentials, the contents of the JSON credentials file."
  type        = "string"
  required    = false
}

parameter {
  key           = "gcp_jwt_exp"
  description   = "The number of minutes a generated JWT should be valid for when using the iam method."
  type          = "int"
  required      = false
  default_value = "15"
}

parameter {
  key         = "gcp_project"
  description = "The project to use, only if it cannot be automatically determined."
  type        = "string"
  required    = false
}

parameter {
  key         = "gcp_role"
  description = "The role to use for GCP authentication.\nThis is required for the `gcp` auth method. This depends on how you configured the Vault [GCP Auth Method](/vault/docs/auth/gcp)."
  type        = "string"
  required    = false
}

parameter {
  key         = "gcp_service_account"
  description = "The service account to use, only if it cannot be automatically determined."
  type        = "string"
  required    = false
}

parameter {
  key         = "gcp_type"
  description = "The type of authentication; must be `gce` or `iam`.\nThis is required for the `gcp` auth method. This depends on how you configured the Vault [GCP Auth Method](/vault/docs/auth/gcp)."
  type        = "string"
  required    = false
}

parameter {
  key         = "kubernetes_role"
  description = "The role to use for Kubernetes authentication.\nThis is required for the `kubernetes` auth method. This is a role that is setup with the [Kubernetes Auth Method in Vault](/vault/docs/auth/kubernetes)."
  type        = "string"
  required    = false
}

parameter {
  key           = "kubernetes_token_path"
  description   = "The path to the Kubernetes service account token.\nIn standard Kubernetes environments, this doesn't have to be set."
  type          = "string"
  required      = false
  default_value = "/var/run/secrets/kubernetes.io/serviceaccount/token"
}

parameter {
  key         = "namespace"
  description = "Default namespace to operate in if you're using Vault namespaces."
  type        = "string"
  required    = false
}

parameter {
  key         = "skip_verify"
  description = "Do not validate the TLS cert presented by the Vault server.\nThis is not recommended unless absolutely necessary."
  type        = "bool"
  required    = false
}

parameter {
  key         = "tls_server_name"
  description = "The TLS server name to verify with the Vault server."
  type        = "string"
  required    = false
}

parameter {
  key         = "token"
  description = "The token to use for communicating to Vault.\nIf you're using a Vault Agent or an `auth_method`, this may not be necessary. If you're using an `auth_method`, this may still be necessary as a minimal token with access to the auth method, but usually these are not protected."
  type        = "string"
  required    = false
}

