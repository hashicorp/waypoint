// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runnerinstall

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	awsinstallutil "github.com/hashicorp/waypoint/internal/installutil/aws"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

const (
	defaultRunnerLogGroup  = "waypoint-runner-logs"
	defaultTaskRuntime     = "FARGATE"
	defaultRunnerTagValue  = "runner-component"
	defaultRunnerIdTagName = "runner-id"
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
        "ecr:CreateRepository",
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
        "ecs:StopTask",
	"elasticloadbalancing:AddTags",
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
	"elasticloadbalancing:DescribeTargetHealth",
        "iam:AttachRolePolicy",
        "iam:CreateRole",
        "iam:GetRole",
        "iam:ListAttachedRolePolicies",
        "iam:PassRole",
        "logs:CreateLogGroup",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:GetLogEvents",
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "*"
    }
  ]
}`

type ECSRunnerInstaller struct {
	Config  EcsConfig
	netInfo *awsinstallutil.NetworkInformation
}

type EcsConfig struct {
	Region            string   `hcl:"region,required"`
	Cluster           string   `hcl:"cluster,required"`
	ExecutionRoleName string   `hcl:"execution_role_name,optional"`
	TaskRoleName      string   `hcl:"task_role_name,optional"`
	CPU               string   `hcl:"runner_cpu,optional"`
	Memory            string   `hcl:"memory_cpu,optional"`
	RunnerImage       string   `hcl:"runner_image,optional"`
	Subnets           []string `hcl:"subnets,optional"`
}

func (i *ECSRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI
	log := opts.Log

	if i.Config.Cluster == "" {
		return errors.New("cluster name not specified")
	} else if i.Config.Region == "" {
		return errors.New("region not specified")
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.Config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	var (
		efsInfo       *awsinstallutil.EfsInformation
		logGroup      string
		executionRole string
		netInfo       *awsinstallutil.NetworkInformation
		taskRole      string
		runSvcArn     *string
	)
	lf := &awsinstallutil.Lifecycle{
		Init: func(ui terminal.UI) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.Config.Region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			// TODO: Add a port here if the user specified the -liveness-tcp-addr flag
			// after `--` on the runner install CLI - currently without this, there
			// is no way to use that flag effectively, since the SG won't allow traffic
			// to the port of the -liveness-tcp-addr flag
			if netInfo, err = awsinstallutil.SetupNetworking(ctx, ui, sess, i.Config.Subnets, []*int64{aws.Int64(int64(2049))}); err != nil {
				return err
			}
			i.netInfo = netInfo

			efsTags := []*efs.Tag{
				{
					Key:   aws.String(defaultRunnerTagName),
					Value: aws.String(defaultRunnerTagValue),
				},
				{
					Key:   aws.String(defaultRunnerIdTagName),
					Value: aws.String(opts.Id),
				},
			}
			if efsInfo, err = awsinstallutil.SetupEFS(ctx, ui, sess, netInfo, efsTags); err != nil {
				return err
			}

			if executionRole, err = awsinstallutil.SetupExecutionRole(ctx, ui, log, sess, i.Config.ExecutionRoleName); err != nil {
				return err
			}

			taskRole, err = i.setupTaskRole(ctx, ui, log, sess, opts.Id)
			if err != nil {
				return err
			}

			logGroup, err = awsinstallutil.SetupLogs(ctx, ui, log, sess, defaultRunnerLogGroup)
			if err != nil {
				return err
			}

			return nil
		},

		Run: func(ui terminal.UI) error {
			runSvcArn, err = launchRunner(
				ctx, ui, log, sess,
				opts.AdvertiseClient.Env(),
				opts.RunnerAgentFlags,
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
		Name:    "ecs-execution-role-name",
		Target:  &i.Config.ExecutionRoleName,
		Usage:   "The name of the execution task IAM Role to associate with the ECS Service.",
		Default: "waypoint-runner-execution-role",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-task-role-name",
		Target:  &i.Config.TaskRoleName,
		Usage:   "IAM Execution Role to assign to the on-demand runner.",
		Default: runnerName,
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
		Default: defaultRunnerImage,
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

// Uninstall deletes the waypoint-runner service from AWS ECS, and its
// associated volume from EFS. The log group, execution role, subnets,
// and ECS cluster are not deleted.
func (i *ECSRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Finding runner in ECS services...")
	defer func() { s.Abort() }()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.Config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	// Find clusterArn which waypoint runner is installed into
	// We check for the serviceName before v0.9 and v0.9+
	ecsSvc := ecs.New(sess)
	serviceNames := []string{
		defaultRunnerTagName,
		installutil.DefaultRunnerName(opts.Id),
	}

	services, err := awsinstallutil.FindServices(serviceNames, ecsSvc, i.Config.Cluster)
	if err != nil {
		log.Debug("Unable to find desired runner; %s", err)
		ui.Output("Could not get list of ECS services: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	for _, service := range services {
		// Delete associated runner service and tasks
		// This does not remove the security group since it may be in use by other
		// runners/waypoint infrastructure.
		s.Update("Deleting runner service")
		_, err = ecsSvc.DeleteService(&ecs.DeleteServiceInput{
			Service: service.ServiceArn,
			Force:   aws.Bool(true),
			Cluster: service.ClusterArn,
		})
		if err != nil {
			log.Debug("error deleting runner service: %s", err)
			ui.Output("Error Deleting Runner service: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return err
		}

		s.Update("Waiting for runner service to be inactive")
		err = ecsSvc.WaitUntilServicesInactive(&ecs.DescribeServicesInput{
			Cluster:  service.ClusterArn,
			Services: []*string{service.ServiceArn},
		})
		if err != nil {
			log.Debug("error waiting for runner to deactivate: %s", err)
			ui.Output("Error waiting for Runner service to deactivate: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return err
		}

		s.Update("Waypoint runner AWS ECS service deleted")
	}
	s.Done()

	// TODO: Still attempt to delete the EFS volume if the ECS service
	// uninstall fails
	s = sg.Add("Deleting runner file system")
	efsSvc := efs.New(sess)
	var fileSystems []*efs.FileSystemDescription
	err = efsSvc.DescribeFileSystemsPages(&efs.DescribeFileSystemsInput{},
		func(page *efs.DescribeFileSystemsOutput, lastPage bool) bool {
			fileSystems = append(fileSystems, page.FileSystems...)
			return !lastPage
		})
	if err != nil {
		return err
	}

	if len(fileSystems) == 0 {
		s.Update("No file systems detected, skipping deletion")
		s.Done()
		return nil
	}

	var fileSystemId *string
	for _, fileSystem := range fileSystems {
		// Check if tags match ID, if so then delete things
		for _, tag := range fileSystem.Tags {
			if *tag.Key == "runner-id" && *tag.Value == opts.Id {
				fileSystemId = fileSystem.FileSystemId
				// This goto skips to the logic for deleting the file system -
				// we know which one we need to delete now, so there's no need
				// to iterate through any additional fileSystems
				goto DeleteFileSystem
			}
		}
	}

	if fileSystemId == nil || *fileSystemId == "" {
		s.Update("File system with tag key `runner-id` and value " + opts.Id + " not detected, skipping deletion")
		s.Done()
		return nil
	}

DeleteFileSystem:
	describeAccessPointsResp, err := efsSvc.DescribeAccessPoints(&efs.DescribeAccessPointsInput{
		FileSystemId: fileSystemId,
	})
	if err != nil {
		return err
	}
	for _, accessPoint := range describeAccessPointsResp.AccessPoints {
		_, err = efsSvc.DeleteAccessPoint(&efs.DeleteAccessPointInput{AccessPointId: accessPoint.AccessPointId})
		if err != nil {
			return err
		}
	}

	describeMountTargetsResp, err := efsSvc.DescribeMountTargets(&efs.DescribeMountTargetsInput{
		FileSystemId: fileSystemId,
	})
	if err != nil {
		return err
	}
	for _, mountTarget := range describeMountTargetsResp.MountTargets {
		_, err = efsSvc.DeleteMountTarget(&efs.DeleteMountTargetInput{MountTargetId: mountTarget.MountTargetId})
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return errors.New("after 5 minutes, the file system could" +
				"not be deleted, because the mount targets weren't deleted")
		default:
			_, err = efsSvc.DeleteFileSystem(&efs.DeleteFileSystemInput{FileSystemId: fileSystemId})
			if err != nil {
				if strings.Contains(err.Error(), "because it has mount targets") {
					// sleep here for 5 seconds to avoid slamming the API
					time.Sleep(5 * time.Second)
					continue
				}
				return err
			}
			// if we reach this point, we're done
			s.Update("Runner file system deleted")
			s.Done()
			return nil
		}
	}
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
	env, extraArgs []string,
	executionRoleArn, taskRoleArn, logGroup, region, cpu, memory, runnerImage, cluster, cookie, id string,
	netInfo *awsinstallutil.NetworkInformation,
	efsInfo *awsinstallutil.EfsInformation,
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

	var args []*string
	for _, arg := range extraArgs {
		args = append(args, aws.String(arg))
	}

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Command: append([]*string{
			aws.String("runner"),
			aws.String("agent"),
			aws.String("-id=" + id),
			aws.String("-liveness-tcp-addr=:1234"),
			aws.String("-cookie=" + cookie),
			aws.String("-state-dir=/data/runner"),
			aws.String("-vv"),
		}, args...),
		Name:  aws.String(runnerName),
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
		Family:                  aws.String(defaultRunnerTagName),
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
				Values: []*string{aws.String(awsinstallutil.DefaultSecurityGroupName)},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var groupId *string
	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
		s.Update("Using existing security group: %s", awsinstallutil.DefaultSecurityGroupName)
	} else {
		return nil, fmt.Errorf("could not find security group (%s)", awsinstallutil.DefaultSecurityGroupName)
	}

	// Check for details of possibly existing cluster `waypoint-server`
	// If server was installed to ECS with `waypoint install` command, we'd expect this
	// query what subnets and vpc information from the server service
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(awsinstallutil.ServerName)},
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
		ServiceName:          aws.String(installutil.DefaultRunnerName(id)),
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

	s.Update("Creating ECS Service (%s)", defaultRunnerTagName)
	svc, err := awsinstallutil.CreateService(createServiceInput, ecsSvc)
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

// setupTaskRole creates an IAM task role for launching on-demand runners
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
		AssumeRolePolicyDocument: aws.String(awsinstallutil.RolePolicy),
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

// OnDemandRunnerConfig implements OnDemandRunnerConfigProvider
func (i *ECSRunnerInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	// Generate some configuration. Some of the OnDemand configurations have
	// defaults so we should be fine to directly use them
	cfgMap := map[string]interface{}{
		"log_group":           defaultRunnerLogGroup,
		"execution_role_name": i.Config.ExecutionRoleName,
		"task_role_name":      i.Config.TaskRoleName,
		"cluster":             i.Config.Cluster,
		"region":              i.Config.Region,
		"odr_cpu":             i.Config.CPU,
		"odr_memory":          i.Config.Memory,
	}

	if i.netInfo != nil {
		var subnets []string
		for _, s := range i.netInfo.Subnets {
			subnets = append(subnets, *s)
		}
		cfgMap["subnets"] = strings.Join(subnets, ",")
		cfgMap["security_group_id"] = i.netInfo.SgID
	}

	// Marshal our config
	cfgJson, err := json.MarshalIndent(cfgMap, "", "\t")
	if err != nil {
		// This shouldn't happen cause we control our input. If it does,
		// just panic cause this will be in a `server install` CLI and
		// we want the user to report a bug.
		panic(err)
	}

	return &pb.OnDemandRunnerConfig{
		Name:         "ecs",
		OciUrl:       installutil.DefaultODRImage,
		PluginType:   "aws-ecs",
		Default:      false,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
	}
}

type Logging struct {
	CreateGroup      bool   `hcl:"create_group,optional"`
	StreamPrefix     string `hcl:"stream_prefix,optional"`
	DateTimeFormat   string `hcl:"datetime_format,optional"`
	MultilinePattern string `hcl:"multiline_pattern,optional"`
	Mode             string `hcl:"mode,optional"`
	MaxBufferSize    string `hcl:"max_buffer_size,optional"`
}
