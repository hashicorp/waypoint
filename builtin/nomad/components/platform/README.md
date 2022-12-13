## nomad (platform)

Deploy to a nomad cluster as a service using docker.

### Interface

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
