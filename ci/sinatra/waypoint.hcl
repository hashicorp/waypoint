project = "sinatra"

app "sinatra" {
  build "pack" {
  }

  registry "docker" {
    image = "localhost:5000/sinatra"
    tag = "latest"
  }
  
  deploy "kubernetes" {
  }

  release "kubernetes" {
    load_balancer = false
    node_port = -1
  }
}
