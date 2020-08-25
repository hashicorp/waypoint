project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env" = "dev"
  }

  build "pack" {
    registry "docker" {
      image = "gcr.io/waypoint-286812/wpmini"
      tag = "latest"
    }
  }

  deploy "google-cloud-run" {
      project = "waypoint-286812"
      region = "europe-north1"

      port = 5000

      env = {
        "NAME": "Nic"
      }

      capacity {
        memory = "128Mi"
        cpu_count = 1
        max_requests_per_container = 10
        request_timeout = 300
      }

      auto_scaling {
        max = 10
      }
  }

  release "google-cloud-run" { }
  
}