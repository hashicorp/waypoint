package apprunner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	if roleArn == "" {
		// TODO(kevinwang): source this from configuration
		roleArn = "arn:aws:iam::000000000000:role/service-role/AppRunnerECRAccessRole"
	}

	mem := int64(p.config.Memory)
	if mem == 0 {
		mem = 2048
	}

	storage := int64(p.config.Cpu)
	if storage == 0 {
		storage = 1024
	}

	step.Done()

	step = sg.Add("App Runner: %s", src.App)
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

	// If we found a previous service, update it.
	if foundServiceSummary != nil {
		// TODO(kevinwang): implement me
		log.Debug("found existing service...", "name", p.config.Name)

		// uso, err :=
		arSvc.UpdateService(&apprunner.UpdateServiceInput{
			ServiceArn: foundServiceSummary.ServiceArn,
			// TODO(kevinwang): update configs
		})

	} else {
		step = sg.Add("App Runner::CreateService")
		// If we didn't find a previous service, create it.
		log.Debug("creating new service...", "name", p.config.Name)
		log.Debug("using image", "image", img.Name)

		cso, err := arSvc.CreateService(&apprunner.CreateServiceInput{
			ServiceName: aws.String(p.config.Name),
			SourceConfiguration: &apprunner.SourceConfiguration{
				AuthenticationConfiguration: &apprunner.AuthenticationConfiguration{
					AccessRoleArn: aws.String(roleArn),
				},
				ImageRepository: &apprunner.ImageRepository{
					ImageRepositoryType: aws.String(apprunner.ImageRepositoryTypeEcr),
					ImageIdentifier:     aws.String(img.Image + ":" + img.Tag),
				},
				AutoDeploymentsEnabled: aws.Bool(false),
			},
		})
		if err != nil {
			log.Error("error creating service", "error", err)
			return nil, err
		}

		step.Done()

		serviceArn = *cso.Service.ServiceArn
		serviceUrl = *cso.Service.ServiceUrl
		serviceStatus = *cso.Service.Status

		// wait for operation status to become `SUCCEEDED`

		step = sg.Add("App Runner::Waiting for Create Service to succeed...")
		d := time.Now().Add(time.Minute * time.Duration(5))
		ctx, cancel := context.WithDeadline(ctx, d)
		defer cancel()
		ticker := time.NewTicker(5 * time.Second)

		opId := *cso.OperationId

		shouldRetry := true
		for shouldRetry {
			loo, err := arSvc.ListOperations(&apprunner.ListOperationsInput{
				ServiceArn: &serviceArn,
			})

			// TODO(kevinwang): better error handling/reporting
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}

			// var foundOperationSummary *apprunner.OperationSummary = nil
			for _, os := range loo.OperationSummaryList {

				// find operation by id from CreateService request
				if *os.Id == opId {
					// update state
					serviceStatus = *os.Status

					switch *os.Status {
					case apprunner.OperationStatusSucceeded:
						// OK — resume
						step.Update("OK!")
						shouldRetry = false
					case apprunner.OperationStatusFailed:
						// Failed — exit
						step.Update("Failed...")
						return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("App Runner responded with status: %s", *os.Status))
					case apprunner.OperationStatusInProgress:
						select {
						case <-ticker.C: // retry
						case <-ctx.Done(): // abort
							step.Update("Timeout...")
							return nil, status.Errorf(codes.Aborted, fmt.Sprintf("Context cancelled from timeout when waiting for App Runner graduate from %s", *os.Status))
						}
					default:
						log.Warn("Unexpected status: %s", *os.Status)
					}
				}
			}
		}
		step.Done()

		serviceArn = *cso.Service.ServiceArn
		serviceUrl = *cso.Service.ServiceUrl
		serviceStatus = *cso.Service.Status
	}

	step.Done()

	deployment.Id = id
	deployment.Region = p.config.Region
	deployment.ServiceName = p.config.Name
	deployment.ServiceArn = serviceArn
	deployment.ServiceUrl = serviceUrl
	deployment.Status = serviceStatus

	return deployment, nil
}

var (
	_ component.Configurable = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
)
