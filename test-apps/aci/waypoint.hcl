project = "aci"

app "wpaci" {
  labels = {
    "service" = "wpaci",
    "env"     = "dev"
  }

  build {
    use "pack" {}
    registry {
      use "docker" {
        image = "nicholasjackson/test_pack"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "azure-container-instance" {
      resource_group="minecraft"
      location = "westeurope"

      ports = [8080]

      static_environment = {
        "NAME": "Nic"
      }

      capacity {
        memory = "1024"
        cpu_count = 4
      }

      volume {
        name = "vol1"
        path = "/consul"
        read_only = true

        git_repo {
          repository = "https://github.com/hashicorp/consul"
          revision = "v1.8.3"
        }
      }
    }
  }
}