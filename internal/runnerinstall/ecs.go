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
	config ecsConfig
}

func (i *ECSRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI
	log := opts.Log

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	var (
		logGroup      string
		executionRole string
		taskRole      string
		runSvcArn     *string
	)
	lf := &aws.Lifecycle{
		Init: func(ui terminal.UI) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.config.Region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			executionRole, err = aws.SetupExecutionRole(ctx, ui, log, sess, i.config.ExecutionRoleName)
			if err != nil {
				return err
			}

			taskRole, err = aws.SetupTaskRole(ctx, ui, log, sess, i.config.TaskRoleName)
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
				i.config.Region,
				i.config.CPU,
				i.config.Memory,
				i.config.RunnerImage,
				i.config.Cluster,
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
		Target: &i.config.Region,
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-exeuction-role-name",
		Target: &i.config.ExecutionRoleName,
		Usage:  "The name of the execution task IAM Role to associate with the ECS Service.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-task-role-name",
		Target: &i.config.TaskRoleName,
		Usage:  "IAM Execution Role to assign to the on-demand runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-cpu",
		Target: &i.config.CPU,
		Usage:  "The amount of CPU to allocate for the Waypoint runner task.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-memory",
		Target: &i.config.Memory,
		Usage:  "The amount of memory to allocate for the Waypoint runner task",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-runner-image",
		Target:  &i.config.RunnerImage,
		Default: "hashicorp/waypoint",
		Usage:   "The Waypoint runner Docker image.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.Cluster,
		Default: "waypoint-server",
		Usage:   "The name of the ECS Cluster to install the Waypoint runner into.",
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

type ecsConfig struct {
	Region            string
	ExecutionRoleName string
	TaskRoleName      string
	CPU               string
	Memory            string
	RunnerImage       string
	Cluster           string
}
