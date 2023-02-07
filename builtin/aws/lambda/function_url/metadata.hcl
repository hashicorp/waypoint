integration {
  name        = "Lambda Function URL"
  description = "The AWS Lambda Function URL plugin releases a function deployed with the AWS Lambda plugin."
  identifier  = "waypoint/lambda-function-url"
  components  = ["release-manager"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
