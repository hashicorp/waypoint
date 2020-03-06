package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/devflow/internal/pkg/status"
	"github.com/mitchellh/devflow/sdk/component"
)

type ECSLauncher struct {
	roleName string
	roleArn  string

	status status.Updater
}

func imageForRuntime(runtime string) string {
	return "robloweco/lambda:ruby2.5"
}

const host = "securetunnel-dev.df.hashicorp.engineering"

func (e *ECSLauncher) updateStatus(str string) {
	e.status.Update(str)
}

func (e *ECSLauncher) SetupCluster(ctx context.Context) error {
	ecsSvc := ecs.New(sess)

	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String("devflow")},
	})

	if err != nil {
		return err
	}

	if len(desc.Clusters) > 1 {
		return nil
	}

	e.updateStatus("creating ECS cluster")
	_, err = ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String("devflow"),
	})

	if err != nil {
		return err
	}

	return nil
}

const rolePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
		  "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

func (e *ECSLauncher) SetupRole(L hclog.Logger, app *component.Source) error {
	svc := iam.New(sess)

	e.roleName = "ecr-" + app.App

	e.updateStatus("setting up IAM role")
	L.Debug("attempting to retrieve existing role", "role-name", e.roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(e.roleName),
	}

	getOut, err := svc.GetRole(queryInput)
	if err == nil {
		e.roleArn = *getOut.Role.Arn
		L.Debug("found existing role", "arn", e.roleArn)
		return nil
	}

	L.Debug("creating new role")

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(e.roleName),
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return err
	}

	e.roleArn = *result.Role.Arn

	L.Debug("created new role", "arn", e.roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(e.roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return err
	}

	L.Debug("attached execution role policy")

	return nil
}

func (e *ECSLauncher) SetupLogs(L hclog.Logger, logGroup string) error {
	e.updateStatus("setting up CloudWatchLogs")

	cwl := cloudwatchlogs.New(sess)
	groups, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              aws.Int64(1),
		LogGroupNamePrefix: aws.String(logGroup),
	})

	if err != nil {
		return err
	}

	if len(groups.LogGroups) == 0 {
		L.Debug("creating log group", "group", logGroup)
		_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroup),
		})
		if err != nil {
			return err
		}
	}

	return nil

}

func (e *ECSLauncher) Launch(ctx context.Context, L hclog.Logger, S status.Updater, app *component.Source, cfg *LambdaConfiguration) (*ConsoleClient, error) {
	e.status = S

	err := e.SetupCluster(ctx)
	if err != nil {
		return nil, err
	}

	err = e.SetupRole(L, app)
	if err != nil {
		return nil, err
	}

	err = e.SetupLogs(L, "devflow-logs")
	if err != nil {
		return nil, err
	}

	e.updateStatus("configuring secure tunnel")
	L.Debug("configuring secure tunnel")
	cc, err := NewConsoleClient(host)
	if err != nil {
		return nil, err
	}

	ecsSvc := ecs.New(sess)

	streamPrefix := fmt.Sprintf("devflow-task-%d", time.Now().Nanosecond())

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Name:      aws.String("devflow-console"),
		Command:   []*string{aws.String("/lambda-runner"), aws.String("-serve")},
		Image:     aws.String(imageForRuntime(cfg.Runtime)),
		Memory:    aws.Int64(512),
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String("awslogs"),
			Options: map[string]*string{
				"awslogs-group":         aws.String("devflow-logs"),
				"awslogs-region":        aws.String("us-west-2"),
				"awslogs-stream-prefix": aws.String(streamPrefix),
			},
		},
	}

	e.updateStatus("registering ECS task definition")
	L.Debug("registring task definition")

	taskOut, err := ecsSvc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{&def},

		ExecutionRoleArn: aws.String(e.roleArn),
		Cpu:              aws.String("256"),
		Memory:           aws.String("512"),
		Family:           aws.String("devflow"),

		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
	})

	if err != nil {
		return nil, err
	}

	e.updateStatus("starting ECS task")
	L.Debug("running task", "arn", *taskOut.TaskDefinition.TaskDefinitionArn)

	taskArn := *taskOut.TaskDefinition.TaskDefinitionArn

	runOut, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		Cluster:        aws.String("devflow"),
		Count:          aws.Int64(1),
		LaunchType:     aws.String("FARGATE"),
		StartedBy:      aws.String("devflow-console"),
		TaskDefinition: aws.String(taskArn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        []*string{aws.String("subnet-2761567d")},
				AssignPublicIp: aws.String("ENABLED"),
				// SecurityGroups: awsStrings(r.SecurityGroups),
			},
		},

		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Name: aws.String("devflow-console"),
					Environment: []*ecs.KeyValuePair{
						{
							Name:  aws.String("DEVFLOW_TUNNEL_TOKEN"),
							Value: aws.String(cc.Tunnel.ServerToken()),
						},
						{
							Name:  aws.String("DEVFLOW_TUNNEL_KEY"),
							Value: aws.String(cc.Tunnel.ServerKey()),
						},
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	L.Debug("task starting", "arn", *runOut.Tasks[0].TaskArn)

	var taskArns []*string

	for _, task := range runOut.Tasks {
		taskArns = append(taskArns, aws.String(*task.TaskArn))
	}

	var status string

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 50; i++ {
		descOut, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: aws.String("devflow"),
			Tasks:   taskArns,
		})

		if err != nil {
			return nil, err
		}

		desc := descOut.Tasks[0]

		if status != *desc.LastStatus {
			status = *desc.LastStatus
			e.updateStatus("task launch status: " + status)

			if status == "RUNNING" {
				break
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			// ok
		}
	}

	e.updateStatus("task running")
	L.Debug("task running", "arn", *runOut.Tasks[0].TaskArn)

	cc.UseApp(cfg)

	return cc, nil
}
