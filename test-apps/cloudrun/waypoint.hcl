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
        image = "gcr.io/waypoint-286812/wpmini"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "google-cloud-run" {
      project  = "waypoint-286812"
      location = "europe-north1"

      port = 5000

      static_environment = {
        "NAME" : "Nic"
      }

      capacity {
        memory                     = 128
        cpu_count                  = 2
        max_requests_per_container = 10
        request_timeout            = 300
      }

      auto_scaling {
        max = 10
      }
    }
  }

  release {
    use "google-cloud-run" {}
  }

}