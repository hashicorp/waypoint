integration {
  name        = "Lambda Function URL"
  description = "The AWS Lambda Function URL plugin releases a function deployed with the AWS Lambda plugin."
  identifier  = "waypoint/lambda-function-url"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "release-manager"
    name = "Lambda Function URL Release Manager"
    slug = "lambda-function-url-release-manager"
  }
}
