package lambda

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// This launches an ECS task to run the given image. It creates an ECS cluster configured
// against Fargate for this purpose, as well as creates an IAM role and security group
// so the task can be accessed via TCP.

type ecsLauncher struct {
	Region       string
	PublicKey    string
	HostKey      string
	Image        string
	DeploymentId string

	LogOutput io.Writer

	roleName string
	roleArn  string

	status terminal.Status
}

const ecsClusterName = "waypoint-lambda-exec"

// SetupCluster creates an ECS cluster if there isn't one.
func (e *ecsLauncher) SetupCluster(sess *session.Session, ctx context.Context, log hclog.Logger) error {
	ecsSvc := ecs.New(sess)

	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(ecsClusterName)},
	})

	if err != nil {
		return err
	}

	if len(desc.Clusters) > 1 {
		log.Info("existing ECS cluster found", "arn", *desc.Clusters[0].ClusterArn)
		return nil
	}

	log.Info("creating ECS cluster")
	out, err := ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(ecsClusterName),
	})

	if err != nil {
		log.Error("error creating new cluster", "error", err)
		return err
	}

	log.Info("created cluster", "arn", *out.Cluster.ClusterArn)

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

// SetupRole creates an IAM role to use for ECS task execution. Reuses a role if there is already one.
func (e *ecsLauncher) SetupRole(L hclog.Logger, sess *session.Session, log hclog.Logger, app *component.Source) error {
	svc := iam.New(sess)

	e.roleName = "ecr-" + app.App

	log.Info("setting up IAM role")
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

// SetupLogs creates a cloudwatch logs LogGroup to send the ECS task logs to. These logs
// are only for debugging.
func (e *ecsLauncher) SetupLogs(L hclog.Logger, sess *session.Session, logGroup string) error {
	L.Info("setting up CloudWatchLogs")

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

type TaskInfo struct {
	IP  string
	Arn string
}

// This is the port the ECS task will be listening on.
const sshPort = 2222

// Launch creates the ECS task and returns it's public IP and ARN
func (e *ecsLauncher) Launch(
	ctx context.Context,
	L hclog.Logger,
	UI terminal.UI,
	app *component.Source,
	cfg *Deployment,
) (*TaskInfo, error) {
	S := UI.Status()
	defer S.Close()
	e.status = S

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: e.Region,
	})
	if err != nil {
		return nil, err
	}

	err = e.SetupCluster(sess, ctx, L)
	if err != nil {
		return nil, err
	}

	err = e.SetupRole(L, sess, L, app)
	if err != nil {
		return nil, err
	}

	logName := ecsClusterName + "-logs"

	err = e.SetupLogs(L, sess, logName)
	if err != nil {
		return nil, err
	}

	ecsSvc := ecs.New(sess)

	streamPrefix := fmt.Sprintf("waypoint-task-%d", time.Now().Nanosecond())

	L.Info("registering task definiton", "image", e.Image)

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Name:      aws.String("waypoint-console"),
		Image:     aws.String(e.Image),
		Memory:    aws.Int64(512),
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String("awslogs"),
			Options: map[string]*string{
				"awslogs-group":         aws.String(logName),
				"awslogs-region":        aws.String(e.Region),
				"awslogs-stream-prefix": aws.String(streamPrefix),
			},
		},
	}

	L.Debug("registring task definition")

	taskOut, err := ecsSvc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{&def},

		ExecutionRoleArn: aws.String(e.roleArn),
		Cpu:              aws.String("256"),
		Memory:           aws.String("512"),
		Family:           aws.String(ecsClusterName),

		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
	})

	if err != nil {
		return nil, err
	}

	L.Debug("running task", "arn", *taskOut.TaskDefinition.TaskDefinitionArn)

	taskArn := *taskOut.TaskDefinition.TaskDefinitionArn

	defaultSubnets, vpc, err := utils.DefaultPublicSubnets(ctx, sess)
	if err != nil {
		return nil, err
	}

	sg, err := utils.CreateSecurityGroup(ctx, sess, "waypoint-lambda-exec", vpc, sshPort)
	if err != nil {
		return nil, err
	}

	runOut, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		Cluster:        aws.String(ecsClusterName),
		Count:          aws.Int64(1),
		LaunchType:     aws.String("FARGATE"),
		StartedBy:      aws.String("waypoint-console"),
		TaskDefinition: aws.String(taskArn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        defaultSubnets,
				AssignPublicIp: aws.String("ENABLED"),
				SecurityGroups: []*string{sg},
			},
		},

		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Name: aws.String("waypoint-console"),
					Environment: []*ecs.KeyValuePair{
						{
							Name:  aws.String("WAYPOINT_EXEC_PLUGIN_SSH"),
							Value: aws.String(strconv.Itoa(sshPort)),
						},
						{
							Name:  aws.String("WAYPOINT_EXEC_PLUGIN_SSH_KEY"),
							Value: aws.String(e.PublicKey),
						},
						{
							Name:  aws.String("WAYPOINT_EXEC_PLUGIN_SSH_HOST_KEY"),
							Value: aws.String(e.HostKey),
						},
						{
							Name:  aws.String("WAYPOINT_EXEC_PLUGIN_SSH_DEPLOYMENT_ID"),
							Value: aws.String(e.DeploymentId),
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

	var ti TaskInfo

	// Watch the task's status for it to run, reporting as it transitions.
	for i := 0; i < 50; i++ {
		descOut, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: aws.String(ecsClusterName),
			Tasks:   taskArns,
		})

		if err != nil {
			return nil, err
		}

		desc := descOut.Tasks[0]

		if status != *desc.LastStatus {
			status = *desc.LastStatus
			L.Info("task launch status", "status", status)

			if status == "RUNNING" {
				ips, err := utils.ECSTaskPublicIPs(sess, descOut.Tasks)
				if err != nil {
					return nil, err
				}

				if len(ips) == 0 {
					L.Error("tasks didn't have any public ips")
					return nil, fmt.Errorf("Unable to calculate public IP of ECS task")
				}

				ti.IP = net.JoinHostPort(ips[0], strconv.Itoa(sshPort))
				rewriteLine(e.LogOutput, "Launching ECS task to provide shell: running")
				break
			} else if status == "STOPPED" {
				rewriteLine(e.LogOutput, "Launching ECS task to provide shell: error")
				L.Error("task stopped before running", "reason", *desc.StoppedReason)
				return nil, fmt.Errorf("task was unable to start")
			} else {
				rewriteLine(e.LogOutput, "Launching ECS task to provide shell: %s", strings.ToLower(status))
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			// ok
		}
	}

	L.Debug("task running", "arn", *runOut.Tasks[0].TaskArn)

	ti.Arn = *runOut.Tasks[0].TaskArn

	return &ti, nil
}
