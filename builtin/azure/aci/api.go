package aci

import (
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func containerInstanceGroupsClient() (*containerinstance.ContainerGroupsClient, error) {

	// get subscription id from the environment
	// TODO: mishra to look into why the authorizer cannot be used
	// to craete a new container group client as it is retrieving
	// credentials from the environment already.
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription isn't set. Use AZURE_SUBSCRIPTION_ID environment variable to set it")
	}

	// create a container groups client
	containerGroupsClient := containerinstance.NewContainerGroupsClient(subscriptionID)

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
