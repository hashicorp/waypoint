project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env"     = "dev"
  }

  build {
    use "pack" {}
    registry {
      use "docker" {
        image = "jacksonnic.azurecr.io/wpmini"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "azure-aci" {
      resource_group="minecraft"
    }
  }
}
