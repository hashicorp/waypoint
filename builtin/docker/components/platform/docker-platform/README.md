<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy a container to Docker, local or remote.

### Interface

- Input: **docker.Image**
- Output: **docker.Deployment**

### Examples

```hcl
deploy {
  use "docker" {
	command      = ["ps"]
	service_port = 3000
	static_environment = {
	  "environment": "production",
	  "LOG_LEVEL": "debug"
	}
  }
}
```

