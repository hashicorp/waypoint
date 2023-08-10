// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package exec

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Platform is the Platform implementation for exec.
type Platform struct {
	config PlatformConfig
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// Deploy deploys an image to exec.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	input *Input,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	args := p.config.Command
	if len(args) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"command must not be empty")
	}

	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()

	// If we have a step set, abort it on exit
	var s terminal.Step
	defer func() {
		if s != nil {
			s.Abort()
		}
	}()

	// Render templates if set
	s = sg.Add("Rendering templates...")

	// Build our template data
	var data tplData
	data.Env = deployConfig.Env()
	data.Workspace = job.Workspace
	data.Populate(input)

	if tpl := p.config.Template; tpl != nil {
		// Render our template
		path, closer, err := p.renderTemplate(tpl, &data)
		if closer != nil {
			defer closer()
		}
		if err != nil {
			return nil, err
		}

		// Replace the template path in the arguments list
		for i, v := range args {
			args[i] = strings.ReplaceAll(v, "<TPL>", path)
		}
	}

	// Render our arguments
	for i, v := range args {
		v, err := p.renderTemplateString(v, &data)
		if err != nil {
			return nil, err
		}

		args[i] = v
	}

	s.Done()
	s = sg.Add("Executing command: %s", strings.Join(args, " "))

	// Ensure we're executing a binary
	if !filepath.IsAbs(args[0]) {
		log.Debug("command is not absolute, will look up on PATH", "command", args[0])
		path, err := exec.LookPath(args[0])
		if err != nil {
			log.Info("failed to find command on PATH", "command", args[0])
			return nil, err
		}

		log.Info("command is not absolute, replaced with value on PATH",
			"old_command", args[0],
			"new_command", path,
		)
		args[0] = path
	}

	// Run our command
	var cmd exec.Cmd
	cmd.Path = args[0]
	cmd.Args = args
	if p.config.Dir != "" {
		cmd.Dir = p.config.Dir
	} else {
		cmd.Dir = src.Path
	}
	cmd.Stdout = s.TermOutput()
	cmd.Stderr = cmd.Stdout

	// Run it
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	s.Done()

	return &Deployment{}, nil
}

func (p *Platform) renderTemplate(tpl *ConfigTemplate, data *tplData) (string, func(), error) {
	fi, err := os.Stat(tpl.Path)
	if err != nil {
		return "", nil, err
	}

	// Create a temporary directory to store our renders
	td, err := ioutil.TempDir("", "waypoint-exec")
	if err != nil {
		return "", nil, err
	}
	closer := func() { os.RemoveAll(td) }

	// Render
	var path string
	if fi.IsDir() {
		path, err = p.renderTemplateDir(tpl, data, td)
	} else {
		path, err = p.renderTemplateFile(tpl, data, td)
	}

	return path, closer, err
}

func (p *Platform) renderTemplateString(v string, data *tplData) (string, error) {
	// Build our template
	tpl, err := template.New("tpl").Parse(v)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *Platform) renderTemplateFile(tplconfig *ConfigTemplate, data *tplData, td string) (string, error) {
	// We'll copy the file into the temporary directory
	path := filepath.Join(td, filepath.Base(tplconfig.Path))

	// Build our template
	tpl, err := template.New(filepath.Base(path)).ParseFiles(tplconfig.Path)
	if err != nil {
		return "", err
	}

	// Create our target path
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return path, tpl.Execute(f, data)
}

func (p *Platform) renderTemplateDir(tpl *ConfigTemplate, data *tplData, td string) (string, error) {
	return td, filepath.Walk(tpl.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := td

		// Determine if we have any directory
		stripped := strings.TrimPrefix(path, tpl.Path)
		if len(stripped) == 0 {
			panic("empty path") // should never happen
		}
		if stripped[0] == '/' || stripped[0] == '\\' {
			// Get rid of any prefix '/' which could happen if tpl.Path doesn't
			// end in a directory sep.
			stripped = stripped[1:]
		}
		if v := filepath.Dir(stripped); v != "." {
			dir = filepath.Join(dir, v)
			if err := os.MkdirAll(dir, 0700); err != nil {
				return err
			}
		}

		// Render
		_, err = p.renderTemplateFile(&ConfigTemplate{Path: path}, data, dir)
		return err
	})
}

// PlatformConfig is the configuration structure for the Platform.
type PlatformConfig struct {
	// The command to execute. The string value "<TPL>" will be replaced
	// with the rendered template. If the template is a file, the value of
	// TPL will be a file. If the template is a directory, TPL will be a path
	// to a directory.
	Command []string `hcl:"command,optional"`

	// Dir is the working directory to set when executing the command.
	// This will default to the path to the application in the Waypoint
	// configuration.
	Dir string `hcl:"dir,optional"`

	// Template is the template to render.
	Template *ConfigTemplate `hcl:"template,block"`
}

type ConfigTemplate struct {
	// Path is the path to the file or directory to template.
	Path string `hcl:"path,attr"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&PlatformConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Execute any command to perform a deploy.

This plugin lets you use almost any pre-existing deployment tool for the
deploy step of Waypoint. This is a great way to take a pre-existing application
and begin using Waypoint. For example, you can wrap "kubectl" calls if you
already have Kubernetes configurations, or "helm" if you use Helm, and so on.

The "exec" plugin is meant to be an escape hatch from Waypoint. In working
this way, you will lose many Waypoint benefits. For example, "waypoint destroy"
functionality will not work with deploys created with the exec plugin.

### Templates

The exec plugin supports templating to access input information about the
artifact. There are two mechanisms for templates:

1. Any argument in "command" is processed as a template.

2. You may specify a file or directory to be processed for templating
using the "template" stanza. Any argument with the value ` + "`<TPL>`" + ` in it
will be replaced with the path to the template.

Templating follows the format of a Go [text/template](https://golang.org/pkg/text/template/)
template. The top of the documentation there has details on the format.

#### Common Values

The following template values are always available:

  - ".Env" (map<string\>string) - These are environment variables that should
    be set on the deployed workload. These enable the entrypoint to work so
    you should set these if able.


  - ".Workspace" (string) - The workspace name that the Waypoint deploy is
    running in. This lets you potentially deploy to different clusters based
    on this value.

#### Docker Image Input

If the build step creates a Docker image, the following template variables
are available:

  - ".Input.DockerImageFull" (string) - The full Docker image name and tag.

  - ".Input.DockerImageName" (string) - The Docker image name, without the tag.

  - ".Input.DockerImageTag" (string) - The Docker image tag, such as "latest".

`)

	doc.Example(`
deploy {
  use "exec" {
    command = ["kubectl", "apply", "-f", "<TPL>"]

    template {
      path = "myapp.yml"
    }
  }
}
`)

	doc.Example(`
deploy {
  use "exec" {
    command = ["docker", "run", "{{.Input.DockerImageFull}}"]
  }
}
`)

	doc.Input("exec.Input")
	doc.Output("exec.Deployment")

	doc.SetField(
		"command",
		"The command to execute for the deploy as a list of strings.",
		docs.Summary(
			"Each value in the list will be rendered as a template, so it",
			"may contain template directives. Additionally, the special string",
			"`<TPL>` will be replaced with the path to the rendered file-based",
			"templates. If your template path was to a file, this will be a path",
			"a file. Otherwise, it will be a path to a directory.",
		),
	)

	doc.SetField(
		"dir",
		"The working directory to use while executing the command.",
		docs.Summary(
			"This will default to the same working directory as the Waypoint execution.",
		),
	)

	doc.SetField(
		"template",
		"A stanza that declares that a file or directory should be template-rendered.",

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"path",
				"The path to the file or directory to render as a template.",
				docs.Summary(
					"Templating uses the following format: https://golang.org/pkg/text/template/",
					"Available template variables depends on the input artifact.",
				),
			)
		}),
	)

	return doc, nil
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
