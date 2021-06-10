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
  description = "This is my description"
  type = string
}