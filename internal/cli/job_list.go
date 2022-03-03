package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobListCommand struct {
	*baseCommand

	flagJson  bool
	flagLimit int
	flagDesc  bool
}

func (c *JobListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	req := &pb.ListJobsRequest{}

	// NOTE(briancain): This is technically not a "public API" function
	resp, err := c.project.Client().XListJobs(ctx, req)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	jobs := resp.Jobs

	if c.flagDesc {
		var reverse []*pb.Job
		for i := len(jobs) - 1; i >= 0; i-- {
			reverse = append(reverse, jobs[i])
		}
		jobs = reverse
	}

	if c.flagLimit > 0 && c.flagLimit <= len(jobs) {
		jobs = jobs[len(jobs)-c.flagLimit:]
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, t := range jobs {
			str, err := m.MarshalToString(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	c.ui.Output("Waypoint Jobs", terminal.WithHeaderStyle())

	tblHeaders := []string{"ID", "Operation", "State", "Time Completed", "Target Runner", "Workspace", "Project", "Application"}
	tbl := terminal.NewTable(tblHeaders...)

	for _, j := range jobs {
		var op string
		// Job_Noop seems to be missing the isJob_operation method
		switch j.Operation.(type) {
		case *pb.Job_Build:
			op = "Build"
		case *pb.Job_Push:
			op = "Push"
		case *pb.Job_Deploy:
			op = "Deploy"
		case *pb.Job_Destroy:
			op = "Destroy"
		case *pb.Job_Release:
			op = "Release"
		case *pb.Job_Validate:
			op = "Validate"
		case *pb.Job_Auth:
			op = "Auth"
		case *pb.Job_Docs:
			op = "Docs"
		case *pb.Job_ConfigSync:
			op = "ConfigSync"
		case *pb.Job_Exec:
			op = "Exec"
		case *pb.Job_Up:
			op = "Up"
		case *pb.Job_Logs:
			op = "Logs"
		case *pb.Job_QueueProject:
			op = "QueueProject"
		case *pb.Job_Poll:
			op = "Poll"
		case *pb.Job_StatusReport:
			op = "StatusReport"
		case *pb.Job_StartTask:
			op = "StartTask"
		case *pb.Job_StopTask:
			op = "StopTask"
		case *pb.Job_Init:
			op = "Init"
		default:
			op = "Unknown"
		}

		var jobState string
		switch j.State {
		case pb.Job_UNKNOWN:
			jobState = "Unknown"
		case pb.Job_QUEUED:
			jobState = "Queued"
		case pb.Job_WAITING:
			jobState = "Waiting"
		case pb.Job_RUNNING:
			jobState = "Running"
		case pb.Job_ERROR:
			jobState = "Error"
		case pb.Job_SUCCESS:
			jobState = "Sucecss"
		default:
			jobState = "Unknown"
		}

		var targetRunner string
		switch target := j.TargetRunner.Target.(type) {
		case *pb.Ref_Runner_Any:
			targetRunner = "*"
		case *pb.Ref_Runner_Id:
			targetRunner = target.Id.Id
		}

		var completeTime string
		if time, err := ptypes.Timestamp(j.CompleteTime); err == nil {
			completeTime = humanize.Time(time)
		}

		tblColumn := []string{
			j.Id,
			op,
			jobState,
			completeTime,
			targetRunner,
			j.Workspace.Workspace,
			j.Application.Project,
			j.Application.Application,
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *JobListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "desc",
			Target:  &c.flagDesc,
			Default: false,
			Usage:   "Output the list of jobs from newest to oldest.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the list of jobs as json.",
		})

		f.IntVar(&flag.IntVar{
			Name:    "limit",
			Target:  &c.flagLimit,
			Default: 0,
			Usage:   "If set, will limit the number of jobs to list.",
		})
	})
}

func (c *JobListCommand) Synopsis() string {
	return "List all jobs in Waypoint"
}

func (c *JobListCommand) Help() string {
	return formatHelp(`
Usage: waypoint job list [options]

  List all known jobs from Waypoint server.

` + c.Flags().Help())
}
