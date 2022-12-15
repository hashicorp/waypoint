## pack (builder)

Create a Docker image using CloudNative Buildpacks.

**This plugin must either be run via Docker or inside an ondemand runner**.

### Interface

- Input: **component.Source**
- Output: **pack.Image**

### Examples

```hcl
build {
  use "pack" {
	builder     = "heroku/buildpacks:20"
	disable_entrypoint = false
  }
}
```

### Mappers

#### Allow pack images to be used as normal docker images

- Input: **pack.Image**
- Output: **docker.Image**
