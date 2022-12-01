package container

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (p *Platform) resourceContainerCreate(
	ctx context.Context,
	containerDeployment *Container,
	log hclog.Logger,
	st terminal.Status,
	deployConfig *component.DeploymentConfig,
	state *Resource_Container,
	img *docker.Image,
	scwContainerAPI *containerSDK.API,
) error {
	// TODO clear steps, add logging
	st.Step(terminal.StatusOK, "Checking for existing container")

	var container *containerSDK.Container
	var err error

	if containerDeployment.Id != "" {
		st.Update("Loading existing container")
		container, err = scwContainerAPI.GetContainer(&containerSDK.GetContainerRequest{
			Region:      scw.Region(containerDeployment.Region),
			ContainerID: containerDeployment.Id,
		})
		if err != nil {
			return fmt.Errorf("failed to find already existing container: %w", err)
		}
		st.Step(terminal.StatusOK, "Loading existing container")
		// TODO update container with new config
	} else {
		st.Update("Creating container")

		req := &containerSDK.CreateContainerRequest{
			Region:                     scw.Region(p.config.Region),
			NamespaceID:                p.config.Namespace,
			Name:                       containerDeployment.Name,
			EnvironmentVariables:       &map[string]string{}, // TODO add user environment
			RegistryImage:              scw.StringPtr(img.Name()),
			Port:                       scw.Uint32Ptr(uint32(p.config.Port)),
			SecretEnvironmentVariables: nil, // TODO
		}

		for key, value := range deployConfig.Env() {
			(*req.EnvironmentVariables)[key] = value
		}
		container, err = scwContainerAPI.CreateContainer(req, scw.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("failed to create container: %w", err)
		}

		st.Step(terminal.StatusOK, "Creating container")
	}
	containerDeployment.Id = container.ID
	containerDeployment.Name = container.Name
	containerDeployment.Image = container.RegistryImage
	containerDeployment.Url = container.DomainName

	state.Id = containerDeployment.Id
	state.Region = containerDeployment.Region

	st.Update("Deploy container")

	_, err = scwContainerAPI.DeployContainer(&containerSDK.DeployContainerRequest{
		Region:      scw.Region(containerDeployment.Region),
		ContainerID: containerDeployment.Id,
	}, scw.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to request container deployment: %w", err)
	}

	st.Update("Waiting for deployment")

	_, err = scwContainerAPI.WaitForContainer(&containerSDK.WaitForContainerRequest{
		ContainerID: containerDeployment.Id,
		Region:      scw.Region(containerDeployment.Region),
	}, scw.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to deploy container: %w", err)
	}

	//st.Step(terminal.StatusOK, fmt.Sprintf("Loaded or created container: %s", containerDeployment.Id))
	return nil
}

func (p *Platform) resourceContainerStatus(
	ctx context.Context,
	sg terminal.StepGroup,
	containerState *Resource_Container,
	sr *resource.StatusResponse,
	scwContainerAPI *containerSDK.API,
) error {
	s := sg.Add("Checking status of deployed Container...")
	defer s.Abort()

	container, err := scwContainerAPI.GetContainer(&containerSDK.GetContainerRequest{
		Region:      scw.Region(containerState.Region),
		ContainerID: containerState.Id,
	}, scw.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to get container (%s/%s): %w", containerState.Region, containerState.Id, err)
	}
	stateJson, err := json.Marshal(container)
	if err != nil {
		return fmt.Errorf("failed to marshal container state: %w", err)
	}

	health := sdk.StatusReport_UNKNOWN

	switch container.Status {
	case containerSDK.ContainerStatusReady:
		health = sdk.StatusReport_READY
	case containerSDK.ContainerStatusDeleting:
		health = sdk.StatusReport_MISSING
	case containerSDK.ContainerStatusError:
		health = sdk.StatusReport_DOWN
	case containerSDK.ContainerStatusLocked:
		health = sdk.StatusReport_DOWN
	case containerSDK.ContainerStatusCreating:
		health = sdk.StatusReport_PARTIAL
	case containerSDK.ContainerStatusPending:
		health = sdk.StatusReport_PARTIAL
	case containerSDK.ContainerStatusCreated:
		health = sdk.StatusReport_PARTIAL
	}

	healthMessage := "container is deploying or deployed"
	if container.ErrorMessage != nil {
		healthMessage = *container.ErrorMessage
	}

	containerResource := sdk.StatusReport_Resource{
		Id:                  container.ID,
		Name:                container.Name,
		Type:                "container",
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE,
		Platform:            "docker", // Make a docker icon appear on web ui
		StateJson:           string(stateJson),
		Health:              health,
		HealthMessage:       healthMessage,
		PlatformUrl: fmt.Sprintf("https://console.scaleway.com/containers/namespaces/fr-par/%s/containers/%s/deployment",
			container.NamespaceID, container.ID),
	}

	sr.Resources = append(sr.Resources, &containerResource)
	s.Done()
	return nil
}

func (p *Platform) resourceContainerDestroy(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	ui terminal.UI,
	containerState *Resource_Container,
	scwContainerAPI *containerSDK.API,
) error {
	s := sg.Add("Destroying container Container...")

	_, err := scwContainerAPI.DeleteContainer(&containerSDK.DeleteContainerRequest{
		Region:      scw.Region(containerState.Region),
		ContainerID: containerState.Id,
	}, scw.WithContext(ctx))
	responseError := &scw.ResponseError{} // TODO move out
	if err != nil && !(errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound) {
		return fmt.Errorf("failed to request container deletion: %w", err)
	}

	s.Done()
	return nil
}
