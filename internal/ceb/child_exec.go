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
	cmd, err := ceb.buildCmd(ctx, cfg.ExecArgs)
	if err != nil {
		return err
	}

	ceb.childCmd = cmd
	return nil
}

// execChildCmd starts the child process, and waits for completion by sending an
// error along the channel.
func (ceb *CEB) execChildCmd(ctx context.Context) <-chan error {
	ch := make(chan error, 1)
	cmd := ceb.childCmd

	// Start our subprocess
	log := ceb.logger.With(
		"cmd", cmd.Path,
		"args", cmd.Args,
	)
	log.Info("starting child process")
	if err := cmd.Start(); err != nil {
		ch <- status.Errorf(codes.Aborted,
			"failed to execute subprocess: %s", err)
		return ch
	}

	// Start a goroutine to wait for completion
	go func() {
		err := cmd.Wait()
		if err == nil {
			log.Info("subprocess gracefully exited")
		} else {
			log.Warn("subprocess exited", "err", err)
		}

		ch <- err
	}()

	return ch
}

func (ceb *CEB) buildCmd(ctx context.Context, args []string) (*exec.Cmd, error) {
	// Avoid a crash below by verifying we got some arguments.
	if len(args) == 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"command was empty")
	}

	// Exec requires a full path to a binary. If we weren't given an absolute
	// path then we need to look it up via the PATH.
	if !filepath.IsAbs(args[0]) {
		path, err := exec.LookPath(args[0])
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"failed to find command %q on PATH: %s", args[0], err)
		}

		args[0] = path
	}

	// Build the subprocess
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
