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
  description = "This is my description"
  type = string
}