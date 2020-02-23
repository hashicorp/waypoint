package terraform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/mitchellh/go-linereader"

	"github.com/mitchellh/devflow/internal/pkg/httpfs"
	"github.com/mitchellh/devflow/sdk/datadir"
)

// Terraform is a wrapper for the Terraform executable and operations
// (Apply, Destroy) for a single workspace/configuration.
type Terraform struct {
	// Context is the context for execution. If this is cancelled then the
	// Terraform apply will be cancelled.
	Context context.Context

	// Logger is the logger to use while executing
	Logger hclog.Logger

	// Dir is the directory to use for storing configuration, state, etc.
	// This should be a scoped directory for this specific Terraform workspace
	// so that the Terraform run doesn't clobber any data.
	Dir datadir.Dir

	// ConfigFS is the http.FileSystem with Terraform configuration files.
	// These files will be copied as necessary.
	//
	// ConfigPath is the path to the directory containing the configuration
	// files within the filesystem. Example: "config/"
	ConfigFS   http.FileSystem
	ConfigPath string

	// Vars are variables to set for the Terraform configuration. The value
	// type should be a valid primitive that can be encoded to JSON for input
	// with Terraform.
	//
	// These will get encoded into a "tfvars" file. This is an important detail
	// to understand the load order.
	Vars map[string]interface{}

	configPath  string // configPath is the path to the Terraform configuration.
	backendPath string // backendPath is the path to backend configuration.
}

// Apply applies Terraform and returns the outputs from the execution.
func (tf *Terraform) Apply() (map[string]interface{}, error) {
	tf.Logger.Debug("setting up Terraform")
	if err := tf.setup(); err != nil {
		tf.Logger.Error("failed setting up Terraform", "error", err)
		return nil, err
	}

	tf.Logger.Info("running Terraform apply")
	return tf.terraformApply()
}

// setup prepares Terraform for execution by calling up to `init`.
func (tf *Terraform) setup() error {
	// Parse and log the Terraform version
	vsn, err := tf.terraformVersion()
	if err != nil {
		return fmt.Errorf("could not determine Terraform version: %s", err)
	}
	tf.Logger.Info("terraform", "version", vsn.String())

	// Download Configuration
	if err := tf.setupConfig(); err != nil {
		return fmt.Errorf("could not set up Terraform config: %s", err)
	}

	// Setup the TF vars
	if err := tf.setupVars(); err != nil {
		return fmt.Errorf("could not set up Terraform vars: %s", err)
	}

	// Setup the backend
	if err := tf.setupBackend(); err != nil {
		return fmt.Errorf("could not set up Terraform backend: %s", err)
	}

	// Initialize
	tf.Logger.Info("initializing Terraform")
	if err := tf.terraformInit(); err != nil {
		return fmt.Errorf("failed to initialize Terraform: %s", err)
	}

	return nil
}

// terraformVersion parses the Terraform version from the binary.
func (tf *Terraform) terraformVersion() (*version.Version, error) {
	var buf bytes.Buffer
	cmd, closer := tf.terraformCmd()
	defer closer()
	cmd.Args = append(cmd.Args, "version")
	cmd.Stdout = io.MultiWriter(&buf, cmd.Stdout)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	fields := strings.Fields(buf.String())
	if len(fields) < 2 {
		return nil, fmt.Errorf("Unrecognized version string: %q", buf.String())
	}

	return version.NewVersion(strings.TrimPrefix(fields[1], "v"))
}

// setupConfig copies the configuration into the target directory.
func (tf *Terraform) setupConfig() error {
	// We put the configuration under the cache directory
	tf.configPath = filepath.Join(tf.Dir.CacheDir(), "config")
	if err := os.RemoveAll(tf.configPath); err != nil {
		return err
	}
	if err := os.MkdirAll(tf.configPath, 0755); err != nil {
		return err
	}

	// Copy the files
	if err := httpfs.Copy(tf.ConfigFS, tf.configPath, tf.ConfigPath); err != nil {
		return err
	}

	return nil
}

// setupVars configures the "tfvars" files based on ApplyParams.
//
// This requires r.configPath to be set.
func (tf *Terraform) setupVars() error {
	// If there are no variables do nothing.
	if len(tf.Vars) == 0 {
		return nil
	}

	// Validate that the config path is set.
	if tf.configPath == "" {
		return fmt.Errorf("internal error: config path is not set on run")
	}

	// Encode them as JSON.
	raw, err := json.Marshal(tf.Vars)
	if err != nil {
		return err
	}

	path := filepath.Join(tf.configPath, "config.auto.tfvars.json")
	return ioutil.WriteFile(path, raw, 0644)
}

// setupBackend configures the backend configuration.
func (tf *Terraform) setupBackend() error {
	typ := "local"
	config := map[string]interface{}{
		"path": filepath.Join(tf.Dir.DataDir(), "terraform.tfstate"),
	}

	// Write the backend type
	path := filepath.Join(tf.configPath, "_backend.tf")
	err := ioutil.WriteFile(path, []byte(fmt.Sprintf(backendTemplate, typ)), 0644)
	if err != nil {
		return fmt.Errorf("could not write backend: %s", err)
	}

	// Encode the backend configuration as JSON
	raw, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not encode backend config: %s", err)
	}

	path = filepath.Join(tf.configPath, "_backend.json")
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		return fmt.Errorf("could not write backend config: %s", err)
	}
	tf.backendPath = path

	return nil
}

// terraformInit calls `terraform init` for the configuration.
func (tf *Terraform) terraformInit() error {
	cmd, closer := tf.terraformCmd()
	defer closer()
	cmd.Args = append(cmd.Args, "init", "-no-color", "-input=false")

	// If we have a backend configuration written then we need to
	// merge that configuration as part of initialization.
	if tf.backendPath != "" {
		path, err := filepath.Rel(cmd.Dir, tf.backendPath)
		if err != nil {
			return err
		}

		cmd.Args = append(cmd.Args, "-backend-config", path)
	}

	return cmd.Run()
}

// terraformApply calls `terraform apply` for the configuration.
func (tf *Terraform) terraformApply() (map[string]interface{}, error) {
	// Run the apply
	cmd, closer := tf.terraformCmd()
	defer closer()
	cmd.Args = append(cmd.Args,
		"apply",
		"-auto-approve",
		"-input=false",
		"-no-color",
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// If the apply succeeded then gather the outputs. We do this by
	// multi-writing the output to a buffer, then parsing that later.
	var buf bytes.Buffer
	cmd, closer = tf.terraformCmd()
	defer closer()
	cmd.Args = append(cmd.Args,
		"output",
		"-no-color",
		"-json",
	)
	// Override stdout: `terraform output -json` emits secrets in plain text. We dont want to log these.
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Parse the outputs
	return parseOutputs(&buf)
}

// terraformCmd returns an *exec.Cmd for executing Terraform.
func (tf *Terraform) terraformCmd() (*exec.Cmd, func()) {
	cmd := exec.CommandContext(tf.Context, "terraform")

	// Working directory is where the configuration is
	cmd.Dir = tf.configPath

	// Set the stdout/stderr to log to different pipes that we then read line-by-line
	// to output to the log with an appropriate level.
	outr, outw := io.Pipe()
	go func() {
		for line := range linereader.New(outr).Ch {
			if line != "" {
				tf.Logger.Debug(
					"terraform output",
					"line", line,
					"source", "terraform stdout",
				)
			}
		}
	}()

	errr, errw := io.Pipe()
	go func() {
		for line := range linereader.New(errr).Ch {
			if line != "" {
				tf.Logger.Error(
					"terraform output",
					"line", line,
					"source", "terraform stderr",
				)
			}
		}
	}()
	cmd.Stdout = outw
	cmd.Stderr = errw

	// Set the environment variables
	cmd.Env = os.Environ()

	return cmd, func() {
		outw.Close()
		outr.Close()
		errw.Close()
		errr.Close()
	}
}

// backendTemplate is the template we use to set the backend type. The
// configuration is passed in via JSON files.
const backendTemplate = `
terraform {
  backend "%s" {}
}
`
