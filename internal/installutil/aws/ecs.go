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
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
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
	defaultHttpPort = "9702"
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

func SetupNetworking(
	ctx context.Context,
	ui terminal.UI,
	sess *session.Session,
	subnet []string,
) (*NetworkInformation, error) {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Setting up networking...")
	defer s.Abort()
	subnets, vpcID, err := subnetInfo(ctx, s, sess, subnet)
	if err != nil {
		return nil, err
	}

	s.Update("Setting up security group...")
	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)
	ports := []*int64{
		aws.Int64(int64(grpcPort)),
		aws.Int64(int64(httpPort)),
		aws.Int64(int64(2049)), // EFS File system port
	}

	sgID, err := createSG(ctx, s, sess, defaultSecurityGroupName, vpcID, ports)
	if err != nil {
		return nil, err
	}
	s.Update("Networking setup")
	s.Done()
	ni := NetworkInformation{
		vpcID:   vpcID,
		subnets: subnets,
		sgID:    sgID,
	}
	return &ni, nil
}

func SetupEFS(
	ctx context.Context,
	ui terminal.UI,
	sess *session.Session,
	netInfo *NetworkInformation,

) (*EfsInformation, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Creating new EFS file system...")
	defer func() { s.Abort() }()

	efsSvc := efs.New(sess)
	ulid, _ := component.Id()

	fsd, err := efsSvc.CreateFileSystem(&efs.CreateFileSystemInput{
		CreationToken: aws.String(ulid),
		Encrypted:     aws.Bool(true),
		Tags: []*efs.Tag{
			{
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = efsSvc.DescribeFileSystems(&efs.DescribeFileSystemsInput{
		CreationToken: aws.String(ulid),
	})
	if err != nil {
		return nil, err
	}
	s.Update("Created new EFS file system: %s", *fsd.FileSystemId)

EFSLOOP:
	for i := 0; i < 10; i++ {
		fsList, err := efsSvc.DescribeFileSystems(&efs.DescribeFileSystemsInput{
			FileSystemId: fsd.FileSystemId,
		})
		if err != nil {
			return nil, err
		}
		if len(fsList.FileSystems) == 0 {
			return nil, fmt.Errorf("file system (%s) not found", *fsd.FileSystemId)
		}
		// check the status of the first one
		fs := fsList.FileSystems[0]
		switch *fs.LifeCycleState {
		case efs.LifeCycleStateDeleted, efs.LifeCycleStateDeleting:
			return nil, fmt.Errorf("files system is deleting/deleted")
		case efs.LifeCycleStateAvailable:
			break EFSLOOP
		}
		time.Sleep(2 * time.Second)
	}

	s.Update("Creating EFS Mount targets...")

	// poll for available
	for _, sub := range netInfo.subnets {
		_, err := efsSvc.CreateMountTarget(&efs.CreateMountTargetInput{
			FileSystemId:   fsd.FileSystemId,
			SecurityGroups: []*string{netInfo.sgID},
			SubnetId:       sub,
			// Mount Targets do not support tags directly
		})
		if err != nil {
			return nil, fmt.Errorf("error creating mount target: %w", err)
		}
	}

	// create EFS access points
	s.Update("Creating EFS Access Point...")
	uid := aws.Int64(int64(100))
	gid := aws.Int64(int64(1000))
	accessPoint, err := efsSvc.CreateAccessPoint(&efs.CreateAccessPointInput{
		FileSystemId: fsd.FileSystemId,
		PosixUser: &efs.PosixUser{
			Uid: uid,
			Gid: gid,
		},
		RootDirectory: &efs.RootDirectory{
			CreationInfo: &efs.CreationInfo{
				OwnerUid:    uid,
				OwnerGid:    gid,
				Permissions: aws.String("755"),
			},
			Path: aws.String("/waypointserverdata"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating access point: %w", err)
	}

	// loop until all mount targets are ready, or the first container can have
	// issues starting
	// TODO: Update to use context instead of sleep
	s.Update("Waiting for EFS mount targets to become available...")
	var available int
	for i := 0; 1 < 30; i++ {
		mtgs, err := efsSvc.DescribeMountTargets(&efs.DescribeMountTargetsInput{
			AccessPointId: accessPoint.AccessPointId,
		})
		if err != nil {
			return nil, err
		}

		for _, m := range mtgs.MountTargets {
			if *m.LifeCycleState == efs.LifeCycleStateAvailable {
				available++
			}
		}
		if available == len(netInfo.subnets) {
			break
		}

		available = 0
		time.Sleep(5 * time.Second)
		continue
	}

	if available != len(netInfo.subnets) {
		return nil, fmt.Errorf("not enough available mount targets found")
	}

	s.Update("EFS ready")
	s.Done()
	return &EfsInformation{
		fileSystemID:  fsd.FileSystemId,
		accessPointID: accessPoint.AccessPointId,
	}, nil
}

type NetworkInformation struct {
	vpcID   *string
	sgID    *string
	subnets []*string
}

type EfsInformation struct {
	fileSystemID  *string
	accessPointID *string
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
	executionRoleArn, taskRoleArn, logGroup, region, cpu, memory, runnerImage, cluster, cookie, id string,
	netInfo *NetworkInformation,
	efsInfo *EfsInformation,
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

	// TODO: Add -state-dir
	// TODO: Add mount from EFS for state
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
		},
		Volumes: []*ecs.Volume{
			{
				EfsVolumeConfiguration: &ecs.EFSVolumeConfiguration{
					AuthorizationConfig: &ecs.EFSAuthorizationConfig{
						AccessPointId: efsInfo.accessPointID,
					},
					FileSystemId:      efsInfo.fileSystemID,
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

	// Check for details of possibly existing cluster `waypoint-server`
	// If server was installed to ECS with `waypoint install` command, we'd expect this
	// query what subnets and vpc information from the server service
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(serverName)},
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
		subnets = netInfo.subnets
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

func subnetInfo(
	ctx context.Context,
	s terminal.Step,
	sess *session.Session,
	subnet []string,
) ([]*string, *string, error) {
	ec2Svc := ec2.New(sess)

	var (
		subnets []*string
		vpcID   *string
	)

	if len(subnet) == 0 {
		s.Update("Using default subnets for Service networking")
		desc, err := ec2Svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("default-for-az"),
					Values: []*string{aws.String("true")},
				},
			},
		})
		if err != nil {
			return nil, nil, err
		}

		for _, subnet := range desc.Subnets {
			subnets = append(subnets, subnet.SubnetId)
		}
		if len(desc.Subnets) == 0 {
			return nil, nil, fmt.Errorf("no default subnet information found")
		}
		vpcID = desc.Subnets[0].VpcId
		return subnets, vpcID, nil
	}

	subnets = make([]*string, len(subnet))
	for j := range subnet {
		subnets[j] = &subnet[j]
	}
	s.Update("Using provided subnets for Service networking")
	subnetInfo, err := ec2Svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: subnets,
	})
	if err != nil {
		return nil, nil, err
	}

	if len(subnetInfo.Subnets) == 0 {
		return nil, nil, fmt.Errorf("no subnet information found for provided subnets")
	}

	vpcID = subnetInfo.Subnets[0].VpcId

	return subnets, vpcID, nil
}

func createSG(
	ctx context.Context,
	s terminal.Step,
	sess *session.Session,
	name string,
	vpcId *string,

	ports []*int64,
) (*string, error) {
	ec2srv := ec2.New(sess)

	dsg, err := ec2srv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(name)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{vpcId},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var groupId *string

	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
		s.Update("Using existing security group: %s", name)
	} else {
		s.Update("Creating security group: %s", name)
		out, err := ec2srv.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			Description: aws.String("created by waypoint"),
			GroupName:   aws.String(name),
			VpcId:       vpcId,
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String(ec2.ResourceTypeSecurityGroup),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String(defaultServerTagName),
							Value: aws.String(defaultServerTagValue),
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		groupId = out.GroupId
		s.Update("Created security group: %s", name)
	}

	s.Update("Authorizing ports to security group")
	// Port 2049 is the port for accessing EFS file systems over NFS
	for _, port := range ports {
		_, err = ec2srv.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			CidrIp:     aws.String("0.0.0.0/0"),
			FromPort:   port,
			ToPort:     port,
			GroupId:    groupId,
			IpProtocol: aws.String("tcp"),
		})
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidPermission.Duplicate":
				// fine, means we already added it.
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return groupId, nil
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
