<!-- This file was generated via `make gen/integrations-hcl` -->
Build a Docker image from a Dockerfile.

If a Docker server is available (either locally or via environment variables
such as "DOCKER_HOST"), then "docker build" will be used to build an image
from a Dockerfile.

### Dockerless Builds

Many hosted environments, such as Kubernetes clusters, don't provide access
to a Docker server. In these cases, it is desirable to perform what is called
a "dockerless" build: building a Docker image without access to a Docker
daemon. Waypoint supports dockerless builds.

Waypoint performs Dockerless builds by leveraging
[Kaniko](https://github.com/GoogleContainerTools/kaniko)
within on-demand launched runners. This should work in all supported
Waypoint installation environments by default and you should not have
to specify any additional configuration.

### Interface

- Output: **docker.Image**

### Examples

```hcl
build {
  use "docker" {
	buildkit    = false
	disable_entrypoint = false
  }
}
```

