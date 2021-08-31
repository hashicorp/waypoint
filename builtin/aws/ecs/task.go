package ecs

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/oklog/ulid"
	"github.com/ryboe/q"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
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

// TaskLauncherConfig is the configuration structure for the task plugin. At
// this time most of these are simply copied from what the Waypoint Server
// installation is using.
type TaskLauncherConfig struct {
	// Cluster is the ECS we're operating in
	Cluster string `hcl:"cluster,optional"`

	// Region is the AWS region we're operating in, e.g. us-west-2, us-east-1
	Region string `hcl:"region,optional"`

	// OdrExecutionRoleName is the name of the AWS IAM role to apply to the
	// task's Execution Role. This is generally the same as the Server Execution
	// Role.
	OdrExecutionRoleName string `hcl:"odr_execution_role_name,optional"`

	// TaskRoleArn is the name of the AWS IAM role to apply to the task. This
	// role determins the privilages the ODR builder has, and must have the correct
	// policies in place to work with the provided registries.
	OdrTaskRoleName string `hcl:"odr_task_role_name,optional"`

	// Subnets are the list of subnets for the cluster. These will match the
	// subnets used for the Cluster
	Subnets string `hcl:"subnets,optional"`

	// SecurityGroupId is the security group used for the Waypoint tasks.
	SecurityGroupId string `hcl:"security_group_id,optional"`

	// LogGroup is the CloudWatch log group name to use.
	LogGroup string `hcl:"log_group,optional"`
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
source authentication information for AWS. If this is running within ECS
itself (typical for a ECS-based installation), it will use the task's
IAM Task Role.
`)

	doc.SetField(
		"odr_execution_role_name",
		"Execution role name to be used for the execution role in the task",
		docs.Summary(
			"Execution role is the name of the AWS IAM Execution Role to use "+
				"as the execution role for the task.",
		),
	)

	doc.SetField(
		"odr_task_role_name",
		"Task role name to be used for the task role in the On-Demand Runner task",
		docs.Summary(
			"Task role name to be used for the task role in the On-Demand Runner task.",
			"This role must have the correct IAM policies to complete it's task.",
		),
	)

	return doc, nil
}

// TaskLauncher implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to docker to stop the container created previously
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	return nil
}

// StartTask creates a docker container for the task.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	q.Q("===================")
	q.Q("=> starting task")
	q.Q("===================")

	// Generate an ID for our pod name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// tli.Arguments represent the command inputs, e.x:
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

	exRoleArn, err := roleArn(p.config.OdrExecutionRoleName, sess)
	if err != nil {
		q.Q(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		q.Q("=> => ex roleArn err: ", err)
		q.Q("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		return nil, err
	}
	q.Q(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	q.Q("ex roleArn: ", exRoleArn)
	q.Q("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	taskRoleArn, err := roleArn(p.config.OdrTaskRoleName, sess)
	if err != nil {
		q.Q(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		q.Q("=> => task roleArn err: ", err)
		q.Q("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		return nil, err
	}
	q.Q(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	q.Q("task roleArn: ", taskRoleArn)
	q.Q("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    []*ecs.ContainerDefinition{&def},
		ExecutionRoleArn:        &exRoleArn,
		TaskRoleArn:             &taskRoleArn,
		Cpu:                     aws.String("1024"),
		Memory:                  aws.String("2048"),
		Family:                  aws.String("waypoint-runner"),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("waypoint-runner")},
		// Tags: []*ecs.Tag{
		// 	{
		// 		Key:   aws.String(defaultRunnerTagName),
		// 		Value: aws.String(defaultRunnerTagValue),
		// 	},
		// },
	}

	taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	// registerTaskDefinition() above ensures taskDef here is non-nil, if the
	// error returned is nil
	taskDefArn := *taskDef.TaskDefinitionArn

	subnetStrings := strings.Split(p.config.Subnets, ",")
	subnets := aws.StringSlice(subnetStrings)

	_, err = ecsSvc.RunTask(&ecs.RunTaskInput{
		LaunchType:           aws.String("FARGATE"),
		Cluster:              &p.config.Cluster,
		Count:                aws.Int64(1),
		TaskDefinition:       &taskDefArn,
		EnableECSManagedTags: aws.Bool(true),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: []*string{&p.config.SecurityGroupId},
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
	})
	if err != nil {
		q.Q(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		q.Q("=> => run task err: ", err)
		q.Q("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		return nil, err
	}

	return &TaskInfo{
		Id: taskName,
	}, nil
}

func roleArn(name string, sess *session.Session) (string, error) {
	iamSvc := iam.New(sess)
	// get the Task Role ARN
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

var _ component.TaskLauncher = (*TaskLauncher)(nil)
