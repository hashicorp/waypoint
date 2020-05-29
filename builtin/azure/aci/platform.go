package aci

import (
	"context"
	"encoding/base64"
	"fmt"
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
	result := &Deployment{
		ContainerGroup: &Deployment_ContainerGroup{
			ResourceGroup: p.config.ResourceGroup,
			Name:          src.App,
		},
	}

	containerGroupsClient, err := containerInstanceGroupsClient()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Create our env vars
	var env []containerinstance.EnvironmentVariable
	for k, v := range deployConfig.Env() {
		env = append(env, containerinstance.EnvironmentVariable{
			Name:  to.StringPtr(k),
			Value: to.StringPtr(v),
		})
	}

	// Add custom ports
	// note: aci currently doesn't support port mapping
	// the default exposed port is 80, but the application
	// cannot bind on it so we are exposing port 8080, we might want to
	// see other ways to expose the application like for example using a sidecar or add option to specify port.
	env = append(env, containerinstance.EnvironmentVariable{
		Name:  to.StringPtr("PORT"),
		Value: to.StringPtr("8080"),
	})

	// todo: mishra to look into the waypoint tagging scheme as aci doesn't allow
	// '<,>,%!!(MISSING),(MISSING)&,\\,?,/' characters in the tags.
	var tags = map[string]*string{
		"waypoint.hashicorp.com": to.StringPtr(time.Now().UTC().Format(time.RFC3339Nano)),
	}

	create := false
	containerGroup, err := containerGroupsClient.Get(ctx, result.ContainerGroup.ResourceGroup, result.ContainerGroup.Name)
	log.Trace("checking if container group already exists", "containergroup")
	st.Update("Checking if container group is already created")
	if err != nil {
		if containerGroup.StatusCode != 404 {
			return nil, err
		}

		log.Trace("container group not found", "containergroup")
		// Create a container group
		create = true
		containerGroup = containerinstance.ContainerGroup{
			Name:     &result.ContainerGroup.Name,
			Location: to.StringPtr("eastus"),
		}
	}

	if create != true {
		log.Info("container group provisioning state", to.String(containerGroup.ContainerGroupProperties.ProvisioningState))
		// There is already a provisioning operation that is creating or pending
		// we would want to error at this point.
		if containerGroup.ProvisioningState == to.StringPtr("Creating") || containerGroup.ProvisioningState == to.StringPtr("Pending") {
			return nil, status.Errorf(codes.Aborted, "container group is \"%s\", cannot create a new deployment", to.String(containerGroup.ProvisioningState))
		}
	}

	// Set container group properties, doesn't matter whether we are creating or updating
	// we are forcing a new revision.
	containerGroup.ContainerGroupProperties = &containerinstance.ContainerGroupProperties{
		IPAddress: &containerinstance.IPAddress{
			Type: containerinstance.Public,
			Ports: &[]containerinstance.Port{
				{
					Port:     to.Int32Ptr(8080),
					Protocol: containerinstance.TCP,
				},
			},
			DNSNameLabel: &result.ContainerGroup.Name,
		},
		OsType: containerinstance.Linux,
		Containers: &[]containerinstance.Container{
			{
				Name: to.StringPtr(result.ContainerGroup.Name),
				ContainerProperties: &containerinstance.ContainerProperties{
					Ports: &[]containerinstance.ContainerPort{
						{
							Port: to.Int32Ptr(8080),
						},
					},
					Image: to.StringPtr(img.GetImage()),
					Resources: &containerinstance.ResourceRequirements{
						Limits: &containerinstance.ResourceLimits{
							MemoryInGB: to.Float64Ptr(1),
							CPU:        to.Float64Ptr(1),
						},
						Requests: &containerinstance.ResourceRequests{
							MemoryInGB: to.Float64Ptr(1),
							CPU:        to.Float64Ptr(1),
						},
					},
					EnvironmentVariables: &env,
				},
			},
		},
	}

	containerGroup.Tags = tags

	if create {
		// Create the container group
		log.Info("creating the container group")
		st.Update("Creating new container group")
	} else {
		// Update
		log.Info("updating a pre-existing container group")
		st.Update("Updating the container group")
	}

	response, err := containerGroupsClient.CreateOrUpdate(ctx, result.ContainerGroup.ResourceGroup, result.ContainerGroup.Name, containerGroup)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// Wait for the container group to be created
	err = response.WaitForCompletionRef(ctx, containerGroupsClient.Client)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	containerGroupResult, err := response.Result(*containerGroupsClient)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// The container group is ready, we will set the id and URL
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	ports := *containerGroupResult.IPAddress.Ports
	var portInt32 int32
	if len(ports) > 0 {
		port := ports[0]
		portInt32 = to.Int32(port.Port)
	}
	result.Url = fmt.Sprintf("http://%s:%d", to.String(containerGroupResult.IPAddress.Fqdn), portInt32)

	// If we have tracing enabled we just dump the full container group as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, _ := containerGroupResult.MarshalJSON()
		log.Trace("container group JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}

	return result, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// ResourceGroup is the resource group to deploy to.
	ResourceGroup string `hcl:"resource_group,attr"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
