package apprunner

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ component.Destroyer = (*Platform)(nil)
)

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return p.Destroy
}

func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	deployment *Deployment,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Deleting Service")
	// We put this in a function because if/when step is reassigned, we want to
	// abort the new value.
	defer func() {
		step.Abort()
	}()

	// Delete app runner by service ARN

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		step.Update(fmt.Sprintf("Failed to get session: %s", err))
		return err
	}

	arSvc := apprunner.New(sess)

	step.Update("Deleting service: %s", deployment.ServiceArn)
	dso, err := arSvc.DeleteService(&apprunner.DeleteServiceInput{
		ServiceArn: aws.String(deployment.ServiceArn),
	})
	step.Done()

	step = sg.Add("App Runner::Waiting for Delete Service to succeed...")
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)

	opId := *dso.OperationId

	shouldRetry := true

	for shouldRetry {
		loo, err := arSvc.ListOperations(&apprunner.ListOperationsInput{
			ServiceArn: &deployment.ServiceArn,
		})

		// TODO(kevinwang): better error reporting
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				return aerr
			}
			return err
		}

		for _, os := range loo.OperationSummaryList {
			// find operation by id from DeleteService request
			if *os.Id == opId {
				switch *os.Status {
				case apprunner.OperationStatusSucceeded:
					// OK — resume
					step.Update("OK!")
					shouldRetry = false
				case apprunner.OperationStatusFailed:
					// Failed — exit
					step.Update("Failed...")
					return status.Error(codes.FailedPrecondition, fmt.Sprintf("App Runner responded with status: %s", *os.Status))
				case apprunner.OperationStatusInProgress:
					select {
					case <-ticker.C: // retry
					case <-ctx.Done(): // abort
						step.Update("Timeout...")
						return status.Errorf(codes.Aborted, fmt.Sprintf("Context cancelled from timeout when waiting for App Runner to graduate from %s", *os.Status))
					}
				default:
					log.Warn("Unexpected status: %s", *os.Status)
				}
			}
		}
	}

	step.Update("Deleted App Runner service: %s", deployment.ServiceName)
	return nil
}
