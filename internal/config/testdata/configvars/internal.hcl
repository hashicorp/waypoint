project = "p"

config {
  internal = {
    value = "V"
  }

  env = {
    "direct"       = config.internal.value
    "interpolated" = "value: ${config.internal.value}"
  }
}
