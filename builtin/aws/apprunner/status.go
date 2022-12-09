package apprunner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusFunc implements component.Status
func (r *Releaser) StatusFunc() interface{} {
	return r.Status
}

func (r *Releaser) Status(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Gathering health report for app runner: %q", release.Url)
	defer step.Abort()

	report := sdk.StatusReport{}
	report.External = true

	if release.Region == "" {
		log.Debug("Region is not available for this release. Unable to determine status.")
		return &report, nil
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: release.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	arSvc := apprunner.New(sess)

	dso, err := arSvc.DescribeService(&apprunner.DescribeServiceInput{
		ServiceArn: &release.ServiceArn,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to describe service %s: %s", release.ServiceName, err)
	}

	report.Health = sdk.StatusReport_READY
	report.HealthMessage = fmt.Sprintf("Service status is: %q", *dso.Service.Status)
	report.Resources = []*sdk.StatusReport_Resource{}
	report.Resources = append(report.Resources, &sdk.StatusReport_Resource{
		Health:        sdk.StatusReport_READY,
		HealthMessage: *dso.Service.Status,
		Name:          *dso.Service.ServiceArn,
	})

	step.Update("Finished building report for AWS/AppRunner platform")
	step.Done()

	return &report, nil
}

var _ component.Status = (*Releaser)(nil)
