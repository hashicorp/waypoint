package runnerinstall

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	installutil "github.com/hashicorp/waypoint/internal/installutil/aws"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

const (
	defaultRunnerLogGroup = "waypoint-runner-logs"
	defaultRunnerTagName  = "waypoint-runner"
	defaultTaskRuntime    = "FARGATE"
	defaultRunnerTagValue = "runner-component"
)

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

type ECSRunnerInstaller struct {
	Config EcsConfig
}

type EcsConfig struct {
	Region            string   `hcl:"region,required"`
	ExecutionRoleName string   `hcl:"execution_role_name,optional"`
	TaskRoleName      string   `hcl:"task_role_name,optional"`
	CPU               string   `hcl:"runner_cpu,optional"`
	Memory            string   `hcl:"memory_cpu,optional"`
	RunnerImage       string   `hcl:"runner_image,optional"`
	Cluster           string   `hcl:"cluster,optional"`
	Subnets           []string `hcl:"subnets,optional"`
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
		efsInfo       *installutil.EfsInformation
		logGroup      string
		executionRole string
		netInfo       *installutil.NetworkInformation
		taskRole      string
		runSvcArn     *string
	)
	lf := &installutil.Lifecycle{
		Init: func(ui terminal.UI) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.Config.Region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			if netInfo, err = installutil.SetupNetworking(ctx, ui, sess, i.Config.Subnets); err != nil {
				return err
			}

			if efsInfo, err = installutil.SetupEFS(ctx, ui, sess, netInfo); err != nil {
				return err
			}

			if executionRole, err = installutil.SetupExecutionRole(ctx, ui, log, sess, i.Config.ExecutionRoleName); err != nil {
				return err
			}

			taskRole, err = i.setupTaskRole(ctx, ui, log, sess, opts.Id)
			if err != nil {
				return err
			}

			logGroup, err = installutil.SetupLogs(ctx, ui, log, sess, defaultRunnerLogGroup)
			if err != nil {
				return err
			}

			return nil
		},

		Run: func(ui terminal.UI) error {
			runSvcArn, err = launchRunner(
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
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Uninstalling Runner...")
	defer func() { s.Abort() }()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.Config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	// Find clusterArn which waypoint runner is installed into
	ecsSvc := ecs.New(sess)
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(i.Config.Cluster),
		Services: []*string{aws.String("waypoint-runner-" + opts.Id)},
	})
	if err != nil {
		s.Update("Could not find runner with ID %s", opts.Id)
		return err
	}
	if len(services.Services) != 1 {
		log.Debug("Unable to uninstall runner. Too many instances")
		s.Update("Expected 1 runner service, found %s.", len(services.Services))

		return fmt.Errorf("Expected 1 runner service, found %d.", len(services.Services))
	}
	clusterArn := services.Services[0].ClusterArn

	// Delete associated runner service and tasks
	// This does not remove the security group since it may be in use by other
	// runners/waypoint infrastructure.
	s.Update("Deleting runner service")
	_, err = ecsSvc.DeleteService(&ecs.DeleteServiceInput{
		Service: services.Services[0].ServiceArn,
		Force:   aws.Bool(true),
		Cluster: clusterArn,
	})
	if err != nil {
		s.Update("Unable to delete runner service.")
		return err
	}

	s.Update("Waiting for runner service to be inactive")
	err = ecsSvc.WaitUntilServicesInactive(&ecs.DescribeServicesInput{
		Cluster:  clusterArn,
		Services: []*string{services.Services[0].ServiceArn},
	})
	if err != nil {
		s.Update("Unable to verify runner uninstalled", len(services.Services))
		return err
	}

	s.Update("Runner uninstalled")
	s.Done()
	return nil
}

func (i *ECSRunnerInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:   "ecs-region",
		Usage:  "AWS region in which to install the Waypoint runner.",
		Target: &i.Config.Region,
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.Config.Cluster,
		Default: "waypoint-server",
		Usage:   "The name of the ECS Cluster to install the Waypoint runner into.",
	})

}

func launchRunner(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,
	env []string,
	executionRoleArn, taskRoleArn, logGroup, region, cpu, memory, runnerImage, cluster, cookie, id string,
	netInfo *installutil.NetworkInformation,
	efsInfo *installutil.EfsInformation,
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

	grpcPort, _ := strconv.Atoi(serverconfig.DefaultGRPCPort)

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
			aws.String("-id=" + id),
			aws.String("-liveness-tcp-addr=:1234"),
			aws.String("-cookie=" + cookie),
			aws.String("-state-dir=/data/runner"),
			aws.String("-vv"),
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
		MountPoints: []*ecs.MountPoint{
			{
				ContainerPath: aws.String("/data/runner"),
				ReadOnly:      aws.Bool(false),
				SourceVolume:  aws.String(defaultRunnerTagName),
			},
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
			{
				Key:   aws.String("runner-id"),
				Value: aws.String(id),
			},
		},
		Volumes: []*ecs.Volume{
			{
				EfsVolumeConfiguration: &ecs.EFSVolumeConfiguration{
					AuthorizationConfig: &ecs.EFSAuthorizationConfig{
						AccessPointId: efsInfo.AccessPointID,
					},
					FileSystemId:      efsInfo.FileSystemID,
					TransitEncryption: aws.String(ecs.EFSTransitEncryptionEnabled),
				},
				Name: aws.String(defaultRunnerTagName),
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
				Values: []*string{aws.String(installutil.DefaultSecurityGroupName)},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var groupId *string
	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
		s.Update("Using existing security group: %s", installutil.DefaultSecurityGroupName)
	} else {
		return nil, fmt.Errorf("could not find security group (%s)", installutil.DefaultSecurityGroupName)
	}

	// Check for details of possibly existing cluster `waypoint-server`
	// If server was installed to ECS with `waypoint install` command, we'd expect this
	// query what subnets and vpc information from the server service
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(installutil.ServerName)},
	})
	if err != nil {
		return nil, err
	}

	// If we found an ECS cluster `waypoint-server` use that
	// If not, use the configs passed to this method
	var clusterArn *string
	var subnets []*string
	if len(services.Services) == 0 {
		// return nil, fmt.Errorf("no waypoint-server service found")
		clusterArn = aws.String(cluster)
		subnets = netInfo.Subnets
	} else {
		service := services.Services[0]
		clusterArn = service.ClusterArn
		subnets = service.NetworkConfiguration.AwsvpcConfiguration.Subnets
	}

	createServiceInput := &ecs.CreateServiceInput{
		Cluster:              clusterArn,
		DesiredCount:         aws.Int64(1),
		LaunchType:           aws.String(defaultTaskRuntime),
		ServiceName:          aws.String(runnerName + "-" + id),
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
			{
				Key:   aws.String("runner-id"),
				Value: aws.String(id),
			},
		},
	}

	s.Update("Creating ECS Service (%s)", runnerName)
	svc, err := installutil.CreateService(createServiceInput, ecsSvc)
	if err != nil {
		return nil, err
	}
	s.Update("Runner service created")
	s.Done()

	return svc.ClusterArn, nil
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

func (i *ECSRunnerInstaller) setupTaskRole(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	sess *session.Session,
	id string,
) (string, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Setting up an IAM task role for ODR runners...")
	defer func() { s.Abort() }()

	svc := iam.New(sess)

	roleName := i.Config.TaskRoleName

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
		AssumeRolePolicyDocument: aws.String(installutil.RolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(defaultRunnerTagName),
				Value: aws.String(defaultRunnerTagValue),
			},
			{
				Key:   aws.String("runner-id"),
				Value: aws.String(id),
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

type Logging struct {
	CreateGroup      bool   `hcl:"create_group,optional"`
	StreamPrefix     string `hcl:"stream_prefix,optional"`
	DateTimeFormat   string `hcl:"datetime_format,optional"`
	MultilinePattern string `hcl:"multiline_pattern,optional"`
	Mode             string `hcl:"mode,optional"`
	MaxBufferSize    string `hcl:"max_buffer_size,optional"`
}
