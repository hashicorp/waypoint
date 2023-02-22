package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const (
	executionRolePolicyArn        = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
	awsCreateRetries              = 30
	awsCreateRetryIntervalSeconds = 2

	// Five minutes of retrying for destroy - destroys are slower than creates in aws.
	awsDestroyRetries              = 60
	awsDestroyRetryIntervalSeconds = 5
	defaultServicePort             = 3000
)

type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	c, ok := config.(*Config)
	if !ok {
		// this should never happen
		return status.Errorf(codes.FailedPrecondition, "invalid configuration, expected *ecs.Config, got %T", config)
	}

	if c.ALB != nil {
		alb := c.ALB
		err := utils.Error(validation.ValidateStruct(alb,
			// ZoneId and FQDN are both used to create a route53 record that points to an ALB.
			// While we could still create this, if a user is managing their own ALB out of band,
			// they probably also want to manage the route53 record themselves
			// too.
			validation.Field(&alb.ZoneId,
				validation.Empty.When(alb.LoadBalancerArn != ""),
				validation.Required.When(alb.FQDN != ""),
			),
			validation.Field(&alb.FQDN,
				validation.Empty.When(alb.LoadBalancerArn != ""),
				validation.Required.When(alb.ZoneId != "").Error("fqdn only valid with zone_id"),
			),

			validation.Field(&alb.InternalScheme,
				validation.Nil.When(alb.LoadBalancerArn != "").Error("internal cannot be used with load_balancer_arn"),
			),
			validation.Field(&alb.LoadBalancerArn,
				validation.Empty.When(
					alb.ZoneId != "" ||
						alb.FQDN != "" ||
						len(alb.SecurityGroupIDs) >= 1).Error("load_balancer_arn can not be used with other options"),
			),
			validation.Field(&alb.SecurityGroupIDs),
		))
		if err != nil {
			return err
		}
	}

	err := utils.Error(validation.ValidateStruct(c,
		validation.Field(&c.Memory, validation.Required, validation.Min(4)),
		validation.Field(&c.MemoryReservation, validation.Min(4), validation.Max(c.Memory)),
	))
	if err != nil {
		return err
	}

	if c.Architecture != "" {
		c.Architecture = strings.ToUpper(c.Architecture)
		cpuArchitectures := make([]interface{}, len(ecs.CPUArchitecture_Values()))
		for i, ca := range ecs.CPUArchitecture_Values() {
			cpuArchitectures[i] = ca
		}
		err := utils.Error(validation.ValidateStruct(c,
			validation.Field(&c.Architecture,
				validation.In(cpuArchitectures...).Error("unsupported CPU architecture"),
			),
		))
		if err != nil {
			return err
		}
	}

	for _, cc := range c.ContainersConfig {
		err := utils.Error(validation.ValidateStruct(cc,
			validation.Field(&cc.Memory, validation.Required, validation.Min(4)),
			validation.Field(&cc.MemoryReservation, validation.Min(4), validation.Max(cc.Memory)),
		))
		if err != nil {
			return err
		}
	}

	if c.HealthCheck != nil {
		if c.HealthCheck.Timeout >= c.HealthCheck.Interval {
			return status.Errorf(codes.FailedPrecondition, "the health check "+
				"timeout must be less than the interval")
		}
	}

	return nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return p.Destroy
}

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return p.ValidateAuth
}

// DestroyWorkspaceFunc implements component.WorkspaceDestroyer
func (p *Platform) DestroyWorkspaceFunc() interface{} {
	return p.DestroyWorkspace
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return p.Auth
}

func (p *Platform) Auth() error {
	return nil
}

func (p *Platform) ValidateAuth() error {
	return nil
}

// StatusFunc implements component.Status
func (p *Platform) StatusFunc() interface{} {
	return p.Status
}

// DefaultReleaserFunc implements component.PlatformReleaser
func (p *Platform) DefaultReleaserFunc() interface{} {
	return func() *Releaser { return &Releaser{p: p} }
}

func (p *Platform) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp, dtr *component.DestroyedResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getSession),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithDestroyedResourcesResp(dtr),
		resource.WithResource(resource.NewResource(
			resource.WithName("cluster"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_Cluster{}),
			resource.WithCreate(p.resourceClusterCreate),
			// TODO: implement destroy when we have better support for globally-scoped resources
			resource.WithStatus(p.resourceClusterStatus),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_OTHER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("execution role"),
			resource.WithType("IAM role"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_ExecutionRole{}),
			resource.WithCreate(p.resourceExecutionRoleCreate),
			// TODO: implement destroy when we have better support for app-scoped resources
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_POLICY),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("task role"),
			resource.WithType("IAM role"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_TaskRole{}),
			resource.WithCreate(p.resourceTaskRoleCreate),
			// TODO: implement destroy when we have better support for app-scoped resources
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_POLICY),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("internal security groups"),
			resource.WithPlatform(platformName),
			resource.WithType("security groups"),
			resource.WithState(&Resource_InternalSecurityGroups{}),
			resource.WithCreate(p.resourceInternalSecurityGroupsCreate),
			// NOTE(izaak/nic): these are destroyed by resourceSecurityGroupsDestroy, registered under external security groups
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_POLICY),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("external security groups"),
			resource.WithPlatform(platformName),
			resource.WithType("security groups"),
			resource.WithState(&Resource_ExternalSecurityGroups{}),
			resource.WithCreate(p.resourceExternalSecurityGroupsCreate),
			resource.WithDestroy(p.resourceSecurityGroupsDestroy),
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_POLICY),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("log group"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_LogGroup{}),
			resource.WithCreate(p.resourceLogGroupCreate),
			// TODO: implement destroy when we have better support for waypoint global resources
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_OTHER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("service subnets"),
			resource.WithType("subnets"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_ServiceSubnets{}),
			resource.WithCreate(p.resourceServiceSubnetsDiscover),
			// We never create subnets, and therefore should never destroy them
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("alb subnets"),
			resource.WithType("subnets"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_AlbSubnets{}),
			resource.WithCreate(p.resourceAlbSubnetsDiscover),
			// We never create subnets, and therefore should never destroy them
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("target group"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_TargetGroup{}),
			resource.WithCreate(p.resourceTargetGroupCreate),
			resource.WithDestroy(p.resourceTargetGroupDestroy),
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_OTHER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("application load balancer"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_Alb{}),
			resource.WithCreate(p.resourceAlbCreate),
			resource.WithDestroy(p.resourceAlbDestroy),
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("alb listener"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_Alb_Listener{}),
			resource.WithCreate(p.resourceAlbListenerCreate),
			resource.WithDestroy(p.resourceAlbListenerDestroy),
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("route53 record"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_Route53Record{}),
			resource.WithCreate(p.resourceRoute53RecordCreate),
			// TODO: implement destroy when we have better support for app-scoped resources
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_ROUTER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("task definition"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_TaskDefinition{}),
			resource.WithCreate(p.resourceTaskDefinitionCreate),
			resource.WithDestroy(p.resourceTaskDefinitionDestroy),
			// TODO: implement status when we have a plan to not hit rate limits
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
		resource.WithResource(resource.NewResource(
			resource.WithName("service"),
			resource.WithPlatform(platformName),
			resource.WithState(&Resource_Service{}),
			resource.WithCreate(p.resourceServiceCreate),
			resource.WithDestroy(p.resourceServiceDestroy),
			resource.WithStatus(p.resourceServiceStatus),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
	)
}

// Below are custom types for various parameters used throughout resource creation.
// They exist as custom types so that argmapper can recognize and wire them
// into the appropriate resource lifecycle functions.

// DeploymentId is a unique ID to be consistently used throughout our deployment
type DeploymentId string

// WorkspaceDestroy marks that we're destroying the entire workspace (not just the deployment), and
// we should destroy global resources too
type WorkspaceDestroy bool

// ExternalIngressPort is the port that the ALB will listen for traffic on
type ExternalIngressPort int64

func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
	dcr *component.DeclaredResourcesResp,
) (*Deployment, error) {
	var result Deployment

	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing deployment...")
	defer s.Abort()

	// Generate a common deployment ID to use in the resources we create.
	// TODO: should include the sequence ID
	ulid, err := component.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate a ULID: %s", err)
	}
	deploymentId := DeploymentId(fmt.Sprintf("%s-%s", src.App, ulid))

	// Set default port - it's used for multiple resources
	if p.config.ServicePort != 0 {
		log.Debug("Using configured service port", "port number", p.config.ServicePort)
	} else {
		log.Debug("Using the default service port", "port number", defaultServicePort)
		p.config.ServicePort = int64(defaultServicePort)
	}

	// Set ALB ingress port - used for multiple resources
	var externalIngressPort ExternalIngressPort
	if p.config.ALB != nil && p.config.ALB.IngressPort != 0 {
		log.Debug("Using configured ingress port", "port number", p.config.ServicePort)
		externalIngressPort = ExternalIngressPort(p.config.ALB.IngressPort)
	} else if p.config.ALB != nil && p.config.ALB.CertificateId != "" {
		log.Debug("ALB config defined and cert configured, using ingress port 443")
		externalIngressPort = ExternalIngressPort(443)
	} else {
		log.Debug("Defaulting external ingress port to 80")
		externalIngressPort = ExternalIngressPort(80)
	}

	// Create our resource manager and create
	rm := p.resourceManager(log, dcr, nil)
	if err := rm.CreateAll(
		ctx, log, sg, ui, deploymentId, externalIngressPort,
		src, img, deployConfig, &result,
	); err != nil {
		return nil, status.Errorf(status.Convert(err).Code(), "failed to create deployment resources: %s", err)
	}

	// Store our resource state
	result.ResourceState = rm.State()

	// Get other state required for older versions to destroy this deployment
	srState := rm.Resource("service").State().(*Resource_Service)
	result.ServiceArn = srState.Arn

	tgState := rm.Resource("target group").State().(*Resource_TargetGroup)
	result.TargetGroupArn = tgState.Arn

	albState := rm.Resource("application load balancer").State().(*Resource_Alb)
	result.LoadBalancerArn = albState.Arn

	listenerState := rm.Resource("alb listener").State().(*Resource_Alb_Listener)
	result.ListenerArn = listenerState.Arn

	cState := rm.Resource("cluster").State().(*Resource_Cluster)
	result.Cluster = cState.Name

	tdState := rm.Resource("task definition").State().(*Resource_TaskDefinition)
	result.TaskArn = tdState.Arn

	s.Update("Deployment resources created")
	s.Done()
	return &result, nil
}

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Gathering health report for ecs deployment...")
	defer s.Abort()

	rm := p.resourceManager(log, nil, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		if err := p.loadResourceManagerState(ctx, rm, deployment, log, sg); err != nil {
			return nil, status.Errorf(codes.Internal, "failed recovering old state into resource manager: %s", err)
		}
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return nil, status.Errorf(codes.Internal, "failed loading state into resource manager: %s", err)
		}
	}

	result, err := rm.StatusReport(ctx, log, sg)
	if err != nil {
		return nil, status.Errorf(status.Convert(err).Code(), "resource manager failed to generate a status report: %s", err)
	}

	s.Update("Finished building report for ecs deployment")
	s.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	st.Update("Determining overall health for ecs deployment...")
	if result.Health == sdk.StatusReport_READY {
		st.Step(terminal.StatusOK, result.HealthMessage)
	} else {
		if result.Health == sdk.StatusReport_PARTIAL {
			st.Step(terminal.StatusWarn, result.HealthMessage)
		} else {
			st.Step(terminal.StatusError, result.HealthMessage)
		}

		// Extra advisory wording to let user know that the deployment could be still starting up
		// if the report was generated immediately after it was deployed or released.
		st.Step(terminal.StatusWarn, mixedHealthWarn)
	}

	return result, nil
}

func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
	dcr *component.DeclaredResourcesResp,
	dtr *component.DestroyedResourcesResp,
) error {

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Destroying ecs deployment...")
	defer s.Abort()

	rm := p.resourceManager(log, dcr, dtr)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		if err := p.loadResourceManagerState(ctx, rm, deployment, log, sg); err != nil {
			return status.Errorf(codes.Internal, "failed recovering old state into resource manager: %s", err)
		}
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return status.Errorf(codes.Internal, "failed loading state into resource manager: %s", err)
		}
	}

	// Destroy app-scoped resources
	// Note: This calls the destroy func for all resources, but resource destroy funcs will use the WorkspaceDestroy
	// param to only actually destroy if they are deployment scoped.
	err := rm.DestroyAll(ctx, log, sg, ui, WorkspaceDestroy(false))
	if err != nil {
		return status.Errorf(status.Convert(err).Code(), "failed to destroy all resources for deployment: %s", err)
	}

	s.Update("Finished destroying ECS deployment")
	s.Done()
	return nil
}

func (p *Platform) DestroyWorkspace(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	app *component.Source,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Destroying ecs workspace...")
	defer s.Abort()

	rm := p.resourceManager(log, nil, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		if err := p.loadResourceManagerState(ctx, rm, deployment, log, sg); err != nil {
			return status.Errorf(codes.Internal, "failed recovering old state into resource manager: %s", err)
		}
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return status.Errorf(codes.Internal, "failed loading state into resource manager: %s", err)
		}
	}

	// Destroy
	// Note: This calls the destroy func for all resources, but resource destroy funcs will use the WorkspaceDestroy
	// param to only actually destroy if they are workspace scoped.
	err := rm.DestroyAll(ctx, log, sg, ui, WorkspaceDestroy(true))
	if err != nil {
		return status.Errorf(status.Convert(err).Code(), "failed to destroy all resources for workspace: %s", err)
	}

	s.Update("Finished destroying ECS deployment")
	s.Done()
	return nil
}

func (p *Platform) resourceClusterCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	state *Resource_Cluster,
) error {
	s := sg.Add("Initiating cluster creation...")
	defer s.Abort()

	cluster := p.config.Cluster
	if cluster == "" {
		cluster = "waypoint"
	}
	state.Name = cluster

	s.Update("Attempting to find existing cluster named %q", cluster)

	ecsSvc := ecs.New(sess)
	desc, err := ecsSvc.DescribeClustersWithContext(ctx, &ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(cluster)},
	})
	if err != nil {
		return err
	}

	for _, c := range desc.Clusters {
		if *c.ClusterName == cluster {
			if *c.Status == "PROVISIONING" {
				s.Update("Existing ecs cluster %q is still provisioning - try again later.", cluster)
			} else if *c.Status == "ACTIVE" {
				s.Update("Using existing ECS cluster %s", cluster)
				if c.ClusterArn != nil {
					state.Arn = *c.ClusterArn
				}
				s.Done()
				return nil
			} else {
				// Warn if we encounter waypoint clusters in other odd states (i.e. DEPROVISIONING, FAILED, etc.)
				// I think it's ok to try to create a new cluster if one exists in a non-active non-provisioning state
				log.Warn("Ignoring cluster named %q in state %q", cluster, *c.Status)
			}
		}
	}

	if p.config.EC2Cluster {
		return status.Errorf(codes.FailedPrecondition, "EC2 clusters can not be automatically created")
	}

	s.Update("No existing cluster found - creating new ECS cluster: %s", cluster)

	c, err := ecsSvc.CreateClusterWithContext(ctx, &ecs.CreateClusterInput{
		ClusterName: aws.String(cluster),
	})

	if err != nil {
		return err
	}

	if c.Cluster != nil && c.Cluster.ClusterArn != nil {
		state.Arn = *c.Cluster.ClusterArn
	}

	s.Update("Created ECS cluster: %s", cluster)
	s.Done()
	return nil
}

func (p *Platform) resourceClusterStatus(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	state *Resource_Cluster,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Checking status of the ecs cluster %q...", state.Name)
	defer s.Abort()

	ecsSvc := ecs.New(sess)
	desc, err := ecsSvc.DescribeClustersWithContext(ctx, &ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(state.Name)},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe cluster named %q (ARN: %q): %s", state.Name, state.Arn, err)
	}

	clusterResource := sdk.StatusReport_Resource{
		Name: state.Name,
	}

	sr.Resources = append(sr.Resources, &clusterResource)

	for _, c := range desc.Clusters {
		if *c.ClusterName == state.Name {
			s.Update("Found existing ECS cluster: %s", state.Name)
			clusterResource.Id = *c.ClusterArn
			switch *c.Status {
			case "ACTIVE":
				clusterResource.Health = sdk.StatusReport_READY
			case "PROVISIONING":
				clusterResource.Health = sdk.StatusReport_ALIVE
			case "DEPROVISIONING", "FAILED", "INACTIVE":
				clusterResource.Health = sdk.StatusReport_DOWN
			default:
				clusterResource.Health = sdk.StatusReport_UNKNOWN
			}
			clusterResource.HealthMessage = *c.Status

			stateJson, err := json.Marshal(c)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to marshal ecs cluster state json: %s", err)
			}
			clusterResource.StateJson = string(stateJson)

			s.Done()
			return nil
		}
	}

	// Failed to find ECS cluster
	clusterResource.Health = sdk.StatusReport_MISSING
	clusterResource.HealthMessage = fmt.Sprintf("No cluster named %q found (expected arn %q)", state.Name, state.Arn)

	s.Update("Done checking ecs cluster status")
	s.Done()
	return nil
}

func (p *Platform) resourceServiceCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	src *component.Source,
	deploymentId DeploymentId,
	state *Resource_Service,

	// Outputs of other resource creation processes
	taskDefinition *Resource_TaskDefinition,
	cluster *Resource_Cluster,
	targetGroup *Resource_TargetGroup,
	subnets *Resource_ServiceSubnets,
	securityGroups *Resource_InternalSecurityGroups,

	_ *Resource_Alb_Listener, // Necessary dependency. service creation will fail unless this exists and the target group has been added.
) error {
	s := sg.Add("Initiating ecs service creation")
	defer s.Abort()

	// Use the common deployment ID as our service name
	serviceName := string(deploymentId)

	// We have to clamp at a length of 32 because the Name field
	// requires that the name is 32 characters or less.
	if len(serviceName) > 32 {
		serviceName = serviceName[:32]
		log.Debug("using a shortened value for service name due to AWS's length limits", "serviceName", serviceName)
	}

	count := int64(p.config.Count)
	if count == 0 {
		count = 1
	}

	securityGroupIds := make([]*string, len(securityGroups.SecurityGroups))
	for i, securityGroup := range securityGroups.SecurityGroups {
		securityGroupIds[i] = &securityGroup.Id
	}

	subnetIds := make([]*string, len(subnets.Subnets.Subnets))
	for i, subnet := range subnets.Subnets.Subnets {
		subnetIds[i] = &subnet.Id
	}

	netCfg := &ecs.AwsVpcConfiguration{
		Subnets:        subnetIds,
		SecurityGroups: securityGroupIds,
	}

	if !p.config.EC2Cluster {
		if p.config.AssignPublicIp == nil {
			log.Debug("AssignPublicIp config value not defined - defaulting to true.")
			p.config.AssignPublicIp = aws.Bool(true)
		}

		if *p.config.AssignPublicIp {
			netCfg.AssignPublicIp = aws.String("ENABLED")
		} else {
			netCfg.AssignPublicIp = aws.String("DISABLED")
		}
	}

	state.Cluster = cluster.Name

	createServiceInput := &ecs.CreateServiceInput{
		Cluster:        &cluster.Name,
		DesiredCount:   aws.Int64(count),
		LaunchType:     &taskDefinition.Runtime,
		ServiceName:    aws.String(serviceName),
		TaskDefinition: aws.String(taskDefinition.Arn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: netCfg,
		},
	}

	if targetGroup.Arn != "" {
		log.Debug("Creating ECS service with a load balancer")
		createServiceInput.SetLoadBalancers([]*ecs.LoadBalancer{{
			ContainerName:  aws.String(src.App),
			ContainerPort:  aws.Int64(targetGroup.Port),
			TargetGroupArn: &targetGroup.Arn,
		}})
	} else {
		log.Debug("No target group specified - skipping load balancer config for ECS service")
	}

	s.Update("Creating ECS Service %s", serviceName)

	ecsSvc := ecs.New(sess)
	// AWS is eventually consistent so even though we probably created the resources that
	// are referenced by the task definition, it can error out if we try to reference those resources
	// too quickly. So we're forced to guard actions which reference other AWS services
	// with loops like this.
	var servOut *ecs.CreateServiceOutput
	var err error
OUTER:
	for i := 0; i <= awsCreateRetries; i++ {
		servOut, err = ecsSvc.CreateServiceWithContext(ctx, createServiceInput)
		if err == nil {
			break
		}

		// if we encounter an unrecoverable error, exit and skip this loop
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "AccessDeniedException", "UnsupportedFeatureException",
				"PlatformUnknownException",
				"PlatformTaskDefinitionIncompatibilityException":
				break OUTER
			}
		}

		s.Update("Failed to register ecs service. Will retry in %d seconds (up to %d more times)\nError: %s", awsCreateRetryIntervalSeconds, awsCreateRetries-i, err)

		// otherwise sleep and try again
		time.Sleep(awsCreateRetryIntervalSeconds * time.Second)
	}
	if err != nil {
		return status.Errorf(codes.Internal, "failed registering ecs service: %s", err)
	}

	state.Name = *servOut.Service.ServiceName
	state.Arn = *servOut.Service.ServiceArn

	s.Update("Created ECS Service %s", serviceName)
	s.Done()
	return nil
}

func (p *Platform) resourceServiceStatus(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	state *Resource_Service,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Determining status of ecs service %s", state.Name)
	defer s.Abort()

	ecsSvc := ecs.New(sess)

	servicesResp, err := ecsSvc.DescribeServicesWithContext(ctx, &ecs.DescribeServicesInput{
		Services: []*string{&state.Name},
		Cluster:  &state.Cluster,
	})
	if _, ok := err.(*ecs.ClusterNotFoundException); ok {
		sr.Resources = append(sr.Resources, &sdk.StatusReport_Resource{
			Name:          state.Name,
			Id:            state.Arn,
			Health:        sdk.StatusReport_MISSING,
			HealthMessage: fmt.Sprintf("Cluster named %q is missing", state.Cluster),
		})
		s.Done()
		return nil
	}
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe service (ARN %q): %s", state.Arn, err)
	}
	if len(servicesResp.Services) == 0 {
		sr.Resources = append(sr.Resources, &sdk.StatusReport_Resource{
			Name:          state.Name,
			Id:            state.Arn,
			Health:        sdk.StatusReport_MISSING,
			HealthMessage: fmt.Sprintf("service %q is missing", state.Name),
		})
		s.Done()
		return nil
	}

	service := servicesResp.Services[0]

	serviceResource := sdk.StatusReport_Resource{
		Name:                *service.ServiceName,
		Id:                  *service.ServiceArn,
		CreatedTime:         timestamppb.New(*service.CreatedAt),
		PlatformUrl:         fmt.Sprintf("https://console.aws.amazon.com/ecs/home?region=%s#/clusters/waypoint/services/%s", p.config.Region, state.Name),
		Type:                "service",
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
		HealthMessage:       fmt.Sprintf("service is %q", *service.Status),
	}
	sr.Resources = append(sr.Resources, &serviceResource)

	if *service.Status == "ACTIVE" {
		serviceResource.Health = sdk.StatusReport_READY
	} else {
		serviceResource.Health = sdk.StatusReport_DOWN
		serviceResource.HealthMessage = fmt.Sprintf("service is %q", *service.Status)
	}

	serviceJson, err := json.Marshal(map[string]interface{}{"service": service})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal service %q (ARN %q) state to json: %s", *service.ServiceName, *service.ServiceArn, err)
	}
	serviceResource.StateJson = string(serviceJson)

	taskArns, err := ecsSvc.ListTasksWithContext(ctx, &ecs.ListTasksInput{
		ServiceName: &state.Name,
		Cluster:     &state.Cluster,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list tasks for service %q in cluster %q: %s", state.Name, state.Cluster, err)
	}

	// Insert missing tasks if necessary
	missingCount := int(*service.DesiredCount) - len(taskArns.TaskArns)
	log.Debug("There are missing tasks. The service may be just starting up.", "missing count", missingCount, "service name", state.Name, "cluster", state.Cluster)
	for i := 0; i < missingCount; i++ {
		sr.Resources = append(sr.Resources, &sdk.StatusReport_Resource{
			Type:                "task",
			Name:                "missing",
			ParentResourceId:    *service.ServiceArn,
			Health:              sdk.StatusReport_MISSING,
			HealthMessage:       fmt.Sprintf("task is missing. The parent service %q may be just starting up at the time of this status check", state.Name),
			CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE,
		})
	}

	if len(taskArns.TaskArns) > 0 {
		tasks, err := ecsSvc.DescribeTasksWithContext(ctx, &ecs.DescribeTasksInput{
			Tasks:   taskArns.TaskArns,
			Cluster: &state.Cluster,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to describe tasks for service %q in cluster %q: %s", state.Name, state.Cluster, err)
		}

		for _, task := range tasks.Tasks {
			// Determine short task ID
			splitArn := strings.Split(*task.TaskArn, "/")
			taskId := splitArn[len(splitArn)-1]

			taskResource := &sdk.StatusReport_Resource{
				Type:                "task",
				Name:                taskId,
				ParentResourceId:    *service.ServiceArn,
				Id:                  *task.TaskArn,
				CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE,
				CreatedTime:         timestamppb.New(*task.CreatedAt),
				PlatformUrl:         fmt.Sprintf("https://console.aws.amazon.com/ecs/home?region=%s#/clusters/waypoint/tasks/%s", p.config.Region, taskId),
			}
			sr.Resources = append(sr.Resources, taskResource)

			switch strings.ToLower(*task.LastStatus) {
			case "running":
				taskResource.Health = sdk.StatusReport_READY
			case "provisioning", "pending", "activating":
				taskResource.Health = sdk.StatusReport_ALIVE
			default:
				taskResource.Health = sdk.StatusReport_DOWN
			}

			taskResource.HealthMessage = fmt.Sprintf("task is %q", *task.LastStatus)

			// Find IP address if possible

			var ipAddress string
			for _, attachment := range task.Attachments {
				for _, detail := range attachment.Details {
					if *detail.Name == "privateIPv4Address" {
						ipAddress = *detail.Value
					}
				}
			}

			stateJson, err := json.Marshal(map[string]interface{}{
				"ipAddress": ipAddress,
				"task":      task,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to marshal task (arn %q) state to json: %s", *task.TaskArn, err)
			}
			taskResource.StateJson = string(stateJson)
		}
	}
	s.Done()
	return nil
}

func (p *Platform) resourceServiceDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	state *Resource_Service,
	workspaceDestroy WorkspaceDestroy,
) error {
	if workspaceDestroy {
		// Services are destroyed when a deployment is destroyed, not a whole workspace.
		return nil
	}

	log.Debug("deleting ecs service", "arn", state.Arn)
	if state.Arn == "" {
		log.Debug("Missing ECS Service ARN - it must not have been created successfully. Skipping delete.")
		return nil
	}

	s := sg.Add("Deleting service %s", state.Name)
	defer s.Abort()

	_, err := ecs.New(sess).DeleteServiceWithContext(ctx, &ecs.DeleteServiceInput{
		Cluster: &state.Cluster,
		Force:   aws.Bool(true),
		Service: &state.Arn,
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok && awserr.Code() == "ServiceNotFoundException" {
			s.Update("Service does not exist - it must have already been deleted. (ARN: %q)", state.Arn)
			s.Done()
			return nil
		}
		return status.Errorf(codes.Internal, "failed to delete ECS service %s (ARN: %q): %s", state.Name, state.Arn, err)
	}

	s.Update("Deleted service %s", state.Name)
	s.Done()
	return nil
}

// resourceAlbListenerCreate finds or creates an ALB listener, and ensures that the
// target group is added to it.
func (p *Platform) resourceAlbListenerCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	alb *Resource_Alb,
	targetGroup *Resource_TargetGroup,
	externalIngressPort ExternalIngressPort,

	state *Resource_Alb_Listener,
) error {
	if p.config.DisableALB {
		log.Debug("ALB disabled - skipping alb listener creation")
		return nil
	}

	s := sg.Add("Initiating ALB Listener")
	defer s.Abort()

	state.TargetGroup = targetGroup

	albConfig := p.config.ALB

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: &targetGroup.Arn,
			Weight:         aws.Int64(0),
		},
	}

	var (
		certs    []*elbv2.Certificate
		protocol = "HTTP"
	)

	if albConfig != nil && albConfig.CertificateId != "" {
		protocol = "HTTPS"
		certs = append(certs, &elbv2.Certificate{
			CertificateArn: &albConfig.CertificateId,
		})
	}

	elbsrv := elbv2.New(sess)

	if alb == nil || alb.Arn == "" {
		return status.Errorf(codes.InvalidArgument, "cannot create ALB listener - no existing ALB defined.")
	}
	log.Info("load-balancer defined", "dns-name", alb.DnsName)

	s.Update("Looking for listeners for ALB %q", alb.Arn)
	listeners, err := elbsrv.DescribeListenersWithContext(ctx, &elbv2.DescribeListenersInput{
		LoadBalancerArn: &alb.Arn,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe listeners for alb (ARN %q): %s", alb.Arn, err)
	}

	// Check for existing listeners on our specified port
	var listener *elbv2.Listener
	for _, l := range listeners.Listeners {
		if *l.Port == int64(externalIngressPort) {
			listener = l
		}
	}

	if listener != nil {
		// Found an existing listener!
		s.Update("Modifying existing ALB Listener to introduce target group")

		def := listener.DefaultActions

		if len(def) > 0 && def[0].ForwardConfig != nil {
			for _, tg := range def[0].ForwardConfig.TargetGroups {
				if *tg.Weight > 0 {
					tgs = append(tgs, tg)
					log.Debug("previous target group", "arn", *tg.TargetGroupArn)
				}
			}
		}

		in := &elbv2.ModifyListenerInput{
			ListenerArn: listener.ListenerArn,
			DefaultActions: []*elbv2.Action{
				{
					ForwardConfig: &elbv2.ForwardActionConfig{
						TargetGroups: tgs,
					},
					Type: aws.String("forward"),
				},
			},

			// Setting these on every deployment so that we pick up any potential changes
			// in the waypoint.hcl config
			Port:         aws.Int64(int64(externalIngressPort)),
			Protocol:     aws.String(protocol),
			Certificates: certs,
		}

		_, err = elbsrv.ModifyListenerWithContext(ctx, in)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to introduce new target group to existing ALB listener: %s", err)
		}

		s.Update("Modified ALB Listener to introduce target group")
	} else {
		s.Update("Creating new ALB Listener")
		state.Managed = true

		tgs[0].Weight = aws.Int64(100)
		lo, err := elbsrv.CreateListenerWithContext(ctx, &elbv2.CreateListenerInput{
			LoadBalancerArn: &alb.Arn,
			Port:            aws.Int64(int64(externalIngressPort)),
			Protocol:        aws.String(protocol),
			Certificates:    certs,
			DefaultActions: []*elbv2.Action{
				{
					ForwardConfig: &elbv2.ForwardActionConfig{
						TargetGroups: tgs,
					},
					Type: aws.String("forward"),
				},
			},
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to create listener: %s", err)
		}

		listener = lo.Listeners[0]

		s.Update("Created ALB Listener")
		log.Debug("Created ALB Listener", "arn", *listener.ListenerArn)
	}

	state.Arn = *listener.ListenerArn

	s.Done()
	return nil
}

// resourceAlbListenerDestroy destroys the ALB listener associated with this deployment
// if it is under waypoint's management, and if the only target group it forwards to
// is this deployment's target group.
func (p *Platform) resourceAlbListenerDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	state *Resource_Alb_Listener,
	workspaceDestroy WorkspaceDestroy,
) error {
	if workspaceDestroy {
		// alb listeners are destroyed when a deployment is destroyed, not a workspace
		return nil
	}
	if !state.Managed {
		log.Debug("Skipping destroy of unmanaged ALB listener with ARN %q", state.Arn)
		return nil
	}
	if state.Arn == "" {
		log.Debug("Missing alb listener ARN - it must not have been created successfully. Skipping delete.")
		return nil
	}

	s := sg.Add("Initiating deletion of ALB Listener (ARN: %q)", state.Arn)
	defer s.Abort()

	elbsrv := elbv2.New(sess)
	s.Update("Describing ALB listener (ARN: %q)", state.Arn)

	listeners, err := elbsrv.DescribeListenersWithContext(ctx, &elbv2.DescribeListenersInput{
		ListenerArns: []*string{&state.Arn},
	})
	if err != nil {
		// There doesn't seem to be an aws error to cast to for this code.
		if strings.Contains(err.Error(), "ListenerNotFound") {
			s.Update("Listener does not exist and must have been destroyed (ARN %q).", state.Arn)
			s.Done()
			return nil
		}
		return status.Errorf(codes.Internal, "failed to describe listener with ARN %q: %s", state.Arn, err)
	}

	if len(listeners.Listeners) == 0 {
		// Could happen if listener was deleted out-of-band
		s.Update("ALB listener does not exist - not deleting (ARN: %q)", state.Arn)
		s.Done()
		return nil
	}

	listener := listeners.Listeners[0]

	log.Debug("listener arn", "arn", *listener.ListenerArn)

	def := listener.DefaultActions

	var tgs []*elbv2.TargetGroupTuple

	// If there is only 1 target group, delete the listener
	if len(def) == 1 && len(def[0].ForwardConfig.TargetGroups) == 1 {
		log.Debug("only 1 target group, deleting listener")

		s.Update("Deleting ALB listener (ARN: %q)", state.Arn)
		_, err = elbsrv.DeleteListenerWithContext(ctx, &elbv2.DeleteListenerInput{
			ListenerArn: listener.ListenerArn,
		})

		if err != nil {
			return status.Errorf(codes.Internal, "failed to delete ALB listener (ARN %q): %s", *listener.ListenerArn, err)
		}
		s.Update("Deleted ALB Listener")
	} else if len(def) > 0 && def[0].ForwardConfig != nil && len(def[0].ForwardConfig.TargetGroups) > 1 {
		// Multiple target groups means we can keep the listener
		var active bool

		for _, tg := range def[0].ForwardConfig.TargetGroups {
			if *tg.TargetGroupArn != state.TargetGroup.Arn {
				tgs = append(tgs, tg)
				if *tg.Weight > 0 {
					active = true
				}
			}
		}

		// If there are no target groups active, then we just activate the first
		// one, otherwise we can't modify the listener.
		if !active && len(tgs) > 0 {
			tgs[0].Weight = aws.Int64(100)
		}

		log.Debug("modifying listener to remove target group", "target-groups", len(tgs))

		s.Update("Deregistering this deployment's target group from ALB listener")
		_, err = elbsrv.ModifyListenerWithContext(ctx, &elbv2.ModifyListenerInput{
			ListenerArn: listener.ListenerArn,
			Port:        listener.Port,
			Protocol:    listener.Protocol,
			DefaultActions: []*elbv2.Action{
				{
					ForwardConfig: &elbv2.ForwardActionConfig{
						TargetGroups: tgs,
					},
					Type: aws.String("forward"),
				},
			},
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to modify listener (ARN %q): %s", *listener.ListenerArn, err)
		}
		s.Update("Deregistered this deployment's target group from ALB listener")
	}

	s.Done()
	return nil
}

func (p *Platform) resourceTargetGroupCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	deploymentId DeploymentId,
	subnets *Resource_AlbSubnets, // Required because we need to know which VPC we're in, and subnets discover it.
	state *Resource_TargetGroup,
	securityGroupsInternal *Resource_InternalSecurityGroups,
	securityGroups *Resource_ExternalSecurityGroups,
) error {
	if p.config.DisableALB {
		log.Debug("ALB disabled - skipping target group creation")
		return nil
	}

	s := sg.Add("Initiating target group creation...")
	defer s.Abort()

	elbsrv := elbv2.New(sess)

	// Use our common deployment ID as the target group name
	targetGroupName := string(deploymentId)

	// We have to clamp at a length of 32 because the Name field
	// requires that the name is 32 characters or less.

	// NOTE(izaak): The random part of ULIDs seems to be near the end, so for long app names, we might not get unique names here.
	// Should use a different source of randomness than component ID
	if len(targetGroupName) > 32 {
		targetGroupName = targetGroupName[:32]
		log.Debug("using a shortened value for service name due to AWS's length limits", "serviceName", targetGroupName)
	}

	if subnets.Subnets.VpcId == "" {
		return status.Error(codes.Internal, "subnets failed to discover a VPC ID - cannot create target group")
	}

	state.Port = p.config.ServicePort

	if p.config.HealthCheck.Protocol == "" {
		p.config.HealthCheck.Protocol = "HTTP"
	}

	createTargetGroupInput := &elbv2.CreateTargetGroupInput{
		HealthCheckEnabled: aws.Bool(true),
		Name:               &targetGroupName,
		Port:               &state.Port,
		Protocol:           aws.String(p.config.HealthCheck.Protocol),
		TargetType:         aws.String("ip"),
		VpcId:              &subnets.Subnets.VpcId,
	}

	if p.config.HealthCheck.Path != "" {
		createTargetGroupInput.HealthCheckPath = aws.String(p.config.HealthCheck.Path)
	}

	if p.config.HealthCheck.Timeout != 0 {
		createTargetGroupInput.HealthCheckTimeoutSeconds = aws.Int64(p.config.HealthCheck.Timeout)
	}

	if p.config.HealthCheck.Interval != 0 {
		createTargetGroupInput.HealthCheckIntervalSeconds = aws.Int64(p.config.HealthCheck.Interval)
	}

	if p.config.HealthCheck.HealthyThresholdCount != 0 {
		createTargetGroupInput.HealthyThresholdCount = aws.Int64(p.config.HealthCheck.HealthyThresholdCount)
	}

	if p.config.HealthCheck.UnhealthyThresholdCount != 0 {
		createTargetGroupInput.UnhealthyThresholdCount = aws.Int64(p.config.HealthCheck.UnhealthyThresholdCount)
	}

	if p.config.HealthCheck.GRPCCode != "" {
		createTargetGroupInput.Matcher.GrpcCode = aws.String(p.config.HealthCheck.GRPCCode)
	}

	if p.config.HealthCheck.HTTPCode != "" {
		createTargetGroupInput.Matcher.HttpCode = aws.String(p.config.HealthCheck.HTTPCode)
	}

	ctg, err := elbsrv.CreateTargetGroupWithContext(ctx, createTargetGroupInput)
	if err != nil || ctg == nil || len(ctg.TargetGroups) == 0 {
		return status.Errorf(codes.Internal, "failed to create target group: %s", err)
	}

	state.Name = *ctg.TargetGroups[0].TargetGroupName
	state.Arn = *ctg.TargetGroups[0].TargetGroupArn

	s.Update("Created target group %s", state.Name)
	s.Done()
	return nil
}

func (p *Platform) resourceTargetGroupDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	state *Resource_TargetGroup,
	workspaceDestroy WorkspaceDestroy,
) error {
	if state.Arn == "" {
		log.Debug("Missing target group ARN - it must not have been created successfully. Skipping delete.")
		return nil
	}
	if workspaceDestroy {
		// The target group is destroyed when a deployment is destroyed, not a whole workspace.
		return nil
	}

	s := sg.Add("Getting details for target group %s", state.Name)
	defer s.Abort()

	elbsrv := elbv2.New(sess)

	groups, err := elbsrv.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: aws.StringSlice([]string{state.Arn}),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe target group %s (ARN: %q): %s", state.Name, state.Arn, err)
	} else if len(groups.TargetGroups) > 0 {
		return status.Errorf(codes.FailedPrecondition, "only one target group should be returned for ARN %q, but found %d matching target groups", state.Arn, len(groups.TargetGroups))
	}

	s.Update("Retrieved target group details")
	s.Done()

	s = sg.Add("Checking if target group %q has an active ALB listener", *groups.TargetGroups[0].TargetGroupName)
	// If there are any load balancers routing to the target group, loop until
	// the listeners are deleted by resourceAlbListenerDestroy, for a max of
	// 5 minutes
	var listenerArns []string
	describeListenersInput := &elbv2.DescribeListenersInput{}
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	listenerDeleted := false
	for !listenerDeleted {
		for _, lb := range groups.TargetGroups[0].LoadBalancerArns {
		CHECK_LISTENERS:
			if len(listenerArns) > 0 {
				describeListenersInput.ListenerArns = aws.StringSlice(listenerArns)
			} else {
				describeListenersInput.LoadBalancerArn = lb
			}

			listeners, err := elbsrv.DescribeListeners(describeListenersInput)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to describe listeners for ALB (ARN: %q): %s", *lb, err)
			}
			for _, listener := range listeners.Listeners {
				for _, defaultAction := range listener.DefaultActions {
					if *defaultAction.TargetGroupArn == state.Arn {
						s.Update("Found active ALB listener with ARN " + *listener.ListenerArn)
						// Save the listener ARN to search by it instead of
						// needlessly looping through all listeners on the
						// LB again. We hard-assign this to the first index
						// of the slice because there should be only one.
						listenerArns[0] = *listener.ListenerArn
						select {
						case <-ticker.C: // wait 5 seconds before checking again
						case <-ctx.Done():
							return status.Errorf(codes.Aborted, "Context "+
								"cancelled from timeout checking if listener for "+
								"target group was deleted: %s", ctx.Err())
						}
						goto CHECK_LISTENERS
					}
				}
			}
		}
		// We've checked all of the possible permutations of ALBs and
		// listeners, but none of the DefaultActions point to our target
		// group - if there was a listener, it's gone now.
		listenerDeleted = true
	}

	s.Update("ALB listener for target group is deleted")
	s.Done()

	s = sg.Add("Deleting target group...")
	// Destroying the listener earlier should have deregistered this target group, so it should be safe
	// to just delete
	_, err = elbsrv.DeleteTargetGroupWithContext(ctx, &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: &state.Arn,
	})
	if err != nil {
		// This doesn't seem to return an error if the target group does not exist.
		return status.Errorf(codes.Internal, "failed to delete target group %s (ARN: %q): %s", state.Name, state.Arn, err)
	}
	s.Update("Target group deleted")
	s.Done()
	return nil
}

func (p *Platform) resourceTaskDefinitionCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	state *Resource_TaskDefinition,

	// Outputs of other resource creation processes
	executionRole *Resource_ExecutionRole,
	taskRole *Resource_TaskRole,
	logGroup *Resource_LogGroup,
) error {
	s := sg.Add("Initiating task definition creation")
	defer s.Abort()

	// Build environment variables
	env := []*ecs.KeyValuePair{
		{
			Name:  aws.String("PORT"),
			Value: aws.String(fmt.Sprint(p.config.ServicePort)),
		},
	}

	for k, v := range p.config.Environment {
		env = append(env, &ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}

	for k, v := range deployConfig.Env() {
		env = append(env, &ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}

	// Build secrets
	var secrets []*ecs.Secret
	for k, v := range p.config.Secrets {
		secrets = append(secrets, &ecs.Secret{
			Name:      aws.String(k),
			ValueFrom: aws.String(v),
		})
	}

	// Build logging options
	defaultStreamPrefix := fmt.Sprintf("waypoint-%d", time.Now().Nanosecond())

	logOptions := buildLoggingOptions(
		p.config.Logging,
		p.config.Region,
		logGroup.Name,
		defaultStreamPrefix,
	)

	// Define app container
	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Name:      aws.String(src.App),
		Image:     aws.String(img.Name()),
		PortMappings: []*ecs.PortMapping{{
			ContainerPort: aws.Int64(p.config.ServicePort),
		}},
		Environment:       env,
		Memory:            utils.OptionalInt64(int64(p.config.Memory)),
		MemoryReservation: utils.OptionalInt64(int64(p.config.MemoryReservation)),
		Secrets:           secrets,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String("awslogs"),
			Options:   logOptions,
		},
	}

	// Define sidecar containers
	var additionalContainers []*ecs.ContainerDefinition
	for _, container := range p.config.ContainersConfig {
		var secrets []*ecs.Secret
		for k, v := range container.Secrets {
			secrets = append(secrets, &ecs.Secret{
				Name:      aws.String(k),
				ValueFrom: aws.String(v),
			})
		}

		var env []*ecs.KeyValuePair
		for k, v := range container.Environment {
			env = append(env, &ecs.KeyValuePair{
				Name:  aws.String(k),
				Value: aws.String(v),
			})
		}

		c := &ecs.ContainerDefinition{
			Essential: aws.Bool(false),
			Name:      aws.String(container.Name),
			Image:     aws.String(container.Image),
			PortMappings: []*ecs.PortMapping{
				{
					ContainerPort: aws.Int64(int64(container.ContainerPort)),
					HostPort:      aws.Int64(int64(container.HostPort)),
					Protocol:      aws.String(container.Protocol),
				},
			},
			Secrets:           secrets,
			Environment:       env,
			Memory:            utils.OptionalInt64(int64(container.Memory)),
			MemoryReservation: utils.OptionalInt64(int64(container.MemoryReservation)),
		}

		if container.HealthCheck != nil {
			c.SetHealthCheck(&ecs.HealthCheck{
				Command:     aws.StringSlice(container.HealthCheck.Command),
				Interval:    aws.Int64(container.HealthCheck.Interval),
				Timeout:     aws.Int64(container.HealthCheck.Timeout),
				Retries:     aws.Int64(container.HealthCheck.Retries),
				StartPeriod: aws.Int64(container.HealthCheck.StartPeriod),
			})
		}

		additionalContainers = append(additionalContainers, c)
	}

	containerDefinitions := append([]*ecs.ContainerDefinition{&def}, additionalContainers...)

	family := "waypoint-" + src.App
	s.Update("Registering Task definition: %s", family)

	var cpuShares int
	runtime := aws.String("FARGATE")
	if p.config.EC2Cluster {
		runtime = aws.String("EC2")
		cpuShares = p.config.CPU
	} else {
		if err := utils.ValidateEcsMemCPUPair(p.config.Memory, p.config.CPU); err != nil {
			return err
		}

		cpuValues := fargateResources[p.config.Memory]

		// at this point we know that config.CPU is either 0, or a valid value
		// for the memory given
		cpuShares = p.config.CPU
		if cpuShares == 0 {
			cpuShares = cpuValues[0]
		}
	}

	cpus := aws.String(strconv.Itoa(cpuShares))
	// on EC2 launch type, `Cpu` is an optional field, so we leave it nil if it is 0
	if p.config.EC2Cluster && cpuShares == 0 {
		cpus = nil
	}
	mems := strconv.Itoa(p.config.Memory)

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: containerDefinitions,

		ExecutionRoleArn: aws.String(executionRole.Arn),
		Cpu:              cpus,
		Memory:           aws.String(mems),
		Family:           aws.String(family),
		RuntimePlatform: &ecs.RuntimePlatform{
			CpuArchitecture: aws.String(p.config.Architecture),
		},

		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{runtime},

		Tags: []*ecs.Tag{
			{
				Key:   aws.String("waypoint-app"),
				Value: aws.String(src.App),
			},
		},
	}

	if taskRole != nil && taskRole.Arn != "" {
		registerTaskDefinitionInput.SetTaskRoleArn(taskRole.Arn)
	}

	ecsSvc := ecs.New(sess)

	var taskOut *ecs.RegisterTaskDefinitionOutput
	var err error
	// AWS is eventually consistent so even though we probably created the resources that
	// are referenced by the task definition, it can error out if we try to reference those resources
	// too quickly. So we're forced to guard actions which reference other AWS services
	// with loops like this.

OUTER:
	for i := 0; i <= awsCreateRetries; i++ {
		taskOut, err = ecsSvc.RegisterTaskDefinitionWithContext(ctx, &registerTaskDefinitionInput)
		if err == nil {
			break
		}

		// if we encounter an unrecoverable error, exit and skip this loop
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "ResourceConflictException" {
			switch aerr.Code() {
			case "ResourceConflictException":
				break OUTER
			}
		}

		s.Update("Failed to register ecs task definition. Will retry in %d seconds (up to %d more times)\nError: %s", awsCreateRetryIntervalSeconds, awsCreateRetries-i, err)
		s.Status(terminal.StatusWarn)

		// otherwise sleep and try again
		time.Sleep(awsCreateRetryIntervalSeconds * time.Second)
	}
	if err != nil {
		return status.Errorf(codes.Internal, "failed registering ecs task definition: %s", err)
	}

	s.Update("Registered Task definition: %s", family)

	state.Runtime = *runtime
	state.Arn = *taskOut.TaskDefinition.TaskDefinitionArn

	s.Done()
	return nil
}

func (p *Platform) resourceTaskDefinitionDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	state *Resource_TaskDefinition,
	taskRole *Resource_TaskRole,
	workspaceDestroy WorkspaceDestroy,
) error {
	if !workspaceDestroy {
		// Do not destroy app-scoped resource unless the whole workspace is being destroyed
		return nil
	}

	ecsSvc := ecs.New(sess)

	if taskRole != nil && taskRole.Arn != "" {
		return nil // user specified Task definition
	}

	s := sg.Add("Deleting ecs task definition")

	taskDefInput := ecs.DeregisterTaskDefinitionInput{TaskDefinition: &state.Arn}
	_, err := ecsSvc.DeregisterTaskDefinition(&taskDefInput)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete ecs task definition: %s", err)
	}

	s.Update("Deleted ECS task definition")
	s.Done()
	return nil
}

func (p *Platform) resourceAlbCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	src *component.Source,
	securityGroupsInternal *Resource_InternalSecurityGroups,
	securityGroups *Resource_ExternalSecurityGroups,
	subnets *Resource_AlbSubnets, // Required because we need to know which VPC we're in, and subnets discover it.
	state *Resource_Alb,
) error {
	s := sg.Add("Initiating ALB creation")
	defer s.Abort()

	if p.config.DisableALB {
		s.Update("ALB disabled - skipping ALB")
		s.Done()
		return nil
	}

	albConfig := p.config.ALB

	if albConfig != nil && albConfig.LoadBalancerArn != "" {
		s.Update("Existing ALB specified - no need to create or discover an ALB")
		s.Done()
		state.Managed = false
		state.Arn = albConfig.LoadBalancerArn
		return nil
	}

	// If not using an existing listener, the load balancer is owned by waypoint
	state.Managed = true

	var certs []*elbv2.Certificate
	if albConfig != nil && albConfig.CertificateId != "" {
		certs = append(certs, &elbv2.Certificate{
			CertificateArn: &albConfig.CertificateId,
		})
	}

	elbsrv := elbv2.New(sess)

	lbName := "waypoint-ecs-" + src.App
	state.Name = lbName

	s.Update("Looking for an existing load balancer named %s", lbName)

	var lb *elbv2.LoadBalancer
	dlb, err := elbsrv.DescribeLoadBalancersWithContext(ctx, &elbv2.DescribeLoadBalancersInput{
		Names: []*string{&lbName},
	})
	if err != nil {
		// If the load balancer wasn't found, we'll create it.
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() != elbv2.ErrCodeLoadBalancerNotFoundException {
			return status.Errorf(codes.Internal, "failed to describe load balancers with name %q: %s", lbName, err)
		}
		log.Debug("load balancer %s was not found - will create it.")
	}

	if dlb != nil && len(dlb.LoadBalancers) > 0 {
		lb = dlb.LoadBalancers[0]
		s.Update("Using existing ALB %s (%s, dns-name: %s)",
			lbName, *lb.LoadBalancerArn, *lb.DNSName)
	} else {
		s.Update("Creating new ALB: %s", lbName)

		scheme := elbv2.LoadBalancerSchemeEnumInternetFacing

		if albConfig != nil && albConfig.InternalScheme != nil && *albConfig.InternalScheme {
			log.Debug("Creating an internal scheme ALB")
			scheme = elbv2.LoadBalancerSchemeEnumInternal
		}

		subnetIds := make([]*string, len(subnets.Subnets.Subnets))
		for i, subnet := range subnets.Subnets.Subnets {
			subnetIds[i] = &subnet.Id
		}

		securityGroupIds := make([]*string, len(securityGroups.SecurityGroups))
		for i, securityGroup := range securityGroups.SecurityGroups {
			securityGroupIds[i] = &securityGroup.Id
		}

		clb, err := elbsrv.CreateLoadBalancerWithContext(ctx, &elbv2.CreateLoadBalancerInput{
			Name:           aws.String(lbName),
			Subnets:        subnetIds,
			SecurityGroups: securityGroupIds,
			Scheme:         &scheme,
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String("waypoint_managed"),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String("waypoint_app"),
					Value: aws.String(src.App),
				},
			},
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to create ALB %q: %s", lbName, err)
		}

		lb = clb.LoadBalancers[0]

		s.Update("Created ALB: %s (dns-name: %s)", lbName, *lb.DNSName)
	}
	state.Arn = *lb.LoadBalancerArn

	state.DnsName = *lb.DNSName
	state.CanonicalHostedZoneId = *lb.CanonicalHostedZoneId

	s.Update("Using Application Load Balancer %q", state.Name)
	s.Done()
	return nil
}

func (p *Platform) resourceAlbDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	state *Resource_Alb,
	workspaceDestroy WorkspaceDestroy,
) error {
	if !workspaceDestroy {
		// Do not destroy app-scoped resource unless the whole workspace is being destroyed
		return nil
	}

	if p.config.DisableALB {
		log.Debug("ALB disabled - skipping target group creation")
		return nil
	}

	s := sg.Add("Initializing ALB deletion")
	defer s.Abort()

	s.Update("Deleting ALB %s", state.DnsName)

	elbsrv := elbv2.New(sess)
	input := elbv2.DeleteLoadBalancerInput{LoadBalancerArn: &state.Arn}

	_, err := elbsrv.DeleteLoadBalancer(&input)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to remove ALB %s: %s", state.DnsName, err)
	}

	s.Update("Deleted ALB %s", state.Arn)
	s.Done()
	return nil
}

func (p *Platform) resourceRoute53RecordCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	alb *Resource_Alb,
	state *Resource_Route53Record,
) error {
	albConfig := p.config.ALB
	if p.config.DisableALB || albConfig == nil || albConfig.ZoneId == "" || albConfig.FQDN == "" {
		log.Debug("Not creating a route53 record")
		return nil
	}

	s := sg.Add("Route53 record is required - checking if one already exists")
	defer s.Abort()

	r53 := route53.New(sess)

	records, err := r53.ListResourceRecordSetsWithContext(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(albConfig.ZoneId),
		StartRecordName: aws.String(albConfig.FQDN),
		StartRecordType: aws.String(route53.RRTypeA),
		MaxItems:        aws.String("1"),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list resource records for alb %q: %s", state.Name, err)
	}

	fqdn := albConfig.FQDN

	// Add trailing period to match Route53 record name
	if fqdn[len(fqdn)-1] != '.' {
		fqdn += "."
	}

	var recordExists bool

	if len(records.ResourceRecordSets) > 0 {
		record := records.ResourceRecordSets[0]
		if aws.StringValue(record.Type) == route53.RRTypeA && aws.StringValue(record.Name) == fqdn {
			s.Update("Found existing Route53 record: %s", aws.StringValue(record.Name))
			log.Debug("found existing record, assuming it's correct")
			recordExists = true
		}
	}

	if !recordExists {
		s.Update("Creating new Route53 record: %s (zone-id: %s)",
			albConfig.FQDN, albConfig.ZoneId)

		log.Debug("creating new route53 record", "zone-id", albConfig.ZoneId)
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action: aws.String(route53.ChangeActionCreate),
						ResourceRecordSet: &route53.ResourceRecordSet{
							Name: aws.String(albConfig.FQDN),
							Type: aws.String(route53.RRTypeA),
							AliasTarget: &route53.AliasTarget{
								DNSName:              &alb.DnsName,
								EvaluateTargetHealth: aws.Bool(true),
								HostedZoneId:         &alb.CanonicalHostedZoneId,
							},
						},
					},
				},
				Comment: aws.String("managed by waypoint"),
			},
			HostedZoneId: aws.String(albConfig.ZoneId),
		}

		result, err := r53.ChangeResourceRecordSetsWithContext(ctx, input)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to create route53 record %q: %s", albConfig.FQDN, err)
		}
		log.Debug("record created", "change-id", *result.ChangeInfo.Id)

		s.Update("Created Route53 record: %s (zone-id: %s)",
			albConfig.FQDN, albConfig.ZoneId)
	}

	s.Done()
	return nil
}

// Returns subnet resource based on the configured subnets to use.
// If no configured subnets are given, will use the default subnets for the VPC
func getSubnetsResource(
	ctx context.Context,
	s terminal.Step,
	log hclog.Logger,
	sess *session.Session,
	configuredSubnets []string,
) (*Resource_Subnets, error) {
	svc := ec2.New(sess)

	var describeSubnetsInput *ec2.DescribeSubnetsInput

	if len(configuredSubnets) == 0 {
		log.Debug("Using default subnets")
		describeSubnetsInput = &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("default-for-az"),
					Values: []*string{aws.String("true")},
				},
			},
		}
	} else {
		log.Debug("Using configured subnets", "subnet ids", strings.Join(configuredSubnets, ", "))
		subnetIdPtrs := make([]*string, len(configuredSubnets))
		for i := range configuredSubnets {
			subnetIdPtrs[i] = &configuredSubnets[i]
		}
		describeSubnetsInput = &ec2.DescribeSubnetsInput{
			SubnetIds: subnetIdPtrs,
		}
	}

	s.Update("Describing subnets")
	desc, err := svc.DescribeSubnetsWithContext(ctx, describeSubnetsInput)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to describe subnets: %s", err)
	}

	if len(desc.Subnets) == 0 {
		if len(configuredSubnets) == 0 {
			return nil, status.Errorf(codes.FailedPrecondition, "No subnets specified, and no default subnets found in the vpc. Either specify subnets explicitly, or assign default subnets for this vpc.")
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "subnet ids %q do not exist", strings.Join(configuredSubnets, ", "))
		}
	}

	var subnets Resource_Subnets
	subnets.VpcId = *desc.Subnets[0].VpcId
	for _, subnet := range desc.Subnets {
		subnets.Subnets = append(subnets.Subnets, &Resource_Subnets_Subnet{Id: *subnet.SubnetId})
	}

	return &subnets, nil
}

func (p *Platform) resourceServiceSubnetsDiscover(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	state *Resource_ServiceSubnets,
) error {
	s := sg.Add("Discovering which subnets to use for the service")
	defer s.Abort()
	subnets, err := getSubnetsResource(ctx, s, log, sess, p.config.Subnets)
	if err != nil {
		return status.Errorf(status.Convert(err).Code(), "failed discovering service subnets: %s", err)
	}
	state.Subnets = subnets
	log.Debug("Service subnets are in vpc", "vpc-name", subnets.VpcId)

	s.Update("Discovered service subnets")
	s.Done()
	return nil
}

func (p *Platform) resourceAlbSubnetsDiscover(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	state *Resource_AlbSubnets,
) error {
	if p.config.DisableALB {
		log.Debug("ALB disabled - skipping alb subnet discovery")
		return nil
	}

	s := sg.Add("Discovering which subnets to use for the alb")
	defer s.Abort()
	var configuredSubnets []string
	if p.config.ALB != nil && len(p.config.ALB.Subnets) > 0 {
		log.Debug("Using configured ALB subnets", "subnet ids", strings.Join(p.config.ALB.Subnets, ", "))
		configuredSubnets = p.config.ALB.Subnets
	}
	subnets, err := getSubnetsResource(ctx, s, log, sess, configuredSubnets)
	if err != nil {
		return status.Errorf(status.Convert(err).Code(), "failed discovering alb subnets: %s", err)
	}
	state.Subnets = subnets
	log.Debug("Alb subnets are in vpc", "vpc-name", subnets.VpcId)

	s.Update("Discovered alb subnets")
	s.Done()
	return nil
}

func (p *Platform) resourceExternalSecurityGroupsCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	src *component.Source,
	subnets *Resource_AlbSubnets, // Required because we need to know which VPC we're in, and subnets discover it.
	externalIngressPort ExternalIngressPort,
	state *Resource_ExternalSecurityGroups,
) error {

	if p.config.DisableALB {
		return nil
	}

	name := fmt.Sprintf("%s-inbound", src.App)
	s := sg.Add("Initiating creation of external security group named %s", name)
	defer s.Abort()

	if p.config.ALB != nil && p.config.ALB.SecurityGroupIDs != nil {
		s.Update("Using specified ALB security group IDs")
		for _, sgId := range p.config.SecurityGroupIDs {
			state.SecurityGroups = append(state.SecurityGroups, &Resource_SecurityGroup{Id: *sgId, Managed: false})
		}
		s.Done()
		return nil
	}

	protocol := "tcp"
	cidr := "0.0.0.0/0"
	cidrDescription := "all traffic"
	port := int64(externalIngressPort)
	perms := []*ec2.IpPermission{{
		IpProtocol: &protocol,
		FromPort:   &port,
		ToPort:     &port,
		IpRanges: []*ec2.IpRange{{
			CidrIp:      &cidr,
			Description: &cidrDescription,
		}},
	}}

	securityGroup, err := upsertSecurityGroup(ctx, sess, s, name, subnets.Subnets.VpcId, perms)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upsert security group %q: %s", name, err)
	}

	state.SecurityGroups = append(state.SecurityGroups, securityGroup)
	s.Update("Using external security group %s", securityGroup.Name)

	s.Done()
	return nil
}

func (p *Platform) resourceSecurityGroupsDestroy(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	log hclog.Logger,
	externalState *Resource_ExternalSecurityGroups,
	internalState *Resource_InternalSecurityGroups,
	workspaceDestroy WorkspaceDestroy,
) error {
	if !workspaceDestroy {
		// Do not destroy app-scoped resource unless the whole workspace is being destroyed
		return nil
	}

	ids := []string{}

	// only delete the managed security groups
	// internal groups have a dependency on the external groups
	// so we must delete these first
	for _, sg := range internalState.SecurityGroups {
		if sg.Managed {
			ids = append(ids, sg.Id)
		}
	}

	for _, sg := range externalState.SecurityGroups {
		if sg.Managed {
			ids = append(ids, sg.Id)
		}
	}

	err := deleteSecurityGroups(ctx, sess, sg, log, ids)
	if err != nil {
		return err
	}

	return nil
}

func (p *Platform) resourceInternalSecurityGroupsCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	src *component.Source,
	subnets *Resource_ServiceSubnets, // Required because we need to know which VPC we're in, and subnets discover it.
	extSecurityGroup *Resource_ExternalSecurityGroups,
	state *Resource_InternalSecurityGroups,
) error {
	s := sg.Add("Initiating security group creation...")
	defer s.Abort()

	if p.config.SecurityGroupIDs != nil {
		s.Update("Using specified security group IDs")
		for _, sgId := range p.config.SecurityGroupIDs {
			state.SecurityGroups = append(state.SecurityGroups, &Resource_SecurityGroup{Id: *sgId, Managed: false})
		}
		s.Done()
		return nil
	}

	name := fmt.Sprintf("%s-inbound-internal", src.App)

	s.Update("No security groups specified - checking for existing security group named %q", name)

	if extSecurityGroup == nil || len(extSecurityGroup.SecurityGroups) == 0 || extSecurityGroup.SecurityGroups[0] == nil {
		return status.Errorf(codes.Internal, "cannot create internal security group without a reference to the external security group ID")
	}

	extSgId := extSecurityGroup.SecurityGroups[0].Id

	protocol := "tcp"
	perms := []*ec2.IpPermission{{
		IpProtocol: &protocol,
		FromPort:   &p.config.ServicePort,
		ToPort:     &p.config.ServicePort,
		UserIdGroupPairs: []*ec2.UserIdGroupPair{{
			GroupId: &extSgId,
		}},
	}}

	securityGroup, err := upsertSecurityGroup(ctx, sess, s, name, subnets.Subnets.VpcId, perms)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upsert security group %q: %s", name, err)
	}

	state.SecurityGroups = append(state.SecurityGroups, securityGroup)
	s.Update("Using internal security group %s", securityGroup.Name)

	s.Done()
	return nil
}

// Finds a security group by name, and creates one if it does not exist.
func upsertSecurityGroup(
	ctx context.Context,
	sess *session.Session,
	s terminal.Step,

	name string,
	vpcId string,
	perms []*ec2.IpPermission,
) (*Resource_SecurityGroup, error) {
	ec2srv := ec2.New(sess)

	s.Update("Looking for existing security group named %q", name)
	dsg, err := ec2srv.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(name)},
			},
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to describe security groups named %q: %s", name, err)
	}

	// We only upsert security groups that we manage
	sg := &Resource_SecurityGroup{Managed: true}

	if len(dsg.SecurityGroups) != 0 {
		sg.Id = *dsg.SecurityGroups[0].GroupId
		sg.Name = *dsg.SecurityGroups[0].GroupName
		s.Update("Using existing security group with ID %s", sg.Id)
	} else {
		s.Update("Creating new security group named %s", name)

		out, err := ec2srv.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("created by waypoint"),
			GroupName:   aws.String(name),
			VpcId:       &vpcId,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create security group %q: %s", name, err)
		}

		sg.Id = *out.GroupId
		sg.Name = name

		s.Update("Authorizing ingress on newly created security group %s", name)
		_, err = ec2srv.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       &sg.Id,
			IpPermissions: perms,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to authorize ingress to security group %q: %s", sg.Name, err)
		}

		s.Update("Created and configured security group %s", name)
	}

	return sg, nil
}

func deleteSecurityGroups(
	ctx context.Context,
	sess *session.Session,
	sg terminal.StepGroup,
	log hclog.Logger,
	ids []string,
) error {
	ec2srv := ec2.New(sess)

	for _, id := range ids {
		req := ec2.DeleteSecurityGroupInput{GroupId: &id}

		// retry until context cancelled as AWS resource deletion is eventually consistent, we need to
		// wait for resources that have been linked to this security group to be removed.
		var returnErr error
		s := sg.Add("Attempting to delete security group %s", id)

		for i := 0; i <= awsDestroyRetries; i++ {
			if err := ctx.Err(); err != nil {
				return returnErr
			}

			s.Update("Attempting to delete security group %s", id)
			_, err := ec2srv.DeleteSecurityGroup(&req)

			if err == nil {
				s.Update("Deleted security group %s", id)
				returnErr = nil
				break
			}

			aerr, isAwsErr := err.(awserr.Error)
			if isAwsErr && aerr.Code() == "DependencyViolation" {
				// Expected case - it takes AWS a while to get around to actually removing the ENI.
				s.Update("Security group %s still has a dependency - waiting for it to be deleted.\nWill retry in %d seconds (up to %d more times)", id, awsDestroyRetryIntervalSeconds, awsDestroyRetries-i)

				// otherwise sleep and try again
				time.Sleep(awsDestroyRetryIntervalSeconds * time.Second)

				returnErr = err // If we time out after hitting this case, return the error
				continue
			} else if isAwsErr && aerr.Code() == "InvalidGroup.NotFound" {
				// Will happen if the security group has been deleted out of band
				s.Update("Security group %s does not exist", id)
				returnErr = nil
				break
			} else {
				// Some unexpected error
				s.Update("unable to delete security group id %s: %s", id, err)
				returnErr = status.Errorf(codes.Internal, "error deleting security group id %s: %s", id, err)
				break
			}

		}

		if returnErr != nil {
			s.Abort()
			return returnErr
		}
		s.Done()
	}

	return nil
}

func (p *Platform) resourceLogGroupCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	sess *session.Session,
	state *Resource_LogGroup,
) error {
	s := sg.Add("Initiating log group creation...")
	defer s.Abort()

	logGroup := p.config.LogGroup
	if logGroup == "" {
		logGroup = "waypoint-logs"
	}

	s.Update("Looking for existing log group named %s", logGroup)

	cwl := cloudwatchlogs.New(sess)
	groups, err := cwl.DescribeLogGroupsWithContext(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              aws.Int64(1),
		LogGroupNamePrefix: aws.String(logGroup),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe log groups: %s", err)
	}

	if len(groups.LogGroups) == 1 {
		s.Update("Using existing log group %s", logGroup)
		lg := groups.LogGroups[0]
		state.Name = *lg.LogGroupName
		state.Arn = *lg.Arn
		s.Done()
		return nil
	}

	s.Update("No existing log group found - creating new CloudWatchLogs group to store logs in: %s", logGroup)

	_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroup),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed creating log group %s: %s", logGroup, err)
	}

	//NOTE(izaak): CreateLogGroup doesn't return the log group ARN.
	state.Name = logGroup

	s.Update("Created CloudWatchLogs group to store logs in: %s", logGroup)
	s.Done()
	return nil
}

func (p *Platform) resourceTaskRoleCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	src *component.Source,
	sess *session.Session,
	state *Resource_TaskRole,
) error {
	roleName := p.config.TaskRoleName

	if roleName == "" && len(p.config.TaskRolePolicyArns) == 0 {
		log.Debug("No task role name or task role policies specified - skipping task role creation")
		return nil
	}

	s := sg.Add("Initiating task role creation...")
	defer s.Abort()

	if roleName == "" {
		roleName = src.App + "-task-role"
	}
	state.Name = roleName

	// role names have to be 64 characters or less, and the client side doesn't validate this.
	if len(roleName) > 64 {
		roleName = roleName[:64]
		log.Debug("using a shortened value for role name due to AWS's length limits", "roleName", roleName)
	}

	s.Update("Attempting to find an existing role named %q", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	svc := iam.New(sess)

	roleFound := true
	getOut, err := svc.GetRoleWithContext(ctx, queryInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NoSuchEntity" {
			roleFound = false
		} else {
			return status.Errorf(codes.Internal, "failed creating execution role %q: %s", roleName, err)
		}
	}

	if roleFound {
		s.Update("Found existing task IAM role: %s", *getOut.Role.RoleName)
		state.Arn = *getOut.Role.Arn
	} else {
		s.Update("No existing task role found - creating role %s", roleName)

		input := &iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(rolePolicy),
			Path:                     aws.String("/"),
			RoleName:                 aws.String(roleName),
		}

		result, err := svc.CreateRoleWithContext(ctx, input)
		if err != nil {
			return status.Errorf(codes.Internal, "failed creating execution role %q: %s", roleName, err)
		}
		state.Arn = *result.Role.Arn
	}

	if len(p.config.TaskRolePolicyArns) > 0 {
		s.Update("Task role policies specified - checking the existing policies on the task role %s", roleName)
		listOut, err := svc.ListAttachedRolePoliciesWithContext(ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: &roleName,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to list attached role policies on role (ARN: %q): %s", roleName, err)
		}

		for _, policyArn := range p.config.TaskRolePolicyArns {
			alreadyAttached := false
			for _, policy := range listOut.AttachedPolicies {
				if policyArn == *policy.PolicyArn {
					log.Debug("Requested policy arn is already attached to role", "policy arn", policyArn, "role name", roleName)
					alreadyAttached = true
					break
				}
			}

			if !alreadyAttached {
				s.Update("Attaching the following policy to role %s: %q", roleName, policyArn)
				_, err = svc.AttachRolePolicyWithContext(ctx, &iam.AttachRolePolicyInput{
					PolicyArn: &policyArn,
					RoleName:  &roleName,
				})
				if err != nil {
					return status.Errorf(codes.Internal, "failed to attach policy (ARN: %q) to role (name: %q): %s", policyArn, roleName, err)
				}
			}
		}
		s.Update("All requested policies are attached to role %s", roleName)
	}
	s.Done()
	return nil
}

func (p *Platform) resourceExecutionRoleCreate(
	ctx context.Context,
	sg terminal.StepGroup,
	log hclog.Logger,
	sess *session.Session,
	src *component.Source,
	state *Resource_ExecutionRole,
) error {
	s := sg.Add("Initiating execution role creation...")
	defer s.Abort()

	roleName := p.config.ExecutionRoleName

	if roleName == "" {
		roleName = "ecr-" + src.App
		state.Managed = true
	} else {
		// If the role name is defined, we're not managing this role, and shouldn't destroy it later.
		state.Managed = false
	}
	// role names have to be 64 characters or less, and the client side doesn't validate this.
	if len(roleName) > 64 {
		roleName = roleName[:64]
		log.Debug("using a shortened value for role name due to AWS's length limits", "roleName", roleName)
	}
	state.Name = roleName

	s.Update("Attempting to find an existing role named %q", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	svc := iam.New(sess)

	getOut, err := svc.GetRoleWithContext(ctx, queryInput)
	if err == nil {
		s.Update("Using existing execution IAM role %q", *getOut.Role.RoleName)

		// NOTE(izaak): We're verifying that the role exists, but not that it has the correct policy attached.
		// It's possible that we failed on that step earlier. We could call AttachRolePolicy every time, but
		// we're trying to minimize per-deployment aws api invocations for to stay under rate limits.

		state.Arn = *getOut.Role.Arn
		s.Done()
		return nil
	}
	// NOTE(izaak): the error returned here is an awserr.requestError, which cannot be cast to the public awserr.Error.
	// So we're forced to do this.
	if !strings.Contains(strings.ToLower(err.Error()), "status code: 404") {
		return status.Errorf(codes.Internal, "failed to get role with name %q: %s", roleName, err)
	}

	s.Update("No existing execution role found: creating IAM role %q", roleName)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
	}

	result, err := svc.CreateRoleWithContext(ctx, input)
	if err != nil {
		return status.Errorf(codes.Internal, "failed creating execution role %q: %s", roleName, err)
	}
	state.Arn = *result.Role.Arn

	s.Update("Attaching default execution policy to role %q", roleName)
	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(executionRolePolicyArn),
	}

	_, err = svc.AttachRolePolicyWithContext(ctx, aInput)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to attach policy %q to role %q: %s", executionRolePolicyArn, roleName, err)
	}

	s.Update("Created execution IAM role: %s", roleName)
	s.Done()
	return nil
}

// getSession is a value provider for resource manager and provides a client
// for use by resources to interact with AWS
func (p *Platform) getSession(log hclog.Logger) (*session.Session, error) {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create aws session: %s", err))
	}
	return sess, nil
}

// Loads state from the old locations on the top-level Deployment struct into resource manager.
// This is only necessary for backwards-compat with deployments created by waypoint pre-0.5.2.
// Newer waypoint deployments will have the complete resource manager state stored in deployment.ResourceState
func (p *Platform) loadResourceManagerState(
	ctx context.Context,
	rm *resource.Manager,
	deployment *Deployment,
	log hclog.Logger,
	sg terminal.StepGroup,
) error {
	log.Debug("Missing deployment resource state - must be an old deployment. Recovering state.")
	s := sg.Add("Gathering deployment resource state")
	defer s.Abort()

	rm.Resource("service").SetState(&Resource_Service{
		Cluster: deployment.Cluster,
		Arn:     deployment.ServiceArn,
	})

	targetGroupResource := Resource_TargetGroup{
		Arn: deployment.TargetGroupArn,
	}
	rm.Resource("target group").SetState(&targetGroupResource)

	// Restore state of ALB listener by inspecting the load balancer
	var listenerResource Resource_Alb_Listener
	listenerResource.TargetGroup = &targetGroupResource

	listenerResource.Managed = false // Impossible to know now if we created this initially, so this is the safe choice
	s.Update("Describing load balancer %s", deployment.LoadBalancerArn)
	sess, err := p.getSession(log)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get aws session: %s", err)
	}
	elbsrv := elbv2.New(sess)

	listeners, err := elbsrv.DescribeListenersWithContext(ctx, &elbv2.DescribeListenersInput{
		LoadBalancerArn: &deployment.LoadBalancerArn,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to describe listeners for ALB %q: %s", deployment.LoadBalancerArn, err)
	}
	if len(listeners.Listeners) == 0 {
		s.Update("No listeners found for ALB %q", deployment.LoadBalancerArn)
	} else {
		listenerResource.Arn = *listeners.Listeners[0].ListenerArn
		s.Update("Found existing listener (ARN: %q)", listenerResource.Arn)
	}
	rm.Resource("alb listener").SetState(&listenerResource)

	s.Update("Finished gathering resource state")
	s.Done()
	return nil
}

// A policy required for ECS roles - roles must be assumable by ECS tasks.
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

// The valid combinations of fargate cpu/memory resources. This will need to be updated
// over time as aws changes allowed combinations. It would be nice to get this from
// aws at runtime instead of hard-coding it in the future.
var fargateResources = map[int][]int{
	512:  {256},
	1024: {256, 512},
	2048: {256, 512, 1024},
	3072: {512, 1024},
	4096: {512, 1024},
	5120: {1024},
	6144: {1024},
	7168: {1024},
	8192: {1024},
}

// Builds logging options for use in an ECS task definition.
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
		// We receive the error `Log driver awslogs disallows options: awslogs-endpoint`
		// when setting `awslogs-endpoint`, so that is not included here of the
		// available options
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

type ALBConfig struct {
	// Certificate ARN to attach to the load balancer
	CertificateId string `hcl:"certificate,optional"`

	// Route53 Zone to setup record in
	ZoneId string `hcl:"zone_id,optional"`

	// Fully qualified domain name of the record to create in the target zone id
	FQDN string `hcl:"domain_name,optional"`

	// When set, waypoint will configure the target group into the load balancer.
	// This allows for usage of existing ALBs.
	LoadBalancerArn string `hcl:"load_balancer_arn"`

	// Indicates, when creating an ALB, that it should be internal rather than
	// internet facing.
	InternalScheme *bool `hcl:"internal,optional"`

	// Internet-facing traffic port. Defaults to 80 if CertificateId is unset, 443 if set.
	IngressPort int64 `hcl:"ingress_port,optional"`

	// Subnets to place the alb into. Defaults to the subnets in the default VPC.
	// This can be used to explicitly place the ALB on public subnets, leaving the service in private subnets.
	Subnets []string `hcl:"subnets,optional"`

	// Security Group ID of existing security group to use for ALB.
	SecurityGroupIDs []string `hcl:"security_group_ids,optional"`
}

type HealthCheckConfig struct {
	// A string array representing the command that the container runs to determine if it is healthy
	Command []string `hcl:"command"`

	// The time period in seconds between each health check execution
	Interval int64 `hcl:"interval,optional"`

	// The time period in seconds to wait for a health check to succeed before it is considered a failure
	Timeout int64 `hcl:"timeout,optional"`

	// The number of times to retry a failed health check before the container is considered unhealthy
	Retries int64 `hcl:"retries,optional"`

	// The optional grace period within which to provide containers time to bootstrap before failed health checks count towards the maximum number of retries
	StartPeriod int64 `hcl:"start_period,optional"`
}

type Logging struct {
	CreateGroup bool `hcl:"create_group,optional"`

	StreamPrefix string `hcl:"stream_prefix,optional"`

	DateTimeFormat string `hcl:"datetime_format,optional"`

	MultilinePattern string `hcl:"multiline_pattern,optional"`

	Mode string `hcl:"mode,optional"`

	MaxBufferSize string `hcl:"max_buffer_size,optional"`
}

type ContainerConfig struct {
	// The name of a container
	Name string `hcl:"name"`

	// The image used to start a container
	Image string `hcl:"image"`

	// The amount (in MiB) of memory to present to the container
	Memory int `hcl:"memory,optional"`

	// The soft limit (in MiB) of memory to reserve for the container
	MemoryReservation int `hcl:"memory_reservation,optional"`

	// The port number on the container
	ContainerPort int `hcl:"container_port,optional"`

	// The port number on the container instance to reserve for your container
	HostPort int `hcl:"host_port,optional"`

	// The protocol used for the port mapping
	Protocol string `hcl:"protocol,optional"`

	// The container health check command
	HealthCheck *HealthCheckConfig `hcl:"health_check,block"`

	// The environment variables to pass to a container
	Environment map[string]string `hcl:"static_environment,optional"`

	// The secrets to pass to a container
	Secrets map[string]string `hcl:"secrets,optional"`
}

type Config struct {
	// AWS Region to deploy into
	Region string `hcl:"region"`

	// Name of the Log Group to store logs into
	LogGroup string `hcl:"log_group,optional"`

	// Name of the ECS cluster to install the service into
	Cluster string `hcl:"cluster,optional"`

	// Name of the execution task IAM Role to associate with the ECS Service
	ExecutionRoleName string `hcl:"execution_role_name,optional"`

	// IAM Policy ARNs for attaching to task role
	TaskRolePolicyArns []string `hcl:"task_role_policy_arns,optional"`

	// Name of the task IAM role to associate with the ECS service
	TaskRoleName string `hcl:"task_role_name,optional"`

	// Subnets to place the service into. Defaults to the subnets in the default VPC.
	Subnets []string `hcl:"subnets,optional"`

	// Security Group IDs of existing security groups to use for ECS.
	SecurityGroupIDs []*string `hcl:"security_group_ids,optional"`

	// How many tasks of the service to run. Default 1.
	Count int `hcl:"count,optional"`

	// How much memory to assign to the containers
	Memory int `hcl:"memory"`

	// The soft limit (in MiB) of memory to reserve for the container
	MemoryReservation int `hcl:"memory_reservation,optional"`

	// How much CPU to assign to the containers
	CPU int `hcl:"cpu,optional"`

	Architecture string `hcl:"architecture,optional"`

	// The environment variables to pass to the main container
	Environment map[string]string `hcl:"static_environment,optional"`

	// The secrets to pass to the main container
	Secrets map[string]string `hcl:"secrets,optional"`

	// Assign each task a public IP. Default true.
	// Note: If private subnets have been specified for ECS tasks, and
	// no NAT gateway is configured, ECS will be unable to pull your image
	// from ECR and your services will be unable to start.
	// NOTE(izaak): If this is not set, we can probably set it to true if
	// the ecs subnet's route table has a nat gw default route, and true if they don't.
	AssignPublicIp *bool `hcl:"assign_public_ip,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to port 3000.
	ServicePort int64 `hcl:"service_port,optional"`

	// Indicate that service should be deployed on an EC2 cluster.
	EC2Cluster bool `hcl:"ec2_cluster,optional"`

	// If set to true, do not create a load balancer assigned to the service
	DisableALB bool `hcl:"disable_alb,optional"`

	// Configuration options for how the ALB will be configured.
	ALB *ALBConfig `hcl:"alb,block"`

	// Configuration options for additional containers
	ContainersConfig []*ContainerConfig `hcl:"sidecar,block"`

	Logging *Logging `hcl:"logging,block"`

	HealthCheck *AppHealthCheck `hcl:"health_check,block"`
}

type AppHealthCheck struct {
	// Protocol is the protocol for the health check to use - TCP, HTTP, HTTPS.
	Protocol string `hcl:"protocol,optional"`

	// Path is the destination of the ping path for target health checks.
	Path string `hcl:"path,optional"`

	// Timeout is the amount of time in seconds during which no target response means a failure.
	Timeout int64 `hcl:"timeout,optional"`

	// Interval is the amount of time in seconds between health checks.
	Interval int64 `hcl:"interval,optional"`

	// HealthyThresholdCount is the number of consecutive successful health checks required before considering an unhealthy target healthy.
	// The range is 2–10. The default is 5.
	HealthyThresholdCount int64 `hcl:"healthy_threshold_count,optional"`

	// UnhealthyThresholdCount is the number of consecutive failed health checks required before considering a target unhealthy.
	// The range is 2–10. The default is 2.
	UnhealthyThresholdCount int64 `hcl:"unhealthy_threshold_count,optional"`

	// GRPCCode is the gRPC codes to use when checking for a successful response from a target.
	// The default value is 12.
	GRPCCode string `hcl:"grpc_code,optional"`

	// HTTPCode is the HTTP codes to use when checking for a successful response from a target.
	// The range is 200 to 599. The default is 200-399.
	HTTPCode string `hcl:"http_code,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy the application into an ECS cluster on AWS")

	doc.Example(
		`
deploy {
  use "aws-ecs" {
    region = "us-east-1"
    memory = 512
  }
}
`)

	doc.Input("docker.Image")
	doc.Output("ecs.Deployment")

	doc.SetField(
		"region",
		"the AWS region for the ECS cluster",
	)

	doc.SetField(
		"log_group",
		"the CloudWatchLogs log group to store container logs into",
		docs.Default("derived from the application name"),
	)

	doc.SetField(
		"cluster",
		"the name of the ECS cluster to deploy into",
		docs.Summary(
			"the ECS cluster that will run the application as a Service.",
			"if there is no ECS cluster with this name, the ECS cluster will be",
			"created and configured to use Fargate to run containers.",
		),
	)

	doc.SetField(
		"execution_role_name",
		"the name of the IAM role to use for ECS execution",
		docs.Default("create a new exeuction IAM role based on the application name"),
	)

	doc.SetField(
		"task_role_name",
		"the name of the task IAM role to assign.",
		docs.Summary(
			"If no role exists and a one or more task role policies are requested,",
			"a role with this name will be created.",
		),
	)

	doc.SetField(
		"task_role_policy_arns",
		"IAM Policy arns for attaching to the task role.",
		docs.Summary(
			"If no task role name is specified a task role with a default name",
			"will be created for this app, and these policies will be attached.",
		),
	)

	doc.SetField(
		"subnets",
		"the VPC subnets to use for the service",
		docs.Default("public subnets in the default VPC"),
		docs.Summary(
			"you may set a list of private subnets here to prevent your tasks from being directly",
			"exposed publicly",
		),
	)

	doc.SetField(
		"security_group_ids",
		"Security Group IDs of existing security groups to use for the ECS service's network access",
		docs.Summary(
			"list of existing group IDs to use for the ECS service's network access.",
			"If none are specified, waypoint will create one. If DisableALB is false (the default), waypoint",
			"will only allow ingress from the ALB's security group",
		),
	)

	doc.SetField(
		"count",
		"how many instances of the application should run",
	)

	doc.SetField(
		"memory",
		"how much memory to assign to the container running the application",
		docs.Summary(
			"when running in Fargate, this must be one of a few values, specified in MB:",
			"512, 1024, 2048, 3072, 4096, 5120, and up to 16384 in increments of 1024.",
			"The memory value also controls the possible values for cpu",
		),
	)

	doc.SetField(
		"ec2_cluster",
		"indicate if the ECS cluster should be EC2 type rather than Fargate",
		docs.Summary(
			"this controls if we should verify the ECS cluster in EC2 type. The cluster",
			"will not be created if it doesn't exist, only that there as existing cluster",
			"this is using EC2 and not Fargate",
		),
	)

	doc.SetField(
		"disable_alb",
		"do not create a load balancer assigned to the service",
	)

	doc.SetField(
		"static_environment",
		"static environment variables to make available",
	)

	doc.SetField(
		"secrets",
		"secret key/values to pass to the ECS container",
	)

	doc.SetField(
		"alb",
		"Provides additional configuration for using an ALB with ECS",
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"certificate",
				"the ARN of an AWS Certificate Manager cert to associate with the ALB",
			)

			doc.SetField(
				"zone_id",
				"Route53 ZoneID to create a DNS record into",
				docs.Summary(
					"set along with alb.domain_name to have DNS automatically setup for the ALB",
				),
			)

			doc.SetField(
				"domain_name",
				"Fully qualified domain name to set for the ALB",
				docs.Summary(
					"set along with zone_id to have DNS automatically setup for the ALB.",
					"this value should include the full hostname and domain name, for instance",
					"app.example.com",
				),
			)

			doc.SetField(
				"internal",
				"Whether or not the created ALB should be internal",
				docs.Summary(
					"used when listener_arn is not set. If set, the created ALB will have a scheme",
					"of `internal`, otherwise by default it has a scheme of `internet-facing`.",
				),
			)

			doc.SetField(
				"load_balancer_arn",
				"the ARN on an existing ALB to configure",
				docs.Summary(
					"when this is set, Waypoint will use this ALB instead of creating",
					"its own. A target group will still be created for each deployment,",
					"and will be added to a listener on the configured ALB port",
					"(Waypoint will the listener if it doesn't exist).",
					"This allows users to configure their ALB outside Waypoint but still ",
					"have Waypoint hook the application to that ALB",
				),
			)

			doc.SetField(
				"ingress_port",
				"Internet-facing traffic port. Defaults to 80 if 'certificate' is unset, 443 if set.",
				docs.Summary("used to set the ALB listener port, and the ALB security group ingress port"),
			)
			doc.SetField(
				"subnets",
				"the VPC subnets to use for the ALB",
				docs.Default("public subnets in the default VPC"),
			)
		}),
	)

	doc.SetField(
		"logging",
		"Provides additional configuration for logging flags for ECS",
		docs.Summary(
			"Part of the ecs task definition.  These configuration flags help",
			"control how the awslogs log driver is configured."),

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"create_group",
				"Enables creation of the aws logs group if not present",
			)

			doc.SetField(
				"region",
				"The region the logs are to be shipped to",
				docs.Default("The same region the task is to be running"),
			)

			doc.SetField(
				"stream_prefix",
				"Prefix for application in cloudwatch logs path",
				docs.Default("Generated based off timestamp"),
			)

			doc.SetField(
				"datetime_format",
				"Defines the multiline start pattern in Python strftime format",
			)

			doc.SetField(
				"multiline_pattern",
				"Defines the multiline start pattern using a regular expression",
			)

			doc.SetField(
				"mode",
				"Delivery method for log messages, either 'blocking' or 'non-blocking'",
			)

			doc.SetField(
				"max_buffer_size",
				"When using non-blocking logging mode, this is the buffer size for message storage",
			)
		}),
	)

	doc.SetField("health_check",
		"Health check settings for the app.",
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField("protocol",
				"The protocol for the health check to use.")

			doc.SetField("path",
				"The destination of the ping path for the target health check.")

			doc.SetField("timeout",
				"The amount of time, in seconds, for which no target response "+
					"means a failure.")

			doc.SetField("interval",
				"The amount of time, in seconds, between health checks.")

			doc.SetField("healthy_threshold_count",
				"The number of consecutive successful health checks required to"+
					"consider a target healthy.",
				docs.Default("5"))

			doc.SetField("unhealthy_threshold_count",
				"The number of consecutive failed health checks required to "+
					"consider a target unhealthy.",
				docs.Default("2"))

			doc.SetField("matcher",
				"The range of HTTP codes to use when checking for a successful response from"+
					"the target.",
			)
		}))

	doc.SetField(
		"sidecar",
		"Additional container to run as a sidecar.",
		docs.Summary(
			"This runs additional containers in addition to the main container that",
			"comes from the build phase.",
		),

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"name",
				"Name of the container",
			)

			doc.SetField(
				"image",
				"Image of the sidecar container",
			)

			doc.SetField(
				"memory",
				"The amount (in MiB) of memory to present to the container",
			)

			doc.SetField(
				"memory_reservation",
				"The soft limit (in MiB) of memory to reserve for the container",
			)

			doc.SetField(
				"container_port",
				"The port number for the container",
			)

			doc.SetField(
				"host_port",
				"The port number on the host to reserve for the container",
			)

			doc.SetField(
				"protocol",
				"The protocol used for port mapping.",
			)

			doc.SetField(
				"static_environment",
				"Environment variables to expose to this container",
			)

			doc.SetField(
				"secrets",
				"Secrets to expose to this container",
			)
		}),
	)

	var memvals []int

	for k := range fargateResources {
		memvals = append(memvals, k)
	}

	sort.Ints(memvals)

	var sb strings.Builder

	for _, mem := range memvals {
		cpu := fargateResources[mem]

		var cpuVals []string

		for _, c := range cpu {
			cpuVals = append(cpuVals, strconv.Itoa(c))
		}

		fmt.Fprintf(&sb, "%dMB: %s\n", mem, strings.Join(cpuVals, ", "))
	}

	doc.SetField(
		"cpu",
		"how many cpu shares the container running the application is allowed",
		docs.Summary(
			"on Fargate, possible values for this are configured by the amount of memory",
			"the container is using. Here is a complete listing of possible values:\n",
			sb.String(),
		),
	)

	doc.SetField(
		"service_port",
		"the TCP port that the application is listening on",
		docs.Default("3000"),
	)

	doc.SetField(
		"assign_public_ip",
		"assign a public ip address to tasks. Defaults to true. Ignored if using an ec2 cluster.",
		docs.Default("true"),
		docs.Summary(
			"If this is set to false, deployments will fail unless tasks are able to egress to the",
			"container registry by some other means (i.e. a subnet default route to a NAT gateway).",
		),
	)

	doc.SetField(
		"architecture",
		"the instruction set CPU architecture that the Amazon ECS supports. Valid values are: \"x86_64\", \"arm64\"",
	)

	return doc, nil
}

var (
	mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current deployment is not ready, however your application
might be available or still starting up.
`)
)
