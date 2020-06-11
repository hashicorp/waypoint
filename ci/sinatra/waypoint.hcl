project = "sinatra"

app "sinatra" {
  build "pack" {
  }

  registry "docker" {
    image = "waypoint.local/sinatra"
    tag = "latest"
    local = true
  }
  
  deploy "kubernetes" {
  }

  release "kubernetes" {
    load_balancer = false
    node_port = -1
  }
}
