<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy the application into an ECS cluster on AWS.

### Interface

- Input: **docker.Image**
- Output: **ecs.Deployment**

### Examples

```hcl
deploy {
  use "aws-ecs" {
    region = "us-east-1"
    memory = 512
  }
}
```

