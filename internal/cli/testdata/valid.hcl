project = "test-project"
app "web" {
    build {
        use "docker" {}
    }

    # Deploy to Docker
    deploy {
        use "docker" {}
    }
}

variable "port" {
  type = number
  default = 1
  # default = 2
}