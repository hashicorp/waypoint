package aci

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Releaser is the ReleaseManager implementation for Azure ACI.
type Releaser struct {
	config ReleaserConfig
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// Release deploys an image to GCR.
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	targets []component.ReleaseTarget,
) (*Release, error) {
	if len(targets) > 1 {
		return nil, fmt.Errorf(
			"The 'aci' release manager does not support traffic splitting.")
	}

	var result Release

	// Get the deployment
	var deploy Deployment
	if err := component.ProtoAnyUnmarshal(targets[0].Deployment, &deploy); err != nil {
		return nil, err
	}

	st := ui.Status()
	defer st.Close()

	containerGroupsClient, err := containerInstanceGroupsClient()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	containerGroup, err := containerGroupsClient.Get(ctx, deploy.ContainerGroup.ResourceGroup, deploy.ContainerGroup.Name)
	log.Trace("checking if container group already exists", "containergroup")
	st.Update("Checking if container group is already created")
	if err != nil {
		if containerGroup.StatusCode == 404 {
			log.Trace("container group not found", "containergroup")
		}

		return nil, err
	}

	log.Trace("container group is found", "containergroup")

	var containerGroupResource containerinstance.Resource
	containerGroupResource.Tags = containerGroup.Tags

	// Set tags
	containerGroupResource.Tags["name"] = to.StringPtr(deploy.Id)

	// Update
	log.Info("updating a pre-existing container group with container group tags")
	st.Update("Updating the container group")

	// We use update instead of createorupdate to allow for tag updates
	containerGroupResult, err := containerGroupsClient.Update(ctx, deploy.ContainerGroup.ResourceGroup, deploy.ContainerGroup.Name, containerGroupResource)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	ports := *containerGroupResult.IPAddress.Ports
	var portInt32 int32
	if len(ports) > 0 {
		port := ports[0]
		portInt32 = to.Int32(port.Port)
	}
	result.Url = fmt.Sprintf("http://%s:%d", to.String(containerGroupResult.IPAddress.Fqdn), portInt32)

	return &result, nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
}

func (r *Release) URL() string { return r.Url }

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Release        = (*Release)(nil)
)
