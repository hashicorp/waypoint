package aci

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/terminal"

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

// Deploy deploys an image to ACI.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	dir *datadir.Component,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// Start building our deployment since we use this information
	deployment := &Deployment{
		ContainerGroup: &Deployment_ContainerGroup{
			ResourceGroup: p.config.ResourceGroup,
			Name:          src.App,
		},
	}

	// if there is no subscription id in the deployment config try and fetch it from the environment
	if p.config.SubscriptionID != "" {
		deployment.ContainerGroup.SubscriptionId = p.config.SubscriptionID
	} else {
		// try and fetch from environment vars
		deployment.ContainerGroup.SubscriptionId = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}

	// if we do not have a subscription id, return an error
	if deployment.ContainerGroup.SubscriptionId == "" {
		return nil, status.Error(
			codes.Aborted,
			"Please set either your Azure subscription ID in the deployment config, or set the environment variable 'AZURE_SUBSCRIPTION_ID'",
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
	l, err := deployment.getLocations(ctx)
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

	log.Info("Checking if container group already exists", "containergroup", deployment.ContainerGroup.Name)
	st.Update("Checking if container group is already created")
	containerGroup, err := deployment.getContainerGroup(ctx)
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

	imageRegistryCredentials := []containerinstance.ImageRegistryCredential{}

	// Set container group properties, doesn't matter whether we are creating or updating
	// we are forcing a new revision.
	containerGroup.ContainerGroupProperties = &containerinstance.ContainerGroupProperties{
		ImageRegistryCredentials: &imageRegistryCredentials,
		IPAddress: &containerinstance.IPAddress{
			Type:         containerinstance.Public,
			DNSNameLabel: &deployment.ContainerGroup.Name,
		},
		OsType: containerinstance.Linux,
		Containers: &[]containerinstance.Container{
			{
				Name: &deployment.ContainerGroup.Name,
				ContainerProperties: &containerinstance.ContainerProperties{
					Image:                &img.Image,
					EnvironmentVariables: &env,
					Resources:            &containerinstance.ResourceRequirements{},
				},
			},
		},
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
				p.config.ManagedIdentity: &containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue{},
			},
		}
	}

	// do we need to add registry credentials for auth?
	if p.config.RegistryCredentials != nil {
		server := parseDockerServer(img.Image)

		imageRegistryCredentials = append(imageRegistryCredentials, containerinstance.ImageRegistryCredential{
			Server:   &server,
			Username: &p.config.RegistryCredentials.Username,
			Password: &p.config.RegistryCredentials.Password,
		})
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

		containers := containerGroup.ContainerGroupProperties.Containers
		(*containers)[0].ContainerProperties.Ports = &containerPorts

		containerGroup.ContainerGroupProperties.IPAddress.Ports = &ports
	}

	// Set the capacity
	if p.config.Capacity != nil {
		containers := containerGroup.ContainerGroupProperties.Containers
		requests := containerinstance.ResourceRequests{}
		limits := containerinstance.ResourceLimits{}

		if p.config.Capacity.Memory > 0 {
			requests.MemoryInGB = to.Float64Ptr(float64(p.config.Capacity.Memory) / 1024.0) // Azure units are 1GiB, our unit is 1MB
			limits.MemoryInGB = to.Float64Ptr(float64(p.config.Capacity.Memory) / 1024.0)

			(*containers)[0].ContainerProperties.Resources.Requests = &requests
			(*containers)[0].ContainerProperties.Resources.Limits = &limits
		}

		if p.config.Capacity.CPUCount > 0 {
			requests.CPU = to.Float64Ptr(float64(p.config.Capacity.CPUCount))
			limits.CPU = to.Float64Ptr(float64(p.config.Capacity.CPUCount))

			(*containers)[0].ContainerProperties.Resources.Requests = &requests
			(*containers)[0].ContainerProperties.Resources.Limits = &limits
		}
	}

	// If we have any volumes, add them
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
	(*containerGroup.Containers)[0].VolumeMounts = &volumeMounts

	if create {
		log.Info("Creating the container group")
		st.Update("Creating new container group")
	} else {
		log.Info("Updating a pre-existing container group")
		st.Update("Updating the container group")
	}

	containerGroupResult, err := deployment.createOrUpdate(ctx, containerGroup)
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
	ManagedIdentity string `hcl:"managed_identity,attr"`

	// RegistryCredentials allow you to set the username and password
	// in the instance the image to deploy is in a private repository and
	// requires authentication.
	RegistryCredentials *RegistryCredentials `hcl:"registry_credentials,block"`

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

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
