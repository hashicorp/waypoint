package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRunnerLogGroup = "waypoint-runner-logs"
	defaultServerLogGroup = "waypoint-server-logs"

	defaultTaskFamily  = "waypoint-server"
	defaultTaskRuntime = "FARGATE"

	defaultSecurityGroupName = "waypoint-server-security-group"
	defaultNLBName           = "waypoint-server-nlb"

	// These tags are used to tag resources as they are created in AWS. Two tags
	// are required, so that the UninstallRunner method(s) can query for runner
	// resources and retrieve the cluster ARN
	defaultServerTagName  = "waypoint-server"
	defaultServerTagValue = "server-component"
	defaultRunnerTagName  = "waypoint-runner"
	defaultRunnerTagValue = "runner-component"

	serverName = "waypoint-server"
	runnerName = "waypoint-runner"

	defaultGrpcPort = "9701"
)

type Lifecycle struct {
	Init    func(terminal.UI) error
	Run     func(terminal.UI) error
	Cleanup func(terminal.UI) error
}

func (lf *Lifecycle) Execute(log hclog.Logger, ui terminal.UI) error {
	if lf.Init != nil {
		log.Debug("lifecycle init")

		err := lf.Init(ui)
		if err != nil {
			return err
		}

	}

	log.Debug("lifecycle run")
	err := lf.Run(ui)
	if err != nil {
		return err
	}

	if lf.Cleanup != nil {
		log.Debug("lifecycle cleanup")

		err = lf.Cleanup(ui)
		if err != nil {
			return err
		}
	}

	return nil
}

type Logging struct {
	CreateGroup      bool   `hcl:"create_group,optional"`
	StreamPrefix     string `hcl:"stream_prefix,optional"`
	DateTimeFormat   string `hcl:"datetime_format,optional"`
	MultilinePattern string `hcl:"multiline_pattern,optional"`
	Mode             string `hcl:"mode,optional"`
	MaxBufferSize    string `hcl:"max_buffer_size,optional"`
}

func buildLoggingOptions(
	lo *Logging,
	region string,
	logGroup string,
	defaultStreamPrefix string,
) map[string]*string {
	result := map[string]*string{
		"awslogs-region":        aws.String(region),
		"awslogs-group":         aws.String(logGroup),
		"awslogs-stream-prefix": aws.String(defaultStreamPrefix),
	}

	if lo != nil {
		// We receive the error `Log driver awslogs disallows options:
		// awslogs-endpoint` when setting `awslogs-endpoint`, so that is not
		// included here of the available options
		// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html
		result["awslogs-datetime-format"] = aws.String(lo.DateTimeFormat)
		result["awslogs-multiline-pattern"] = aws.String(lo.MultilinePattern)
		result["mode"] = aws.String(lo.Mode)
		result["max-buffer-size"] = aws.String(lo.MaxBufferSize)

		if lo.CreateGroup {
			result["awslogs-create-group"] = aws.String("true")
		}
		if lo.StreamPrefix != "" {
			result["awslogs-stream-prefix"] = aws.String(lo.StreamPrefix)
		}
	}

	for k, v := range result {
		if *v == "" {
			delete(result, k)
		}
	}

	return result
}

func LaunchRunner(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,
	env []string,
	executionRoleArn, taskRoleArn, logGroup, region, cpu, memory, runnerImage, cluster string,
) (*string, error) {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Installing Waypoint runner into ECS...")
	defer func() { s.Abort() }()

	defaultStreamPrefix := fmt.Sprintf("waypoint-runner-%d", time.Now().Nanosecond())
	logOptions := buildLoggingOptions(
		nil,
		region,
		logGroup,
		defaultStreamPrefix,
	)

	grpcPort, _ := strconv.Atoi(defaultGrpcPort)

	envs := []*ecs.KeyValuePair{}
	for _, line := range env {
		idx := strings.Index(line, "=")
		if idx == -1 {
			// Should never happen but let's not crash.
			continue
		}

		key := line[:idx]
		value := line[idx+1:]
		envs = append(envs, &ecs.KeyValuePair{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Command: []*string{
			aws.String("runner"),
			aws.String("agent"),
			aws.String("-vv"),
			aws.String("-liveness-tcp-addr=:1234"),
		},
		Name:  aws.String("waypoint-runner"),
		Image: aws.String(runnerImage),
		PortMappings: []*ecs.PortMapping{
			{
				ContainerPort: aws.Int64(int64(grpcPort)),
				HostPort:      aws.Int64(int64(grpcPort)),
			},
		},
		Environment: envs,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String(ecs.LogDriverAwslogs),
			Options:   logOptions,
		},
	}

	s.Update("Registering Task definition: waypoint-runner")

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{&def},

		ExecutionRoleArn:        aws.String(executionRoleArn),
		Cpu:                     aws.String(cpu),
		Memory:                  aws.String(memory),
		Family:                  aws.String(runnerName),
		TaskRoleArn:             &taskRoleArn,
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(defaultTaskRuntime)},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(defaultRunnerTagName),
				Value: aws.String(defaultRunnerTagValue),
			},
		},
	}

	ecsSvc := ecs.New(sess)
	taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc, log)
	if err != nil {
		return nil, err
	}

	taskDefArn := *taskDef.TaskDefinitionArn
	s.Update("Creating Service...")
	log.Debug("creating service", "arn", *taskDef.TaskDefinitionArn)

	// find the default security group to use
	ec2srv := ec2.New(sess)
	dsg, err := ec2srv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(defaultSecurityGroupName)},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var groupId *string
	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
		s.Update("Using existing security group: %s", defaultSecurityGroupName)
	} else {
		return nil, fmt.Errorf("could not find security group (%s)", defaultSecurityGroupName)
	}

	// query what subnets and vpc information from the server service
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(serverName)},
	})
	if err != nil {
		return nil, err
	}

	// should only find one
	service := services.Services[0]
	if service == nil {
		return nil, fmt.Errorf("no waypoint-server service found")
	}

	clusterArn := service.ClusterArn
	subnets := service.NetworkConfiguration.AwsvpcConfiguration.Subnets

	createServiceInput := &ecs.CreateServiceInput{
		Cluster:              clusterArn,
		DesiredCount:         aws.Int64(1),
		LaunchType:           aws.String(defaultTaskRuntime),
		ServiceName:          aws.String(runnerName),
		EnableECSManagedTags: aws.Bool(true),
		TaskDefinition:       aws.String(taskDefArn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: []*string{groupId},
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(defaultRunnerTagName),
				Value: aws.String(defaultRunnerTagValue),
			},
		},
	}

	s.Update("Creating ECS Service (%s)", runnerName)
	svc, err := createService(createServiceInput, ecsSvc)
	if err != nil {
		return nil, err
	}
	s.Update("Runner service created")
	s.Done()

	return svc.ClusterArn, nil
}

func createService(serviceInput *ecs.CreateServiceInput, ecsSvc *ecs.ECS) (*ecs.Service, error) {
	// AWS is eventually consistent so even though we probably created the
	// resources that are referenced by the service, it can error out if we try to
	// reference those resources too quickly. So we're forced to guard actions
	// which reference other AWS services with loops like this.
	var (
		servOut *ecs.CreateServiceOutput
		err     error
	)
	for i := 0; i < 30; i++ {
		servOut, err = ecsSvc.CreateService(serviceInput)
		if err == nil {
			break
		}

		// if we encounter an unrecoverable error, exit now.
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "AccessDeniedException", "UnsupportedFeatureException",
				"PlatformUnknownException",
				"PlatformTaskDefinitionIncompatibilityException":
				return nil, err
			}
		}

		// otherwise sleep and try again
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}
	return servOut.Service, nil
}

func SetupExecutionRole(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,

	executionRoleName string,
) (string, error) {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Setting up an IAM execution role...")
	defer func() { s.Abort() }()

	svc := iam.New(sess)

	roleName := executionRoleName

	// role names have to be 64 characters or less, and the client side doesn't
	// validate this.
	if len(roleName) > 64 {
		roleName = roleName[:64]
		log.Debug("using a shortened value for role name due to AWS's length limits", "roleName", roleName)
	}

	log.Debug("attempting to retrieve existing role", "role-name", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	getOut, err := svc.GetRole(queryInput)
	if err == nil {
		s.Update("Found existing IAM role to use: %s", roleName)
		s.Done()
		return *getOut.Role.Arn, nil
	}

	log.Debug("creating new role")
	s.Update("Creating IAM role: %s", roleName)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
			},
		},
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return "", err
	}

	roleArn := *result.Role.Arn

	log.Debug("created new execution role", "arn", roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return "", err
	}

	log.Debug("attached IAM execution role policy")

	s.Update("Created IAM execution role: %s", roleName)
	s.Done()
	return roleArn, nil
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

// odrRolePolicy represents the minimum policies required for an On-Demand
// Runner task to successfully build and deploy a Waypoint application to ECS.
// We chose to enumerate the minimum policies to avoid being over privileged.
// This list may not be exhaustive or complete to deploy to all platforms (EC2,
// Lambda), but represent a reasonable minimum. To add additional policies,
// users can create their own role and use it as the Task Role with the
// -ecs-task-role-name server installation flag, using these policies as a
// starting point.
const odrRolePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateSecurityGroup",
        "ec2:DescribeSecurityGroups",
        "ec2:DeleteSecurityGroup",
        "ec2:DescribeSubnets",
        "ecr:BatchGetImage",
        "ecr:CompleteLayerUpload",
        "ecr:DescribeImages",
        "ecr:DescribeRepositories",
        "ecr:GetAuthorizationToken",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:InitiateLayerUpload",
        "ecr:ListImages",
        "ecr:ListTagsForResource",
        "ecr:PutImage",
        "ecr:PutImageTagMutability",
        "ecr:ReplicateImage",
        "ecr:TagResource",
        "ecr:UntagResource",
        "ecr:UploadLayerPart",
        "ecs:CreateCluster",
        "ecs:CreateService",
        "ecs:DeleteService",
        "ecs:DescribeClusters",
        "ecs:DescribeServices",
        "ecs:ListTasks",
        "ecs:DescribeTasks",
        "ecs:RegisterTaskDefinition",
        "ecs:DeregisterTaskDefinition",
        "ecs:RunTask",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateRule",
        "elasticloadbalancing:CreateTargetGroup",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:DeleteRule",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:DescribeListeners",
        "elasticloadbalancing:DescribeLoadBalancers",
        "elasticloadbalancing:DescribeRules",
        "elasticloadbalancing:DescribeTargetGroups",
        "elasticloadbalancing:ModifyListener",
        "iam:AttachRolePolicy",
        "iam:CreateRole",
        "iam:GetRole",
        "iam:ListAttachedRolePolicies",
        "iam:PassRole",
        "logs:CreateLogGroup",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "*"
    }
  ]
}`

func SetupTaskRole(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,

	taskRoleName string,
) (string, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Setting up an IAM task role for ODR runners...")
	defer func() { s.Abort() }()

	svc := iam.New(sess)

	roleName := taskRoleName

	// role names have to be 64 characters or less, and the client side doesn't
	// validate this.
	if len(roleName) > 64 {
		roleName = roleName[:64]
		log.Debug("using a shortened value for role name due to AWS's length limits", "roleName", roleName)
	}

	log.Debug("attempting to retrieve existing role", "role-name", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	getOut, err := svc.GetRole(queryInput)
	if err == nil {
		s.Update("Found existing IAM role to use: %s", roleName)
		s.Done()
		return *getOut.Role.Arn, nil
	}

	log.Debug("creating new role")
	s.Update("Creating IAM task role: %s", roleName)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(defaultRunnerTagName),
				Value: aws.String(defaultRunnerTagValue),
			},
		},
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return "", err
	}

	roleArn := *result.Role.Arn

	log.Debug("created new task role", "arn", roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return "", err
	}

	log.Debug("attached IAM task role policy")

	// ODR specific policies
	_, err = svc.PutRolePolicy(&iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyDocument: aws.String(odrRolePolicy),
		PolicyName:     aws.String("waypoint-odr-policy"),
	})
	if err != nil {
		return "", err
	}

	s.Update("Created IAM task role: %s", roleName)
	s.Done()
	return roleArn, nil
}

func SetupLogs(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,

	logGroup string,
) (string, error) {
	cwl := cloudwatchlogs.New(sess)

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Examining existing CloudWatchLogs groups...")
	defer func() { s.Abort() }()

	groups, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              aws.Int64(1),
		LogGroupNamePrefix: aws.String(logGroup),
	})
	if err != nil {
		return "", err
	}

	if len(groups.LogGroups) == 0 {
		s.Update("Creating CloudWatchLogs group to store logs in: %s", logGroup)

		log.Debug("creating log group", "group", logGroup)
		_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroup),
		})
		if err != nil {
			return "", err
		}

		s.Update("Created CloudWatchLogs group to store logs in: %s", logGroup)
	} else {
		s.Update("Using existing log group")
	}

	s.Done()
	return logGroup, nil
}
