// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aci

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-11-01/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)
var locations []string // accessible locations for the current account

func (d *Deployment) containerInstanceGroupsClient(auth autorest.Authorizer) (*containerinstance.ContainerGroupsClient, error) {
	// create a container groups client
	containerGroupsClient := containerinstance.NewContainerGroupsClient(d.ContainerGroup.SubscriptionId)
	containerGroupsClient.Authorizer = auth

	return &containerGroupsClient, nil
}

// init sets up the authorizer and fetches the locations
func (d *Deployment) authenticate(ctx context.Context, log hclog.Logger) (autorest.Authorizer, error) {
	// create an authorizer from env vars or Azure Managed Service Identity
	//authorizer, err := auth.NewAuthorizerFromCLI()

	// first try and create an environment
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Warn("unable to create subscriptions client", "error", err)
	}

	// we need to timeout this request as this request never fails when we have
	// invalid credentials
	timeoutContext, cf := context.WithTimeout(ctx, 15*time.Second)
	defer cf()

	_, err = d.getLocations(timeoutContext, authorizer)
	if err == nil {
		return authorizer, nil
	}

	timeoutContext, cf2 := context.WithTimeout(ctx, 15*time.Second)
	defer cf2()

	// the environment variable auth has failed fall back to CLI auth
	log.Info("attempting CLI auth")
	authorizer, err = auth.NewAuthorizerFromCLI()
	if err != nil {
		return authorizer, err
	}
	_, err = d.getLocations(timeoutContext, authorizer)
	if err == nil {
		return authorizer, nil
	}

	return nil, fmt.Errorf(
		"unable to authenticate with the Azure API, ensure you have your credentials set as environment variables, " +
			"or you have logged in using the 'az' command line tool",
	)
}

func (d *Deployment) getLocations(ctx context.Context, auth autorest.Authorizer) ([]string, error) {
	// create a account client
	subscriptionClient := subscriptions.NewClient()
	subscriptionClient.Authorizer = auth

	llr, err := subscriptionClient.ListLocations(ctx, d.ContainerGroup.SubscriptionId)
	if err != nil {
		return nil, fmt.Errorf("Unable to list locations for this subscription: %s", err)
	}

	locs := []string{}
	for _, v := range *llr.Value {
		locs = append(locs, *v.Name)
	}

	return locs, nil
}

func (d *Deployment) getContainerGroup(ctx context.Context, auth autorest.Authorizer) (containerinstance.ContainerGroup, error) {
	c, err := d.containerInstanceGroupsClient(auth)
	if err != nil {
		return containerinstance.ContainerGroup{}, fmt.Errorf("Unable to create Container Groups client: %s", err)
	}

	return c.Get(ctx, d.ContainerGroup.ResourceGroup, d.ContainerGroup.Name)
}

func (d *Deployment) createOrUpdate(ctx context.Context, auth autorest.Authorizer, cg containerinstance.ContainerGroup) (containerinstance.ContainerGroup, error) {
	c, err := d.containerInstanceGroupsClient(auth)
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
