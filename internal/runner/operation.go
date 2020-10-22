package runner

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/datadir"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	configpkg "github.com/hashicorp/waypoint/internal/config2"
	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// executeJob executes an assigned job. This will source the data (if necessary),
// setup the project, execute the job, and return the outcome.
func (r *Runner) executeJob(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
	wd string,
) (*pb.Job_Result, error) {
	// Eventually we'll need to extract the data source. For now we're
	// just building for local exec so it is the working directory.
	path := configpkg.Filename
	if wd != "" {
		path = filepath.Join(wd, path)
	}

	// Determine the evaluation context we'll be using
	log.Trace("reading configuration", "path", path)
	cfg, err := configpkg.Load(path, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	// Setup our project data directory.
	projDir, err := datadir.NewProject(".waypoint")
	if err != nil {
		return nil, err
	}

	// Find all our plugins
	factories, err := r.pluginFactories(log, cfg.Plugins(), wd)
	if err != nil {
		return nil, err
	}

	// Build our job info
	jobInfo := &component.JobInfo{
		Id:    job.Id,
		Local: r.local,
	}

	// Create our project
	log.Trace("initializing project", "project", cfg.Project)
	project, err := core.NewProject(ctx,
		core.WithLogger(log),
		core.WithUI(ui),
		core.WithComponents(factories),
		core.WithClient(r.client),
		core.WithConfig(cfg),
		core.WithDataDir(projDir),
		core.WithLabels(job.Labels),
		core.WithWorkspace(job.Workspace.Workspace),
		core.WithJobInfo(jobInfo),
	)
	if err != nil {
		return nil, err
	}
	defer project.Close()

	// Execute the operation
	log.Info("executing operation", "type", fmt.Sprintf("%T", job.Operation))
	switch job.Operation.(type) {
	case *pb.Job_Noop_:
		if r.noopCh != nil {
			select {
			case <-r.noopCh:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		log.Debug("noop job success")
		return nil, nil

	case *pb.Job_Build:
		return r.executeBuildOp(ctx, job, project)

	case *pb.Job_Push:
		return r.executePushOp(ctx, job, project)

	case *pb.Job_Deploy:
		return r.executeDeployOp(ctx, job, project)

	case *pb.Job_Destroy:
		return r.executeDestroyOp(ctx, job, project)

	case *pb.Job_Release:
		return r.executeReleaseOp(ctx, log, job, project)

	case *pb.Job_Validate:
		return r.executeValidateOp(ctx, job, project)

	case *pb.Job_Auth:
		return r.executeAuthOp(ctx, log, job, project)

	case *pb.Job_Docs:
		return r.executeDocsOp(ctx, log, job, project)

	default:
		return nil, status.Errorf(codes.Aborted, "unknown operation %T", job.Operation)
	}
}

func (r *Runner) pluginFactories(
	log hclog.Logger,
	plugins []*configpkg.Plugin,
	wd string,
) (map[component.Type]*factory.Factory, error) {
	// Copy all our base factories first
	result := map[component.Type]*factory.Factory{}
	for k, f := range r.factories {
		result[k] = f.Copy()
	}

	// Get our plugin search paths
	pluginPaths, err := plugin.DefaultPaths(wd)
	if err != nil {
		return nil, err
	}
	log.Debug("plugin search path", "path", pluginPaths)

	// Search for all of our plugins
	var perr error
	for _, pluginCfg := range plugins {
		plog := log.With("plugin_name", pluginCfg.Name)
		plog.Debug("searching for plugin")

		// Find our plugin.
		cmd, err := plugin.Discover(pluginCfg, pluginPaths)
		if err != nil {
			plog.Warn("error searching for plugin", "err", err)
			perr = multierror.Append(perr, err)
			continue
		}

		// If the plugin was not found, it is only an error if
		// we don't have it already registered.
		if cmd == nil {
			if _, ok := plugin.Builtins[pluginCfg.Name]; !ok {
				perr = multierror.Append(perr, fmt.Errorf(
					"plugin %q not found",
					pluginCfg.Name))
				plog.Warn("plugin not found")
			} else {
				plog.Debug("plugin found as builtin")
				for _, t := range pluginCfg.Types() {
					result[t].Register(pluginCfg.Name, plugin.BuiltinFactory(pluginCfg.Name, t))
				}
			}

			continue
		}

		// Register the command
		plog.Debug("plugin found as external binary", "path", cmd.Path)
		for _, t := range pluginCfg.Types() {
			result[t].Register(pluginCfg.Name, plugin.Factory(cmd, t))
		}
	}

	return result, perr
}
