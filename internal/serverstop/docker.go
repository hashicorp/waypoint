package serverstop

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

const (
	containerLabel = "waypoint-type=server"
	volumeId       = "waypoint-server"
)

// StopDocker stops and removes the Waypoint Docker container along
// with its Docker volume.
func StopDocker(ctx context.Context, status terminal.Status) error {
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	defer func() {
		_ = dockerCli.Close()
	}()

	dockerCli.NegotiateAPIVersion(ctx)

	containers, err := dockerCli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerLabel,
		}),
	})

	if err != nil {
		return err
	}

	if len(containers) < 1 {
		return fmt.Errorf("cannot find a Waypoint Docker container")
	}

	// Pick the first container regardless of the number of matching
	// containers, as there should be only one.
	containerId := containers[0].ID

	status.Update("Stopping Waypoint Docker container...")

	// Stop the container gracefully, respecting the Engine's default timeout.
	if err := dockerCli.ContainerStop(ctx, containerId, nil); err != nil {
		return err
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := dockerCli.ContainerRemove(ctx, containerId, removeOptions); err != nil {
		return err
	}

	volumeExists, err := volumeExists(ctx, dockerCli)
	if err != nil {
		return err
	}

	// If the Waypoint Docker volume does not exist, return. This normally
	// shouldn't happen.
	if !volumeExists {
		status.Update("Couldn't find Waypoint Docker volume")
		return nil
	}

	status.Update("Removing Waypoint Docker volume...")

	if err := dockerCli.VolumeRemove(ctx, volumeId, true); err != nil {
		return err
	}

	return nil
}

// volumeExists determines whether the Waypoint Docker volume exists.
func volumeExists(ctx context.Context, dockerCli *client.Client) (bool, error) {
	listBody, err := dockerCli.VolumeList(ctx, filters.NewArgs(filters.KeyValuePair{
		Key:   "name",
		Value: volumeId,
	}))

	if err != nil {
		return false, err
	}

	exists := len(listBody.Volumes) > 0

	return exists, nil
}
