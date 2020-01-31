variable "number" {
    type = number
}

output "double" {
    value = var.number * 2
}
