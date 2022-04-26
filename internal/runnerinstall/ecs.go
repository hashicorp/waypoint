package runnerinstall

import (
	"context"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/internal/installutil/aws"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

const (
	defaultRunnerLogGroup = "waypoint-runner-logs"
)

type ECSRunnerInstaller struct {
	Config EcsConfig
}

func (i *ECSRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI
	log := opts.Log

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.Config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	var (
		efsInfo       *aws.EfsInformation
		logGroup      string
		executionRole string
		netInfo       *aws.NetworkInformation
		taskRole      string
		runSvcArn     *string
	)
	lf := &aws.Lifecycle{
		Init: func(ui terminal.UI) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.Config.Region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			if netInfo, err = aws.SetupNetworking(ctx, ui, sess, i.Config.Subnets); err != nil {
				return err
			}

			if efsInfo, err = aws.SetupEFS(ctx, ui, sess, netInfo); err != nil {
				return err
			}

			if executionRole, err = aws.SetupExecutionRole(ctx, ui, log, sess, i.Config.ExecutionRoleName); err != nil {
				return err
			}

			taskRole, err = aws.SetupTaskRole(ctx, ui, log, sess, i.Config.TaskRoleName, opts.Id)
			if err != nil {
				return err
			}

			logGroup, err = aws.SetupLogs(ctx, ui, log, sess, defaultRunnerLogGroup)
			if err != nil {
				return err
			}

			return nil
		},

		Run: func(ui terminal.UI) error {
			runSvcArn, err = aws.LaunchRunner(
				ctx, ui, log, sess,
				opts.AdvertiseClient.Env(),
				executionRole,
				taskRole,
				logGroup,
				i.Config.Region,
				i.Config.CPU,
				i.Config.Memory,
				i.Config.RunnerImage,
				i.Config.Cluster,
				opts.Cookie,
				opts.Id,
				netInfo,
				efsInfo,
			)
			return err
		},

		Cleanup: func(ui terminal.UI) error { return nil },
	}

	if err := lf.Execute(log, ui); err != nil {
		return err
	}

	log.Debug("runner service started", "arn", *runSvcArn)

	return nil
}

func (i *ECSRunnerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "ecs-region",
		Usage:  "AWS region in which to install the Waypoint runner.",
		Target: &i.Config.Region,
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-exeuction-role-name",
		Target:  &i.Config.ExecutionRoleName,
		Usage:   "The name of the execution task IAM Role to associate with the ECS Service.",
		Default: "waypoint-runner-execution-role",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-task-role-name",
		Target:  &i.Config.TaskRoleName,
		Usage:   "IAM Execution Role to assign to the on-demand runner.",
		Default: "waypoint-runner",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-cpu",
		Target:  &i.Config.CPU,
		Usage:   "The amount of CPU to allocate for the Waypoint runner task.",
		Default: "512",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-memory",
		Target:  &i.Config.Memory,
		Usage:   "The amount of memory to allocate for the Waypoint runner task",
		Default: "2048",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-runner-image",
		Target:  &i.Config.RunnerImage,
		Default: "hashicorp/waypoint",
		Usage:   "The Waypoint runner Docker image.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.Config.Cluster,
		Default: "waypoint-server",
		Usage:   "The name of the ECS Cluster to install the Waypoint runner into.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "ecs-subnets",
		Target: &i.Config.Subnets,
		Usage:  "Subnets to install the Waypoint runner into.",
	})
}

func (i *ECSRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *ECSRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}

type EcsConfig struct {
	Region            string
	ExecutionRoleName string
	TaskRoleName      string
	CPU               string
	Memory            string
	RunnerImage       string
	Cluster           string
	Subnets           []string
}
