package client

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

func (c *App) Exec(ctx context.Context, op *pb.Job_ExecOp, mon chan pb.Job_State) error {
	if op == nil {
		op = &pb.Job_ExecOp{}
	}

	// Build our job
	job := c.job()
	job.Operation = &pb.Job_Exec{
		Exec: op,
	}

	// Execute it
	_, err := c.doJobMonitored(ctx, job, mon)
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

func (a *App) Logs(ctx context.Context) (pb.Waypoint_GetLogStreamClient, error) {
	log := a.project.logger.Named("logs")

	// First we attempt to query the server for logs for this deployment.
	log.Info("requesting log stream")
	client, err := a.project.client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_Application_{
			Application: &pb.GetLogStreamRequest_Application{
				Application: a.Ref(),
				Workspace:   a.project.WorkspaceRef(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Build our log viewer
	return client, nil
}

func (c *App) ConfigSync(ctx context.Context, op *pb.Job_ConfigSyncOp) (*pb.Job_ConfigSyncResult, error) {
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

	return result.ConfigSync, nil
}
