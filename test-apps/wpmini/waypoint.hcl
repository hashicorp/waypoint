project = "wpmini"

app "wpmini" {
  build "pack" {
    registry "docker" {
      image = "localhost:5000/wpmini"
      tag = "latest"
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

server {
  address = "localhost:9701"
  address_internal = "waypoint:9701"
  insecure = true
  require_auth = false
}
