package docker

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"

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
	defer sg.Wait()
	var step terminal.Step
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	step = sg.Add("Preparing Docker configuration...")
	env := os.Environ()
	if path, err := TempDockerConfig(log, target, r.config.EncodedAuth); err != nil {
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

	// If we're a local image, then we don't want to push, just tag.
	if r.config.Local {
		return nil
	}

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
