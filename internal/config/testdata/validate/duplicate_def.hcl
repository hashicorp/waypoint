project = "foo"

app "web" {
  config {
    env = {
      static = "hello"
    }
  }

  build {}

  deploy {}
}

variable "dinosaur" {
  default = "longneck"
  type = string
}

variable "dinosaur" {
  default = "bigtooth"
  type = string
}