parameter {
  key         = "config_key"
  description = <<EOT
Return a value from the config source configuration.
This looks up the given key in the `values` configuration for the config sourcer. This can be used to actually test pulling a dynamic value, except the dynamic value is just Waypoint server-stored. This is useful for learning about and experimenting with config sourcer configuration with Waypoint.

EOT
  type        = "string"
  required    = false

}

parameter {
  key         = "static_value"
  description = <<EOT
A static value to use for the dynamic configuration.
This just returns the value given as the dynamic configuration. This isn't very "dynamic" but it helps to exercise the full dynamic configuration code paths which can be useful for experimentation or testing This is not expected to be used in a real-world production system.

EOT
  type        = "string"
  required    = false

}

