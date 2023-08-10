// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/datasource"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectApplyCommand struct {
	*baseCommand

	flagDataSource            string
	flagGitURL                string
	flagGitPath               string
	flagGitRef                string
	flagGitAuthType           string
	flagGitUsername           string
	flagGitPassword           string
	flagGitKeyPath            string
	flagGitKeyPassword        string
	flagGitRecurseSubmodules  int
	flagFromWaypointHcl       string
	flagWaypointHcl           string
	flagPoll                  *bool
	flagPollInterval          string
	flagAppStatusPoll         *bool
	flagAppStatusPollInterval string
}

func (c *ProjectApplyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()
	ctx := c.Ctx

	if len(args) != 1 {
		c.ui.Output("Single argument required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Checking for an existing project named: %s", name)
	defer func() { s.Abort() }()

	// Check for an existing project of the same name.
	resp, err := c.project.Client().GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: name,
		},
	})
	if status.Code(err) == codes.NotFound {
		// If the error is a not found error, act as though there is no error
		// and the project is nil so that we can handle that later.
		resp = nil
		err = nil
	}
	if err != nil {
		c.ui.Output(
			"Error checking for project: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Setup our project that we're going to override
	var proj *pb.Project
	var updated bool
	if resp != nil {
		s.Update("Updating project %q...", name)
		updated = true
		proj = resp.Project
	} else {
		s.Update("Creating project %q...", name)
		proj = &pb.Project{Name: name}
	}

	// If we were specified a file then we're going to load that up.
	if c.flagFromWaypointHcl != "" {
		path, err := filepath.Abs(c.flagFromWaypointHcl)
		if err != nil {
			c.ui.Output(
				"Error loading HCL file specified with the -from-waypoint-hcl flag: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		cfg, err := configpkg.Load(path, &configpkg.LoadOptions{
			Pwd:       filepath.Dir(path),
			Workspace: c.refWorkspace.Workspace,
		})
		if err != nil {
			c.ui.Output(
				"Error loading HCL file specified with the -from-waypoint-hcl flag: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		// Load the data source configuration
		if dscfg := cfg.Runner.DataSource; dscfg != nil {
			factory, ok := datasource.FromString[dscfg.Type]
			if !ok {
				c.ui.Output(
					"Data source type specified in the HCL file is unknown: %q", dscfg.Type,
					terminal.WithErrorStyle(),
				)

				return 1
			}

			source := factory()
			ds, err := source.ProjectSource(dscfg.Body, cfg.HCLContext())
			if err != nil {
				c.ui.Output(
					"Error loading HCL file specified with the -from-waypoint-hcl flag: %s", clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)

				return 1
			}

			// Override the project data source with this configuration.
			proj.DataSource = ds

			// If the data source flag is not set, set it to our type. This
			// lets users override datasource values without specifying the
			// type flag since it is specified in the HCL file.
			if c.flagDataSource == "" {
				c.flagDataSource = dscfg.Type
			}
		}
	}

	// The logic throughout below has a ton of if checks cause we want to
	// merge information rather than set it all new.

	switch strings.ToLower(c.flagDataSource) {
	case "git":
		// Set remote enabled to true so users can run remote ops
		proj.RemoteEnabled = true
		if c.flagPoll != nil {
			if proj.DataSourcePoll == nil {
				proj.DataSourcePoll = &pb.Project_Poll{}
			}

			proj.DataSourcePoll.Enabled = *c.flagPoll
			proj.DataSourcePoll.Interval = c.flagPollInterval
		}

		if c.flagAppStatusPoll != nil {
			if proj.StatusReportPoll == nil {
				proj.StatusReportPoll = &pb.Project_AppStatusPoll{}
			}

			proj.StatusReportPoll.Enabled = *c.flagAppStatusPoll
			proj.StatusReportPoll.Interval = c.flagAppStatusPollInterval
		}

		// If the project existing datasource is Git, then we're overriding.
		// If the existing datasource is not Git or not set, then we set it
		// to Git and create new.
		var gitInfo *pb.Job_Git
		if proj.DataSource != nil {
			if v, ok := proj.DataSource.Source.(*pb.Job_DataSource_Git); ok {
				gitInfo = v.Git
			}
		}
		if gitInfo == nil {
			gitInfo = &pb.Job_Git{}
			proj.DataSource = &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{Git: gitInfo},
			}
		}

		if v := c.flagGitURL; v != "" {
			gitInfo.Url = v
		}
		if v := c.flagGitPath; v != "" {
			gitInfo.Path = v
		}
		if v := c.flagGitRef; v != "" {
			gitInfo.Ref = v
		}
		if v := c.flagGitRecurseSubmodules; v > 0 {
			gitInfo.RecurseSubmodules = uint32(v)
		}

		switch strings.ToLower(c.flagGitAuthType) {
		case "basic":
			authInfo, ok := gitInfo.Auth.(*pb.Job_Git_Basic_)
			if !ok {
				authInfo = &pb.Job_Git_Basic_{Basic: &pb.Job_Git_Basic{}}
				gitInfo.Auth = authInfo
			}

			if v := c.flagGitUsername; v != "" {
				authInfo.Basic.Username = v
			}
			if v := c.flagGitPassword; v != "" {
				authInfo.Basic.Password = v
			}

		case "ssh":
			authInfo, ok := gitInfo.Auth.(*pb.Job_Git_Ssh)
			if !ok {
				authInfo = &pb.Job_Git_Ssh{Ssh: &pb.Job_Git_SSH{}}
				gitInfo.Auth = authInfo
			}

			if v := c.flagGitKeyPath; v != "" {
				bs, err := ioutil.ReadFile(v)
				if err != nil {
					c.ui.Output(
						"Error reading private key specified with -git-private-key-path: %s", err,
						terminal.WithErrorStyle(),
					)

					return 1
				}

				authInfo.Ssh.PrivateKeyPem = bs
			}
			if v := c.flagGitKeyPassword; v != "" {
				authInfo.Ssh.Password = v
			}

		case "":
			gitInfo.Auth = nil

		default:
			c.ui.Output(
				"Unknown value for -git-auth-type set. Must be either 'basic' or 'ssh' or unset.",
				terminal.WithErrorStyle(),
			)
		}

	case "local":
		// Disable polling cause this never works with local
		if proj.DataSourcePoll != nil {
			proj.DataSourcePoll.Enabled = false
		}

		// Set the data source to local if it isn't set.
		var localInfo *pb.Job_Local
		if proj.DataSource != nil {
			if v, ok := proj.DataSource.Source.(*pb.Job_DataSource_Local); ok {
				localInfo = v.Local
			}
		}
		if localInfo == nil {
			localInfo = &pb.Job_Local{}
			proj.DataSource = &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{Local: localInfo},
			}
		}

		// We don't use localInfo yet since there is nothing to set.

	case "":
		// Do nothing, we aren't updating this information for this project.

		// Some basic error handling if polling is requested but was not explicit
		// about a data source for the project
		if c.flagPoll != nil && *c.flagPoll {
			c.ui.Output(
				"To enable polling, you must specify a git data source for the project with -data-source=git",
				terminal.WithErrorStyle(),
			)
			return 1
		}

		// Some basic error handling if app status polling is requested but was not explicit
		// about a data source for the project
		if c.flagAppStatusPoll != nil && *c.flagAppStatusPoll {
			c.ui.Output(
				"To enable application status polling, you must specify a git data "+
					"source for the project with -data-source=git",
				terminal.WithErrorStyle(),
			)
			return 1
		}

	default:
		s.Abort()

		c.ui.Output(
			"Unknown value for -data-source set. Must be either 'git' or 'local' or unset.",
			terminal.WithErrorStyle(),
		)
	}

	// Setup our default waypoint.hcl if it was given
	if v := c.flagWaypointHcl; v != "" {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			c.ui.Output(
				"Error reading HCL file specified with the -waypoint-hcl flag: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		var format pb.Hcl_Format
		switch filepath.Ext(v) {
		case ".hcl":
			format = pb.Hcl_HCL
			_, diag := hclsyntax.ParseConfig(bs, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -waypoint-hcl: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		case ".json":
			format = pb.Hcl_JSON
			_, diag := hcljson.Parse(bs, "<waypoint-hcl>")
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -waypoint-hcl: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		default:
			c.ui.Output(
				"File specified via -waypoint-hcl must end in '.hcl' or '.json'",
				terminal.WithErrorStyle(),
			)

			return 1
		}

		proj.WaypointHcl = bs
		proj.WaypointHclFormat = format
	}

	// Upsert
	_, err = c.project.Client().UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: proj,
	})
	if err != nil {
		c.ui.Output(
			"Error upserting project: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if updated {
		s.Update("Project %q updated", name)
	} else {
		s.Update("Project %q created", name)
	}
	s.Done()

	return 0
}

func (c *ProjectApplyCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "from-waypoint-hcl",
			Target:  &c.flagFromWaypointHcl,
			Default: "",
			Usage: "waypoint.hcl formatted file to load settings from. This can be used " +
				"to read settings from a file. Additional flags will override values found " +
				"in the file. Note that any settings in the file will NOT be merged with " +
				"what is already in the server; they will overwrite the server.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "waypoint-hcl",
			Target:  &c.flagWaypointHcl,
			Default: "",
			Usage: "Path to a waypoint.hcl file to associate with this project. This " +
				"is only necessary if a waypoint.hcl is not committed alongside the project " +
				"source code. If a waypoint.hcl file does not exist in the project source " +
				"then this waypoint.hcl file will be used. This file will not be validated " +
				"until an operation is run against the project; this is done on purpose since " +
				"the waypoint.hcl file may depend on files in the source repository.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "data-source",
			Target:  &c.flagDataSource,
			Default: "",
			Usage: "The data source type to use (currently only supports 'git' or 'local'). " +
				"Associated data source settings (such as flags starting with '-git') will not " +
				"take effect unless the appropriate data source is set with this flag.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-url",
			Target:  &c.flagGitURL,
			Default: "",
			Usage:   "URL of the Git repository to clone. This can be an HTTP or SSH URL.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-path",
			Target:  &c.flagGitPath,
			Default: "",
			Usage: "Path is a subdirectory within the checked out repository to " +
				"go into for the project's configuration. This must be a relative path " +
				"and may not contain '..'",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-ref",
			Target:  &c.flagGitRef,
			Default: "",
			Usage:   "Git ref (i.e. branch, tag, commit) to clone on new operations.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-auth-type",
			Target:  &c.flagGitAuthType,
			Default: "",
			Usage: "Authentication type for Git. If set, must be one of 'basic' or 'ssh'. " +
				"Basic auth is username/password and SSH uses an SSH key.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-username",
			Target:  &c.flagGitUsername,
			Default: "",
			Usage: "Username for authentication when git-auth-type is 'basic'. " +
				"For GitHub, this can be any value but it must be non-empty.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-password",
			Target:  &c.flagGitPassword,
			Default: "",
			Usage: "Password for authentication when git-auth-type is 'basic'. " +
				"For GitHub, this should be a personal access token (PAT).",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-private-key-path",
			Target:  &c.flagGitKeyPath,
			Default: "",
			Usage:   "Path to a PEM-encoded private key for 'ssh'-based auth.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "git-private-key-password",
			Target:  &c.flagGitKeyPassword,
			Default: "",
			Usage: "Password for the private key specified by git-private-key-path " +
				"if the key requires a password to decode. This is not required if " +
				"the private key doesn't require a password.",
		})

		f.IntVar(&flag.IntVar{
			Name:    "git-recurse-submodules",
			Target:  &c.flagGitRecurseSubmodules,
			Default: 0,
			Usage: "The maximum depth to recursively clone submodules. A value of " +
				"zero disables cloning any submodules recursively.",
		})

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "poll",
			Target: &c.flagPoll,
			Usage: "Enable polling. This is only valid if a Git data source is supplied. " +
				"This will watch the repo for changes and trigger a remote 'up'.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "poll-interval",
			Target:  &c.flagPollInterval,
			Default: "30s",
			Usage:   "Interval between polling if polling is enabled.",
		})

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "app-status-poll",
			Target: &c.flagAppStatusPoll,
			Usage: "Enable polling to continuously generate status reports for apps. " +
				"This is only valid if a Git data source is supplied.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "app-status-poll-interval",
			Target:  &c.flagAppStatusPollInterval,
			Default: "5m",
			Usage:   "Interval between polling to generate status reports if polling is enabled.",
		})
	})
}

func (c *ProjectApplyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectApplyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectApplyCommand) Synopsis() string {
	return "Create or update a project."
}

func (c *ProjectApplyCommand) Help() string {
	return formatHelp(`
Usage: waypoint project apply [options] NAME

  Create or update a project.

  This will create a new project with the given options. If a project with
  the same name already exists, this will update the existing project using
  the fields that are set.

  This command should be used to create a new project pointing to a VCS
  repo. If you have a "waypoint.hcl" file and a local repository, you can
  also use "waypoint init" in the directory of the project.

  You may create a project from a waypoint.hcl file and optionally overwrite
  some fields using flags by specifying the -waypoint-hcl flag.

` + c.Flags().Help())
}
