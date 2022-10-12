package apprunner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Releaser struct {
	// ReleaseFunc should return the method handle for the "release" operation.
	config ReleaserConfig
}

type ReleaserConfig struct {
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	dep *Deployment,
) (*Release, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	// Create our resource manager and create
	rm := r.resourceManager(log)
	if err := rm.CreateAll(
		ctx, log, sg, ui, src,
		dep,
	); err != nil {
		log.Info("Error creating resources", "error", err)
		return nil, err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	arSvc := apprunner.New(sess)

	// wait for operation status to become `SUCCEEDED`

	operationId := dep.OperationId

	if operationId == "" {
		// Likely an attempt was made to update the app runner
		// service, but there was no configuration change to be
		// made.

		// do nothing as there is no operation to wait for
	} else {
		// Wait for 5 minutes
		d := time.Now().Add(time.Minute * time.Duration(5))
		ctx, cancel := context.WithDeadline(ctx, d)
		defer cancel()

		// Poll every 10 seconds
		ticker := time.NewTicker(10 * time.Second)
		shouldRetry := true
		for shouldRetry {
			loo, err := arSvc.ListOperations(&apprunner.ListOperationsInput{
				ServiceArn: &dep.ServiceArn,
			})

			// TODO(kevinwang): better error handling/reporting
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}

			var opsum *apprunner.OperationSummary = nil
			// Find our operation by id
			for _, os := range loo.OperationSummaryList {
				// find operation by id from CreateService request
				if *os.Id == operationId {

					opsum = os
					break
				}
			}

			switch *opsum.Status {
			case apprunner.OperationStatusSucceeded:
				// OK — resume
				shouldRetry = false
			case apprunner.OperationStatusFailed:
				// Failed — exit
				return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("App Runner responded with status: %s", *opsum.Status))
			// case apprunner.OperationStatusInProgress:
			default:
				select {
				case <-ticker.C: // retry
				case <-ctx.Done(): // abort
					return nil, status.Errorf(codes.Aborted, fmt.Sprintf("Context cancelled from timeout when waiting for App Runner graduate from %s", *opsum.Status))
				}
			}
		}
	}

	return &Release{
		Url:         "https://" + dep.ServiceUrl,
		ServiceArn:  dep.ServiceArn,
		ServiceName: dep.ServiceName,
	}, nil
}

func (r *Releaser) getSession(
	_ context.Context,
	log hclog.Logger,
	dep *Deployment,
) (*session.Session, error) {
	return utils.GetSession(&utils.SessionConfig{
		Region: dep.Region,
		Logger: log,
	})
}

func (r *Releaser) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(r.getSession),

		// resource.WithResource(resource.NewResource(
		// 	resource.WithName("function_permission"),
		// 	resource.WithCreate(r.resourceFunctionPermissionCreate),
		// )),

		// resource.WithResource(resource.NewResource(
		// 	resource.WithName("function_url"),
		// 	resource.WithState(&Resource_FunctionUrl{}),
		// 	resource.WithCreate(r.resourceFunctionUrlCreate),
		// )),
	)
}

func (p *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&PlatformConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description(`
This releaser is a lightweight layer that polls App Runner until a deployment
reaches the ` + "`SUCCEEDED` status." + `

~> **Note:** Using the ` + "`-prune=false`" + ` flag is recommended for this releaser. By default,
Waypoint prunes and destroys all unreleased deployments and keeps only one previous
deployment. Therefore, if ` + "`-prune=false`" + ` is not set, Waypoint will delete the single
service that App Runner manages upon a ` + "second `waypoint up`." + `

See [deployment pruning](/docs/lifecycle/release#deployment-pruning) for more information.
`)

	return doc, nil
}

var (
	_ component.Configurable   = (*Releaser)(nil)
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
)
