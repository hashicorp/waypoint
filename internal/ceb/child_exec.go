// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ceb

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// initChildCmd initializes the child command that we'll execute when
// we run. This just sets the `childCmd` field on the CEB structure. This
// does not have any side effecting behavior.
func (ceb *CEB) initChildCmd(ctx context.Context, cfg *config) error {
	// Prepare our base command. This validates that the path and everything is found.
	cmd, err := ceb.buildCmd(ctx, cfg.ExecArgs)
	if err != nil {
		return err
	}
	ceb.childCmdBase = cmd

	// Setup our channels and goroutine to prepare to execute this. This
	// won't actually start any child commands, but it starts the watcher
	// goroutine that will eventually run them.
	childCmdCh := make(chan *exec.Cmd, 5)
	doneCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	go ceb.watchChildCmd(ctx, childCmdCh, doneCh, sigCh)
	ceb.childCmdCh = childCmdCh
	ceb.childDoneCh = doneCh
	ceb.childSigCh = sigCh

	return nil
}

// markChildCmdReady will allow watchChildCmd to begin executing commands.
// This should be called once and should not be called concurrently.
func (ceb *CEB) markChildCmdReady() {
	ceb.setState(&ceb.stateChildReady, true)
}

// watchChildCmd should be started in a goroutine. This will run the child
// command sent on the channel. This ensures only one child command is run
// at a time.
func (ceb *CEB) watchChildCmd(
	ctx context.Context,
	cmdCh <-chan *exec.Cmd,
	doneCh chan<- error,
	sigCh <-chan os.Signal,
) {
	// We always close the done channel when we exit so that callers can
	// detect this and also exit.
	defer close(doneCh)

	log := ceb.logger.Named("child")

	// We need to wait for stateChildReady. This is set by ceb/init.go
	// or ceb/config.go when we're ready to begin processing. We have to do
	// this because we want to try to connect to the server to get our initial
	// config values before executing our child command. But if that fails, we
	// execute the child command anyways.
	log.Debug("waiting for stateChildReady to flip to true")
	if ceb.waitState(&ceb.stateChildReady, true) {
		// Early exit request
		log.Warn("exit state received before child was ready to start")
		return
	}

	log.Debug("starting child command watch loop")
	var currentCh <-chan error
	var currentCmd *exec.Cmd
	for {
		select {
		case <-ctx.Done():
			log.Warn("request to exit")

			// If we have a child, we need to exit that first.
			if currentCh != nil {
				log.Info("terminating current child process")
				err := ceb.termChildCmd(log, currentCmd, currentCh, false, false)
				log.Info("child process termination result", "err", err)
			}

			return

		case sig := <-sigCh:
			if currentCmd != nil {
				err := currentCmd.Process.Signal(sig)
				if err != nil {
					log.Error("error sending signal to process",
						"error", err,
						"signal", sig,
						"pid", currentCmd.Process.Pid,
					)
				} else {
					log.Info("sent signal to process",
						"signal", sig,
						"pid", currentCmd.Process.Pid,
					)
				}
			}
		case cmd := <-cmdCh:
			log.Debug("child command received")

			// If we have an existing process, we need to exit that first.
			if currentCh != nil {
				log.Info("terminating current child process for restart")
				err := ceb.termChildCmd(log, currentCmd, currentCh, false, false)
				if err != nil {
					// In the event terminating the child fails, we exit
					// the whole CEB because we can't be sure we're not just
					// going to fork bomb ourselves.
					log.Info("child process termination error", "err", err)
					doneCh <- err
					return
				}
			}

			// Drain the cmdCh. During graceful termination as well as initial
			// ready state waiting, we may accumulate a small buffer of command
			// changes. Let's drain that.
		CMD_DRAIN:
			for {
				select {
				case cmd = <-cmdCh:
				default:
					break CMD_DRAIN
				}
			}

			// Start our new command.
			currentCh = ceb.startChildCmd(log, cmd)
			currentCmd = cmd

		case err := <-currentCh:
			// Our child process exited, we're done. This function does not
			// restart any crashing child process, it only restarts if there
			// are changes to configuration. We assume a higher level process
			// such as a scheduler handles full crash restarts.
			log.Info("child process exited on its own", "err", err)
			doneCh <- err
			return
		}
	}
}

// startChildCmd starts the child process, and waits for completion by sending an
// error along the channel.
func (ceb *CEB) startChildCmd(
	log hclog.Logger,
	cmd *exec.Cmd,
) <-chan error {
	ch := make(chan error, 1)

	// Start our subprocess
	log = log.With(
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

// termChildCmd terminates the child command.
//
// If force is set to true, this will send a SIGKILL.
//
// If force is false, this will send a SIGTERM and wait up to 30 seconds
// for the child process to gracefully exit. If the process does not gracefully
// exit in 30 seconds, we will send a SIGKILL.
func (ceb *CEB) termChildCmd(
	log hclog.Logger,
	cmd *exec.Cmd,
	childErrCh <-chan error, // error channel from startChildCmd
	force bool,
	returnExitErr bool, // if true, doesn't treat exec.ExitError as nil
) error {
	log = log.With("pid", cmd.Process.Pid)

	// If we're not forcing, try a SIGTERM first.
	if !force {
		log.Debug("sending SIGTERM")
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Warn("error sending SIGTERM, will proceed to SIGKILL", "err", err)
		} else {
			log.Debug("SIGTERM sent, waiting for child process to end or timeout")
			select {
			case err := <-childErrCh:
				// If we got an exit error then everything worked propertly so
				// we just set error to nil.
				if _, ok := err.(*exec.ExitError); ok && !returnExitErr {
					err = nil
				}

				// Child successfully exited.
				log.Info("child process exited", "wait_err", err)
				return err

			case <-time.After(30 * time.Second):
				// Timeout, fall through to SIGKILL
				log.Warn("graceful termination failed, will send SIGKILL")
			}
		}
	}

	// SIGKILL. We send the signal to the negative value of the pid so
	// that it goes to the entire process group, therefore killing all
	// grandchildren of our child process as well.
	log.Debug("sending SIGKILL")
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		log.Warn("error sending SIGKILL", "err", err)
		return err
	}

	// Wait for the process to die
	log.Debug("waiting for child process to end")
	err := <-childErrCh
	if _, ok := err.(*exec.ExitError); ok && !returnExitErr {
		// If we got an exit error then everything worked propertly so
		// we just set error to nil.
		err = nil
	}
	log.Debug("child process exited", "wait_err", err)
	return err
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

	// Create a new process group so we can kill this child and all its
	// grandchildren when the time comes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd, nil
}

// copyCmd creates a new copy of the given command.
func (ceb *CEB) copyCmd(cmd *exec.Cmd) *exec.Cmd {
	var new exec.Cmd
	new.Path = cmd.Path
	new.Args = cmd.Args // not a deep copy
	new.Env = cmd.Env   //not a deep copy
	new.Dir = cmd.Dir
	new.Stdin = cmd.Stdin
	new.Stdout = cmd.Stdout
	new.Stderr = cmd.Stderr
	return &new
}
