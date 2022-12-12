## docker-ref (builder)

Use an existing, pre-built Docker image without modifying it.

### Interface

- Input: **component.Source**
- Output: **docker.Image**

### Examples

```hcl
build {
  use "docker-ref" {
    image = "gcr.io/my-project/my-image"
    tag   = "abcd1234"
  }
}
```
