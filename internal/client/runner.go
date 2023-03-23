// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/hashicorp/go-hclog"
	empty "google.golang.org/protobuf/types/known/emptypb"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// startRunner initializes and starts a local runner.
// It stores it on the parent struct, and will deactivate it
// when the parent is closed.
func (c *Project) startRunner(ctx context.Context) error {
	c.logger.Debug("starting runner to process local jobs")

	// Initialize our runner
	r, err := runner.New(
		runner.WithClient(c.client),
		runner.WithLogger(c.logger.Named("runner")),
		runner.ByIdOnly(),      // We'll direct target this
		runner.WithLocal(c.UI), // Local mode
	)
	if err != nil {
		return err
	}

	// Start the runner
	if err := r.Start(ctx); err != nil {
		return err
	}

	c.activeRunner = r

	// We spin up the job processing here. Anything that spawns jobs (either locally spawned
	// or server spawned) will be processed by this runner ONLY if the runner is directly targeted.
	// Because this runner's lifetime is bound to a CLI context and therefore transient, we don't
	// want to accept jobs that aren't related to local activities (jobs queued or RPCs made)
	// because they'll hang the CLI randomly as those jobs run (it's also a security issue).
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		r.AcceptMany(c.bg)
	}()

	return nil
}

// remoteOpPreferred attempts to determine if the current waypoint infrastructure will be successful
// performing a remote operation against this project. It verifies the project's datasource,
// it's ODR runner profile, and detects if a remote runner is currently registered.
// If an operation can occur successfully remotely, we prefer the remote environment for consistency
// and security reasons.
//
// Note that this cannot guarantee that an operation will succeed remotely - the remote environment
// may not have the right auth configured, the right plugins configured, etc.
func remoteOpPreferred(ctx context.Context, client pb.WaypointClient, project *pb.Project, runnerCfgs []*configpkg.Runner, log hclog.Logger) (bool, error) {
	if !project.RemoteEnabled {
		log.Debug("Remote operations are disabled for this project - operation cannot occur remotely")
		return false, nil
	}

	if project.DataSource == nil {
		log.Debug("Project has no datasource configured - operation cannot occur remotely")
		// This is probably going to be fatal somewhere downstream
		return false, nil
	}

	var hasRemoteDataSource bool
	switch project.DataSource.GetSource().(type) {
	case *pb.Job_DataSource_Local:
		hasRemoteDataSource = false
	default:
		hasRemoteDataSource = true
	}

	if !hasRemoteDataSource {
		log.Debug("Project does not have a remote data source - operation cannot occur remotely")
		return false, nil
	}

	// We know the project can handle remote ops at this point - but do we have runners?

	runnersResp, err := client.ListRunners(ctx, &pb.ListRunnersRequest{})
	if err != nil {
		return false, err
	}
	hasRemoteRunner := false
	for _, runner := range runnersResp.Runners {
		if _, ok := runner.Kind.(*pb.Runner_Remote_); ok {
			// NOTE(izaak): There is currently no way to distinguish between a remote runner and a CLI runner.
			// So if some other waypoint client is performing an operation at this moment, we will interpret
			// that as a remote runner, and this will return a false positive.

			// Also note that this is designed to run before se start our own CLI runner.
			hasRemoteRunner = true
			break
		}
	}
	if !hasRemoteRunner {
		log.Debug("No remote runner detected - operation cannot occur remotely")
		return false, nil
	}

	// For now, if any app has an ODR profile set, we'll prefer remote for every op
	// NOTE: this means that it isn't possible to have one app in a project execute
	// locally only, and another execute remotely.
	// TODO(izaak): it's possible for us to fix this by invoking this once per app, instead of once per project
	for _, runnerCfg := range runnerCfgs {
		if runnerCfg == nil {
			continue
		}
		if runnerCfg.Profile != "" {
			log.Warn("An explicit ODR profile is set - choosing remote operations for all app operations.")
			return true, nil
		}
	}

	// Check to see if we have a global default ODR profile
	// TODO: it would be more efficient if we had an arg to filter to just get default profiles.
	configsResp, err := client.ListOnDemandRunnerConfigs(ctx, &empty.Empty{})
	if err != nil {
		return false, err
	}

	defaultRunnerProfileExists := false
	for _, odrConfig := range configsResp.Configs {
		if odrConfig.Default {
			defaultRunnerProfileExists = true
			break
		}
	}

	if defaultRunnerProfileExists {
		log.Debug("Default runner profile exists - operation is possible remotely.")
		return true, nil
	}

	log.Debug("No runner profile is set for this project and no global default exists - operation should happen locally")

	// The operation here _could_ still happen remotely - executed on the remote runner itself without ODR.
	// If it's a container build op it will probably fail (because no kaniko), and if it's a deploy/release op it
	// very well might fail do to incorrect/insufficient permissions. Because it probably won't work, we won't try,
	// but the user could force it to happen locally by setting -local=false.
	return false, nil
}
