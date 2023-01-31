package apprunner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
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

	sess, err := utils.GetSession(&utils.SessionConfig{
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	arSvc := apprunner.New(sess)

	step := sg.Add("Waiting for App Runner service to be deployed: %q", dep.ServiceName)
	defer func() {
		step.Abort()
	}()
	step.Done()

	// wait for operation status to become `SUCCEEDED`
	operationId := dep.OperationId

	if operationId == "" {
		// Likely an attempt was made to update the app runner
		// service, but there was no configuration change to be
		// made.

		// do nothing as there is no operation to wait for
		step = sg.Add("No configuration change was made. Moving on...")
		step.Done()
	} else {
		step = sg.Add("App Runner is in state: ...")
		// Wait for X-minutes
		d := time.Now().Add(DEFAULT_TIMEOUT)
		ctx, cancel := context.WithDeadline(ctx, d)
		defer cancel()

		// Poll every 10 seconds
		ticker := time.NewTicker(10 * time.Second)
		shouldRetry := true

		now := time.Now()
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

			step.Update("App Runner is in state: %s. Time elapsed: %s", *opsum.Status, time.Since(now).Round(time.Second))

			switch *opsum.Status {
			case apprunner.OperationStatusSucceeded:
				// OK — ready to proceed
				shouldRetry = false
			case apprunner.OperationStatusFailed:
				// Failed — exit
				return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("App Runner responded with status: %s", *opsum.Status))
			default:
				select {
				case <-ticker.C: // retry
				case <-ctx.Done(): // abort
					return nil, status.Errorf(codes.Aborted, fmt.Sprintf("Context cancelled from timeout when waiting for App Runner graduate from %s", *opsum.Status))
				}
			}
		}
		step.Done()
	}

	step = sg.Add("App Runner service is ready!")
	step.Done()

	return &Release{
		Url:         "https://" + dep.ServiceUrl,
		ServiceArn:  dep.ServiceArn,
		ServiceName: dep.ServiceName,
		Region:      dep.Region,
	}, nil
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

const DEFAULT_TIMEOUT = time.Minute * time.Duration(10)

var (
	_ component.Configurable   = (*Releaser)(nil)
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
)
