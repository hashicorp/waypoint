integration {
  name        = "AWS AMI"
  description = "The AWS AMI plugin searches for an returns an existing AMI, to be deployed as an EC2."
  identifier  = "waypoint/aws-ami"
  components  = ["builder"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
