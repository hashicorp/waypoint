package cli

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
	"sort"
	"time"
)

type StatusCommand struct {
	*baseCommand
}

type deploymentNameSequenceSorted []deploymentStatusReport

func (d deploymentNameSequenceSorted) Len() int {
	return len(d)
}

func (d deploymentNameSequenceSorted) Less(i, j int) bool {
	if d[i].Deployment.Application.Application == d[j].Deployment.Application.Application {
		return d[i].Sequence > d[j].Sequence // We want newer sequences first
	}

	return d[i].Deployment.Application.Application < d[j].Deployment.Application.Application
}

func (d deploymentNameSequenceSorted) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type deploymentStatusReport struct {
	*pb.Deployment
	*pb.StatusReport
}

var _ sort.Interface = (*deploymentNameSequenceSorted)(nil)

func (c *StatusCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListStatusReports(c.Ctx, &pb.ListStatusReportsRequest{
		Application: c.refApp,
		Workspace:   c.refWorkspace,
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Desc:  true,
		},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Deployments Status", terminal.WithHeaderStyle())
	err = c.listDeploymentsStatus(resp)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *StatusCommand) listDeploymentsStatus(resp *pb.ListStatusReportsResponse) error {
	table := terminal.NewTable("Name", "Version", "Health", "Message", "Last Check")
	var dsr []deploymentStatusReport

	for _, sr := range resp.StatusReports {
		if sr.TargetId == nil {
			continue
		}

		// Get Deployment
		dep, err := c.project.Client().GetDeployment(c.Ctx, &pb.GetDeploymentRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: sr.GetDeploymentId()},
			},
		})

		if err != nil {
			return fmt.Errorf(
				"unable to get deployment: %v",
				err,
			)
		}

		dsr = append(dsr, deploymentStatusReport{
			Deployment:   dep,
			StatusReport: sr,
		})

	}

	sort.Sort(deploymentNameSequenceSorted(dsr))

	for _, dsrVal := range dsr {
		table.Rich([]string{
			dsrVal.StatusReport.Application.Application,
			fmt.Sprintf("v%d", dsrVal.Deployment.Sequence),
			dsrVal.StatusReport.Health.HealthStatus,
			dsrVal.StatusReport.Health.HealthMessage,
			humanize.Time(
				time.Unix(
					dsrVal.StatusReport.Status.CompleteTime.Seconds,
					int64(dsrVal.StatusReport.Status.CompleteTime.Nanos),
				),
			),
		}, nil)
	}

	c.ui.Table(table)
	return nil
}

func (c *StatusCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *StatusCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *StatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StatusCommand) Synopsis() string {
	return "List status for the current project."
}

func (c *StatusCommand) Help() string {
	return formatHelp(`
Usage: waypoint status

  Show the status for the current project.

  This shows a detailed structure of the project resources and
  the current reported health associated with them.

`)
}
