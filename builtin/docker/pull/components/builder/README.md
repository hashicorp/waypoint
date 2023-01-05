<!-- This file was generated via `make gen/integrations-hcl` -->
Use an existing, pre-built Docker image.

This builder will automatically inject the Waypoint entrypoint. You
can disable this with the "disable_entrypoint" configuration.

If you wish to rename or retag an image, use this along with the
"docker" registry option which will rename/retag the image and then
push it to the specified registry.

If Docker isn't available (the Docker daemon isn't running or a DOCKER_HOST
isn't set), a daemonless solution will be used instead.

If "disable_entrypoint" is set to true and the Waypoint configuration
has no registry, this builder will not physically pull the image. This enables
Waypoint to work in environments where the image is built outside of Waypoint
(such as in a CI pipeline).

### Interface

- Input: **component.Source**
- Output: **docker.Image**

### Examples

```hcl
build {
  use "docker-pull" {
    image = "gcr.io/my-project/my-image"
    tag   = "abcd1234"
  }
}
```

