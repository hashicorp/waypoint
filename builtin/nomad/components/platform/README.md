<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy to a nomad cluster as a service using Docker.

### Interface

- Input: **docker.Image**
- Output: **nomad.Deployment**

### Examples

```hcl
deploy {
        use "nomad" {
          region = "global"
          datacenter = "dc1"
          auth {
            username = "username"
            password = "password"
          }
          static_environment = {
            "environment": "production",
            "LOG_LEVEL": "debug"
          }
          service_port = 3000
          replicas = 1
        }
}
```

