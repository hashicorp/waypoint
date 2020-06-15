project = "sinatra"

app "sinatra" {
  build "pack" {
    registry "docker" {
      image = "localhost:5000/sinatra"
      tag = "latest"
    }
  }

  deploy "kubernetes" {
    probe_path = "/"
  }

  release "kubernetes" {
    load_balancer = false
    node_port = -1
  }
}
