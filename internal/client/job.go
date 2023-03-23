// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/jobstream"
	"github.com/hashicorp/waypoint/internal/pkg/gitdirty"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/grpcmetadata"
)

// job returns the basic job skeleton pre-populated with the correct
// defaults based on how the client is configured. For example, for local
// operations, this will already have the targeting for the local runner.
func (c *Project) job() *pb.Job {
	job := &pb.Job{
		TargetRunner: c.runner,
		Labels:       c.labels,
		Variables:    c.variables,
		Workspace:    c.workspace,
		Application: &pb.Ref_Application{
			Project: c.project.Project,
		},
		DataSourceOverrides: c.dataSourceOverrides,

		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},
	}

	return job
}

// DoJobDangerously executes the given job and returns the result. The
// "Dangerously" suffix is because this function isn't meant to be generally
// used; it is dangerous because it doesn't perform many validation steps.
// In almost all cases, callers should use a more focused function such as
// Build or Deploy, or write a new one.
func (c *Project) DoJobDangerously(ctx context.Context, job *pb.Job) (*pb.Job_Result, error) {
	return c.doJob(ctx, job, c.UI)
}

// doJob will queue and execute the job, and target the proper runner.
func (c *Project) doJob(ctx context.Context, job *pb.Job, ui terminal.UI) (*pb.Job_Result, error) {
	return c.doJobMonitored(ctx, job, ui, nil)
}

// setupLocalJobSystem does the pre-work required to run jobs locally.
// This includes:
//   - figure out if jobs should be executed locally or remotely.
//   - if job should be executed locally, start a local runner
//   - if jobs will be executed remotely, but local VCS is present and
//     dirty, warn.
//
// This lives separately from DoJob because the logs and exec commands
// need to conditionally warm up the local job infrastructure, but
// don't actually create a job (the server does).
func (c *Project) setupLocalJobSystem(ctx context.Context) (isLocal bool, newCtx context.Context, err error) {
	log := c.logger.Named("setupLocalJobSystem")
	defer func() {
		log.Debug("result", "isLocal", isLocal)
	}()

	// Automatically determine if we should use a local or a remote
	// runner
	newCtx = ctx

	// We may use this in multiple places, so we save the result if we
	// obtain it
	// NOTE(izaak): If in the future we need the full project in other
	// places in this codepath, we should probably cache it early on
	// the parent struct.
	var project *pb.Project

	// A nil useLocalRunner means the option was not set explicitly
	// when this client was created. We'll decide a value for it here,
	// and set it for future runs.
	if c.useLocalRunner == nil {
		log.Debug("determining if a local or remote runner should be used for this and future jobs")

		getProjectResp, err := c.client.GetProject(ctx, &pb.GetProjectRequest{Project: c.project})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return false, newCtx, fmt.Errorf("Project %q was not found! Please ensure that 'waypoint init' was run with this project.", c.project.Project)
			} else {
				return false, newCtx, errors.Wrapf(err, "failed to get project %s", c.project.Project)
			}
		}
		project = getProjectResp.Project

		var runnerCfgs []*configpkg.Runner
		// Note(XX): temp (?) workaround the issue where runner is only upserted to profile on the first `waypoint init`
		if c.waypointHCL != nil {
			runnerCfgs = append(runnerCfgs, c.waypointHCL.ConfigRunner())
			for _, app := range project.Applications {
				runnerCfgs = append(runnerCfgs, c.waypointHCL.ConfigAppRunner(app.Name))
			}
		}

		remotePreferred, err := remoteOpPreferred(ctx, c.client, project, runnerCfgs, c.logger)
		if err != nil {
			return false, newCtx, errors.Wrapf(err, "failed to determine if job should run locally or remotely")
		}

		// Store this for later operations on this same project
		useLocalRunner := !remotePreferred
		c.useLocalRunner = &useLocalRunner
	}

	if *c.useLocalRunner {
		if c.activeRunner == nil {
			// we need a local runner and we haven't started it yet
			if err := c.startRunner(ctx); err != nil {
				return false, newCtx, errors.Wrapf(err, "failed to start local runner for job %s", err)
			}
		}
		// Inject the metadata about the client, such as the runner id
		// if it is running a local runner.
		newCtx = grpcmetadata.AddRunner(ctx, c.activeRunner.Id())
	} else {
		// We're about to run a remote op. We should check if we have
		// a dirty local vcs, because the user may expect their local
		// changes to be reflected in the remote op execution, and
		// they won't.
		gitDirtyErr := func() error {
			// Running this inside of an anonymous func so that we can
			// return early
			if c.configPath == "" {
				// No local project dir, so nothing is dirty!
				return nil
			}

			if !gitdirty.GitInstalled() {
				return errors.New("git is not installed - unable to check if local git directory is dirty for warning purposes")
			}

			repoRoot, err := gitdirty.RepoTopLevelPath(log, c.configPath)
			if err != nil {
				return errors.Wrapf(err, "failed to find the top level of the repository that contains %q", c.configPath)
			}

			// Get the project if we haven't already
			if project == nil {
				getProjectResp, err := c.client.GetProject(ctx, &pb.GetProjectRequest{Project: c.project})
				if err != nil {
					if status.Code(err) == codes.NotFound {
						return fmt.Errorf("project %q was not found! Please ensure that 'waypoint init' was run with this project.", c.project.Project)
					} else {
						return errors.Wrapf(err, "failed to get project %s", c.project.Project)
					}
				}
				project = getProjectResp.Project
			}

			if project.DataSource == nil || project.DataSource.Source == nil {
				return fmt.Errorf("no valid data source configured for Project %q", c.project.Project)
			}

			gitDs, ok := project.DataSource.Source.(*pb.Job_DataSource_Git)

			if !ok {
				// The remote op will likely fail anyway, because it
				// needs a remote-capable datasource.
				log.Debug("local config directory is a git repo, but project has non-remote datasource type. Will not attempt dirty git warning.",
					"project datasource type", fmt.Sprintf("%t", project.DataSource.Source),
				)
				return nil
			}

			var dirty bool
			if gitDs.Git.Path != "" {
				diffPath := filepath.Join(repoRoot, gitDs.Git.Path)
				dirty, err = gitdirty.PathIsDirty(log, repoRoot, gitDs.Git.Url, gitDs.Git.Ref, diffPath)
				if err != nil {
					return errors.Wrapf(err, "failed to diff repo at %q subpath %q against remote with url %q ref %q",
						repoRoot, diffPath, gitDs.Git.Url, gitDs.Git.Ref,
					)
				}
			} else {
				dirty, err = gitdirty.RepoIsDirty(log, repoRoot, gitDs.Git.Url, gitDs.Git.Ref)
				return errors.Wrapf(err, "failed to diff repo at %q against remote with url %q ref %q",
					repoRoot, gitDs.Git.Url, gitDs.Git.Ref,
				)
			}
			if dirty {
				c.UI.Output(warnGitDirty, terminal.WithWarningStyle())
			}

			return nil
		}()
		if gitDirtyErr != nil {
			log.Warn("failed to determine if local vcs is dirty", "err", gitDirtyErr)
		}
	}
	return *c.useLocalRunner, newCtx, nil
}

// Same as doJob, but with the addition of a mon channel that can be
// used to monitor the job status as it changes.
// The receiver must be careful to not block sending to mon as it will
// block the job state processing loop.
func (c *Project) doJobMonitored(ctx context.Context, job *pb.Job, ui terminal.UI, monCh chan pb.Job_State) (*pb.Job_Result, error) {
	isLocal, ctx, err := c.setupLocalJobSystem(ctx)
	if err != nil {
		return nil, err
	}

	// Be sure that the monitor is closed so the receiver knows for
	// sure the job isn't going anymore.
	if monCh != nil {
		defer close(monCh)
	}

	// In local mode we have to target the local runner.
	if isLocal {
		// If we're local, we set a local data source. Otherwise, it
		// defaults to whatever the project has remotely.
		job.DataSource = &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Local{
				Local: &pb.Job_Local{},
			},
		}

		// Modify the job to target this runner and use the local data
		// source. The runner will have been started when we created
		// the Project value and be used for all local jobs.
		job.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: c.activeRunner.Id(),
				},
			},
		}
	} else if c.waypointHCL != nil {
		var configRunner *configpkg.Runner
		// Find runner configuration on the app
		if job.Application != nil {
			configRunner = c.waypointHCL.ConfigAppRunner(job.Application.Application)
		}
		// If not on app, try to find it on the project
		if configRunner == nil {
			configRunner = c.waypointHCL.ConfigRunner()
		}
		// If runner config is found, assign to job
		if configRunner != nil {
			if configRunner.Profile != "" {
				job.OndemandRunner = &pb.Ref_OnDemandRunnerConfig{
					Name: configRunner.Profile,
				}
			}
		}
	}

	return c.queueAndStreamJob(ctx, job, ui, monCh, isLocal)
}

// queueAndStreamJob will queue the job. If the client is configured to watch the job,
// it'll also stream the output to the configured UI.
func (c *Project) queueAndStreamJob(
	ctx context.Context,
	job *pb.Job,
	ui terminal.UI,
	monCh chan pb.Job_State,
	localJob bool,
) (*pb.Job_Result, error) {
	log := c.logger

	// When local, we set an expiration here in case we can't gracefully
	// cancel in the event of an error. This will ensure that the jobs don't
	// remain queued forever. This is only for local ops.
	expiration := ""
	if localJob {
		expiration = "30s"
	}

	// Queue the job
	log.Debug("queueing job", "operation", fmt.Sprintf("%T", job.Operation))
	queueResp, err := c.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job:       job,
		ExpiresIn: expiration,
	})
	if err != nil {
		return nil, err
	}
	log = log.With("job_id", queueResp.JobId)

	// Stream
	return jobstream.Stream(ctx, queueResp.JobId,
		jobstream.WithClient(c.client),
		jobstream.WithLogger(log),
		jobstream.WithUI(ui),
		jobstream.WithCancelOnError(localJob),
		jobstream.WithIgnoreTerminal(localJob),
		jobstream.WithStateCh(monCh),
	)
}

// The time here is meant to encompass the typical case for an operation to begin.
// With the introduction of ondemand runners, we bumped it up from 1500 to 3000
// to accommodate the additional time before the job was picked up when testing in
// local Docker.
const stateEventPause = 3000 * time.Millisecond

var warnGitDirty = strings.TrimSpace(`
There are local changes that do not match the remote repository. By default,
Waypoint will perform this operation using a remote runner that will use the
remote repository’s git ref and not these local changes. For these changes
to be used for future operations, either commit and push, or run the operation
locally with the -local flag.
`)
