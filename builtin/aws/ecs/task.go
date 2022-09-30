package ecs

import (
	"context"
	"crypto/rand"
	"fmt"
	"google.golang.org/grpc/status"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/oklog/ulid"
	"google.golang.org/grpc/codes"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// TaskLauncher implements the TaskLauncher plugin interface to support
// launching on-demand tasks for the Waypoint server.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// StartTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StartTaskFunc() interface{} {
	return p.StartTask
}

// StopTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StopTaskFunc() interface{} {
	return p.StopTask
}

// WatchTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) WatchTaskFunc() interface{} {
	return p.WatchTask
}

// TaskLauncherConfig is the configuration structure for the task plugin. At
// this time all these are simply copied from what the Waypoint Server
// installation is using, with the only exception being the TaskRoleName.
type TaskLauncherConfig struct {
	// Cluster is the ECS we're operating in
	Cluster string `hcl:"cluster,optional"`

	// Region is the AWS region we're operating in, e.g. us-west-2, us-east-1
	Region string `hcl:"region,optional"`

	// ExecutionRoleName is the name of the AWS IAM role to apply to the task's
	// Execution Role. At this time we reuse the same Role as the Server
	// Execution Role.
	ExecutionRoleName string `hcl:"execution_role_name,optional"`

	// TaskRoleName is the name of the AWS IAM role to apply to the task. This
	// role determines the privileges the ODR builder has, and must have the correct
	// policies in place to work with the provided registries.
	TaskRoleName string `hcl:"task_role_name,optional"`

	// Subnets are the list of subnets for the cluster. These will match the
	// subnets used for the Cluster
	Subnets string `hcl:"subnets"`

	// SecurityGroupId is the security group used for the Waypoint tasks.
	SecurityGroupId string `hcl:"security_group_id"`

	// LogGroup is the CloudWatch log group name to use.
	LogGroup string `hcl:"log_group,optional"`

	// ODR Resource configuration
	OdrMemory string `hcl:"odr_memory,optional"`
	OdrCPU    string `hcl:"odr_cpu,optional"`
}

func (p *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(p.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`
Launch an ECS task for on-demand tasks from the Waypoint server.

This will use the standard AWS environment variables and IAM Role information to
source authentication information for AWS, using the configured task role.
If no task role name is specified, Waypoint will create one with the required
permissions.
`)

	doc.SetField(
		"odr_image",
		"Docker image for the Waypoint On-Demand Runners",
		docs.Summary(`
		Docker image for the Waypoint On-Demand Runners. This will
default to the server image with the name (not label) suffixed with '-odr'."
`,
		),
	)

	doc.SetField(
		"task_role_name",
		"Task role name to be used for the task role in the On-Demand Runner task",
		docs.Summary(`
This role must have the correct IAM policies to complete its task.
If this IAM role does not already exist, a role will be created with the correct
permissions"
`,
		),
	)

	doc.SetField(
		"cluster",
		"Cluster name to place On-Demand runner tasks in",
		docs.Summary(
			"ECS Cluster to place On-Demand runners in. This defaults to the cluster",
			"used by the Waypoint server",
		),
	)

	doc.SetField(
		"region",
		"AWS Region to use",
		docs.Summary(
			"AWS region to use. Defaults to the region used for the Waypoint Server.",
		),
	)

	doc.SetField(
		"execution_role_name",
		"The name of the AWS IAM role to apply to the task's Execution Role",
		docs.Summary(
			"ExecutionRoleName is the name of the AWS IAM role to apply to the task's",
			"Execution Role. At this time we reuse the same Role as the Waypoint",
			"server Execution Role.",
		),
	)

	doc.SetField(
		"task_role_name",
		"The name of the AWS IAM role to apply to the task's Task Role",
		docs.Summary(
			"TaskRoleName is the name of the AWS IAM role to apply to the task.",
			"This role determines the privileges the ODR builder. If no role",
			"name is given, an IAM role will be created with the required",
			"policies",
		),
	)

	doc.SetField(
		"subnets",
		"List of subnets to place the On-Demand Runner task in.",
		docs.Summary(
			"List of subnets to place the On-Demand Runner task in. This defaults",
			"to the list of subnets configured for the Waypoint server and ",
			"must be either identical or a subset of the subnets used by the ",
			"Waypoint server",
		),
	)

	doc.SetField(
		"security_group_id",
		"Security Group ID to place the On-Demand Runner task in",
		docs.Summary(
			"Security Group ID to place the On-Demand Runner task in. This defaults ",
			"to the security group used for the Waypoint server",
		),
	)

	doc.SetField(
		"log_group",
		"Cloud Watch Log Group to use for On-Demand Runners",
		docs.Summary(
			"Cloud Watch Log Group to use for On-Demand Runners. Defaults to the ",
			"log group used for runners (waypoint-runner).",
		),
	)

	doc.SetField(
		"odr_cpu",
		"CPU to use for the On-Demand runners.",
		docs.Summary(
			"Configure the CPU for the On-Demand runners. The default is 512. ",
			"See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html ",
			"for valid values",
		),
	)
	doc.SetField(
		"odr_memory",
		"Memory to use for the On-Demand runners.",
		docs.Summary(
			"Configure the memory for the On-Demand runners. The default is 1024. ",
			"See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html ",
			"for valid values",
		),
	)

	return doc, nil
}

// TaskLauncher implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to AWS ECS to stop the container created previously.
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
	})
	if err != nil {
		return err
	}

	ecsSvc := ecs.New(sess)
	tasks, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster: aws.String(p.config.Cluster),
		Family:  aws.String("waypoint-runner"),
	})
	if err != nil {
		return err
	}
	for _, taskArn := range tasks.TaskArns {
		taskResp, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: aws.String(p.config.Cluster),
			Tasks:   aws.StringSlice([]string{*taskArn}),
			Include: aws.StringSlice([]string{"TAGS"}),
		})
		if err != nil {
			return err
		} else if len(taskResp.Tasks) != 1 {
			return status.Errorf(codes.Internal, "there should be only 1 task, but there are %d", len(taskResp.Tasks))
		}
		for _, tag := range taskResp.Tasks[0].Tags {
			if *tag.Key == "waypoint-odr-task-name" && *tag.Value == ti.Id {
				log.Info("stopping ECS task", "task_id", ti.Id)
				_, err := ecsSvc.StopTask(&ecs.StopTaskInput{
					Cluster: aws.String(p.config.Cluster),
					Reason:  aws.String("Waypoint ODR job is complete"),
					Task:    taskArn,
				})
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	// If we reach this point, then we could not find our task in ECS
	// so there's nothing to clean up. Normally, ECS will automatically
	// clean up the task.
	return nil
}

// WatchTask implements TaskLauncher
func (p *TaskLauncher) WatchTask(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	ti *TaskInfo,
) (*component.TaskResult, error) {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
	})
	if err != nil {
		return nil, err
	}

	// this channel will receive the status of the ECS task
	taskStatusCh := make(chan string)
	go p.taskStatus(ctx, 5*time.Minute, &taskStatusCh, log, sess, ti.Id)

	var logStreamName string
	// We wait 5 minutes for the log stream to be available
	logStreamContext, cancel := context.WithTimeout(ctx, time.Minute*5)
	cwl := cloudwatchlogs.New(sess)
	// loop until task is ready
	defer cancel()
	for logStreamName == "" {
		select {
		case <-logStreamContext.Done():
			log.Error("timeout waiting for log stream to become ready", terminal.WithErrorStyle())
			return nil, err
		case taskStatus := <-taskStatusCh:
			if taskStatus == "RUNNING" {
				// use the task prefix "waypoint-odr-task-<id>" to filter the log streams
				logStreamsResp, err := cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
					LogGroupName:        aws.String(p.config.LogGroup),
					LogStreamNamePrefix: aws.String(ti.Id),
				})
				if err != nil {
					log.Error("error describing log streams", "err", err)
					ui.Output("Error describing log streams: %s", err)
					return nil, err
				}

				for _, logStream := range logStreamsResp.LogStreams {
					if strings.Contains(*logStream.LogStreamName, ti.Id) {
						logStreamName = *logStream.LogStreamName
					}
				}
			}
		default:
			continue
		}
	}

	token := ""
	var resp *cloudwatchlogs.GetLogEventsOutput
	for {
		select {
		default:
			// loop until log stream is ready
			getLogsInput := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(p.config.LogGroup),
				LogStreamName: aws.String(logStreamName),
				StartFromHead: aws.Bool(true),
			}
			if resp != nil {
				token = *resp.NextForwardToken
				getLogsInput.NextToken = aws.String(token)
			}
			resp, err = cwl.GetLogEvents(getLogsInput)
			if err != nil {
				return nil, err
			}

			if *resp.NextForwardToken == token {
				// if the task is done, AND we're at the end of the log
				// stream, we exit
				taskStatus := <-taskStatusCh
				if taskStatus == "DELETED" || taskStatus == "DEPROVISIONING" || taskStatus == "STOPPED" || taskStatus == "DEACTIVATING" {
					close(taskStatusCh)
					return &component.TaskResult{ExitCode: 0}, nil
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}

			for _, event := range resp.Events {
				log.Info(*event.Message)
				ui.Output(*event.Message, terminal.WithInfoStyle())
			}

			token = *resp.NextForwardToken
			// Sleep for half a second to slam the CloudWatch API less
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// StartTask runs an ECS Task to perform the requested job.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	// Generate an ID for our task name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Arguments here represent the command inputs used when executing the
	// waypoint runner, e.x:
	//   "runner, agent, -vvv, -id, <some id>, -odr"
	cmd := aws.StringSlice(tli.Arguments)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
	})
	if err != nil {
		return nil, err
	}
	ecsSvc := ecs.New(sess)

	envs := []*ecs.KeyValuePair{}
	for k, v := range tli.EnvironmentVariables {
		envs = append(envs, &ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}
	// This environment variable informs Kaniko that it is operating in a
	// container. Without this environment variable, Kaniko will fail with an
	// error
	envs = append(envs, &ecs.KeyValuePair{
		Name:  aws.String("container"),
		Value: aws.String("docker"),
	})

	taskName := fmt.Sprintf("waypoint-odr-task-%s", id.String())

	logOptions := buildLoggingOptions(
		nil,
		p.config.Region,
		p.config.LogGroup,
		taskName,
	)

	def := ecs.ContainerDefinition{
		Name:        &taskName,
		Image:       &tli.OciUrl,
		Command:     cmd,
		Environment: envs,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String(ecs.LogDriverAwslogs),
			Options:   logOptions,
		},
	}

	exRoleArn, err := roleArn(p.config.ExecutionRoleName, sess)
	if err != nil {
		return nil, err
	}

	taskRoleArn, err := roleArn(p.config.TaskRoleName, sess)
	if err != nil {
		return nil, err
	}

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    []*ecs.ContainerDefinition{&def},
		ExecutionRoleArn:        &exRoleArn,
		TaskRoleArn:             &taskRoleArn,
		Cpu:                     aws.String(p.config.OdrCPU),
		Memory:                  aws.String(p.config.OdrMemory),
		Family:                  aws.String("waypoint-runner"),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(ecs.LaunchTypeFargate)},
	}

	taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc, log)
	if err != nil {
		return nil, err
	}

	taskDefArn := *taskDef.TaskDefinitionArn

	subnetStrings := strings.Split(p.config.Subnets, ",")
	subnets := aws.StringSlice(subnetStrings)

	if _, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		LaunchType:           aws.String(ecs.LaunchTypeFargate),
		Cluster:              &p.config.Cluster,
		Count:                aws.Int64(1),
		TaskDefinition:       &taskDefArn,
		EnableECSManagedTags: aws.Bool(true),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: []*string{&p.config.SecurityGroupId},
				// without a public IP we cannot reach out to ECR or other
				// registries
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String("waypoint-odr-task-name"),
				Value: aws.String(taskName),
			},
		},
	}); err != nil {
		return nil, err
	}

	return &TaskInfo{
		Id: taskName,
	}, nil
}

func roleArn(name string, sess *session.Session) (string, error) {
	iamSvc := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: &name,
	}

	roleOut, err := iamSvc.GetRole(input)
	if err != nil {
		return "", err
	}
	if roleOut.Role == nil {
		return "", fmt.Errorf("no role found for (%s) role name", name)
	}
	// Arn is a required field of Role so we assume it will be populated and
	// prevent any nil dereference here
	return *roleOut.Role.Arn, nil
}

// taskStatus gets the status of an ECS task and provides it to the caller via a channel
// the caller is responsible for closing the channel
func (p *TaskLauncher) taskStatus(ctx context.Context, d time.Duration, taskStatusCh *chan string, log hclog.Logger, sess *session.Session, taskId string) {
	ecsSvc := ecs.New(sess)

	taskContext, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	for {
		select {
		case <-taskContext.Done():
			log.Error("timeout waiting for task to become ready", terminal.WithErrorStyle())
			return
		default:
		}
		tasks, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
			Cluster: aws.String(cluster),
			Family:  aws.String("waypoint-runner"),
		})
		if err != nil {
			log.Error("error listing ECS tasks", "err", err)
			return
		}
		var odrTask *ecs.Task
		for _, taskArn := range tasks.TaskArns {
			taskResp, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
				Cluster: aws.String(cluster),
				Tasks:   aws.StringSlice([]string{*taskArn}),
				Include: aws.StringSlice([]string{"TAGS"}),
			})
			if err != nil {
				log.Error("error describing ECS tasks", "err", err)
				return
			} else if len(taskResp.Tasks) != 1 {
				log.Error("there should be only 1 task", "num_tasks", len(taskResp.Tasks))
				return
			}
			for _, tag := range taskResp.Tasks[0].Tags {
				if *tag.Key == "waypoint-odr-task-name" && *tag.Value == taskId {
					odrTask = taskResp.Tasks[0]
					break
				}
			}
		}

		if odrTask == nil {
			log.Info("ODR Task not found")
			*taskStatusCh <- "DELETED"
			return
		} else {
			log.Debug("ODR task status", "status", *odrTask.LastStatus)
		}
		*taskStatusCh <- *odrTask.LastStatus
		time.Sleep(500 * time.Millisecond)
	}
}

var _ component.TaskLauncher = (*TaskLauncher)(nil)
