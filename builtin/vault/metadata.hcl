integration {
  name        = "Vault"
  description = "The Vault plugin reads configuration values from Vault."
  identifier  = "waypoint/vault"
  components  = ["config-sourcer"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
