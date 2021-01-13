package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *Registry) pushWithImg(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	source *Image,
	target *Image,
) error {
	sg := ui.StepGroup()
	var step terminal.Step
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	step = sg.Add("Preparing Docker configuration...")
	env := os.Environ()
	if path, err := r.createDockerConfig(log, target); err != nil {
		return err
	} else if path != "" {
		defer os.RemoveAll(path)
		env = append(env, "DOCKER_CONFIG="+path)
	}

	step.Done()
	step = sg.Add("Tagging the image from %s => %s", source.Name(), target.Name())

	// Tag
	cmd := exec.CommandContext(ctx,
		"img",
		"tag",
		source.Name(),
		target.Name(),
	)
	cmd.Env = env
	cmd.Stdout = step.TermOutput() // Command output should go to the step
	cmd.Stderr = cmd.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	step.Done()
	step = sg.Add("Pushing image...")

	// Push the image
	var buf bytes.Buffer
	cmd = exec.CommandContext(ctx,
		"img",
		"push",
		target.Name(),
	)
	cmd.Env = env
	cmd.Stdout = io.MultiWriter(&buf, step.TermOutput()) // Command output should go to the step
	cmd.Stderr = cmd.Stdout
	buf.Reset()
	if err := cmd.Run(); err != nil {
		return status.Errorf(codes.Internal,
			"Failure while pushing image: %s\n\n%s", err, buf.String())
	}

	step.Done()
	return nil
}

// createDockerConfig creates a new Docker configuration with the
// configured auth in it. It saves this Docker config to a temporary path
// and returns the path to that Docker file.
//
// We have to do this because `img` doesn't support setting auth for
// a single operation. Therefore, we must set auth in the Docker config,
// but we don't want to pollute any concurrent runs or the main file. So
// we create a copy.
//
// This can return ("", nil) if there is no custom Docker config necessary.
//
// Callers should defer file deletion for this temporary file.
func (r *Registry) createDockerConfig(
	log hclog.Logger,
	target *Image,
) (string, error) {
	if r.config.EncodedAuth == "" {
		return "", nil
	}

	// Create a reader that base64 decodes our encoded auth and then
	// JSON decodes that.
	var authCfg types.AuthConfig
	var rdr io.Reader = strings.NewReader(r.config.EncodedAuth)
	rdr = base64.NewDecoder(base64.URLEncoding, rdr)
	dec := json.NewDecoder(rdr)
	if err := dec.Decode(&authCfg); err != nil {
		return "", status.Errorf(codes.FailedPrecondition,
			"Failed to decode encoded_auth: %s", err)
	}

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}
	host := reference.Domain(ref)
	log.Trace("auth host", "host", host)

	// Parse our old Docker config and add the auth.
	log.Trace("loading Docker configuration")
	file, err := config.Load(config.Dir())
	if err != nil {
		return "", err
	}

	if file.AuthConfigs == nil {
		file.AuthConfigs = map[string]types.AuthConfig{}
	}
	file.AuthConfigs[host] = authCfg

	// Create a temporary directory for our config
	td, err := ioutil.TempDir("", "wp-docker-config")
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary directory for Docker config: %s", err)
	}

	// Create a temporary file and write our Docker config to it
	f, err := os.Create(filepath.Join(td, "config.json"))
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}
	defer f.Close()
	if err := file.SaveToWriter(f); err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}

	log.Info("temporary Docker config created for auth",
		"auth_host", host,
		"path", td,
	)

	return td, nil
}
