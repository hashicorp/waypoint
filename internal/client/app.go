// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// App is used for application-specific operations.
type App struct {
	UI terminal.UI

	project     *Project
	application *pb.Ref_Application
	runner      *configpkg.Runner
}

// App returns the app-specific operations client.
func (c *Project) App(n string) *App {
	app := &App{
		UI:      c.UI,
		project: c,
		application: &pb.Ref_Application{
			Project:     c.project.Project,
			Application: n,
		},
	}
	if c.waypointHCL != nil {
		app.runner = c.waypointHCL.ConfigAppRunner(n)
	}

	return app
}

// Ref returns the application reference that this client is using.
func (c *App) Ref() *pb.Ref_Application {
	return c.application
}

// job is the same as Project.job except this also sets the application
// reference.
func (c *App) job() *pb.Job {
	job := c.project.job()
	job.Application = c.application
	if c.runner != nil && c.runner.Profile != "" {
		job.OndemandRunner = &pb.Ref_OnDemandRunnerConfig{
			Name: c.runner.Profile,
		}
	}
	return job
}

// doJob is the same as Project.doJob except we set the proper app-specific UI.
func (c *App) doJob(ctx context.Context, job *pb.Job) (*pb.Job_Result, error) {
	return c.project.doJob(ctx, job, c.UI)
}

// doJob is the same as Project.doJob except we set the proper app-specific UI and can
// monitor the job status.
func (c *App) doJobMonitored(ctx context.Context, job *pb.Job, monCh chan pb.Job_State) (*pb.Job_Result, error) {
	return c.project.doJobMonitored(ctx, job, c.UI, monCh)
}
