// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"strconv"

	"github.com/hashicorp/waypoint/internal/server/execclient"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (c *Project) Validate(ctx context.Context, op *pb.Job_ValidateOp) (*pb.Job_ValidateResult, error) {
	if op == nil {
		op = &pb.Job_ValidateOp{}
	}

	// Validate our job
	job := c.job()
	job.Operation = &pb.Job_Validate{
		Validate: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job, c.UI)
	if err != nil {
		return nil, err
	}

	return result.Validate, nil
}

func (c *Project) DestroyProject(ctx context.Context, op *pb.Job_DestroyProjectOp) (*pb.Job_ProjectDestroyResult, error) {
	if op == nil {
		op = &pb.Job_DestroyProjectOp{}
	}
	// Build our job
	job := c.job()
	job.Operation = &pb.Job_DestroyProject{
		DestroyProject: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job, c.UI)
	if err != nil {
		return nil, err
	}
	return result.ProjectDestroy, nil
}

func (c *App) Auth(ctx context.Context, op *pb.Job_AuthOp) (*pb.Job_AuthResult, error) {
	if op == nil {
		op = &pb.Job_AuthOp{}
	}

	// Auth our job
	job := c.job()
	job.Operation = &pb.Job_Auth{
		Auth: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.Auth, nil
}

func (c *App) Docs(ctx context.Context, op *pb.Job_DocsOp) (*pb.Job_DocsResult, error) {
	if op == nil {
		op = &pb.Job_DocsOp{}
	}

	job := c.job()
	job.Operation = &pb.Job_Docs{
		Docs: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.Docs, nil
}

func (c *App) Up(ctx context.Context, op *pb.Job_UpOp) (*pb.Job_Result, error) {
	if op == nil {
		op = &pb.Job_UpOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Up{
		Up: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	// Return the full result struct since Up populates multiple fields.
	return result, nil
}

func (c *App) Build(ctx context.Context, op *pb.Job_BuildOp) (*pb.Job_BuildResult, error) {
	if op == nil {
		op = &pb.Job_BuildOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Build{
		Build: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.Build, nil
}

func (c *App) PushBuild(ctx context.Context, op *pb.Job_PushOp) error {
	if op == nil {
		op = &pb.Job_PushOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Push{
		Push: op,
	}

	// Execute it
	_, err := c.doJob(ctx, job)
	return err
}

func (c *App) Deploy(ctx context.Context, op *pb.Job_DeployOp) (*pb.Job_DeployResult, error) {
	if op == nil {
		op = &pb.Job_DeployOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Deploy{
		Deploy: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.Deploy, nil
}

func (c *App) Destroy(ctx context.Context, op *pb.Job_DestroyOp) error {
	if op == nil {
		op = &pb.Job_DestroyOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Destroy{
		Destroy: op,
	}

	// Execute it
	_, err := c.doJob(ctx, job)
	return err
}

func (c *App) Exec(ctx context.Context, ec *execclient.Client) (exitCode int, err error) {

	// Depending on which deployments are at play, and which plugins those deployments
	// correspond to, we may need a local runner. It'll be up to the server to actually
	// create the job, but we'll need to create the local runner if necessary, error
	// if vcs is dirty, etc.
	_, ctx, err = c.project.setupLocalJobSystem(ctx)
	if err != nil {
		return 0, err
	}

	ec.Context = ctx
	return ec.Run()
}

func (c *App) Release(ctx context.Context, op *pb.Job_ReleaseOp) (*pb.Job_ReleaseResult, error) {
	if op == nil {
		op = &pb.Job_ReleaseOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Release{
		Release: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.Release, nil
}

func (a *App) Logs(ctx context.Context, deploySeq string) (pb.Waypoint_GetLogStreamClient, error) {
	log := a.project.logger.Named("logs")

	// Depending on which deployments are at play, and which plugins those deployments
	// correspond to, we may need a local runner. It'll be up to the server to actually
	// create the job, but we'll need to create the local runner if necessary, error
	// if vcs is dirty, etc.
	_, ctx, err := a.project.setupLocalJobSystem(ctx)
	if err != nil {
		return nil, err
	}
	var logStreamRequest *pb.GetLogStreamRequest

	if deploySeq != "" {
		i, err := strconv.ParseUint(deploySeq, 10, 64)
		if err != nil {
			return nil, err
		}
		deploy, err := a.project.client.GetDeployment(ctx, &pb.GetDeploymentRequest{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Sequence{
					Sequence: &pb.Ref_OperationSeq{
						Application: a.Ref(),
						Number:      i,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		logStreamRequest = &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_DeploymentId{
				DeploymentId: deploy.Id,
			},
		}
	} else {
		logStreamRequest = &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_Application_{
				Application: &pb.GetLogStreamRequest_Application{
					Application: a.Ref(),
					Workspace:   a.project.WorkspaceRef(),
				},
			}}
	}

	// First we attempt to query the server for logs for this deployment.
	log.Info("requesting log stream")
	client, err := a.project.client.GetLogStream(ctx, logStreamRequest)
	if err != nil {
		return nil, err
	}

	// Build our log viewer
	return client, nil
}

func (c *App) ConfigSync(ctx context.Context, op *pb.Job_ConfigSyncOp) (*pb.Job_Result, error) {
	if op == nil {
		op = &pb.Job_ConfigSyncOp{}
	}

	job := c.job()
	job.Operation = &pb.Job_ConfigSync{
		ConfigSync: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *App) StatusReport(ctx context.Context, op *pb.Job_StatusReportOp) (*pb.Job_StatusReportResult, error) {
	if op == nil {
		op = &pb.Job_StatusReportOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_StatusReport{
		StatusReport: op,
	}

	// Execute it
	result, err := c.doJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return result.StatusReport, nil
}
