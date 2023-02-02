<!-- This file was generated via `make gen/integrations-hcl` -->
Push a Docker image to a Docker compatible registry.

### Interface

- Input: **docker.Image**
- Output: **docker.Image**

### Examples

```hcl
build {
  use "docker" {}
  registry {
    use "docker" {
      image = "hashicorp/http-echo"
      tag   = "latest"
    }
  }
}
```

