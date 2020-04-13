package ceb

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// initChildCmd initializes the child command that we'll execute when
// we run. This just sets the `childCmd` field on the CEB structure. This
// does not have any side effecting behavior.
func (ceb *CEB) initChildCmd(ctx context.Context, cfg *config) error {
	args := cfg.ExecArgs

	// Exec requires a full path to a binary. If we weren't given an absolute
	// path then we need to look it up via the PATH.
	if !filepath.IsAbs(args[0]) {
		path, err := exec.LookPath(args[0])
		if err != nil {
			return status.Errorf(codes.InvalidArgument,
				"failed to find command %q on PATH: %s", args[0], err)
		}

		args[0] = path
	}

	// Start building our subprocess. Even if we are just going to
	// exec into it (syscall.Exec), we use an exec.Cmd to store it along
	// the way.
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ceb.childCmd = cmd

	return nil
}
