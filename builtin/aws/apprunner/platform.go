package apprunner

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apprunner"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

type Platform struct {
	config PlatformConfig
}

type PlatformConfig struct {
	Region  string `hcl:"region,optional"`
	RoleArn string `hcl:"role_arn,optional"`
	// InstanceConfiguration map[string]int `hcl:"source_configuration,optional"`

	Name   string `hcl:"name"`
	Memory int    `hcl:"memory,optional"`
	Cpu    int    `hcl:"cpu,optional"`
	// Once created, Port cannot be modified
	Port int `hcl:"port,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	SourceConfiguration map[string]string `hcl:"source_configuration,optional"`
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// Deploy deploys an image to AWS Lambda.
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

	// We put this in a function because if/when step is reassigned, we want to
	// abort the new value.
	defer func() {
		step.Abort()
	}()

	// Start building our deployment since we use this information
	deployment := &Deployment{}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	roleArn := p.config.RoleArn
	port := p.config.Port
	if port == 0 {
		// App Runner default port
		port = 8080
	}

	if roleArn == "" {
		// if role arn is not specified, we should create a service role,
		// similar to how the AWS console
		roleArn = "arn:aws:iam::111122223333:role/service-role/AppRunnerECRAccessRole"
	}

	mem := int64(p.config.Memory)
	if mem == 0 {
		mem = 2048
	}

	cpu := int64(p.config.Cpu)
	if cpu == 0 {
		cpu = 1024
	}

	step.Done()

	step = sg.Add("App Runner::%s", src.App)
	step.Done()

	arSvc := apprunner.New(sess)

	// The operator's waypoint.hcl will realistically only have a svc name, not ARN,
	// since that is provided by AWS.
	//
	// While `DescribeService` only supports a service ARN, our best effort approach
	// is to list all services and manually match by user-provided name.
	//
	// - If we find a match, we should update
	// - If no match, we should create

	// TODO(kevinwang): pagination
	step = sg.Add("App Runner::ListServices")
	lso, err := arSvc.ListServices(&apprunner.ListServicesInput{
		// max is 20
		MaxResults: aws.Int64(20),
	})
	if err != nil {
		log.Error("error listing services")

		// TODO(kevinwang): specific error handling
		if aerr, ok := err.(awserr.Error); ok {
			log.Error("AWS Error", "code", aerr.Code(), "message", aerr.Message())
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		}
		log.Error("Non-AWS Error", "error", err)
		return nil, err
	}
	step.Done()

	log.Info("found services", "service count", len(lso.ServiceSummaryList))

	// Did we find a previous apprunner service by name?
	var foundServiceSummary *apprunner.ServiceSummary = nil

	serviceArn := ""
	serviceUrl := ""
	serviceStatus := ""

	for _, ss := range lso.ServiceSummaryList {
		if *ss.ServiceName == p.config.Name {
			foundServiceSummary = ss
			serviceArn = *ss.ServiceArn
			serviceUrl = *ss.ServiceUrl
			serviceStatus = *ss.Status
			break
		}
	}

	operationId := ""

	// If we found a previous service, update it.
	if foundServiceSummary != nil {
		// TODO(kevinwang): implement me
		log.Debug("found existing service...", "name", p.config.Name)

		step = sg.Add("App Runner::UpdateService %s", *foundServiceSummary.ServiceName)

		uso, err := arSvc.UpdateService(&apprunner.UpdateServiceInput{
			ServiceArn: foundServiceSummary.ServiceArn,
			InstanceConfiguration: &apprunner.InstanceConfiguration{
				Cpu:    aws.String(strconv.FormatInt(cpu, 10)),
				Memory: aws.String(strconv.FormatInt(mem, 10)),
			},
		})

		if err != nil {
			step.Update("App Runner::UpdateService Failed: %s", err)
			return nil, err
		}

		// If no configuration is changed, no Operation is triggered,
		// and no OperationId is returned
		if uso.OperationId != nil {
			operationId = *uso.OperationId
		} else {
			// step.Update("uso.OperationId was null %+v", uso)
			// TODO(Kevinwang): handle this; Maybe fail?
		}

		step.Done()

	} else {
		step = sg.Add("App Runner::CreateService")
		// If we didn't find a previous service, create it.
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
						Port: aws.String(strconv.Itoa(port)),
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

	step = sg.Add("Assigning Deployment Values")
	deployment.Id = id
	deployment.Region = p.config.Region
	deployment.ServiceName = p.config.Name
	deployment.ServiceArn = serviceArn
	deployment.ServiceUrl = serviceUrl
	deployment.Status = serviceStatus
	// possibly empty string when no configuration change is detected
	deployment.OperationId = operationId
	step.Done()

	return deployment, nil
}

var (
	_ component.Configurable = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
)

// {
// Service: {
//   AutoScalingConfigurationSummary: {
//     AutoScalingConfigurationArn: "arn:aws:apprunner:us-east-1:058050752201:autoscalingconfiguration/DefaultConfiguration/1/00000000000000000000000000000001"
//     AutoScalingConfigurationName: "DefaultConfiguration",
// 		AutoScalingConfigurationRevision: 1
// 	},
// 	CreatedAt: 2022-09-29 06:08:12 +0000 UTC,
// 	HealthCheckConfiguration: {
// 		HealthyThreshold: 1,
// 		Interval: 5,
// 		Path: "/",
// 		 Protocol: "TCP",
// 		 Timeout: 2,
// 		 UnhealthyThreshold: 5
// 		},
// 		InstanceConfiguration: {
// 			Cpu: "1024",
// 			Memory: "2048"
// 		},
// 		NetworkConfiguration: {
// 			EgressConfiguration: {
// 				EgressType: "DEFAULT"
// 			}
// 		},
// 		ServiceArn: "arn:aws:apprunner:us-east-1:058050752201:service/my-apprunner-service3/4992ba7a6832496e9289bce260d65b8a",
// 		ServiceId: "4992ba7a6832496e9289bce260d65b8a",
// 		ServiceName: "my-apprunner-service3",
// 		ServiceUrl: "vi4rav3dvx.us-east-1.awsapprunner.com",
// 		SourceConfiguration: {
// 			AuthenticationConfiguration: {
// 				AccessRoleArn: "arn:aws:iam::058050752201:role/service-role/AppRunnerECRAccessRole"
// 			},
// 			AutoDeploymentsEnabled: false,
// 			ImageRepository: {
// 				ImageConfiguration: {
// 					Port: "7531",
// 					RuntimeEnvironmentVariables: {
// 						Test: "1"
// 					}
// 				},
// 				ImageIdentifier: "058050752201.dkr.ecr.us-east-1.amazonaws.com/python-fastapi:latest",
// 				ImageRepositoryType: "ECR"
// 			}
// 		},
// 		Status: "RUNNING",
// 		UpdatedAt: 2022-09-29 06:08:12 +0000 UTC
// 	}
// }
