variable "teeth" {
  default = configdynamic("static", {
    json = <<-EOF
      {"k1":"v1", "k2":"v2"}
EOF
  })
  type = map(string)
}
