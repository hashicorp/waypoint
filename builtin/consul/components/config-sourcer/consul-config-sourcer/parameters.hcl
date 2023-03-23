# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

parameter {
  key         = "key"
  description = <<EOT
The KV path to retrieve.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "allow_stale"
  description = <<EOT
Whether to perform a stale query for retrieving the KV data.

If not set this will default to true. It must explicitly be set to false in order to use consistent queries.
EOT
  type        = "bool"
  required    = false
}

parameter {
  key         = "datacenter"
  description = <<EOT
The datacenter to load the KV value from.

If not specified then it will default to the plugin's global datacenter configuration. If that is also not specified then Consul will default the datacenter like it would any other request.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "namespace"
  description = <<EOT
The namespace to load the KV value from.

If not specified then it will default to the plugin's global namespace configuration. If that is also not specified then Consul will default the namespace like it would any other request.
EOT
  type        = "string"
  required    = false
}

parameter {
  key         = "partition"
  description = <<EOT
The partition to load the KV value from.

If not specified then it will default to the plugin's global partition configuration. If that is also not specified then Consul will default the partition like it would any other request.
EOT
  type        = "string"
  required    = false
}

