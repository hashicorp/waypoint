project = "sinatra"

app "sinatra" {
  build {
    use "pack" {}

    registry {
      use "docker" {
        image = "localhost:5000/sinatra"
        tag = "latest"
      }
    }
  }

  deploy {
    use "kubernetes" {
      probe_path = "/"
    }
  }

  release {
    use "kubernetes" {
      load_balancer = false
      node_port = -1
    }
  }
}
