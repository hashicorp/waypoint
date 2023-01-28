package apprunner

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

type Platform struct {
	config PlatformConfig
}

type PlatformConfig struct {
	Region string `hcl:"region,optional"`

	Name string `hcl:"name"`

	Memory int `hcl:"memory,optional"`
	Cpu    int `hcl:"cpu,optional"`

	// Once created, Port cannot be modified
	Port int `hcl:"port,optional"`
	// Once created, RoleName cannot be modified
	RoleName string `hcl:"role_name,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// Deploy creates an AppRunner service, updates an
// AppRunner service, or no-ops if zero configuration
// changes are detected.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *ecr.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {

	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Connecting to AWS")
	defer func() {
		step.Abort()
	}()

	deployment := &Deployment{}
	if id, err := component.Id(); err != nil {
		return nil, err
	} else {
		deployment.Id = id
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	port := p.config.Port
	if port == 0 {
		// App Runner default port
		port = 8080
	}

	roleName := p.config.RoleName
	if roleName == "" {
		// if role name is not specified, we'll fall back to
		// getting or creating a service role, with the same
		// default roleName that AWS uses in the console
		roleName = defaultRoleName
	}

	// roleArn only applies to the initial create
	var roleArn string
	if arn, err := p.getOrCreateIamRoleArn(sess, log, roleName); err != nil {
		return nil, err
	} else {
		roleArn = arn
	}

	mem := int64(p.config.Memory)
	if mem == 0 {
		mem = defaultMemory
	}

	cpu := int64(p.config.Cpu)
	if cpu == 0 {
		cpu = defaultCpu
	}

	envVars := make(map[string]*string)
	for k, v := range p.config.StaticEnvVars {
		envVars[k] = aws.String(v)
	}

	step.Done()

	step = sg.Add("Deploying App Runner service %q", p.config.Name)
	step.Done()

	rm := p.resourceManager(log)

	// The operator's waypoint.hcl will realistically only have a svc name, not ARN,
	// since that is provided by AWS.
	//
	// While `DescribeService` only supports a service ARN, our best effort approach
	// is to list all services and manually match by user-provided name.
	step = sg.Add("Checking for existing service...")
	_service, err := p.getServiceSummaryByName(sess, log, p.config.Name)
	if err != nil {
		return nil, err
	}

	// Save state to resource manager
	rm.Resource("apprunner_service_summary").SetState(_service)
	step.Done()

	// update or create app runner server
	summary := rm.Resource("apprunner_service_summary").State().(*ServiceSummary)
	arSvc := apprunner.New(sess)

	operationId := ""
	serviceArn := ""
	serviceUrl := ""
	serviceStatus := ""

	// If we found a previous service, update it.
	if summary != nil {
		step = sg.Add("Found! Updating service %q", summary.Name)
		serviceArn = summary.Arn
		serviceUrl = summary.Url
		serviceStatus = summary.Status

		uso, err := arSvc.UpdateService(&apprunner.UpdateServiceInput{
			ServiceArn: aws.String(summary.Arn),
			InstanceConfiguration: &apprunner.InstanceConfiguration{
				Cpu:    aws.String(strconv.FormatInt(cpu, 10)),
				Memory: aws.String(strconv.FormatInt(mem, 10)),
			},
			SourceConfiguration: &apprunner.SourceConfiguration{
				AuthenticationConfiguration: &apprunner.AuthenticationConfiguration{
					AccessRoleArn: aws.String(roleArn),
				},
				ImageRepository: &apprunner.ImageRepository{
					ImageRepositoryType: aws.String(apprunner.ImageRepositoryTypeEcr),
					ImageIdentifier:     aws.String(img.Name()),
					ImageConfiguration: &apprunner.ImageConfiguration{
						Port:                        aws.String(strconv.Itoa(port)),
						RuntimeEnvironmentVariables: envVars,
					},
				},
				AutoDeploymentsEnabled: aws.Bool(false),
			},
		})

		if err != nil {
			step.Update("Failed to update service: %s", err)
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == apprunner.ErrCodeInvalidStateException {
					// If state is CREATE_FAILED, the service likely should be manually deleted,
					// and then re-created.

					// A substring check for "CREATE_FAILED" is our best effort to detect this
					if strings.Contains(aerr.Message(), "CREATE_FAILED") {
						// for visibility, here are examples for the  error fields:
						// aerr.Code() "InvalidStateException"
						// aerr.Error() "InvalidStateException: Service cannot be updated in the current state: CREATE_FAILED."
						// aerr.Message() "Service cannot be updated in the current state: CREATE_FAILED."
						step.Update("Service %q is in a failed state. Please delete the service and try running `waypoint up` again.", summary.Name)
						step.Abort()
					}
				}
			}

			return nil, err
		}

		// If no configuration is changed, no Operation is triggered,
		// and no OperationId is returned
		if uso.OperationId != nil {
			operationId = *uso.OperationId
		} else {
			log.Warn("No operationId was returned. This likely means no configuration change was detected.")
		}

		step.Done()

	} else {
		// If we didn't find a previous service, create it.
		step = sg.Add("Creating new App Runner service")
		log.Debug("creating new service...", "name", p.config.Name)
		log.Debug("using image", "image", img.Name)

		// Warning: App Runner will crash with a "exec format error"
		// when running Arm64 images.
		cso, err := arSvc.CreateService(&apprunner.CreateServiceInput{
			ServiceName: aws.String(p.config.Name),
			InstanceConfiguration: &apprunner.InstanceConfiguration{
				Cpu:    aws.String(strconv.Itoa(int(cpu))),
				Memory: aws.String(strconv.Itoa(int(mem))),
			},
			HealthCheckConfiguration: &apprunner.HealthCheckConfiguration{},
			SourceConfiguration: &apprunner.SourceConfiguration{
				AuthenticationConfiguration: &apprunner.AuthenticationConfiguration{
					AccessRoleArn: aws.String(roleArn),
				},
				ImageRepository: &apprunner.ImageRepository{
					ImageRepositoryType: aws.String(apprunner.ImageRepositoryTypeEcr),
					ImageIdentifier:     aws.String(img.Name()),
					ImageConfiguration: &apprunner.ImageConfiguration{
						Port:                        aws.String(strconv.Itoa(port)),
						RuntimeEnvironmentVariables: envVars,
					},
				},
				AutoDeploymentsEnabled: aws.Bool(false),
			},
		})
		if err != nil {
			log.Error("error creating service", "error", err)
			return nil, err
		}

		step.Done()

		operationId = *cso.OperationId
		serviceArn = *cso.Service.ServiceArn
		serviceUrl = *cso.Service.ServiceUrl
		serviceStatus = *cso.Service.Status
	}

	deployment.Region = p.config.Region
	deployment.ServiceName = p.config.Name
	deployment.ServiceArn = serviceArn
	deployment.ServiceUrl = serviceUrl
	deployment.Status = serviceStatus
	// possibly empty string when no configuration change is detected
	deployment.OperationId = operationId

	return deployment, nil
}

type ServiceSummary struct {
	Name   string
	Arn    string
	Url    string
	Status string
}

// a helper func to find an apprunner service by name
func (p *Platform) getServiceSummaryByName(
	sess *session.Session,
	log hclog.Logger,
	name string,
) (*ServiceSummary, error) {
	arSvc := apprunner.New(sess)
	lso, err := arSvc.ListServices(&apprunner.ListServicesInput{})
	if err != nil {
		log.Error("error listing services")
		return nil, err
	}

	log.Debug("found services", "service count", len(lso.ServiceSummaryList))

	var serviceSummary *ServiceSummary = nil
	for _, ss := range lso.ServiceSummaryList {
		if *ss.ServiceName == p.config.Name {

			serviceSummary = &ServiceSummary{
				Name:   *ss.ServiceName,
				Arn:    *ss.ServiceArn,
				Url:    *ss.ServiceUrl,
				Status: *ss.Status,
			}
			break
		}
	}

	if serviceSummary != nil {
		log.Debug("Found matching apprunner service", "name", serviceSummary.Name, "arn", serviceSummary.Arn)
	} else {
		log.Warn("No matching apprunner service was found")
	}

	return serviceSummary, nil
}

const defaultMemory = 2048
const defaultCpu = 1024
const defaultRoleName = `AppRunnerECRAccessRole`

const rolePolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "build.apprunner.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}`

// AWS-managed service policy
const AWSAppRunnerServicePolicyForECRAccess = `arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess`

func (p *Platform) getOrCreateIamRoleArn(
	sess *session.Session,
	log hclog.Logger,
	name string,
) (string, error) {
	iamSvc := iam.New(sess)
	ro, err := iamSvc.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException: // if "NoSuchEntity", create the role
				log.Info("CreateRole...")
				if _, err := iamSvc.CreateRole(&iam.CreateRoleInput{
					RoleName:                 aws.String(name),
					Path:                     aws.String("/service-role/"),
					AssumeRolePolicyDocument: aws.String(rolePolicy),
					Description:              aws.String("This role gives App Runner permission to access ECR"),
					// Tags: ,
				}); err != nil {
					log.Info("CreateRole ERROR")
					return "", err
				}
				log.Info("CreateRole OK")

				log.Info("AttachRolePolicy...")
				if _, err := iamSvc.AttachRolePolicy(&iam.AttachRolePolicyInput{
					RoleName:  aws.String(name),
					PolicyArn: aws.String(AWSAppRunnerServicePolicyForECRAccess),
				}); err != nil {
					log.Info("AttachRolePolicy ERROR")
					return "", err
				}
				log.Info("AttachRolePolicy OK")
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
	}
	return *ro.Role.Arn, nil
}

// For Resource manager
func (p *Platform) getSession(
	log hclog.Logger,
) (*session.Session, error) {
	return utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
}

// for a reference see builtin/aws/alb/release.go > resourceSecurityGroupCreate()'s args
func (p *Platform) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getSession),
		resource.WithResource(resource.NewResource(
			resource.WithName("apprunner_service_summary"),
			resource.WithPlatform("aws-apprunner"),
			resource.WithState(&ServiceSummary{}),
		)),
	)
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&PlatformConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Creates an App Runner service, with the specified configuration 
from ` + "`waypoint.hcl`.")

	doc.Example(`
deploy {
	use "aws-apprunner" {
		name = "go-gin"
		region = var.region
		port = 8080
		static_environment = {
			"PORT" = "8080"
		}
	}
}
`)

	doc.Input("ecr.Image")
	doc.Output("apprunner.Deployment")

	doc.SetField(
		"name",
		"the name of the service. This is required.",
	)

	doc.SetField(
		"port",
		"the port to listen on",
		docs.Default("8080 — This is AppRunner's default port"),
	)

	doc.SetField(
		"cpu",
		"the number of vCPU units to reserve for the service. Example: 1024 = 1 vCPU",
		docs.Default(fmt.Sprintf("%d", defaultCpu)),
	)

	doc.SetField(
		"memory",
		"the amount of memory to reserve for the service. Example: 2048 = 2 GB",
		docs.Default(fmt.Sprintf("%d", defaultMemory)),
	)

	doc.SetField(
		"role_name",
		"the name of the IAM role to use for the service",
		docs.Default(defaultRoleName),
	)

	doc.SetField(
		"static_environment",
		"environment variables to set for the service. A key-value map or string to string",
	)

	return doc, nil
}

var (
	_ component.Configurable = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
	_ component.Documented   = (*Platform)(nil)
)
