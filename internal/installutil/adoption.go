package installutil

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func AdoptRunner(ctx context.Context, ui terminal.UI, client pb.WaypointClient, id string, addr string) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Waiting for runner to connect to server at %s...", addr)
	defer func() { s.Abort() }()

	// Waits 5 minutes for the server to detect the new runner before timing out
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
LOOP:
	for {
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
			ui.Output(runnerFailedToConnectToServer, clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
		for _, myRunner := range runners.Runners {
			if myRunner.Id == id {
				break LOOP
			}
		}
	}
	s.Update("Runner detected by server")
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Adopting runner...")
	_, err := client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
		RunnerId: id,
		Adopt:    true,
	})
	if err != nil {
		ui.Output("Error adopting runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}
	s.Update("Runner %q adopted successfully.", id)
	s.Status(terminal.StatusOK)
	s.Done()
	return nil
}

var (
	runnerFailedToConnectToServer = strings.TrimSpace(`
The Waypoint runner was unable to connect to Waypoint server. Maybe the
-server-addr specified is not accessible from the Waypoint runner?
`)
)
