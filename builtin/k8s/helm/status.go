// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"helm.sh/helm/v3/pkg/release"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/builtin/k8s/internal/k8sstatus"
	"github.com/hashicorp/waypoint/builtin/k8s/internal/manifest"
)

func (p *Platform) StatusFunc() interface{} {
	return p.Status
}

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	actionConfig, err := p.actionInit(log)
	if err != nil {
		return nil, err
	}

	rel, err := getRelease(actionConfig, deployment.Release)
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, status.Errorf(codes.NotFound, "Helm release not found")
	}

	// parse the manifest and turn that into a status report
	m, err := manifest.Parse(strings.NewReader(rel.Manifest))
	if err != nil {
		return nil, err
	}

	// Get our K8S API
	cs, _, _, err := k8s.Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	// Build our report
	report, err := k8sstatus.FromManifest(ctx, cs, m)
	if err != nil {
		return nil, err
	}

	// The health is PROBABLY unknown. If it is, then we try to set some
	// basic health based on the status of the chart.
	if report.Health == sdk.StatusReport_UNKNOWN {
		switch rel.Info.Status {
		case release.StatusDeployed:
			// If it is deployed, just consider it ready.
			report.Health = sdk.StatusReport_READY

		case release.StatusFailed:
			// If we failed the helm install, we probably won't have a
			// status report but mark it as unhealthy.
			report.Health = sdk.StatusReport_DOWN

		case release.StatusPendingInstall,
			release.StatusPendingUpgrade,
			release.StatusPendingRollback,
			release.StatusUninstalling:
			// Pending or active operations just get a "partial" status.
			report.Health = sdk.StatusReport_PARTIAL
		}

		report.HealthMessage = fmt.Sprintf(
			"%s: %s", rel.Info.Status, rel.Info.Description)
	}

	// For our resources, set a created time to the FirstDeployed of
	// this release if there is none. This lets us have SOMETHING.
	for _, r := range report.Resources {
		if r.CreatedTime == nil {
			r.CreatedTime = timestamppb.New(rel.Info.FirstDeployed.Time)
		}
	}

	return report, nil
}
