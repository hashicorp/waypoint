package docker

import (
	"context"
	"os/exec"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// HasImg returns true if "img" is available on the PATH.
//
// This doesn't do any fancy checking that "img" is the "img" we expect.
// We can make the checking here more advanced later.
func HasImg() bool {
	_, err := exec.LookPath("img")
	return err == nil
}

func (b *Builder) buildWithImg(
	ctx context.Context,
	ui terminal.UI,
	sg terminal.StepGroup,
	dockerfilePath string,
	contextDir string,
	tag string,
) error {
	step := sg.Add("Building Docker image with img...")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	// NOTE(mitchellh): we can probably use the img Go pkg directly one day.
	cmd := exec.CommandContext(ctx,
		"img",
		"build",
		"-f", dockerfilePath,
		"-t", tag,
		".",
	)

	// Working directory to directory with build context
	cmd.Dir = contextDir

	// Command output should go to the step
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	step.Done()
	return nil
}
