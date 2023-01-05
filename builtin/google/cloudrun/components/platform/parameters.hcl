# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "auto_scaling"
  description = "Configuration to control the auto scaling parameters for Cloud Run."
  type        = "category"
  required    = true
}

parameter {
  key           = "auto_scaling.max"
  description   = "Maximum number of Cloud Run instances. When the maximum requests per container is exceeded, Cloud Run will create an additional container instance to handle load.\n\t\t\t\tThis parameter controls the maximum number of instances that can be created."
  type          = "int"
  required      = true
  default_value = "1000"
}

parameter {
  key         = "capacity"
  description = "CPU, Memory, and resource limits for each Cloud Run instance."
  type        = "category"
  required    = true
}

parameter {
  key           = "capacity.cpu_count"
  description   = "Number of CPUs to allocate the Cloud Run instance, min 1, max 2."
  type          = "int"
  required      = true
  default_value = "1"
}

parameter {
  key           = "capacity.max_requests_per_container"
  description   = "Maximum number of concurrent requests each instance can handle. When the maximum requests are exceeded, Cloud Run will create an additional instance."
  type          = "int"
  required      = true
  default_value = "80"
}

parameter {
  key           = "capacity.memory"
  description   = "Memory to allocate the Cloud Run instance specified in MB, min 128, max 4096."
  type          = "int"
  required      = true
  default_value = "128"
}

parameter {
  key           = "capacity.request_timeout"
  description   = "Maximum time a request can take before timing out, max 900."
  type          = "int"
  required      = true
  default_value = "300"
}

parameter {
  key         = "cloudsql_instances"
  description = "Specify list of CloudSQL instances that the Cloud Run instance will have access to."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "location"
  description = "GCP location, e.g. europe-north-1."
  type        = "string"
  required    = true
}

parameter {
  key         = "port"
  description = "The port your application listens on."
  type        = "int"
  required    = false
}

parameter {
  key         = "project"
  description = "GCP project ID where the Cloud Run instance will be deployed."
  type        = "string"
  required    = true
}

parameter {
  key         = "service_account_name"
  description = "Specify a service account email that Cloud Run will use to run the service. You must have the `iam.serviceAccounts.actAs` permission on the service account."
  type        = "string"
  required    = false
}

parameter {
  key         = "static_environment"
  description = "Additional environment variables to be added to the Cloud Run instance."
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "unauthenticated"
  description = "Is public unauthenticated access allowed for the Cloud Run instance?"
  type        = "bool"
  required    = false
}

parameter {
  key         = "vpc_access"
  description = "VPCAccess details"
  type        = "category"
  required    = true
}

parameter {
  key         = "vpc_access.connector"
  description = "Set VPC Access Connector for the Cloud Run instance."
  type        = "string"
  required    = false
}

parameter {
  key         = "vpc_access.egress"
  description = "Set VPC egress. Supported values are 'all' and 'private-ranges-only'."
  type        = "string"
  required    = false
}

