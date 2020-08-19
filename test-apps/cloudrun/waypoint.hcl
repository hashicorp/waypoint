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
  }
}
