package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/datadir"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/appconfig"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/config/variables"
	"github.com/hashicorp/waypoint/internal/config/variables/formatter"
	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// executeJob executes an assigned job. This will source the data (if necessary),
// setup the project, execute the job, and return the outcome.
func (r *Runner) executeJob(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	assignment *pb.RunnerJobStreamResponse_JobAssignment,
	wd string,
	clientMutex *sync.Mutex,
	client pb.Waypoint_RunnerJobStreamClient,
) (*pb.Job_Result, error) {
	job := assignment.Job

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

	// We're going to slowly accumulate configuration info that we use to
	// update our job later.
	var configInfo pb.Job_Config

	// Try to load the waypoint.hcl from the working directory first.
	configInfo.Source = pb.Job_Config_FILE
	path, err := configpkg.FindPath(wd, "", false)
	if err != nil {
		return nil, err
	}

	// Waypoint.hcl set on the job overrides all, even if we found a
	// waypoint.hcl in the working directory above.
	if job.WaypointHcl != nil && len(job.WaypointHcl.Contents) > 0 {
		log.Info("using waypoint.hcl associated with the project in the server")

		// ext has the extra extension information for the file. We add
		// ".json" if this is JSON-formatted.
		ext := ""
		if job.WaypointHcl.Format == pb.Hcl_JSON {
			ext = ".json"
		}

		// We just write this into the working directory.
		path = filepath.Join(wd, configpkg.Filename+ext)
		if err := ioutil.WriteFile(path, job.WaypointHcl.Contents, 0644); err != nil {
			return nil, status.Errorf(codes.Internal,
				"Failed to write waypoint.hcl from job metadata: %s", err)
		}

		configInfo.Source = pb.Job_Config_JOB
	}

	// If we still have no path, try to load from the project.
	if path == "" {
		log.Trace("waypoint.hcl not found in downloaded data, looking for default in server")
		if v := resp.Project.WaypointHcl; len(v) > 0 {
			log.Info("using waypoint.hcl associated with the project in the server")

			// ext has the extra extension information for the file. We add
			// ".json" if this is JSON-formatted.
			ext := ""
			if resp.Project.WaypointHclFormat == pb.Hcl_JSON {
				ext = ".json"
			}

			// We just write this into the working directory.
			path = filepath.Join(wd, configpkg.Filename+ext)
			if err := ioutil.WriteFile(path, v, 0644); err != nil {
				return nil, status.Errorf(codes.Internal,
					"Failed to write waypoint.hcl from server: %s", err)
			}

			configInfo.Source = pb.Job_Config_SERVER
		} else {
			log.Trace("waypoint.hcl not found in server data")
		}
	}

	if path == "" {
		// Only warn if operation is project destroy
		if _, ok := job.Operation.(*pb.Job_DestroyProject); ok {
			log.Warn("A waypoint.hcl was not found.")
		} else {
			// No waypoint.hcl file is found.
			return nil, status.Errorf(codes.FailedPrecondition,
				"A waypoint.hcl was not found. Please either add a waypoint.hcl to "+
					"the project source or in the project settings in the Waypoint UI.")
		}
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

	// Send our first update so we know where the configuration came from.
	log.Trace("sending initial config load information back to server")
	clientMutex.Lock()
	err = client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_ConfigLoad_{
			ConfigLoad: &pb.RunnerJobStreamRequest_ConfigLoad{
				Config: &configInfo,
			},
		},
	})
	clientMutex.Unlock()
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
	if _, err := cfg.Validate(); err != nil {
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

	// Load variable values from the environment.
	envVars, diags := variables.LoadEnvValues(cfg.InputVariables)
	if diags.HasErrors() {
		return nil, diags
	}

	// Here we'll load our values from auto vars files and the server/UI, and
	// combine them with any values set on the job
	// The order values are added to our final pbVars slice is the order
	// of precedence
	vcsVars, diags := variables.LoadAutoFiles(wd)
	if diags.HasErrors() {
		return nil, diags
	}

	// Combine all our known variables into a list of variables. This
	// determines the precedence order of values.
	pbVars := resp.Project.GetVariables()
	pbVars = append(pbVars, envVars...)
	pbVars = append(pbVars, vcsVars...)
	pbVars = append(pbVars, job.Variables...)

	log.Debug("looking to see if there are dynamic variable default values to load")

	// Load any dynamic default values. This happens after the above
	// because we only load dynamic default values for variables we do
	// not have values for.
	if variables.NeedsDynamicDefaults(pbVars, cfg.InputVariables) {
		// If we're a local runner, we can't support this. We don't
		// support dynamic config sourcing on local runners since it
		// requires config sourcing plugins, auth, etc.
		if r.local {
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"Variables with dynamic defaults cannot be used with local "+
					"runners. Projects using variables with dynamic defaults must "+
					"be executed remotely using a remote runner.",
			)
		}

		log.Debug("loading default values for dynamic variables")

		dynamicVars, diags := variables.LoadDynamicDefaults(
			ctx,
			log,
			pbVars,
			assignment.ConfigSources,
			cfg.InputVariables,
			appconfig.WithLogger(log),
			appconfig.WithPlugins(r.configPlugins),
			appconfig.WithDynamicEnabled(true), // true because we've already determined variables need dynamic defaults
		)
		if diags.HasErrors() {
			log.Warn("failed to load dynamic defaults for variables", "diags", diags)
			return nil, diags
		}

		if len(dynamicVars) > 0 {
			log.Debug("dynamic variables discovered, adding to project variables")
			// If we have dynamic variable values, we _prepend_ them so that
			// they have the lowest precedence. In reality, this shouldn't
			// matter because we only grab dynamic values for vars that have
			// no value set, but we might as well be careful.
			pbVars = append(dynamicVars, pbVars...)
		} else {
			log.Debug("no dynamic variables found")
		}
	}

	// Evaluate all variables against the variable blocks we just decoded:
	// We grab the server cookie here to pass along for the variables
	// evaluation to use a salt for sensitive values
	clientResp, err := r.client.GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	var serverCookie string
	if clientResp != nil && clientResp.Config != nil {
		serverCookie = clientResp.Config.Cookie
	} else {
		panic("server config does not exist")
	}
	// We set both inputVars and jobVars on the project.
	// inputVars is the set of cty.Values to use in our hcl evaluation
	// and jobVars is the matching set of variable refs to store on the job that
	// has sensitive values obfuscated and is used for user-facing feedback/output.
	inputVars, jobVars, diags := variables.EvaluateVariables(log, pbVars, cfg.InputVariables, serverCookie)
	if diags.HasErrors() {
		return nil, diags
	}
	// Update the job with the final set of variable values
	log.Debug("setting final set of variable values on the job")
	clientMutex.Lock()
	err = client.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_VariableValuesSet_{
			VariableValuesSet: &pb.RunnerJobStreamRequest_VariableValuesSet{
				FinalValues: jobVars,
			},
		},
	})
	clientMutex.Unlock()
	if err != nil {
		return nil, err
	}
	// log outputtable values
	output := formatter.ValuesForOutput(jobVars)
	for name, value := range output {
		log.Debug("set variable", "name", name, "value", value.Value, "type", value.Type, "source", value.Source)
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

	case *pb.Job_DestroyProject:
		return r.executeDestroyProjectOp(ctx, log, job, project, cfg)

	case *pb.Job_StatusReport:
		return r.executeStatusReportOp(ctx, log, job, project)

	case *pb.Job_Init:
		return r.executeInitOp(ctx, log, project)

	case *pb.Job_PipelineStep:
		return r.executePipelineStepOp(ctx, log, job, project)

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
		reattachPluginConfigs, err = plugin.ParseReattachPlugins(reattachPluginsStr)
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

// operationNoDataFunc is the function type for operations that are
// executed without data downloaded.
type operationNoDataFunc func(*Runner, context.Context, hclog.Logger, *pb.Job) (*pb.Job_Result, error)

var operationsNoData = map[reflect.Type]operationNoDataFunc{
	reflect.TypeOf((*pb.Job_Poll)(nil)): nil,
}
