project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env" = "dev"
  }

  build "pack" {}

  deploy "docker" {}
}
