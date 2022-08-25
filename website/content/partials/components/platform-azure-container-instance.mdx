## azure-container-instance (platform)

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

### Required Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### capacity (category)

The capacity details for the container.

##### capacity.cpu

Number of CPUs to allocate the container, min 1, max based on resource availability of the region.

##### capacity.cpu_count

- Type: **int**

##### capacity.memory

Memory to allocate the container specified in MB, min 1024, max based on resource availability of the region.

- Type: **int**

#### resource_group

The resource group to deploy the container to.

- Type: **string**

#### volume (category)

The volume details for a container.

##### volume.azure_file_share

The details for the Azure file share volume.

- Type: **aci.AzureFileShareVolume**

##### volume.git_repo

The details for GitHub repo to mount as a volume.

- Type: **aci.GitRepoVolume**

##### volume.name

The name of the volume to mount into the container.

- Type: **string**

##### volume.path

The path to mount the volume to in the container.

- Type: **string**

##### volume.read_only

Specify if the volume is read only.

- Type: **bool**

### Optional Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### location

The resource location to deploy the container instance to.

- Type: **string**
- **Optional**

#### managed_identity

The managed identity assigned to the container group.

- Type: **string**
- **Optional**

#### ports

The ports the container is listening on, the first port in this list will be used by the entrypoint binary to direct traffic to your application.

- Type: **list of int**
- **Optional**

#### static_environment

Environment variables to control broad modes of the application.

Environment variables that are meant to configure the application in a static way. This might be control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.

- Type: **map of string to string**
- **Optional**

#### subscription_id

The Azure subscription id.

If not set uses the environment variable AZURE_SUBSCRIPTION_ID.

- Type: **string**
- **Optional**
- Environment Variable: **AZURE_SUBSCRIPTION_ID**

### Output Attributes

Output attributes can be used in your `waypoint.hcl` as [variables](/docs/waypoint-hcl/variables) via [`artifact`](/docs/waypoint-hcl/variables/artifact) or [`deploy`](/docs/waypoint-hcl/variables/deploy).

#### container_group

- Type: **aci.Deployment_ContainerGroup**

#### id

- Type: **string**

#### url

- Type: **string**
