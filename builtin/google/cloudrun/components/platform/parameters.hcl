parameter {
  key         = "auto_scaling"
  description = <<EOT
Configuration to control the auto scaling parameters for Cloud Run.
EOT
  type        = "category"
  required    = true
}

parameter {
  key         = "auto_scaling.max"
  description = <<EOT
Maximum number of Cloud Run instances. When the maximum requests per container is exceeded, Cloud Run will create an additional container instance to handle load.
This parameter controls the maximum number of instances that can be created.
EOT
  type        = "int"
  required    = true
}

parameter {
  key         = "capacity"
  description = <<EOT
CPU, Memory, and resource limits for each Cloud Run instance.
EOT
  type        = "category"
  required    = true
}

parameter {
  key         = "capacity.cpu_count"
  description = <<EOT
Number of CPUs to allocate the Cloud Run instance, min 1, max 2.
EOT
  type        = "int"
  required    = true
}

parameter {
  key         = "capacity.max_requests_per_container"
  description = <<EOT
Maximum number of concurrent requests each instance can handle. When the maximum requests are exceeded, Cloud Run will create an additional instance.
EOT
  type        = "int"
  required    = true
}

parameter {
  key         = "capacity.memory"
  description = <<EOT
Memory to allocate the Cloud Run instance specified in MB, min 128, max 4096.
EOT
  type        = "int"
  required    = true
}

parameter {
  key         = "capacity.request_timeout"
  description = <<EOT
Maximum time a request can take before timing out, max 900.
EOT
  type        = "int"
  required    = true
}

parameter {
  key         = "location"
  description = <<EOT
GCP location, e.g. europe-north-1.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "project"
  description = <<EOT
GCP project ID where the Cloud Run instance will be deployed.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "cloudsql_instances"
  description = <<EOT
Specify list of CloudSQL instances that the Cloud Run instance will have access to.
EOT
  type        = "list of string"
  required    = false
}

parameter {
  key         = "port"
  description = <<EOT
The port your application listens on.
EOT
  type        = "int"
  required    = false
}

parameter {
  key         = "service_account_name"
  description = <<EOT
Specify a service account email that Cloud Run will use to run the service. You must have the `iam.serviceAccounts.actAs` permission on the service account.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "static_environment"
  description = <<EOT
Additional environment variables to be added to the Cloud Run instance.
EOT
  type        = "map of string to string"
  required    = false
}

parameter {
  key         = "unauthenticated"
  description = <<EOT
Is public unauthenticated access allowed for the Cloud Run instance?.
EOT
  type        = "bool"
  required    = false
}

parameter {
  key         = "vpc_access"
  description = <<EOT
VPCAccess details.
EOT
  type        = "category"
  required    = false
}

parameter {
  key         = "vpc_access.connector"
  description = <<EOT
Set VPC Access Connector for the Cloud Run instance.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "vpc_access.egress"
  description = <<EOT
Set VPC egress. Supported values are 'all' and 'private-ranges-only'.
EOT
  type        = "string"
  required    = false
}

