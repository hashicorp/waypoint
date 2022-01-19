variable "teeth" {
  default = configdynamic("static", {
    value = "hello"
  })
  type = string
}
