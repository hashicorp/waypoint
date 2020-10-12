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

// Platform is the Platform implementation for Kubernetes.
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

// Deploy deploys an image to Kubernetes.
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
	cmd.Dir = src.Path
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
	tpl, err := template.New("tpl").ParseFiles(tplconfig.Path)
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
	return "", nil
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

	doc.Description("Deploy a container to Docker, local or remote")

	doc.Example(`
deploy {
  use "docker" {
	command      = "ps"
	service_port = 3000
	static_environment = {
	  "environment": "production",
	  "LOG_LEVEL": "debug"
	}
  }
}
`)

	doc.Input("docker.Image")
	doc.Output("docker.Deployment")

	doc.SetField(
		"command",
		"the command to run to start the application in the container",
		docs.Default("the image entrypoint"),
	)

	doc.SetField(
		"scratch_path",
		"a path within the container to store temporary data",
		docs.Summary(
			"docker will mount a tmpfs at this path",
		),
	)

	doc.SetField(
		"static_environment",
		"environment variables to expose to the application",
		docs.Summary(
			"these environment variables should not be run of the mill",
			"configuration variables, use waypoint config for that.",
			"These variables are used to control over all container modes,",
			"such as configuring it to start a web app vs a background worker",
		),
	)

	doc.SetField(
		"service_port",
		"port that your service is running on in the container",
		docs.Default("3000"),
	)

	return doc, nil
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
