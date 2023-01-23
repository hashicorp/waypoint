integration {
  name        = "Terraform Cloud"
  description = "The Terraform Cloud plugin reads Terraform state outputs from Terraform Cloud."
  identifier  = "waypoint/terraform-cloud"
  components  = ["config-sourcer"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
