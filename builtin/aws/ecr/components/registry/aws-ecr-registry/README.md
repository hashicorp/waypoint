<!-- This file was generated via `make gen/integrations-hcl` -->
Store a docker image within an Elastic Container Registry on AWS.

### Interface

- Input: **docker.Image**
- Output: **ecr.Image**

### Examples

```hcl
registry {
    use "aws-ecr" {
      region = "us-east-1"
      tag = "latest"
    }
}
```

### Mappers

#### Allow an ECR Image to be used as a standard docker.Image

- Input: **ecr.Image**
- Output: **docker.Image**

