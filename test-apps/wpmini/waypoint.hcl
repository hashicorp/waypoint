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
