project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env" = "dev"
  }

  build "pack" {
    registry "docker" {
      image = "waypoint-example.local/wpmini"
      tag = "latest"
      local = true
    }
  }

  deploy "kubernetes" {
    probe_path = "/"
  }

  release "kubernetes" {
    load_balancer = true
    port = 8080
  }
}
