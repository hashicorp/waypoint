package serverinstall

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

const (
	defaultRunnerLogGroup = "waypoint-runner-logs"
	defaultServerLogGroup = "waypoint-server-logs"

	defaultTaskCPU     = 512
	defaultTaskMemory  = 1024
	defaultTaskFamily  = "waypoint-server"
	defaultTaskRuntime = "FARGATE"

	defaultSecurityGroupName = "waypoint-server-security-group"
	defaultNLBName           = "waypoint-server-nlb"
)

var ulid string

func init() {
	// generate a ULID for use in FileSystem creation and unique tag value for
	// this installation
	ulid, _ = component.Id()
}

type ECSInstaller struct {
	config ecsConfig
}

type ecsConfig struct {
	// serverImage is the image/tag of the Waypoint server to use. Default is
	// hashicorp/waypoint:latest
	serverImage string `hcl:"server_image,optional"`

	// region defines which AWS region to use.
	region string `hcl:"region,optional"`

	// Name of the ECS cluster to install the service into. Defaults to
	// waypoint-server
	cluster string `hcl:"cluster,optional"`

	// Name of the execution task IAM Role to associate with the ECS Service
	executionRoleName string `hcl:"execution_role_name,optional"`

	// Subnets to place the service into. Defaults to the public subnets in the
	// default VPC.
	subnets []string `hcl:"subnets,optional"`
}

// Install is a method of ECSInstaller and implements the Installer interface to
// register a waypoint-server in a ecs cluster
func (i *ECSInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, error) {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Installing Waypoint server")
	defer func() { s.Abort() }()
	log.Info("Starting lifecycle")

	var (
		efsInfo                                *efsInformation
		executionRole, cluster, serverLogGroup string
		netInfo                                *networkInformation
		server                                 *ecsServer
		sess                                   *session.Session

		err error
	)

	lf := &Lifecycle{
		Init: func(s LifecycleStatus) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.config.region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			if netInfo, err = i.SetupNetworking(ctx, s, sess, ulid); err != nil {
				return err
			}

			if cluster, err = i.SetupCluster(ctx, s, sess, ulid); err != nil {
				return err
			}

			if efsInfo, err = i.SetupEFS(ctx, s, sess, netInfo, ulid); err != nil {
				return err
			}

			if executionRole, err = i.SetupExecutionRole(ctx, s, log, sess, ulid); err != nil {
				return err
			}

			if serverLogGroup, err = i.SetupLogs(ctx, s, log, sess, ulid, defaultServerLogGroup); err != nil {
				return err
			}

			return nil
		},

		Run: func(s LifecycleStatus) error {
			server, err = i.Launch(ctx, s, log, ui, sess, efsInfo, netInfo, executionRole, cluster, serverLogGroup, ulid)
			return err
		},

		Cleanup: func(s LifecycleStatus) error { return nil },
	}

	if err := lf.Execute(log, ui); err != nil {
		return nil, err
	}

	// Set our connection information
	grpcAddr := fmt.Sprintf("%s:%s", server.Url, grpcPort)
	httpAddr := fmt.Sprintf("%s:%s", server.Url, httpPort)
	// Set our advertise address
	advertiseAddr := pb.ServerConfig_AdvertiseAddr{
		Addr:          grpcAddr,
		Tls:           true,
		TlsSkipVerify: true,
	}
	contextConfig := clicontext.Config{
		Server: serverconfig.Client{
			Address:       grpcAddr,
			Tls:           true,
			TlsSkipVerify: true, // always for now
			Platform:      "ecs",
		},
	}
	s.Done()
	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Launch takes the previously created resource and launches the Waypoint server
// service
func (i *ECSInstaller) Launch(
	ctx context.Context,
	s LifecycleStatus,
	log hclog.Logger,
	ui terminal.UI,
	sess *session.Session,
	efsInfo *efsInformation,
	netInfo *networkInformation,
	executionRoleArn, clusterName, logGroup, ulid string,
) (*ecsServer, error) {
	s.Status("Installing Waypoint server into ECS...")

	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)
	s.Status("Creating Network resources...")
	nlb, err := createNLB(
		ctx, s, log, sess,
		netInfo.vpcID,
		aws.Int64(int64(grpcPort)),
		aws.Int64(int64(httpPort)),
		netInfo.subnets,
		ulid,
	)
	if err != nil {
		return nil, err
	}
	s.Update("Created Network resources")

	defaultStreamPrefix := fmt.Sprintf("waypoint-server-%d", time.Now().Nanosecond())
	logOptions := buildLoggingOptions(
		nil,
		i.config.region,
		logGroup,
		defaultStreamPrefix,
	)

	cmd := []*string{
		aws.String("server"),
		aws.String("run"),
		aws.String("-accept-tos"),
		aws.String("-vvv"),
		aws.String("-db=/waypoint-data/data.db"),
		aws.String(fmt.Sprintf("-listen-grpc=0.0.0.0:%d", grpcPort)),
		aws.String(fmt.Sprintf("-listen-http=0.0.0.0:%d", httpPort)),
	}

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Command:   cmd,
		Name:      aws.String(serverName),
		Image:     aws.String(i.config.serverImage),
		PortMappings: []*ecs.PortMapping{
			{
				ContainerPort: aws.Int64(int64(httpPort)),
			},
			{
				ContainerPort: aws.Int64(int64(grpcPort)),
			},
		},
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String(ecs.LogDriverAwslogs),
			Options:   logOptions,
		},
		MountPoints: []*ecs.MountPoint{
			{
				SourceVolume:  aws.String("waypointdata"),
				ContainerPath: aws.String("/waypoint-data"),
			},
		},
	}

	// Create mount points for the EFS file system. The EFS mount targets need to
	// existin in a 1:1 pair with the subnets in use.
	s.Status("Creating ECS Service and Tasks...")
	log.Debug("registering task definition", "ulid", ulid)

	cpus := aws.String(strconv.Itoa(defaultTaskCPU))
	mems := strconv.Itoa(defaultTaskMemory)

	s.Update("Registering Task definition: %s", defaultTaskFamily)

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    []*ecs.ContainerDefinition{&def},
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Cpu:                     cpus,
		Memory:                  aws.String(mems),
		Family:                  aws.String(defaultTaskFamily),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(defaultTaskRuntime)},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
		Volumes: []*ecs.Volume{
			{
				Name: aws.String("waypointdata"),
				EfsVolumeConfiguration: &ecs.EFSVolumeConfiguration{
					TransitEncryption: aws.String(ecs.EFSTransitEncryptionEnabled),
					FileSystemId:      efsInfo.fileSystemID,
					AuthorizationConfig: &ecs.EFSAuthorizationConfig{
						AccessPointId: efsInfo.accessPointID,
					},
				},
			},
		},
	}

	ecsSvc := ecs.New(sess)
	taskDef, err := registerTaskDefinition(&registerTaskDefinitionInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	taskDefArn := *taskDef.TaskDefinitionArn

	// Create the service
	s.Update("Creating Service...")
	log.Debug("creating service", "arn", *taskDef.TaskDefinitionArn)

	createServiceInput := &ecs.CreateServiceInput{
		Cluster:                       &clusterName,
		DesiredCount:                  aws.Int64(1),
		LaunchType:                    aws.String(defaultTaskRuntime),
		ServiceName:                   aws.String(serverName),
		TaskDefinition:                aws.String(taskDefArn),
		HealthCheckGracePeriodSeconds: aws.Int64(int64(600)),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        netInfo.subnets,
				SecurityGroups: []*string{netInfo.sgID},
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
		LoadBalancers: []*ecs.LoadBalancer{
			{
				ContainerName:  aws.String("waypoint-server"),
				ContainerPort:  aws.Int64(int64(httpPort)),
				TargetGroupArn: aws.String(nlb.httpTgArn),
			},
			{
				ContainerName:  aws.String("waypoint-server"),
				ContainerPort:  aws.Int64(int64(grpcPort)),
				TargetGroupArn: aws.String(nlb.grpcTgArn),
			},
		},
	}

	s.Status("Creating ECS Service (%s, cluster-name: %s)", serviceName, clusterName)

	service, err := createService(createServiceInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	s.Update("Created ECS Service (%s, cluster-name: %s)", serviceName, clusterName)
	log.Debug("service started", "arn", service.ServiceArn)

	// after the service is created with the specified target groups, the load
	// balancer will start making health checks. Initial registration and health
	// checks can regularly take upwards of 5 minutes.
	s.Update("Waiting for target group to be healthy...")
	elbsrv := elbv2.New(sess)
	var healthy bool
	for i := 0; i < 80; i++ {
		health, err := elbsrv.DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{
			TargetGroupArn: &nlb.httpTgArn,
		})
		if err != nil {
			return nil, err
		}
		// it's possible no health descriptions are available yet
		if len(health.TargetHealthDescriptions) > 0 {
			// grab the first, most recent
			hd := health.TargetHealthDescriptions[0]

			if hd.TargetHealth.State != nil && *hd.TargetHealth.State == elbv2.TargetHealthStateEnumHealthy {
				healthy = true
				break
			}
		}
		time.Sleep(5 * time.Second)
	}

	if !healthy {
		return nil, fmt.Errorf("no healthy target group found")
	}
	s.Status("Service launched!")

	return &ecsServer{
		Url:                nlb.publicDNS,
		Cluster:            clusterName,
		TaskArn:            taskDefArn,
		HttpTargetGroupArn: nlb.httpTgArn,
		ServiceArn:         *service.ServiceArn,
	}, nil
}

// Upgrade is a method of ECSInstaller and implements the Installer interface to
// upgrade a waypoint-server in a ecs cluster
func (i *ECSInstaller) Upgrade(
	ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (
	*InstallResults, error,
) {
	ui := opts.UI
	log := opts.Log
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting ecs cluster...")
	defer s.Abort()

	// inspect current service - looking for image used in Task
	// Get Task definition
	var clusterArn string
	cluster := i.config.cluster
	ecsSvc := ecs.New(sess)

	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(cluster)},
	})
	if err != nil {
		return nil, err
	}

	var found bool
	for _, c := range desc.Clusters {
		if *c.ClusterName == cluster && strings.ToLower(*c.Status) == "active" {
			clusterArn = *c.ClusterArn
			found = true
			s.Update("Found existing ECS cluster: %s", cluster)
		}
	}
	if !found {
		return nil, fmt.Errorf("error: could not find ecs cluster")
	}
	// list the services to find the task descriptions
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(i.config.cluster),
		Services: []*string{aws.String(serverName)},
	})
	if err != nil {
		return nil, err
	}
	// should only find one
	serverSvc := services.Services[0]
	if serverSvc == nil {
		return nil, fmt.Errorf("no waypoint-server service found")
	}

	def, err := ecsSvc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		Include:        []*string{aws.String("TAGS")},
		TaskDefinition: serverSvc.TaskDefinition,
	})
	if err != nil {
		return nil, err
	}

	// assume only 1 task running here
	taskDef := def.TaskDefinition
	taskTags := def.Tags
	containerDef := taskDef.ContainerDefinitions[0]

	upgradeImg := defaultServerImage
	if i.config.serverImage != "" {
		upgradeImg = i.config.serverImage
	}
	// assume upgrade to latest
	if *containerDef.Image == defaultServerImage {
		// we can just update/force-deploy the service
		_, err := ecsSvc.UpdateService(&ecs.UpdateServiceInput{
			ForceNewDeployment:            aws.Bool(true),
			Cluster:                       &clusterArn,
			Service:                       serverSvc.ServiceName,
			HealthCheckGracePeriodSeconds: aws.Int64(int64(600)),
		})
		if err != nil {
			return nil, err
		}
		err = ecsSvc.WaitUntilServicesStable(&ecs.DescribeServicesInput{
			Cluster:  &clusterArn,
			Services: []*string{serverSvc.ServiceName},
		})
		if err != nil {
			return nil, err
		}
	} else {
		containerDef.Image = &upgradeImg
		// update task definition

		taskDef.SetContainerDefinitions([]*ecs.ContainerDefinition{containerDef})
		registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
			ContainerDefinitions:    taskDef.ContainerDefinitions,
			Cpu:                     taskDef.Cpu,
			ExecutionRoleArn:        taskDef.ExecutionRoleArn,
			Family:                  taskDef.Family,
			Memory:                  taskDef.Memory,
			NetworkMode:             aws.String("awsvpc"),
			RequiresCompatibilities: []*string{aws.String(defaultTaskRuntime)},
			Tags:                    taskTags,
			Volumes:                 taskDef.Volumes,
		}

		ecsSvc := ecs.New(sess)
		taskDef, err := registerTaskDefinition(&registerTaskDefinitionInput, ecsSvc)
		if err != nil {
			return nil, err
		}

		_, err = ecsSvc.UpdateService(&ecs.UpdateServiceInput{
			Cluster:        &clusterArn,
			TaskDefinition: taskDef.TaskDefinitionArn,
			Service:        serverSvc.ServiceName,
		})
		if err != nil {
			return nil, err
		}
		err = ecsSvc.WaitUntilServicesStable(&ecs.DescribeServicesInput{
			Cluster:  &clusterArn,
			Services: []*string{serverSvc.ServiceName},
		})
		if err != nil {
			return nil, err
		}
	}

	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	advertiseAddr.Addr = serverCfg.Address
	advertiseAddr.Tls = true
	advertiseAddr.TlsSkipVerify = true
	contextConfig = clicontext.Config{
		Server: serverCfg,
	}
	httpAddr := strings.Replace(serverCfg.Address, "9701", "9702", 1)

	s.Done()
	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Uninstall is a method of ECSInstaller and implements the Installer interface
// to remove a waypoint-server statefulset and the associated PVC and service
// from a ecs cluster
func (i *ECSInstaller) Uninstall(
	ctx context.Context,
	opts *InstallOpts,
) error {
	ui := opts.UI
	log := opts.Log
	s := new(lStatus)
	s.ui = ui
	defer s.Close()

	s.Status("Uninstalling Server resources...")

	// Get list of resources created with either the waypoint-server, or
	// waypoint-runner tag
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.region,
		Logger: log,
	})
	if err != nil {
		return err
	}
	rgSvc := resourcegroups.New(sess)

	query := fmt.Sprintf(defaultResourceQuery, serverName)
	results, err := rgSvc.SearchResources(&resourcegroups.SearchResourcesInput{
		ResourceQuery: &resourcegroups.ResourceQuery{
			Type:  aws.String(resourcegroups.QueryTypeTagFilters10),
			Query: aws.String(query),
		},
	})
	if err != nil {
		return err
	}

	resources := results.ResourceIdentifiers

	// Start destroying things. Some cannot be destroyed before others. The
	// general order to destroy things:
	// - ECS Service
	// - ECS Cluster
	// - Cloudwatch Log Group
	// - ELB Target Groups
	// - ELB Network Load Balancer
	// - EFS File System

	s.Status("Deleting ECS resources...")
	if err := deleteEcsResources(ctx, s, sess, resources); err != nil {
		return err
	}
	s.Status("Deleting Cloud Watch Log Group resources...")
	if err := deleteCWLResources(ctx, s, sess, defaultServerLogGroup); err != nil {
		return err
	}
	s.Status("Deleting EFS resources...")
	if err := deleteEFSResources(ctx, s, sess, resources); err != nil {
		return err
	}
	s.Status("Deleting Network resources...")
	if err := deleteNLBResources(ctx, s, sess, resources); err != nil {
		return err
	}

	s.Status("Server resources deleted")
	return nil
}

func deleteEFSResources(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	resources []*resourcegroups.ResourceIdentifier,
) error {
	// 	"AWS::EFS::FileSystem",
	var id string
	for _, r := range resources {
		if *r.ResourceType == "AWS::EFS::FileSystem" {
			id = nameFromArn(*r.ResourceArn)
			break
		}
	}
	efsSvc := efs.New(sess)
	mtgs, err := efsSvc.DescribeMountTargets(&efs.DescribeMountTargetsInput{
		FileSystemId: &id,
	})
	if err != nil {
		return err
	}
	s.Update("Deleting EFS Mount Targets...")
	for _, mt := range mtgs.MountTargets {
		_, err := efsSvc.DeleteMountTarget(&efs.DeleteMountTargetInput{
			MountTargetId: mt.MountTargetId,
		})
		if err != nil {
			return err
		}
	}

	for i := 0; 1 < 30; i++ {
		mtgs, err := efsSvc.DescribeMountTargets(&efs.DescribeMountTargetsInput{
			FileSystemId: &id,
		})
		if err != nil {
			return err
		}

		var deleted int
		mtgCount := len(mtgs.MountTargets)

		for _, m := range mtgs.MountTargets {
			if *m.LifeCycleState == efs.LifeCycleStateDeleted {
				deleted++
			}
		}
		if mtgCount == 0 {
			break
		}

		if deleted == mtgCount {
		}

		time.Sleep(5 * time.Second)
		continue
	}

	s.Update("Deleting EFS File System...")
	_, err = efsSvc.DeleteFileSystem(&efs.DeleteFileSystemInput{
		FileSystemId: &id,
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteNLBResources(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	resources []*resourcegroups.ResourceIdentifier,
) error {
	elbSvc := elbv2.New(sess)
	for _, r := range resources {
		if *r.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer" {
			s.Update("Deleting Network Load Balancer LISTENERS", *r.ResourceArn)
			results, err := elbSvc.DescribeListeners(&elbv2.DescribeListenersInput{
				LoadBalancerArn: r.ResourceArn,
			})
			if err != nil {
				return err
			}
			for _, l := range results.Listeners {
				_, err := elbSvc.DeleteListener(&elbv2.DeleteListenerInput{
					ListenerArn: l.ListenerArn,
				})
				if err != nil {
					return err
				}
			}
			s.Update("Deleting Network Load Balancer ", *r.ResourceArn)
			_, err = elbSvc.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: r.ResourceArn,
			})
			if err != nil {
				return err
			}
		}
	}
	s.Update("Deleting Target Groups...")
	for _, r := range resources {
		if *r.ResourceType == "AWS::ElasticLoadBalancingV2::TargetGroup" {
			_, err := elbSvc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
				TargetGroupArn: r.ResourceArn,
			})
			if err != nil {
				return err
			}
		}
	}

	ec2Svc := ec2.New(sess)
	results, err := ec2Svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String("waypoint-server")},
			},
		},
	})
	if err != nil {
		return err
	}
	if len(results.SecurityGroups) > 0 {
		for _, g := range results.SecurityGroups {
			for i := 0; i < 20; i++ {
				s.Update("Deleting Security Group (%s)", *g.GroupId)
				_, err := ec2Svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
					GroupId: g.GroupId,
				})
				// if we encounter an unrecoverable error, exit now.
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "DependencyViolation":
						time.Sleep(2 * time.Second)
						continue
					default:
						return err
					}
				}
				return err
			}
		}
	}

	return nil
}

func nameFromArn(arn string) string {
	parts := strings.Split(arn, ":")
	last := parts[len(parts)-1]
	parts = strings.Split(last, "/")
	return parts[len(parts)-1]
}

func deleteCWLResources(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	logGroup string,
) error {
	cwlSvc := cloudwatchlogs.New(sess)

	s.Update("Deleting Log Group %s", logGroup)
	_, err := cwlSvc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroup),
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteEcsResources(
	ctx context.Context,
	s LifecycleStatus,
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
	if err := deleteEcsCommonResources(ctx, s, sess, clusterArn, resources); err != nil {
		return err
	}

	s.Update("Deleting ECS cluster: %s", clusterArn)
	_, err := ecsSvc.DeleteCluster(&ecs.DeleteClusterInput{
		Cluster: &clusterArn,
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteEcsCommonResources(
	ctx context.Context,
	s LifecycleStatus,
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
		s.Update("No runner ECS services found")
		return nil
	}

	s.Update("Deleting ECS service: %s", serviceArn)
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
			s.Update("Deregistering ECS task: %s", *r.ResourceArn)
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

// InstallRunner implements Installer.
func (i *ECSInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	ui := opts.UI
	log := opts.Log

	sg := ui.StepGroup()
	defer sg.Wait()
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	s := sg.Add("Starting Runner installation")
	defer func() { s.Abort() }()
	var (
		logGroup      string
		executionRole string
		runSvcArn     *string
	)
	lf := &Lifecycle{
		Init: func(s LifecycleStatus) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.config.region,
				Logger: log,
			})
			if err != nil {
				return err
			}
			executionRole, err = i.SetupExecutionRole(ctx, s, log, sess, ulid)
			if err != nil {
				return err
			}

			logGroup, err = i.SetupLogs(ctx, s, log, sess, ulid, defaultRunnerLogGroup)
			if err != nil {
				return err
			}

			return nil
		},

		Run: func(s LifecycleStatus) error {
			runSvcArn, err = i.LaunchRunner(
				ctx, s, log, ui, sess,
				opts.AdvertiseClient.Env(),
				executionRole,
				logGroup,
				ulid,
			)
			return err
		},

		Cleanup: func(s LifecycleStatus) error { return nil },
	}

	if err := lf.Execute(log, ui); err != nil {
		return err
	}

	s.Update("Created ECS Service (%s)", runnerName)
	log.Debug("runner service started", "arn", *runSvcArn)

	s.Done()

	return nil
}

var defaultResourceQuery = "{\"ResourceTypeFilters\":[\"AWS::AllSupported\"],\"TagFilters\":[{\"Key\":\"%s\",\"Values\":[]}]}"

func (i *ECSInstaller) UninstallRunner(
	ctx context.Context,
	opts *InstallOpts,
) error {
	ui := opts.UI
	log := opts.Log
	s := new(lStatus)
	s.ui = ui
	defer s.Close()
	s.Status("Uninstalling Runner resources...")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.region,
		Logger: log,
	})
	if err != nil {
		return err
	}
	rgSvc := resourcegroups.New(sess)

	query := fmt.Sprintf(defaultResourceQuery, runnerName)
	results, err := rgSvc.SearchResources(&resourcegroups.SearchResourcesInput{
		ResourceQuery: &resourcegroups.ResourceQuery{
			Type:  aws.String(resourcegroups.QueryTypeTagFilters10),
			Query: aws.String(query),
		},
	})
	if err != nil {
		return err
	}

	resources := results.ResourceIdentifiers
	var clusterArn string
	for _, r := range resources {
		if *r.ResourceType == "AWS::ECS::Cluster" {
			clusterArn = *r.ResourceArn
		}
	}
	s.Status("Deleting ECS resources...")
	if err := deleteEcsCommonResources(ctx, s, sess, clusterArn, resources); err != nil {
		return err
	}
	s.Status("Deleting Cloud Watch Log Group resources...")
	if err := deleteCWLResources(ctx, s, sess, defaultRunnerLogGroup); err != nil {
		return err
	}
	return nil
}

// HasRunner implements Installer.
func (i *ECSInstaller) HasRunner(
	ctx context.Context,
	opts *InstallOpts,
) (bool, error) {
	log := opts.Log
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.region,
		Logger: log,
	})
	if err != nil {
		return false, err
	}
	ecsSvc := ecs.New(sess)
	// query what subnets and vpc information from the server service
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(i.config.cluster),
		Services: []*string{aws.String(runnerName)},
	})
	if err != nil {
		return false, err
	}

	return len(services.Services) > 0, nil
}

func (i *ECSInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.cluster,
		Usage:   "Configures the Cluster to install into.",
		Default: "waypoint-server",
	})
	set.StringVar(&flag.StringVar{
		Name:    "region",
		Target:  &i.config.region,
		Usage:   "Configures the region specific things.",
		Default: "us-west-2",
	})
	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "subnets",
		Target: &i.config.subnets,
		Usage:  "Subnets to install server into.",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-execution-role-name",
		Target:  &i.config.executionRoleName,
		Usage:   "Configures the Execution role name to use.",
		Default: "waypoint-server-execution-role",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})
}

func (i *ECSInstaller) UpgradeFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.cluster,
		Usage:   "Configures the Cluster to upgrade.",
		Default: "waypoint-server",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})
}

func (i *ECSInstaller) UninstallFlags(set *flag.Set) {
}

type Lifecycle struct {
	Init    func(LifecycleStatus) error
	Run     func(LifecycleStatus) error
	Cleanup func(LifecycleStatus) error
}

type lStatus struct {
	ui   terminal.UI
	sg   terminal.StepGroup
	step terminal.Step
}

func (l *lStatus) Status(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Done()
		l.step = nil
	}

	l.step = l.sg.Add(str, args...)
}

func (l *lStatus) Update(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Update(str, args...)
	} else {
		l.step = l.sg.Add(str, args)
	}
}

func (l *lStatus) Error(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Update(str, args...)
		l.step.Abort()
	} else {
		l.step = l.sg.Add(str, args)
		l.step.Abort()
	}

	l.step = nil
}

func (l *lStatus) Abort() error {
	if l.step != nil {
		l.step.Abort()
		l.step = nil
	}

	if l.sg != nil {
		l.sg.Wait()
		l.sg = nil
	}

	return nil
}

func (l *lStatus) Close() error {
	if l.step != nil {
		l.step.Done()
		l.step = nil
	}

	if l.sg != nil {
		l.sg.Wait()
		l.sg = nil
	}

	return nil
}

func (lf *Lifecycle) Execute(log hclog.Logger, ui terminal.UI) error {
	var l lStatus
	l.ui = ui

	defer l.Close()

	if lf.Init != nil {
		log.Debug("lifecycle init")

		err := lf.Init(&l)
		if err != nil {
			l.Abort()
			return err
		}

	}

	log.Debug("lifecycle run")
	err := lf.Run(&l)
	if err != nil {
		l.Abort()
		return err
	}

	if lf.Cleanup != nil {
		log.Debug("lifecycle cleanup")

		err = lf.Cleanup(&l)
		if err != nil {
			l.Abort()
			return err
		}
	}

	return nil
}

type LifecycleStatus interface {
	Status(str string, args ...interface{})
	Update(str string, args ...interface{})
	Error(str string, args ...interface{})
}

func (i *ECSInstaller) SetupNetworking(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	ulid string,
) (*networkInformation, error) {

	s.Status("Establishing subnets and security group...")
	subnets, vpcID, err := i.subnetInfo(ctx, s, sess)
	if err != nil {
		return nil, err
	}

	s.Status("Setting up security group...")
	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)
	ports := []*int64{
		aws.Int64(int64(grpcPort)),
		aws.Int64(int64(httpPort)),
		aws.Int64(int64(2049)), // EFS File system port
	}

	sgID, err := createSG(ctx, s, sess, defaultSecurityGroupName, vpcID, ulid, ports)
	if err != nil {
		return nil, err
	}
	return &networkInformation{
		vpcID:   vpcID,
		subnets: subnets,
		sgID:    sgID,
	}, nil
}

func (i *ECSInstaller) SetupCluster(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	ulid string,
) (string, error) {
	ecsSvc := ecs.New(sess)

	s.Status("Inspecting existing ECS clusters...")
	cluster := i.config.cluster

	// re-use an existing cluster if we have one
	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(cluster)},
	})
	if err != nil {
		return "", err
	}

	for _, c := range desc.Clusters {
		if *c.ClusterName == cluster && strings.ToLower(*c.Status) == "active" {
			s.Status("Found existing ECS cluster: %s", cluster)
			return cluster, nil
		}
	}

	s.Status("Creating new ECS cluster: %s", cluster)

	_, err = ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(cluster),
		// we need to tag with both the server and runner names, so we can properly
		// cleanup
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
			{
				Key:   aws.String(runnerName),
				Value: aws.String(ulid),
			},
		},
	})

	if err != nil {
		return "", err
	}

	s.Update("Created new ECS cluster: %s", cluster)
	return cluster, nil
}

func (i *ECSInstaller) SetupEFS(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	netInfo *networkInformation,
	ulid string,
) (*efsInformation, error) {
	efsSvc := efs.New(sess)

	s.Status("Creating new EFS file system...")
	fsd, err := efsSvc.CreateFileSystem(&efs.CreateFileSystemInput{
		CreationToken: aws.String(ulid),
		Encrypted:     aws.Bool(true),
		Tags: []*efs.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
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

	s.Status("Creating EFS Mount targets...")

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
		Tags: []*efs.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating access point: %w", err)
	}

	// loop until all mount targets are ready, or the first container can have
	// issues starting
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

	return &efsInformation{
		fileSystemID:  fsd.FileSystemId,
		accessPointID: accessPoint.AccessPointId,
	}, nil
}

type ecsServer struct {
	Url                string
	TaskArn            string
	ServiceArn         string
	HttpTargetGroupArn string
	GRPCTargetGroupArn string
	LoadBalancerArn    string
	Cluster            string
}

type networkInformation struct {
	vpcID   *string
	sgID    *string
	subnets []*string
}

type efsInformation struct {
	fileSystemID  *string
	accessPointID *string
}

type nlb struct {
	lbArn     string
	httpTgArn string
	grpcTgArn string
	publicDNS string
}

func createSG(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	name string,
	vpcId *string,
	ulid string,
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
							Key:   aws.String(serverName),
							Value: aws.String(ulid),
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

func (i *ECSInstaller) SetupLogs(
	ctx context.Context,
	s LifecycleStatus,
	log hclog.Logger,
	sess *session.Session,
	ulid string,
	logGroup string,
) (string, error) {
	cwl := cloudwatchlogs.New(sess)

	s.Status("Examining existing CloudWatchLogs groups...")

	groups, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              aws.Int64(1),
		LogGroupNamePrefix: aws.String(logGroup),
	})
	if err != nil {
		return "", err
	}

	if len(groups.LogGroups) == 0 {
		s.Status("Creating CloudWatchLogs group to store logs in: %s", logGroup)

		log.Debug("creating log group", "group", logGroup)
		_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroup),
		})
		if err != nil {
			return "", err
		}

		s.Update("Created CloudWatchLogs group to store logs in: %s", logGroup)
	}

	return logGroup, nil
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

type Logging struct {
	CreateGroup      bool   `hcl:"create_group,optional"`
	StreamPrefix     string `hcl:"stream_prefix,optional"`
	DateTimeFormat   string `hcl:"datetime_format,optional"`
	MultilinePattern string `hcl:"multiline_pattern,optional"`
	Mode             string `hcl:"mode,optional"`
	MaxBufferSize    string `hcl:"max_buffer_size,optional"`
}

func (i *ECSInstaller) SetupExecutionRole(
	ctx context.Context,
	s LifecycleStatus,
	log hclog.Logger,
	sess *session.Session,
	ulid string,
) (string, error) {
	svc := iam.New(sess)

	roleName := i.config.executionRoleName

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
		s.Status("Found existing IAM role to use: %s", roleName)
		return *getOut.Role.Arn, nil
	}

	log.Debug("creating new role")
	s.Status("Creating IAM role: %s", roleName)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return "", err
	}

	roleArn := *result.Role.Arn

	log.Debug("created new role", "arn", roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return "", err
	}

	log.Debug("attached execution role policy")

	s.Update("Created IAM role: %s", roleName)
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

// creates a network load balancer for grpc and http
func createNLB(
	ctx context.Context,
	s LifecycleStatus,
	log hclog.Logger,
	sess *session.Session,
	vpcId *string,
	grpcPort *int64,
	httpPort *int64,
	subnets []*string,
	ulid string,
) (serverNLB *nlb, err error) {

	s.Update("Creating NLB target groups")
	elbsrv := elbv2.New(sess)

	ctgGPRC, err := elbsrv.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:                    aws.String("waypoint-server-grpc"),
		Port:                    grpcPort,
		Protocol:                aws.String("TCP"),
		TargetType:              aws.String("ip"),
		HealthyThresholdCount:   aws.Int64(int64(2)),
		UnhealthyThresholdCount: aws.Int64(int64(2)),
		VpcId:                   vpcId,
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	htgGPRC, err := elbsrv.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:                    aws.String("waypoint-server-http"),
		Port:                    httpPort,
		Protocol:                aws.String("TCP"),
		TargetType:              aws.String("ip"),
		VpcId:                   vpcId,
		HealthCheckProtocol:     aws.String(elbv2.ProtocolEnumHttps),
		HealthCheckPath:         aws.String("/auth"),
		HealthyThresholdCount:   aws.Int64(int64(2)),
		UnhealthyThresholdCount: aws.Int64(int64(2)),
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	httpTgArn := htgGPRC.TargetGroups[0].TargetGroupArn
	grpcTgArn := ctgGPRC.TargetGroups[0].TargetGroupArn

	// Create the load balancer OR modify the existing one to have this new target
	// group but with a weight of 0

	htgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: httpTgArn,
		},
	}
	gtgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: grpcTgArn,
		},
	}

	var certs []*elbv2.Certificate

	var lb *elbv2.LoadBalancer

	dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(defaultNLBName)},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				// fine, means we'll create it.
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if dlb != nil && len(dlb.LoadBalancers) > 0 {
		lb = dlb.LoadBalancers[0]
		s.Update("Using existing NLB %s (%s, dns-name: %s)",
			defaultNLBName, *lb.LoadBalancerArn, *lb.DNSName)
	} else {
		s.Update("Creating new NLB: %s", defaultNLBName)

		scheme := elbv2.LoadBalancerSchemeEnumInternetFacing

		clb, err := elbsrv.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
			Name:    aws.String(defaultNLBName),
			Subnets: subnets,
			// SecurityGroups: []*string{sgWebId},
			Scheme: &scheme,
			Type:   aws.String(elbv2.LoadBalancerTypeEnumNetwork),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(serverName),
					Value: aws.String(ulid),
				},
			},
		})
		if err != nil {
			return nil, err
		}

		s.Update("Waiting on NLB to be active...")
		lb = clb.LoadBalancers[0]
		for i := 0; 1 < 70; i++ {
			clbd, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: []*string{lb.LoadBalancerArn},
			})
			if err != nil {
				return nil, err
			}
			lb = clbd.LoadBalancers[0]
			if lb.State != nil && *lb.State.Code == elbv2.LoadBalancerStateEnumActive {
				break
			}
			if lb.State != nil && *lb.State.Code == elbv2.LoadBalancerStateEnumFailed {
				return nil, fmt.Errorf("failed to create NLB")
			}

			time.Sleep(5 * time.Second)
		}

		if *lb.State.Code != elbv2.LoadBalancerStateEnumActive {
			return nil, fmt.Errorf("failed to create NLB in time, last state: (%s)", *lb.State.Code)
		}

		s.Update("Created new NLB: %s (dns-name: %s)", defaultNLBName, *lb.DNSName)
	}

	s.Update("Creating new NLB Listener")

	log.Info("load-balancer defined", "dns-name", *lb.DNSName)

	_, err = elbsrv.CreateListener(&elbv2.CreateListenerInput{
		LoadBalancerArn: lb.LoadBalancerArn,
		Port:            grpcPort,
		Protocol:        aws.String("TCP"),
		Certificates:    certs,
		DefaultActions: []*elbv2.Action{
			{
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: gtgs,
				},
				Type: aws.String("forward"),
			},
		},
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = elbsrv.CreateListener(&elbv2.CreateListenerInput{
		LoadBalancerArn: lb.LoadBalancerArn,
		Port:            aws.Int64(int64(9702)),
		Protocol:        aws.String("TCP"),
		Certificates:    certs,
		DefaultActions: []*elbv2.Action{
			{
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: htgs,
				},
				Type: aws.String("forward"),
			},
		},
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(serverName),
				Value: aws.String(ulid),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &nlb{
		lbArn:     *lb.LoadBalancerArn,
		httpTgArn: *httpTgArn,
		grpcTgArn: *grpcTgArn,
		publicDNS: *lb.DNSName,
	}, nil
}

func (i *ECSInstaller) subnetInfo(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
) ([]*string, *string, error) {
	ec2Svc := ec2.New(sess)

	var (
		subnets []*string
		vpcID   *string
	)

	if len(i.config.subnets) == 0 {
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

	subnets = make([]*string, len(i.config.subnets))
	for j := range i.config.subnets {
		subnets[j] = &i.config.subnets[j]
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

func registerTaskDefinition(def *ecs.RegisterTaskDefinitionInput, ecsSvc *ecs.ECS) (*ecs.TaskDefinition, error) {
	// AWS is eventually consistent so even though we probably created the
	// resources that are referenced by the task definition, it can error out if
	// we try to reference those resources too quickly. So we're forced to guard
	// actions which reference other AWS services with loops like this.
	var taskOut *ecs.RegisterTaskDefinitionOutput
	var err error
	for i := 0; i < 30; i++ {
		taskOut, err = ecsSvc.RegisterTaskDefinition(def)
		if err == nil {
			break
		}

		// if we encounter an unrecoverable error, exit now.
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "ResourceConflictException" {
				return nil, err
			}
		}

		// otherwise sleep and try again
		time.Sleep(2 * time.Second)
	}
	return taskOut.TaskDefinition, nil
}

func (i *ECSInstaller) LaunchRunner(
	ctx context.Context,
	s LifecycleStatus,
	log hclog.Logger,
	ui terminal.UI,
	sess *session.Session,
	env []string,
	executionRoleArn, logGroup, ulid string,
) (*string, error) {

	s.Status("Installing Waypoint runner into ECS...")

	defaultStreamPrefix := fmt.Sprintf("waypoint-runner-%d", time.Now().Nanosecond())
	logOptions := buildLoggingOptions(
		nil,
		i.config.region,
		logGroup,
		defaultStreamPrefix,
	)

	cpu := defaultTaskCPU
	mem := defaultTaskMemory
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
			aws.String("-vvv"),
			aws.String("-liveness-tcp-addr=:1234"),
		},
		Name:  aws.String("waypoint-runner"),
		Image: aws.String(i.config.serverImage),
		PortMappings: []*ecs.PortMapping{
			{
				ContainerPort: aws.Int64(int64(grpcPort)),
			},
		},
		Environment: envs,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String(ecs.LogDriverAwslogs),
			Options:   logOptions,
		},
	}

	cpus := aws.String(strconv.Itoa(cpu))
	mems := strconv.Itoa(mem)

	s.Update("Registering Task definition: waypoint-runner")

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{&def},

		ExecutionRoleArn: aws.String(executionRoleArn),
		Cpu:              cpus,
		Memory:           aws.String(mems),
		Family:           aws.String(runnerName),

		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(defaultTaskRuntime)},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(runnerName),
				Value: aws.String(runnerName),
			},
		},
	}

	ecsSvc := ecs.New(sess)
	taskDef, err := registerTaskDefinition(&registerTaskDefinitionInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	taskDefArn := *taskDef.TaskDefinitionArn
	s.Status("Creating Service...")
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
		Cluster:  aws.String(i.config.cluster),
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
		Cluster:        clusterArn,
		DesiredCount:   aws.Int64(1),
		LaunchType:     aws.String(defaultTaskRuntime),
		ServiceName:    aws.String(runnerName),
		TaskDefinition: aws.String(taskDefArn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: []*string{groupId},
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(runnerName),
				Value: aws.String(runnerName),
			},
		},
	}

	s.Update("Creating ECS Service (%s)", runnerName)

	svc, err := createService(createServiceInput, ecsSvc)
	if err != nil {
		return nil, err
	}

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
