package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/datadir"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/config/variables"
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
	// NOTE(mitchellh; krantzinator): For now, we query the project directly here
	// since we use it only in case of a missing local waypoint.hcl, and to
	// collect input variable values set on the server. I can see us moving this
	// to accept() eventually though if other data is used.
	resp, err := r.client.GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: job.Application.Project,
		},
	})
	if err != nil {
		return nil, err
	}

	// Eventually we'll need to extract the data source. For now we're
	// just building for local exec so it is the working directory.
	path, err := configpkg.FindPath(wd, "", false)
	if err != nil {
		return nil, err
	}
	if path == "" {
		// If no waypoint.hcl file is found in the downloaded data, look for
		// a default waypoint HCL.
		log.Trace("waypoint.hcl not found in downloaded data, looking for default in server")
		if v := resp.Project.WaypointHcl; len(v) > 0 {
			log.Info("using waypoint.hcl associated with the project in the server")

			// ext has the extra extension information for the file. We add
			// ".json" if this is JSON-formatted.
			ext := ""
			if resp.Project.WaypointHclFormat == pb.Project_JSON {
				ext = ".json"
			}

			// We just write this into the working directory.
			path = filepath.Join(wd, configpkg.Filename+ext)
			if err := ioutil.WriteFile(path, v, 0644); err != nil {
				return nil, status.Errorf(codes.Internal,
					"Failed to write waypoint.hcl from server: %s", err)
			}
		} else {
			log.Trace("waypoint.hcl not found in server data")
		}
	}

	if path == "" {
		// No waypoint.hcl file is found.
		return nil, status.Errorf(codes.FailedPrecondition,
			"A waypoint.hcl was not found. Please either add a waypoint.hcl to "+
				"the project source or in the project settings in the Waypoint UI.")
	}

	// Determine the evaluation context we'll be using
	log.Trace("reading configuration", "path", path)
	cfg, err := configpkg.Load(path, &configpkg.LoadOptions{
		Pwd:       filepath.Dir(path),
		Workspace: job.Workspace.Workspace,
	})
	if err != nil {
		return nil, err
	}

	// If we have a project specified on the job, override the configuration
	// project with that. This allows the same Waypoint configuration to
	// be shared by multiple projects, which is very possible in the UI,
	// and less useful when stored as a file in the repo.
	if v := job.Application.Project; v != "" {
		cfg.Project = v
	}

	// Validate our configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Setup our project data directory.
	projDir, err := datadir.NewProject(filepath.Join(wd, ".waypoint"))
	if err != nil {
		return nil, err
	}

	// Find all our plugins
	factories, err := r.pluginFactories(log, cfg.Plugins(), wd)
	if err != nil {
		return nil, err
	}

	// Here we'll load our values from auto vars files and the server/UI, and
	// combine them with any values set on the job
	// The order values are added to our final pbVars slice is the order
	// of precedence
	vcsVars, diags := variables.LoadAutoFiles(wd)
	if diags.HasErrors() {
		return nil, diags
	}

	pbVars := resp.Project.GetVariables()
	pbVars = append(pbVars, vcsVars...)
	pbVars = append(pbVars, job.Variables...)

	// evaluate all variables against the variable blocks we just decoded
	inputVars, diags := variables.EvaluateVariables(pbVars, cfg.InputVariables, log)
	if diags.HasErrors() {
		return nil, diags
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
		core.WithVariables(inputVars),
		core.WithWorkspace(job.Workspace.Workspace),
		core.WithJobInfo(jobInfo),
	)
	if err != nil {
		return nil, err
	}
	defer project.Close()

	// Execute the operation
	//
	// Note some operation types don't require downloaded data. These are
	// not executed here but are executed in accept.go.
	log.Info("executing operation")
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

	case *pb.Job_Up:
		return r.executeUpOp(ctx, log, job, project)

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

	case *pb.Job_ConfigSync:
		return r.executeConfigSyncOp(ctx, log, job, project)

	case *pb.Job_Exec:
		return r.executeExecOp(ctx, job, project)

	case *pb.Job_Logs:
		return r.executeLogsOp(ctx, job, project)

	case *pb.Job_QueueProject:
		return r.executeQueueProjectOp(ctx, log, job, project)

	case *pb.Job_StatusReport:
		return r.executeStatusReportOp(ctx, job, project)

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

	// Look for any reattach plugins
	var reattachPluginConfigs map[string]*goplugin.ReattachConfig
	reattachPluginsStr := os.Getenv("WP_REATTACH_PLUGINS")
	if reattachPluginsStr != "" {
		var err error
		reattachPluginConfigs, err = parseReattachPlugins(reattachPluginsStr)
		if err != nil {
			return nil, err
		}
	}

	// Search for all of our plugins
	var perr error
	for _, pluginCfg := range plugins {
		plog := log.With("plugin_name", pluginCfg.Name)
		plog.Debug("searching for plugin")

		if reattachConfig, ok := reattachPluginConfigs[pluginCfg.Name]; ok {
			plog.Debug(fmt.Sprintf("plugin %s is declared as running for reattachment", pluginCfg.Name))
			for _, t := range pluginCfg.Types() {
				if err := result[t].Register(pluginCfg.Name, plugin.ReattachPluginFactory(reattachConfig, t)); err != nil {
					return nil, err
				}
			}
			continue
		}

		// Find our plugin.
		cmd, err := plugin.Discover(&plugin.Config{
			Name:     pluginCfg.Name,
			Checksum: pluginCfg.Checksum,
		}, pluginPaths)
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
					plog.Info("register", "type", t.String(), "nil", result[t] == nil)
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

// parse information on reattaching to plugins out of a
// JSON-encoded environment variable.
func parseReattachPlugins(in string) (map[string]*goplugin.ReattachConfig, error) {
	reattachConfigs := map[string]*goplugin.ReattachConfig{}
	if in != "" {
		in = strings.TrimRight(in, "'")
		in = strings.TrimLeft(in, "'")
		var m map[string]sdk.ReattachConfig
		err := json.Unmarshal([]byte(in), &m)
		if err != nil {
			return reattachConfigs, fmt.Errorf("Invalid format for WP_REATTACH_PROVIDERS: %w", err)
		}
		for p, c := range m {
			var addr net.Addr
			switch c.Addr.Network {
			case "unix":
				addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
				if err != nil {
					return reattachConfigs, fmt.Errorf("Invalid unix socket path %q for %q: %w", c.Addr.String, p, err)
				}
			case "tcp":
				addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
				if err != nil {
					return reattachConfigs, fmt.Errorf("Invalid TCP address %q for %q: %w", c.Addr.String, p, err)
				}
			default:
				return reattachConfigs, fmt.Errorf("Unknown address type %q for %q", c.Addr.String, p)
			}
			reattachConfigs[p] = &goplugin.ReattachConfig{
				Protocol:        goplugin.Protocol(c.Protocol),
				ProtocolVersion: c.ProtocolVersion,
				Pid:             c.Pid,
				Test:            c.Test,
				Addr:            addr,
			}
		}
	}
	return reattachConfigs, nil
}

// operationNoDataFunc is the function type for operations that are
// executed without data downloaded.
type operationNoDataFunc func(*Runner, context.Context, hclog.Logger, *pb.Job) (*pb.Job_Result, error)

var operationsNoData = map[reflect.Type]operationNoDataFunc{
	reflect.TypeOf((*pb.Job_Poll)(nil)): nil,
}
