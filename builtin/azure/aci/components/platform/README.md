<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy a container to Azure Container Instances.

### Interface

- Input: **docker.Image**
- Output: **aci.Deployment**

### Examples

```hcl
deploy {
  use "azure-container-instance" {
    resource_group = "resource-group-name"
    location       = "westus"
    ports          = [8080]

    capacity {
      memory = "1024"
      cpu_count = 4
    }

    volume {
      name = "vol1"
      path = "/consul"
      read_only = true

      git_repo {
        repository = "https://github.com/hashicorp/consul"
        revision = "v1.8.3"
      }
    }
  }
}
```

