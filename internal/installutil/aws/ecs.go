package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

const (
	defaultServerTagName     = "waypoint-server"
	defaultServerTagValue    = "server-component"
	ServerName               = "waypoint-server"
	DefaultStaticRunnerName  = "waypoint-static-runner"
	DefaultSecurityGroupName = "waypoint-server-security-group"
	RolePolicy               = `{
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
)

type Lifecycle struct {
	Init    func(terminal.UI) error
	Run     func(terminal.UI) error
	Cleanup func(terminal.UI) error
}

type NetworkInformation struct {
	VpcID   *string
	SgID    *string
	Subnets []*string
}

type EfsInformation struct {
	FileSystemID  *string
	AccessPointID *string
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

// SetupNetworking retrieves subnet information and creates the security group
// necessary for Waypoint.
func SetupNetworking(
	ctx context.Context,
	ui terminal.UI,
	sess *session.Session,
	subnet []string,
	ports []*int64,
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
	sgID, err := createSG(ctx, s, sess, DefaultSecurityGroupName, vpcID, ports)
	if err != nil {
		return nil, err
	}
	s.Update("Networking setup")
	s.Done()
	ni := NetworkInformation{
		VpcID:   vpcID,
		Subnets: subnets,
		SgID:    sgID,
	}
	return &ni, nil
}

func SetupEFS(
	ctx context.Context,
	ui terminal.UI,
	sess *session.Session,
	netInfo *NetworkInformation,
	efsTags []*efs.Tag,
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
		Tags:          efsTags,
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
	for _, sub := range netInfo.Subnets {
		_, err := efsSvc.CreateMountTarget(&efs.CreateMountTargetInput{
			FileSystemId:   fsd.FileSystemId,
			SecurityGroups: []*string{netInfo.SgID},
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
	// TODO: Change path to not always include "server"
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
		if available == len(netInfo.Subnets) {
			break
		}

		available = 0
		time.Sleep(5 * time.Second)
		continue
	}

	if available != len(netInfo.Subnets) {
		return nil, fmt.Errorf("not enough available mount targets found")
	}

	s.Update("EFS ready")
	s.Done()
	return &EfsInformation{
		FileSystemID:  fsd.FileSystemId,
		AccessPointID: accessPoint.AccessPointId,
	}, nil
}

func CreateService(serviceInput *ecs.CreateServiceInput, ecsSvc *ecs.ECS) (*ecs.Service, error) {
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

// TODO: Add runner ID as tag
// SetupExecutionRole creates the necessary IAM execution role for Waypoint, if it does not exist
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
		AssumeRolePolicyDocument: aws.String(RolePolicy),
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

// subnetInfo returns the subnets and VPC to the caller. If no subnets
// were provided as input, then the default subnets are returned.
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

// DeleteEcsCommonResources deletes the provided ECS service and task definition
func DeleteEcsCommonResources(
	ctx context.Context,
	sess *session.Session,
	clusterArn string,
	resources []*resourcegroups.ResourceIdentifier,
) error {
	ecsSvc := ecs.New(sess)

	var serviceArn string
	for _, r := range resources {
		if *r.ResourceType == "AWS::ECS::Service" {
			serviceArn = *r.ResourceArn
		}
	}
	if serviceArn == "" {
		return nil
	}

	_, err := ecsSvc.DeleteService(&ecs.DeleteServiceInput{
		Service: &serviceArn,
		Force:   aws.Bool(true),
		Cluster: &clusterArn,
	})
	if err != nil {
		return err
	}

	runningTasks, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &clusterArn,
		DesiredStatus: aws.String(ecs.DesiredStatusRunning),
	})
	if err != nil {
		return err
	}

	for _, task := range runningTasks.TaskArns {
		_, err := ecsSvc.StopTask(&ecs.StopTaskInput{
			Cluster: &clusterArn,
			Task:    task,
		})
		if err != nil {
			return err
		}
	}

	err = ecsSvc.WaitUntilServicesInactive(&ecs.DescribeServicesInput{
		Cluster:  &clusterArn,
		Services: []*string{&serviceArn},
	})
	if err != nil {
		return err
	}
	for _, r := range resources {
		if *r.ResourceType == "AWS::ECS::TaskDefinition" {
			_, err := ecsSvc.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
				TaskDefinition: r.ResourceArn,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteCWLResources(
	ctx context.Context,
	sess *session.Session,
	logGroup string,
) error {
	cwlSvc := cloudwatchlogs.New(sess)

	_, err := cwlSvc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroup),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceNotFoundException":
				// the log group has already been destroyed
				return nil
			}
		}
		return err
	}
	return nil
}

func DeleteEcsResources(
	ctx context.Context,
	sess *session.Session,
	resources []*resourcegroups.ResourceIdentifier,
) error {
	ecsSvc := ecs.New(sess)

	var clusterArn string
	for _, r := range resources {
		if *r.ResourceType == "AWS::ECS::Cluster" {
			clusterArn = *r.ResourceArn
		}
	}
	if err := DeleteEcsCommonResources(ctx, sess, clusterArn, resources); err != nil {
		return err
	}

	_, err := ecsSvc.DeleteCluster(&ecs.DeleteClusterInput{
		Cluster: &clusterArn,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				// the cluster has already been destroyed
				return nil
			}
		}
		return err
	}

	return nil
}

func FindServices(serviceNames []string, ecsSvc *ecs.ECS, cluster string) ([]*ecs.Service, error) {
	var services []*ecs.Service
	for _, serviceName := range serviceNames {
		ss, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: []*string{aws.String(serviceName)},
		})
		if err != nil {
			return nil, err
		}
		if len(ss.Services) > 0 {
			services = append(services, ss.Services...)
		}
	}
	return services, nil
}
