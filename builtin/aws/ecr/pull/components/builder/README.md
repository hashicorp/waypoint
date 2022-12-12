## aws-ecr-pull (builder)

Use an existing, pre-built AWS ECR image.

This builder attempts to find an image by repository and tag in the
specified region. If found, it will pass along the image information
to the next step.

This builder will not modify the image.

If you wish to rename or retag an image, please use the "docker-pull" component
in conjunction with the "aws-ecr" registry option.

### Interface

- Input: **component.Source**
- Output: **ecr.Image**

### Examples

```hcl
build {
  use "aws-ecr-pull" {
    region     = "us-east-1"
    repository = "deno-http"
    tag        = "latest"
  }
}
```
