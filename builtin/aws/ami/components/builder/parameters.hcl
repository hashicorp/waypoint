# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "filters"
  description = "DescribeImage specific filters to search with\nthe filters are always name => [value]"
  type        = "map of string to list of string"
  required    = false
}

parameter {
  key         = "name"
  description = "the name of the AMI to search for, supports wildcards"
  type        = "string"
  required    = false
}

parameter {
  key         = "owners"
  description = "the set of AMI owners to restrict the search to"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "region"
  description = "the AWS region to search in"
  type        = "string"
  required    = true
}

