// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverinstall

import (
	"context"
	json "encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/waypoint/internal/clierrors"

	"github.com/hashicorp/waypoint/internal/installutil"

	"github.com/hashicorp/waypoint/internal/runnerinstall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/internal/clicontext"
	awsinstallutil "github.com/hashicorp/waypoint/internal/installutil/aws"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

const (
	defaultRunnerLogGroup = "waypoint-runner-logs"
	defaultServerLogGroup = "waypoint-server-logs"

	defaultTaskFamily  = "waypoint-server"
	defaultTaskRuntime = "FARGATE"

	defaultNLBName = "waypoint-server-nlb"

	// These tags are used to tag resources as they are created in AWS. Two tags
	// are required, so that the UninstallRunner method(s) can query for runner
	// resources and retrieve the cluster ARN
	defaultServerTagName  = "waypoint-server"
	defaultServerTagValue = "server-component"
	defaultRunnerTagName  = "waypoint-runner"
	defaultRunnerTagValue = "runner-component"
)

type ECSInstaller struct {
	config ecsConfig
	// netInfo stores information needed to setup on-demand runners between
	// installation and on demand runner setup, so that we don't need to query
	// AWS again to re-establish all the information.
	netInfo *awsinstallutil.NetworkInformation
}

type ecsConfig struct {
	// ServerImage is the image/tag of the Waypoint server to use. Default is
	// hashicorp/waypoint:latest
	ServerImage string `hcl:"server_image,optional"`

	// Region defines which AWS region to use.
	Region string `hcl:"region,optional"`

	// Cluster is the name of the ECS Cluster to install the service into.
	// Defaults to waypoint-server
	Cluster string `hcl:"cluster,optional"`

	// ExecutionRoleName is the name of the execution task IAM Role to associate
	// with the ECS Service
	ExecutionRoleName string `hcl:"execution_role_name,optional"`

	// Subnets to place the service into. Defaults to the public Subnets in the
	// default VPC.
	Subnets []string `hcl:"subnets,optional"`

	// CPU configures the default amount of CPU for the task
	CPU string `hcl:"cpu,optional"`
	// Memory configures the default amount of memory for the task
	Memory string `hcl:"memory,optional"`

	// IAM Execution Role to assign to the on-demand runner
	TaskRoleName string `hcl:"task_role_name,optional"`

	// On-Demand Runner docker image. Defaults to hashicorp/waypoint-odr
	OdrImage string `hcl:"odr_image,optional"`

	// On-Demand Runner
	OdrCPU string `hcl:"odr_cpu,optional"`

	// On-Demand Runner
	OdrMemory string `hcl:"odr_memory,optional"`
}

// Install is a method of ECSInstaller and implements the Installer interface to
// register a waypoint-server in a ecs cluster
func (i *ECSInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, string, error) {
	ui := opts.UI
	log := opts.Log

	log.Info("Starting lifecycle")

	var (
		efsInfo                                *awsinstallutil.EfsInformation
		executionRole, cluster, serverLogGroup string
		netInfo                                *awsinstallutil.NetworkInformation
		server                                 *ecsServer
		sess                                   *session.Session

		err error
	)

	// validate we have a memory/cpu combination that ECS will accept. See
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-cpu-memory-error.html
	// for more information on valid combinations
	mem, err := strconv.Atoi(i.config.Memory)
	if err != nil {
		return nil, "", err
	}
	cpu, err := strconv.Atoi(i.config.CPU)
	if err != nil {
		return nil, "", err
	}

	if err := utils.ValidateEcsMemCPUPair(mem, cpu); err != nil {
		return nil, "", err
	}

	// we need to validate the given ODR mem/cpu at install time to verify the
	// ODR will be able to launch without adjusting the configuration
	// post-install
	odrMem, err := strconv.Atoi(i.config.OdrMemory)
	if err != nil {
		return nil, "", err
	}
	odrCpu, err := strconv.Atoi(i.config.OdrCPU)
	if err != nil {
		return nil, "", err
	}
	if err := utils.ValidateEcsMemCPUPair(odrMem, odrCpu); err != nil {
		return nil, "", err
	}

	lf := &Lifecycle{
		Init: func(ui terminal.UI) error {
			sess, err = utils.GetSession(&utils.SessionConfig{
				Region: i.config.Region,
				Logger: log,
			})
			if err != nil {
				return err
			}

			grpcPort, _ := strconv.Atoi(serverconfig.DefaultGRPCPort)
			httpPort, _ := strconv.Atoi(serverconfig.DefaultHTTPPort)
			ports := []*int64{
				aws.Int64(int64(grpcPort)),
				aws.Int64(int64(httpPort)), // TODO: Not needed for runner install
				aws.Int64(int64(2049)),     // EFS File system port
			}
			if netInfo, err = awsinstallutil.SetupNetworking(ctx, ui, sess, i.config.Subnets, ports); err != nil {
				return err
			}
			i.netInfo = netInfo

			if cluster, err = i.SetupCluster(ctx, ui, sess); err != nil {
				return err
			}

			efsTags := []*efs.Tag{
				{
					Key:   aws.String(defaultServerTagName),
					Value: aws.String(defaultServerTagValue),
				},
			}
			if efsInfo, err = awsinstallutil.SetupEFS(ctx, ui, sess, netInfo, efsTags); err != nil {
				return err
			}

			if executionRole, err = awsinstallutil.SetupExecutionRole(ctx, ui, log, sess, i.config.ExecutionRoleName); err != nil {
				return err
			}

			if serverLogGroup, err = awsinstallutil.SetupLogs(ctx, ui, log, sess, defaultServerLogGroup); err != nil {
				return err
			}

			return nil
		},

		Run: func(ui terminal.UI) error {
			server, err = i.Launch(ctx, log, ui, sess, efsInfo, netInfo, executionRole, cluster, serverLogGroup, opts.ServerRunFlags)
			return err
		},

		Cleanup: func(ui terminal.UI) error { return nil },
	}

	if err := lf.Execute(log, ui); err != nil {
		return nil, "", err
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
	return &InstallResults{
		Context:       &contextConfig,
		AdvertiseAddr: &advertiseAddr,
		HTTPAddr:      httpAddr,
	}, "", nil
}

// Launch takes the previously created resource and launches the Waypoint server
// service
func (i *ECSInstaller) Launch(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	sess *session.Session,
	efsInfo *awsinstallutil.EfsInformation,
	netInfo *awsinstallutil.NetworkInformation,
	executionRoleArn, clusterName, logGroup string, rawRunFlags []string,
) (*ecsServer, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Creating Network resources...")
	defer func() { s.Abort() }()

	grpcPort, _ := strconv.Atoi(serverconfig.DefaultGRPCPort)
	httpPort, _ := strconv.Atoi(serverconfig.DefaultHTTPPort)
	nlb, err := createNLB(
		ctx, s, log, sess,
		netInfo.VpcID,
		aws.Int64(int64(grpcPort)),
		aws.Int64(int64(httpPort)),
		netInfo.Subnets,
	)
	if err != nil {
		return nil, err
	}
	s.Update("Network load balancer created")
	s.Done()
	s = sg.Add("")

	defaultStreamPrefix := fmt.Sprintf("waypoint-server-%d", time.Now().Nanosecond())
	logOptions := buildLoggingOptions(
		nil,
		i.config.Region,
		logGroup,
		defaultStreamPrefix,
	)

	cmd := []*string{
		aws.String("server"),
		aws.String("run"),
		aws.String("-accept-tos"),
		aws.String("-vv"),
		aws.String("-db=/waypoint-data/data.db"),
		aws.String(fmt.Sprintf("-listen-grpc=0.0.0.0:%d", grpcPort)),
		aws.String(fmt.Sprintf("-listen-http=0.0.0.0:%d", httpPort)),
	}
	for _, f := range rawRunFlags {
		cmd = append(cmd, aws.String(f))
	}

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Command:   cmd,
		Name:      aws.String(serverName),
		Image:     aws.String(i.config.ServerImage),
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
	log.Debug("registering task definition")

	s.Update("Registering Task definition: %s", defaultTaskFamily)

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    []*ecs.ContainerDefinition{&def},
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Cpu:                     aws.String(i.config.CPU),
		Memory:                  aws.String(i.config.Memory),
		Family:                  aws.String(defaultTaskFamily),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(defaultTaskRuntime)},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
			},
		},
		Volumes: []*ecs.Volume{
			{
				Name: aws.String("waypointdata"),
				EfsVolumeConfiguration: &ecs.EFSVolumeConfiguration{
					TransitEncryption: aws.String(ecs.EFSTransitEncryptionEnabled),
					FileSystemId:      efsInfo.FileSystemID,
					AuthorizationConfig: &ecs.EFSAuthorizationConfig{
						AccessPointId: efsInfo.AccessPointID,
					},
				},
			},
		},
	}

	ecsSvc := ecs.New(sess)
	taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc, log)
	if err != nil {
		return nil, err
	}

	// registerTaskDefinition() above ensures taskDef here is non-nil, if the
	// error returned is nil
	taskDefArn := *taskDef.TaskDefinitionArn

	// Create the service
	s.Update("Creating server service...")
	log.Debug("creating service", "arn", *taskDef.TaskDefinitionArn)

	createServiceInput := &ecs.CreateServiceInput{
		Cluster:                       &clusterName,
		DesiredCount:                  aws.Int64(1),
		LaunchType:                    aws.String(defaultTaskRuntime),
		ServiceName:                   aws.String(serverName),
		TaskDefinition:                aws.String(taskDefArn),
		EnableECSManagedTags:          aws.Bool(true),
		HealthCheckGracePeriodSeconds: aws.Int64(int64(60)),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        netInfo.Subnets,
				SecurityGroups: []*string{netInfo.SgID},
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
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

	service, err := awsinstallutil.CreateService(createServiceInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	s.Update("Created ECS Service (%s, cluster-name: %s)", serviceName, clusterName)
	s.Done()
	s = sg.Add("")
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
	s.Done()
	s = sg.Add("")
	s.Update("Service launched!")
	s.Done()

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

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting ecs cluster...")
	defer s.Abort()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	// inspect current service - looking for image used in Task
	// Get Task definition
	var clusterArn string
	cluster := i.config.Cluster
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
	s.Done()
	s = sg.Add("Inspecting load balancer target groups")
	// list the services to find the task descriptions
	services, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(i.config.Cluster),
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

	// Set or update the deregistration delay in the target groups. Waypoint
	// server installations pre-0.9.1 used the default 300 second delay, which
	// we need to lower. Note: "serverSvc.LoadBalancers" below is a slice of
	// load balancer + target group pairs. We expect two target groups attached
	// to the same load balancer.
	// See
	// https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeServices.html
	for _, lb := range serverSvc.LoadBalancers {
		if lb.TargetGroupArn != nil {
			s.Update("Updating deregistration delay and termination for Target Groups")
			if err := modifyTargetGroups(elbv2.New(sess), *lb.TargetGroupArn); err != nil {
				s.Update("Error updating Target Group %s: %w", *lb, err)
			}
		}
	}

	s.Done()
	s = sg.Add("Updating task definition")
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

	upgradeImg := installutil.DefaultServerImage
	if i.config.ServerImage != "" {
		upgradeImg = i.config.ServerImage
	}

	s.Done()
	s = sg.Add("Updating task definition")
	defer func() { s.Abort() }()

	if containerDef != nil && *containerDef.Image == installutil.DefaultServerImage && upgradeImg == installutil.DefaultServerImage {
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
		taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc, log)
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
	}
	// after updating the service, we need to stop the existing task so that the
	// new task that is created will be able to open the db and respond to
	// health checks
	tasks, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:     &clusterArn,
		ServiceName: serverSvc.ServiceName,
	})
	if err != nil {
		s.Update("failed to list tasks for service %s, error: %w", *serverSvc.ServiceName, err)
	}

	if len(tasks.TaskArns) > 1 {
		s.Update("Warning: multiple running server tasks detected; there can" +
			"be only 1 active server task at a time")
	}

	// STOP any running tasks. Ideally this completes before the service starts
	// the new task
	for _, taskArn := range tasks.TaskArns {
		_, err := ecsSvc.StopTask(&ecs.StopTaskInput{
			Cluster: &clusterArn,
			Reason:  aws.String("Waypoint server upgrade"),
			Task:    taskArn,
		})
		if err != nil {
			s.Update("failed to stop task %s, error: %w", *taskArn, err)
		}
	}

	s.Update("Waiting until service is stable")
	// WaitUntil waits for ~10 minutes
	err = ecsSvc.WaitUntilServicesStable(&ecs.DescribeServicesInput{
		Cluster:  &clusterArn,
		Services: []*string{serverSvc.ServiceName},
	})
	if err != nil {
		return nil, err
	}

	s.Done()
	s = sg.Add("Updating context...")

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

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Uninstalling Server resources...")
	defer func() { s.Abort() }()

	// Get list of resources created with either the waypoint-server, or
	// waypoint-runner tag
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: i.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	rgSvc := resourcegroups.New(sess)
	var resources []*resourcegroups.ResourceIdentifier
	query := fmt.Sprintf(serverResourceQuery, defaultServerTagName)
	searchInput := resourcegroups.SearchResourcesInput{
		MaxResults: aws.Int64(20),
		ResourceQuery: &resourcegroups.ResourceQuery{
			Type:  aws.String(resourcegroups.QueryTypeTagFilters10),
			Query: aws.String(query),
		},
	}

	// The Resource Group Search results can sometimes be limited to a few
	// results at a time and may not include all resources tagged. Use the
	// pagination function to retrieve the complete list.
	err = rgSvc.SearchResourcesPages(&searchInput,
		func(page *resourcegroups.SearchResourcesOutput, _ bool) bool {
			resources = append(resources, page.ResourceIdentifiers...)
			return page.NextToken != nil
		})

	if err != nil {
		return fmt.Errorf("error retrieving tag search results: %w", err)
	}

	if len(resources) == 0 {
		return fmt.Errorf("no server resources found with tag (%s)", defaultServerTagName)
	}

	// Start destroying things. Some cannot be destroyed before others. The
	// general order to destroy things:
	// - ECS Service
	// - ECS Cluster
	// - Cloudwatch Log Group
	// - ELB Target Groups
	// - ELB Network Load Balancer
	// - EFS File System

	s.Update("Deleting ECS resources...")
	if err := awsinstallutil.DeleteEcsResources(ctx, sess, resources); err != nil {
		return err
	}
	s.Done()

	s.Update("Deleting Cloud Watch Log Group resources...")
	if err := awsinstallutil.DeleteCWLResources(ctx, sess, defaultServerLogGroup); err != nil {
		return err
	}
	s.Done()

	s.Update("Deleting EFS resources...")
	if err := deleteEFSResources(ctx, sess, resources); err != nil {
		return err
	}

	s.Update("Deleting Network resources...")
	if err := deleteNLBResources(ctx, sess, resources); err != nil {
		return err
	}

	s.Update("Server resources deleted")
	s.Done()
	return nil
}

func deleteEFSResources(
	ctx context.Context,
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
			break
		}

		time.Sleep(5 * time.Second)
		continue
	}

	_, err = efsSvc.DeleteFileSystem(&efs.DeleteFileSystemInput{
		FileSystemId: &id,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "FileSystemNotFound":
				// the file system has already been destroyed
				return nil
			}
		}
		return err
	}
	return nil
}

func deleteNLBResources(
	ctx context.Context,
	sess *session.Session,
	resources []*resourcegroups.ResourceIdentifier,
) error {
	elbSvc := elbv2.New(sess)
	for _, r := range resources {
		if *r.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer" {
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

			_, err = elbSvc.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: r.ResourceArn,
			})
			if err != nil {
				return err
			}
		}
	}

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
				Values: []*string{aws.String(defaultServerTagName)},
			},
		},
	})
	if err != nil {
		return err
	}
	if len(results.SecurityGroups) > 0 {
		for _, g := range results.SecurityGroups {
			for i := 0; i < 60; i++ {
				_, err := ec2Svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
					GroupId: g.GroupId,
				})
				// if we encounter an unrecoverable error, exit now.
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "DependencyViolation":
						time.Sleep(3 * time.Second)
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

// InstallRunner implements Installer.
func (i *ECSInstaller) InstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.ECSRunnerInstaller{Config: runnerinstall.EcsConfig{
		Region:            i.config.Region,
		ExecutionRoleName: i.config.ExecutionRoleName,
		TaskRoleName:      i.config.TaskRoleName,
		CPU:               i.config.CPU,
		Memory:            i.config.Memory,
		RunnerImage:       i.config.ServerImage,
		Cluster:           i.config.Cluster,
		Subnets:           i.config.Subnets,
	}}
	err := runnerInstaller.Install(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

var (
	serverResourceQuery = "{\"ResourceTypeFilters\":[\"AWS::AllSupported\"],\"TagFilters\":[{\"Key\":\"%s\",\"Values\":[]}]}"
	runnerResourceQuery = "{\"ResourceTypeFilters\":[\"AWS::AllSupported\"],\"TagFilters\":[{\"Key\":\"%s\",\"Values\":[\"%s\"]}]}"
)

func (i *ECSInstaller) UninstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.ECSRunnerInstaller{Config: runnerinstall.EcsConfig{
		Region:  i.config.Region,
		Cluster: i.config.Cluster,
	}}

	err := runnerInstaller.Uninstall(ctx, opts)
	if err != nil {
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
		Region: i.config.Region,
		Logger: log,
	})
	if err != nil {
		return false, err
	}
	serviceNames := []string{
		defaultRunnerTagName,
		installutil.DefaultRunnerName("static"),
	}
	ecsSvc := ecs.New(sess)
	services, err := awsinstallutil.FindServices(serviceNames, ecsSvc, i.config.Cluster)
	if err != nil {
		opts.UI.Output("Could not get list of ECS services: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return false, err
	}

	return len(services) > 0, nil
}

func (i *ECSInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.Cluster,
		Usage:   "Configures the Cluster to install into.",
		Default: "waypoint-server",
	})
	set.StringVar(&flag.StringVar{
		Name:   "ecs-region",
		Target: &i.config.Region,
		Usage:  "Configures which AWS region to install into.",
	})
	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "ecs-subnets",
		Target: &i.config.Subnets,
		Usage:  "Subnets to install server into.",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-execution-role-name",
		Target:  &i.config.ExecutionRoleName,
		Usage:   "Configures the IAM Execution role name to use.",
		Default: "waypoint-server-execution-role",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-server-image",
		Target:  &i.config.ServerImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cpu",
		Target:  &i.config.CPU,
		Usage:   "Configures the requested CPU amount for the Waypoint server task in ECS.",
		Default: "512",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-mem",
		Target:  &i.config.Memory,
		Usage:   "Configures the requested memory amount for the Waypoint server task in ECS.",
		Default: "1024",
	})
	set.StringVar(&flag.StringVar{
		Name:   "ecs-task-role-name",
		Target: &i.config.TaskRoleName,
		Usage: "IAM Execution Role to assign to the on-demand runner. If this is blank, " +
			"an IAM role will be created automatically with the default permissions.",
		Default: "waypoint-runner",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-odr-image",
		Target: &i.config.OdrImage,
		Usage: "Docker image for the Waypoint On-Demand Runners. This will " +
			"default to the server image with the name (not label) suffixed with '-odr'.",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-odr-mem",
		Target:  &i.config.OdrMemory,
		Usage:   "Configures the requested memory amount for the Waypoint On-Demand runner in ECS.",
		Default: "2048",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-odr-cpu",
		Target:  &i.config.OdrCPU,
		Usage:   "Configures the requested CPU amount for the Waypoint On-Demand runner in ECS.",
		Default: "512",
	})
}

func (i *ECSInstaller) UpgradeFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.Cluster,
		Usage:   "Configures the Cluster to upgrade.",
		Default: "waypoint-server",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-server-image",
		Target:  &i.config.ServerImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})
	set.StringVar(&flag.StringVar{
		Name:   "ecs-region",
		Target: &i.config.Region,
		Usage:  "Configures which AWS region to install into.",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cpu",
		Target:  &i.config.CPU,
		Usage:   "Configures the requested CPU amount for the Waypoint server task in ECS.",
		Default: "512",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-mem",
		Target:  &i.config.Memory,
		Usage:   "Configures the requested memory amount for the Waypoint server task in ECS.",
		Default: "1024",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-odr-image",
		Target: &i.config.OdrImage,
		Usage: "Docker image for the Waypoint On-Demand Runners. This will " +
			"default to the server image with the name (not label) suffixed with '-odr'.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "ecs-task-role-name",
		Target: &i.config.TaskRoleName,
		Usage: "IAM Execution Role to assign to the on-demand runner. If this is blank, " +
			"an IAM role will be created automatically with the default permissions.",
		Default: "waypoint-runner",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-execution-role-name",
		Target:  &i.config.ExecutionRoleName,
		Usage:   "Configures the IAM Execution role name to use.",
		Default: "waypoint-server-execution-role",
	})

	set.StringVar(&flag.StringVar{
		Name:    "ecs-odr-mem",
		Target:  &i.config.OdrMemory,
		Usage:   "Configures the requested memory amount for the Waypoint On-Demand runner in ECS.",
		Default: "2048",
	})
	set.StringVar(&flag.StringVar{
		Name:    "ecs-odr-cpu",
		Target:  &i.config.OdrCPU,
		Usage:   "Configures the requested CPU amount for the Waypoint On-Demand runner in ECS.",
		Default: "512",
	})
}

func (i *ECSInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "ecs-cluster",
		Target:  &i.config.Cluster,
		Usage:   "Configures the Cluster to uninstall.",
		Default: "waypoint-server",
	})
	set.StringVar(&flag.StringVar{
		Name:   "ecs-region",
		Target: &i.config.Region,
		Usage:  "Configures which AWS region to uninstall from.",
	})
}

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

func (i *ECSInstaller) SetupCluster(
	ctx context.Context,
	ui terminal.UI,
	sess *session.Session,
) (string, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Inspecting existing ECS clusters...")
	defer func() { s.Abort() }()

	cluster := i.config.Cluster

	ecsSvc := ecs.New(sess)
	// re-use an existing cluster if we have one
	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(cluster)},
	})
	if err != nil {
		return "", err
	}

	for _, c := range desc.Clusters {
		if *c.ClusterName == cluster && strings.ToLower(*c.Status) == "active" {
			s.Update("Found existing ECS cluster: %s", cluster)
			s.Done()
			return cluster, nil
		}
	}

	s.Update("Creating new ECS cluster: %s", cluster)

	_, err = ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(cluster),
		// we need to tag with both the server and runner names, so we can properly
		// cleanup
		Tags: []*ecs.Tag{
			{
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
			},
			{
				Key:   aws.String(defaultRunnerTagName),
				Value: aws.String(defaultRunnerTagValue),
			},
		},
	})

	if err != nil {
		return "", err
	}

	s.Update("Created new ECS cluster: %s", cluster)
	s.Done()

	return cluster, nil
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

// creates a network load balancer for grpc and http
func createNLB(
	ctx context.Context,
	s terminal.Step,
	log hclog.Logger,
	sess *session.Session,
	vpcId *string,
	grpcPort *int64,
	httpPort *int64,
	subnets []*string,
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
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
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
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	httpTgArn := htgGPRC.TargetGroups[0].TargetGroupArn
	grpcTgArn := ctgGPRC.TargetGroups[0].TargetGroupArn

	s.Update("Updating deregistration delay and termination for Target Groups")
	for _, tgArn := range []*string{httpTgArn, grpcTgArn} {
		if err := modifyTargetGroups(elbsrv, *tgArn); err != nil {
			s.Update("Error updating Target Group %s: %w", *tgArn, err)
		}
	}

	// Create the load balancer OR modify the existing one to have this new
	// target group. Note: Network Load Balancers do not use the weight
	// attribute of TargetGroupTuple.
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
					Key:   aws.String(defaultServerTagName),
					Value: aws.String(defaultServerTagValue),
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
				Key:   aws.String(defaultServerTagName),
				Value: aws.String(defaultServerTagValue),
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

// OnDemandRunnerConfig implements OnDemandRunnerConfigProvider
func (i *ECSInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	// Generate some configuration. Some of the OnDemand configurations have
	// defaults so we should be fine to directly use them
	cfgMap := map[string]interface{}{
		"log_group":           defaultRunnerLogGroup,
		"execution_role_name": i.config.ExecutionRoleName,
		"task_role_name":      i.config.TaskRoleName,
		"cluster":             i.config.Cluster,
		"region":              i.config.Region,
		"odr_cpu":             i.config.OdrCPU,
		"odr_memory":          i.config.OdrMemory,
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
		Name:         "aws-ecs",
		OciUrl:       i.config.OdrImage,
		PluginType:   "aws-ecs",
		Default:      true,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
	}
}

// modifyTargetGroups modifies the target group  to support a shorter
// deregistration period. This helps with upgrades and fail-overs, to route
// traffic more quickly to any new server tasks.
func modifyTargetGroups(elbsrv *elbv2.ELBV2, targetGroupArn string) error {
	_, err := elbsrv.ModifyTargetGroupAttributes(&elbv2.ModifyTargetGroupAttributesInput{
		TargetGroupArn: &targetGroupArn,
		Attributes: []*elbv2.TargetGroupAttribute{
			{
				Key:   aws.String("deregistration_delay.timeout_seconds"),
				Value: aws.String("5"),
			},
			{
				Key:   aws.String("deregistration_delay.connection_termination.enabled"),
				Value: aws.String("true"),
			},
		},
	})
	return err
}
