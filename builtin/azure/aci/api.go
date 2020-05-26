package aci

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func containerInstanceGroupsClient() (*containerinstance.ContainerGroupsClient, error) {

	// get subscription id from the environment
	// TODO: mishra to look into why the authorizer cannot be used
	// to craete a new container group client as it is retrieving
	// credentials from the environment already.
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// create a container groups client
	containerGroupsClient := containerinstance.NewContainerGroupsClient(subscriptionID)

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	containerGroupsClient.Authorizer = authorizer

	return &containerGroupsClient, nil
}
