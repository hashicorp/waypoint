variable "rate" {
  default = configdynamic("vault", {})
  type = number
}
