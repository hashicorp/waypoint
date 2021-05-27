project = "foo"

app "web" {
    build {}
    deploy {}
}

variable "bees" {
  default = "buzz"
  description = "This is my description"
  type = string
}

variable "dinosaur" {
  default = "longneck"
  type = string
}