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

variable "bees" {
  default = "buzz"
  type = bool
}