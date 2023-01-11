## AWS ECR Pull

The AWS ECR Pull plugin references an existing image, if found, in an AWS
[Elastic Container Registry](https://aws.amazon.com/ecr/getting-started/).
The image information can be used to push an image to a new registry, or be 
deployed to [AWS ECS](https://aws.amazon.com/ecs/getting-started/).

### Components

1. [Builder](./components/builder/README.md)

### Related Plugins

1. [AWS ECR](../README.md)
2. [AWS Lambda](../../lambda/README.md)
2. [Docker](../../../docker/README.md)