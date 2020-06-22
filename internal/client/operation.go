package client

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logviewer"
	"github.com/hashicorp/waypoint/sdk/component"
)

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

func (c *App) DestroyDeploy(ctx context.Context, op *pb.Job_DestroyDeployOp) error {
	if op == nil {
		op = &pb.Job_DestroyDeployOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_DestroyDeploy{
		DestroyDeploy: op,
	}

	// Execute it
	_, err := c.doJob(ctx, job)
	return err
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

func (a *App) Logs(ctx context.Context, d *pb.Deployment) (component.LogViewer, error) {
	log := a.project.logger.Named("logs")

	// First we attempt to query the server for logs for this deployment.
	log.Info("requesting log stream", "deployment_id", d.Id)
	client, err := a.project.client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		DeploymentId: d.Id,
	})
	if err != nil {
		return nil, err
	}

	// Build our log viewer
	return &logviewer.Viewer{Stream: client}, nil
}
