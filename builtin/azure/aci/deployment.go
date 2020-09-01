package aci

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-11-01/subscriptions"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/waypoint/sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)

func (d *Deployment) containerInstanceGroupsClient() (*containerinstance.ContainerGroupsClient, error) {
	// create a container groups client
	containerGroupsClient := containerinstance.NewContainerGroupsClient(d.ContainerGroup.SubscriptionId)

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	// todo: mishra to create an option to set polling time out in the waypoint configuration.
	// this is to fix long provisioning times for aci.
	containerGroupsClient.PollingDuration = 60 * time.Minute
	containerGroupsClient.Authorizer = authorizer

	return &containerGroupsClient, nil
}

func (d *Deployment) getLocations(ctx context.Context) ([]string, error) {
	// create a account client
	subscriptionClient := subscriptions.NewClient()

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	subscriptionClient.PollingDuration = 60 * time.Minute
	subscriptionClient.Authorizer = authorizer
	if err != nil {
		return nil, fmt.Errorf("Unable to create subscriptions client: %s", err)
	}

	llr, err := subscriptionClient.ListLocations(ctx, d.ContainerGroup.SubscriptionId)
	if err != nil {
		return nil, fmt.Errorf("Unable to list locations for this subscription: %s", err)
	}

	locs := []string{}
	for _, l := range *llr.Value {
		locs = append(locs, *l.Name)
	}

	return locs, nil
}

func (d *Deployment) getContainerGroup(ctx context.Context) (containerinstance.ContainerGroup, error) {
	c, err := d.containerInstanceGroupsClient()
	if err != nil {
		return containerinstance.ContainerGroup{}, fmt.Errorf("Unable to create Container Groups client: %s", err)
	}

	return c.Get(ctx, d.ContainerGroup.ResourceGroup, d.ContainerGroup.Name)
}

func (d *Deployment) createOrUpdate(ctx context.Context, cg containerinstance.ContainerGroup) (containerinstance.ContainerGroup, error) {
	c, err := d.containerInstanceGroupsClient()
	if err != nil {
		return containerinstance.ContainerGroup{}, fmt.Errorf("Unable to create Container Groups client: %s", err)
	}

	response, err := c.CreateOrUpdate(ctx, d.ContainerGroup.ResourceGroup, d.ContainerGroup.Name, cg)
	if err != nil {
		return containerinstance.ContainerGroup{}, fmt.Errorf("Unable to create or update container group: %s", err)
	}

	// Wait for the container group to be created
	err = response.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return containerinstance.ContainerGroup{}, fmt.Errorf("Error waiting for container group creation to complete: %s", err)
	}

	return response.Result(*c)
}
