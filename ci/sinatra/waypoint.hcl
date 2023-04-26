project = "sinatra"

app "sinatra" {
  build {
    use "pack" {}

    registry {
      use "docker" {
        image = "registry.localhost:5000/sinatra"
        tag   = "latest"
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
      node_port     = 30000 // can only be 30000-32767 in k8s
    }
  }
}
