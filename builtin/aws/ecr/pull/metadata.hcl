integration {
  name        = "AWS ECR Pull"
  description = "The AWS ECR Pull plugin references an existing image, if found, in an AWS Elastic Container Registry. The image information can be used to push an image to a new registry, or be deployed to AWS ECS."
  identifier  = "waypoint/aws-ecr-pull"
  components  = ["builder"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
