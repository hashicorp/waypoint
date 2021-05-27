project = "foo"

app "web" {
    build {}
    deploy {}
}

variable "yellow" {
    default = "balloon"
    type = int
}