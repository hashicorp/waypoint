package apprunner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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

	s := sg.Add("Gathering health report for app runner: %q", release.Url)
	defer s.Done()

	// TODO
	report := sdk.StatusReport{}
	report.External = true
	report.Health = sdk.StatusReport_READY
	report.HealthMessage = "IMPLEMENT ME"
	report.Resources = []*sdk.StatusReport_Resource{}

	return &report, nil
}

var _ component.Status = (*Releaser)(nil)
