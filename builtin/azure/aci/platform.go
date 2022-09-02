package aci

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	reflect "reflect"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/go-autorest/autorest/to"
)

// Platform is the Platform implementation for Azure ACI.
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	c, ok := config.(*Config)
	if !ok {
		// this should never happen
		return fmt.Errorf("Invalid configuration, expected *cloudrun.Config, got %s", reflect.TypeOf(config))
	}

	return validateConfig(*c)
}

// Deploy deploys an image to ACI.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {

	// if there is no subscription id in the deployment config try and fetch it from the environment
	if p.config.SubscriptionID == "" {
		p.config.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}

	// if we do not have a subscription id, return an error
	if p.config.SubscriptionID == "" {
		return nil, status.Error(
			codes.FailedPrecondition,
			"Please set either your Azure subscription ID in the deployment config, or set the environment variable 'AZURE_SUBSCRIPTION_ID'",
		)
	}

	// Start building our deployment since we use this information
	deployment := &Deployment{
		ContainerGroup: &Deployment_ContainerGroup{
			ResourceGroup:  p.config.ResourceGroup,
			Name:           src.App,
			SubscriptionId: p.config.SubscriptionID,
		},
	}

	auth, err := deployment.authenticate(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			err.Error(),
		)
	}

	create := false
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	if p.config.Location == "" {
		// set the default location eastus
		p.config.Location = "eastus"
	}

	// validate that the region for the deployment is valid
	l, err := deployment.getLocations(ctx, auth)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to list locations for subscription: %s", err)
	}

	err = validateLocationAvailable(p.config.Location, l)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Create our env vars
	var env []containerinstance.EnvironmentVariable
	for k, v := range deployConfig.Env() {
		env = append(env, containerinstance.EnvironmentVariable{
			Name:  to.StringPtr(k),
			Value: to.StringPtr(v),
		})
	}

	// if we have a port we need to set the PORT env var so that the CEB binary can direct
	// traffic to the correct service
	if len(p.config.Ports) > 0 {
		env = append(env, containerinstance.EnvironmentVariable{
			Name:  to.StringPtr("PORT"),
			Value: to.StringPtr(fmt.Sprintf("%d", p.config.Ports[0])),
		})
	}

	log.Info("Checking if container group already exists", "containergroup", deployment.ContainerGroup.Name)
	st.Update("Checking if container group is already created")
	containerGroup, err := deployment.getContainerGroup(ctx, auth)
	if err != nil {
		if containerGroup.StatusCode != 404 {
			return nil, status.Errorf(codes.Internal, "Unable to check if container group already exists: %s", err)
		}

		// Create a new container group
		log.Info("Container group not found, creating new container group", "containergroup", deployment.ContainerGroup.Name)
		create = true
		containerGroup = containerinstance.ContainerGroup{
			Name:     &deployment.ContainerGroup.Name,
			Location: &p.config.Location,
		}
	}

	// If we have a container group already, can it be updated?
	if create != true {
		// There is already a provisioning operation that is creating or pending we should return an error
		// as it is not possible to update a container group in this state
		if *containerGroup.ProvisioningState == "Creating" || *containerGroup.ProvisioningState == "Pending" {
			log.Debug("Container group provisioning state", *containerGroup.ContainerGroupProperties.ProvisioningState)
			return nil, status.Errorf(codes.AlreadyExists, "Container group state is '%s', unable to create new deployment", *containerGroup.ProvisioningState)
		}

		// If we are updating check that we are not changing the region, this is not possible
		// for an existing container and will NOOP.
		if *containerGroup.Location != p.config.Location {
			return nil, status.Errorf(
				codes.InvalidArgument,
				`Waypoint is unable to change the location of an existing container instance as Azure does not allow the location of
existing container instances to be changed.
The container instance '%s' is currently deployed to the location '%s', but your configuration sets a location of '%s'.
To update the location you will need to manually destroy and recreate the resource.`,
				src.App,
				*containerGroup.Location,
				p.config.Location,
			)
		}
	}

	// Set container group properties, doesn't matter whether we are creating or updating
	// we are forcing a new revision.
	containerGroup.ContainerGroupProperties = &containerinstance.ContainerGroupProperties{
		IPAddress: &containerinstance.IPAddress{
			Type:         containerinstance.Public,
			DNSNameLabel: &deployment.ContainerGroup.Name,
		},
		OsType: containerinstance.Linux,
	}

	// Add the tags
	var tags = map[string]*string{
		"_waypoint_hashicorp_com_nonce": to.StringPtr(time.Now().UTC().Format(time.RFC3339Nano)),
	}
	containerGroup.Tags = tags

	// Add static environment variables
	for k, v := range p.config.StaticEnvVars {
		env = append(env, containerinstance.EnvironmentVariable{
			Name:  to.StringPtr(k),
			Value: to.StringPtr(v),
		})
	}

	// Add the managed identity if set
	if p.config.ManagedIdentity != "" {
		containerGroup.Identity = &containerinstance.ContainerGroupIdentity{
			Type: containerinstance.UserAssigned,
			UserAssignedIdentities: map[string]*containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue{
				p.config.ManagedIdentity: {},
			},
		}
	}

	// do we need to add registry credentials for auth?
	registryUser := os.Getenv("REGISTRY_USERNAME")
	registryPass := os.Getenv("REGISTRY_PASSWORD")
	if registryUser != "" && registryPass != "" {
		server := parseDockerServer(img.Image)

		containerGroup.ImageRegistryCredentials = &[]containerinstance.ImageRegistryCredential{
			{
				Server:   &server,
				Username: &registryUser,
				Password: &registryPass,
			},
		}
	}

	// define a container
	container := containerinstance.Container{
		Name: &deployment.ContainerGroup.Name,
		ContainerProperties: &containerinstance.ContainerProperties{
			Image:                &img.Image,
			EnvironmentVariables: &env,
			Resources:            &containerinstance.ResourceRequirements{},
		},
	}

	// Add the ports
	if len(p.config.Ports) > 0 {
		ports := []containerinstance.Port{}
		containerPorts := []containerinstance.ContainerPort{}
		for _, port := range p.config.Ports {
			ports = append(
				ports,
				containerinstance.Port{
					Port:     to.Int32Ptr(int32(port)),
					Protocol: containerinstance.TCP,
				},
			)

			containerPorts = append(
				containerPorts,
				containerinstance.ContainerPort{
					Port:     to.Int32Ptr(int32(port)),
					Protocol: containerinstance.ContainerNetworkProtocolTCP,
				},
			)
		}

		// set the ports on the container and the parent container group
		container.ContainerProperties.Ports = &containerPorts
		containerGroup.ContainerGroupProperties.IPAddress.Ports = &ports
	}

	// Set the capacity
	if p.config.Capacity != nil {
		requests := containerinstance.ResourceRequests{}
		limits := containerinstance.ResourceLimits{}

		if p.config.Capacity.Memory > 0 {
			requests.MemoryInGB = to.Float64Ptr(float64(p.config.Capacity.Memory) / 1024.0) // Azure units are 1GiB, our unit is 1MB
			limits.MemoryInGB = to.Float64Ptr(float64(p.config.Capacity.Memory) / 1024.0)

			container.ContainerProperties.Resources.Requests = &requests
			container.ContainerProperties.Resources.Limits = &limits
		}

		if p.config.Capacity.CPUCount > 0 {
			requests.CPU = to.Float64Ptr(float64(p.config.Capacity.CPUCount))
			limits.CPU = to.Float64Ptr(float64(p.config.Capacity.CPUCount))

			container.ContainerProperties.Resources.Requests = &requests
			container.ContainerProperties.Resources.Limits = &limits
		}
	}

	// If we have any volumes, add them
	if len(p.config.Volumes) > 0 {
		volumes := []containerinstance.Volume{}
		volumeMounts := []containerinstance.VolumeMount{}

		for _, v := range p.config.Volumes {
			func(v Volume) {
				vol := containerinstance.Volume{
					Name: &v.Name,
				}

				volMount := containerinstance.VolumeMount{
					Name:      &v.Name,
					MountPath: &v.Path,
					ReadOnly:  &v.ReadOnly,
				}

				if v.AzureFileShare != nil {
					vol.AzureFile = &containerinstance.AzureFileVolume{
						ShareName:          &v.AzureFileShare.Name,
						ReadOnly:           &v.ReadOnly,
						StorageAccountName: &v.AzureFileShare.StorageAccountName,
						StorageAccountKey:  &v.AzureFileShare.StorageAccountKey,
					}
				}

				if v.GitRepoVolume != nil {
					vol.GitRepo = &containerinstance.GitRepoVolume{
						Repository: &v.GitRepoVolume.Repository,
						Directory:  &v.GitRepoVolume.Directory,
						Revision:   &v.GitRepoVolume.Revision,
					}
				}

				volumes = append(volumes, vol)
				volumeMounts = append(volumeMounts, volMount)
			}(v)
		}

		containerGroup.ContainerGroupProperties.Volumes = &volumes
		container.VolumeMounts = &volumeMounts
	}

	// set the container on the container group
	containerGroup.Containers = &[]containerinstance.Container{container}

	if create {
		log.Info("Creating the container group")
		st.Update("Creating new container group")
	} else {
		log.Info("Updating a pre-existing container group")
		st.Update("Updating the container group")
	}

	containerGroupResult, err := deployment.createOrUpdate(ctx, auth, containerGroup)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to create or update container instance: %s", err)
	}

	// The container group is ready, we will set the id and URL
	id, err := component.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to generate ID for the deployment: %s", err)
	}

	deployment.Id = id

	ports := *containerGroupResult.IPAddress.Ports
	if len(ports) > 0 {
		// Only set the URL if there is a port
		deployment.Url = fmt.Sprintf("http://%s:%d", *containerGroupResult.IPAddress.Fqdn, *ports[0].Port)

		// Clear the status before we print the url
		st.Close()
		// Show the container group url
		ui.Output("\nURL: %s", deployment.Url, terminal.WithSuccessStyle())
	}

	// If we have tracing enabled we just dump the full container group as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, err := containerGroupResult.MarshalJSON()
		if err != nil {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}

		log.Trace("container group JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}

	return deployment, nil
}

// Config is the configuration structure for the Platform.
// In addition to HCL defined configuration the following environment variables
// are also valid
// AZURE_SUBSCRIPTION_ID = Subscription ID for your Azure account [required]
// REGISTRY_USERNAME = Username for container registry, required when using a private registry
// REGISTRY_PASSWORD = Password for container registry, required when using a private registry
type Config struct {
	// ResourceGroup is the resource group to deploy to.
	ResourceGroup string `hcl:"resource_group,attr"`

	// Region to deploy the container instance to
	Location string `hcl:"location,optional"`

	// Azure subscription id, if not set plugin will attempt to use the environment variable
	// AZURE_SUBSCRIPTION_ID
	SubscriptionID string `hcl:"subscription_id,optional"`

	// ManagedIdentity assigned to the container group
	// use managed identity to enable containers to access other
	// Azure resources
	// (https://docs.microsoft.com/en-us/azure/container-instances/container-instances-managed-identity#:~:text=Enable%20a%20managed%20identity&text=Azure%20Container%20Instances%20supports%20both,or%20both%20types%20of%20identities.)
	// Note: ManagedIdentity can not be used to authorize Container Instances to pull from private Container registries in Azure
	ManagedIdentity string `hcl:"managed_identity,optional"`

	// Port the applications is listening on.
	Ports []int `hcl:"ports,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Capacity details for cloud run container.
	Capacity *Capacity `hcl:"capacity,block"`

	Volumes []Volume `hcl:"volume,block" validate:"dive"`
}

// RegistryCredentials are the user credentials needed to
// authenticate with a container registry.
type RegistryCredentials struct {
	Username string `hcl:"username"`
	Password string `hcl:"password"`
}

type Capacity struct {
	// Memory to allocate to the container specified in MB, min 512, max 16384.
	// Default value of 0 sets memory to 1536MB which is default container instance value
	Memory int `hcl:"memory,attr" validate:"eq=0|gte=512,lte=16384"`
	// CPUCount is the number CPUs to allocate to a container instance
	CPUCount int `hcl:"cpu_count,attr" validate:"gte=0,lte=4"`
}

// Volume defines a volume mount for the container.
// Supported types are Azure file share or GitHub repository
// Only one type can be set
type Volume struct {
	// Name of the Volume to mount
	Name string `hcl:"name,attr"`
	// Filepath where the volume will be mounted in the container
	Path string `hcl:"path,attr"`
	// Is the volume read only?
	ReadOnly bool `hcl:"read_only,attr"`
	// Details for  an Azure file share volume
	AzureFileShare *AzureFileShareVolume `hcl:"azure_file_share,block"`
	// Details for  an GitHub repo volume
	GitRepoVolume *GitRepoVolume `hcl:"git_repo,block"`
}

// AzureFileShareVolume allows you to mount an Azure container storage
// fileshare into the container
type AzureFileShareVolume struct {
	// Name of the FileShare in Azure storage
	Name string `hcl:"name,attr"`
	// Storage account name
	StorageAccountName string `hcl:"storage_account_name,attr"`
	// Storage account key to access the storage
	StorageAccountKey string `hcl:"storage_account_key,attr"`
}

// GitRepoVolume allows the mounting of a Git repository
// into the container
type GitRepoVolume struct {
	// GitHub repository to mount as a volume
	Repository string `hcl:"repository,attr" validate:"url"`
	// Branch, Tag or Commit SHA
	Revision string `hcl:"revision,optional"`
	// Directory name to checkout repo to, defaults to repository name
	Directory string `hcl:"directory,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy a container to Azure Container Instances")

	doc.Example(`
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
`)

	doc.SetField(
		"resource_group",
		"the resource group to deploy the container to",
	)

	doc.SetField(
		"location",
		"the resource location to deploy the container instance to",
	)

	doc.SetField(
		"subscription_id",
		"the Azure subscription id",
		docs.Summary("if not set uses the environment variable AZURE_SUBSCRIPTION_ID"),
		docs.EnvVar("AZURE_SUBSCRIPTION_ID"),
	)

	doc.SetField(
		"managed_identity",
		"the managed identity assigned to the container group",
	)

	doc.SetField(
		"ports",
		"the ports the container is listening on, the first port in this list will be used by the entrypoint binary to direct traffic to your application",
	)

	doc.SetField(
		"static_environment",
		"environment variables to control broad modes of the application",
		docs.Summary(
			"environment variables that are meant to configure the application in a static",
			"way. This might be control an image that has multiple modes of operation,",
			"selected via environment variable. Most configuration should use the waypoint",
			"config commands.",
		),
	)

	doc.SetField(
		"capacity",
		"the capacity details for the container",

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"memory",
				"memory to allocate the container specified in MB, min 1024, max based on resource availability of the region.",
				docs.Default("1024"),
			)

			doc.SetField(
				"cpu",
				"number of CPUs to allocate the container, min 1, max based on resource availability of the region.",
				docs.Default("1"),
			)
		}),
	)

	doc.SetField(
		"volume",
		"the volume details for a container",

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"name",
				"the name of the volume to mount into the container",
			)

			doc.SetField(
				"path",
				"the path to mount the volume to in the container",
			)

			doc.SetField(
				"read_only",
				"specify if the volume is read only",
			)

			doc.SetField(
				"azure_file_share",
				"the details for the Azure file share volume",
			)

			doc.SetField(
				"git_repo",
				"the details for GitHub repo to mount as a volume",
			)
		}),
	)

	doc.Input("docker.Image")
	doc.Output("aci.Deployment")

	return doc, nil
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Documented   = (*Platform)(nil)
)
