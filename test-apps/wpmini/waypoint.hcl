project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env"     = "dev"
  }

  build {
    use "pack" {}
    registry {
      use "docker" {
        image = "waypoint-example.local/wpmini"
        tag   = "latest"
        local = true
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
      load_balancer = true
      port          = 8080
    }
  }
}
