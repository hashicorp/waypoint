package installutil

import (
	"context"
	"errors"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"time"
)

func AdoptRunner(ctx context.Context, ui terminal.UI, client pb.WaypointClient, id string) error {
	ui.Output("Waiting for runner to connect to server...")
	// Waits 5 minutes for the server to detect the new runner before timing out
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	// TODO: Something safer than for true
	for true {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			ui.Output("Cancelled.",
				terminal.WithErrorStyle(),
			)
			return errors.New("context canceled")
		}
		// Use runner list API to check if runner is reporting to server yet
		// If it's found, adopt it. Otherwise, try until deadline.
		runners, err := client.ListRunners(ctx, &pb.ListRunnersRequest{})
		if err != nil {
			ui.Output("Error getting runners: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
		found := false
		for _, myRunner := range runners.Runners {
			if myRunner.Id == id {
				found = true
				break
			}
		}
		if found {
			_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
				RunnerId: id,
				Adopt:    true,
			})
			if err != nil {
				ui.Output("Error adopting runner: %s", clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return err
			}
			ui.Output("Runner %s adopted successfully.", id,
				terminal.WithSuccessStyle(),
			)
			return nil
		}
	}
	return nil
}
