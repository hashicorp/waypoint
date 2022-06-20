package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobListCommand struct {
	*baseCommand

	flagJson               bool
	flagLimit              int
	flagDesc               bool
	flagState              []string
	flagTargetRunner       string
	flagTargetRunnerLabels map[string]string
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

	if c.flagWorkspace != "" {
		req.Workspace = &pb.Ref_Workspace{
			Workspace: c.flagWorkspace,
		}
	}

	if c.flagProject != "" {
		req.Project = &pb.Ref_Project{
			Project: c.flagProject,
		}
	}

	if c.flagApp != "" {
		req.Application = &pb.Ref_Application{
			Application: c.flagApp,
			Project:     c.flagProject,
		}
	}

	if len(c.flagState) > 0 {
		var states []pb.Job_State
		for _, s := range c.flagState {
			js, ok := pb.Job_State_value[strings.ToUpper(s)]
			if !ok {
				// this shouldn't happen given the State flag is an enum var, but protect
				// against it anyway
				c.ui.Output("Undefined job state value: "+s, terminal.WithErrorStyle())
				return 1
			} else {
				states = append(states, pb.Job_State(js))
			}
		}

		req.JobState = states
	}

	if len(c.flagTargetRunnerLabels) > 0 && c.flagTargetRunner != "" {
		c.ui.Output("Cannot define both 'target-runner-id' and 'target-runner-label' flags.\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagTargetRunner != "" {
		if c.flagTargetRunner == "*" {
			req.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}}
		} else {
			req.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: c.flagTargetRunner,
				},
			},
			}
		}
	} else if len(c.flagTargetRunnerLabels) > 0 {
		req.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Labels{
			Labels: &pb.Ref_RunnerLabels{
				Labels: c.flagTargetRunnerLabels,
			},
		}}
	}

	resp, err := c.project.Client().ListJobs(ctx, req)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	jobs := resp.Jobs

	// sort by complete time
	if c.flagDesc {
		sort.Slice(jobs, func(i, j int) bool {
			if jobs[i].CompleteTime == nil {
				return false
			} else if jobs[j].CompleteTime == nil {
				return true
			} else {
				return jobs[i].CompleteTime.AsTime().Before(jobs[j].CompleteTime.AsTime())
			}
		})
	} else {
		sort.Slice(jobs, func(i, j int) bool {
			if jobs[i].CompleteTime == nil {
				return true
			} else if jobs[j].CompleteTime == nil {
				return false
			} else {
				return jobs[i].CompleteTime.AsTime().After(jobs[j].CompleteTime.AsTime())
			}
		})
	}

	// limit to the first n jobs
	if c.flagLimit > 0 && c.flagLimit <= len(jobs) {
		jobs = jobs[:c.flagLimit]
	}

	if c.flagJson {
		m := protojson.MarshalOptions{
			Indent: "\t",
		}
		for _, t := range jobs {
			data, err := m.Marshal(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(string(data))
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
		case *pb.Job_WatchTask:
			op = "WatchTask"
		case *pb.Job_Init:
			op = "Init"
		case *pb.Job_PipelineStep:
			op = "PipelineStep"
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
			jobState = "Success"
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
		if j.CompleteTime != nil {
			completeTime = humanize.Time(j.CompleteTime.AsTime())
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

var jobStateValues = []string{"Success", "Error", "Running", "Waiting", "Queued", "Unknown"}

func (c *JobListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.EnumVar(&flag.EnumVar{
			Name:   "state",
			Target: &c.flagState,
			Values: jobStateValues,
			Usage:  "List jobs that only match the requested state. Can be repeated multiple times.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "target-runner-id",
			Target:  &c.flagTargetRunner,
			Default: "",
			Usage:   "List jobs that were only assigned to the target runner by id.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "target-runner-label",
			Target: &c.flagTargetRunnerLabels,
			Usage: "List jobs that were only assigned to the target runner by labels. " +
				"Can be repeated multiple times.",
		})

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
