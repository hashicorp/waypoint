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

      capacity {
        memory = "256Mi"
        cpu_count = 2
        max_requests_per_container = 10
        request_timeout = 300
      }
  }
}
