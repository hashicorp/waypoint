parameter {
  key         = "region"
  description = <<EOT
The AWS region to search in.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "filters"
  description = <<EOT
DescribeImage specific filters to search with.

The filters are always name => [value].
EOT
  type        = "map of string to list of string"
  required    = false
}

parameter {
  key         = "name"
  description = <<EOT
The name of the AMI to search for, supports wildcards.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "owners"
  description = <<EOT
The set of AMI owners to restrict the search to.
EOT
  type        = "list of string"
  required    = false
}
