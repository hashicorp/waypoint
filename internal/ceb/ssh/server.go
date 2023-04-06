// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/hashicorp/go-hclog"
	gossh "golang.org/x/crypto/ssh"
)

// RunExecSSHServer starts up an ssh server on the given port +sport+. The server
// will use +hostkey+ as the host key and only accepts a connection when the client
// has authenticated with +key+.
// The SSH server will run one command and then exit, as it's designed to only be used
// for one-shot commands via an exec plugin.
func RunExecSSHServer(
	ctx context.Context,
	logger hclog.Logger,
	sport string,
	hostkey gossh.Signer,
	key gossh.PublicKey,
) error {
	var server *ssh.Server

	port, err := strconv.Atoi(sport)
	if err != nil {
		return err
	}

	check := func(ctx ssh.Context, inputKey ssh.PublicKey) bool {
		if ssh.KeysEqual(inputKey, key) {
			return true
		}

		logger.Error("keys did not match")
		return false
	}

	logger.Info("starting ssh listener...")

	err = ssh.ListenAndServe(
		fmt.Sprintf(":%d", port),
		createHandler(ctx, logger, &server),
		ssh.Option(func(serv *ssh.Server) error {
			server = serv
			serv.PublicKeyHandler = check
			serv.AddHostKey(hostkey)
			return nil
		}),
	)

	if err != ssh.ErrServerClosed {
		return err
	}

	return nil
}

// This is pulled out to make it easier to test.
func createHandler(ctx context.Context, logger hclog.Logger, server **ssh.Server) ssh.Handler {
	return func(s ssh.Session) {
		// Build the subprocess

		args := s.Command()

		cmd := exec.CommandContext(s.Context(), args[0], args[1:]...)
		cmd.Env = s.Environ()

		var (
			ptyFile *os.File
			err     error
		)

		logger.Debug("executing command", "command", args)

		outDoneCh := make(chan struct{})
		if ptyInfo, winCh, ok := s.Pty(); ok {
			logger.Debug("running command in a PTY")

			// If we're setting a pty we'll be overriding our stdin/out/err
			// so we need to get access to the original gRPC writers so we can
			// copy later.
			stdin := s
			stdout := s

			// Set our TERM value
			if ptyInfo.Term != "" {
				cmd.Env = append(cmd.Env, "TERM="+ptyInfo.Term)
			}

			// pty.StartWithSize sets "setsid" which is mutually exclusive to
			// Setpgid. They both result in a new process group being created with
			// the process group ID equal to the PID, which is the behavior we
			// expect when terminating processes.
			if cmd.SysProcAttr != nil {
				cmd.SysProcAttr.Setpgid = false
			}

			// Start with a pty
			ptyFile, err = pty.StartWithSize(cmd, &pty.Winsize{
				X: uint16(ptyInfo.Window.Width),
				Y: uint16(ptyInfo.Window.Height),
			})
			if err != nil {
				fmt.Fprintf(s, "Error occured: %s\r\n", err)
				return
			}

			defer ptyFile.Close()

			// Copy stdin to the pty
			go func() {
				io.Copy(ptyFile, stdin)
				logger.Debug("ssh client closed stdin")
			}()

			go func() {
				defer close(outDoneCh)
				io.Copy(stdout, ptyFile)
				logger.Debug("command closed stdout")
			}()

			go func() {
				for {
					select {
					case <-s.Context().Done():
						return

					case win := <-winCh:
						sz := pty.Winsize{
							X: uint16(win.Width),
							Y: uint16(win.Height),
						}

						if err := pty.Setsize(ptyFile, &sz); err != nil {
							logger.Warn("error changing window size, this doesn't quit the stream",
								"err", err)
						}
					}
				}
			}()
		} else {
			logger.Debug("executing command without pty")

			// You might think "Hey, why make these pipes, can't we just
			// assign s to Stdin and Stdout directly?" Well originally this
			// code did that exact thing, BUUUUT it doesn't work right.
			// When Wait() attempts to cleanup the stdin go routine it launches
			// in the background, it ends up hanging because it's waiting for
			// s.Read to because finish but it's blocked inside SSH.
			// But if we make normal pipes and copy between them, everything is
			// fine, so here we are.

			stdin, err := cmd.StdinPipe()
			if err != nil {
				fmt.Fprintf(s, "Error occured: %s\r\n", err)
				return
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Fprintf(s, "Error occured: %s\r\n", err)
				return
			}

			cmd.Stderr = cmd.Stdout

			// Create a new process group so we can kill this child and all its
			// grandchildren when the time comes.
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			if err := cmd.Start(); err != nil {
				fmt.Fprintf(s, "Error occured: %s\r\n", err)
				return
			}

			go io.Copy(stdin, s)
			go func() {
				defer close(outDoneCh)
				io.Copy(s, stdout)
			}()
		}

		exitCh := make(chan error, 1)
		go func() {
			logger.Debug("waiting for command to finish")
			err := cmd.Wait()
			logger.Debug("command has finished", "error", err)
			exitCh <- err
		}()

		breakCh := make(chan bool, 1)
		s.Break(breakCh)

		signalsCh := make(chan ssh.Signal, 1)
		s.Signals(signalsCh)

		logger.Debug("waiting on events")

		for {
			select {
			case <-ctx.Done():
				logger.Debug("context done, aborting loop")
				cmd.Process.Kill()
				return
			case err := <-exitCh:
				logger.Debug("command has exited, waiting for output to complete")

				select {
				case <-outDoneCh:
					logger.Trace("output copy complete")
				case <-time.After(1 * time.Second):
					// We don't want to block on poorly behaved commands. They
					// should all finish copying relatively quickly.
					logger.Trace("output copy timeout, just forcing exit")
				}

				if err != nil {
					if exiterr, ok := err.(*exec.ExitError); ok {
						s.Exit(exiterr.ExitCode())
					} else {
						logger.Error("error waiting on command", "error", err)
						s.Exit(1)
					}
				} else {
					s.Exit(0)
				}

				if server != nil {
					go (*server).Shutdown(ctx)
				}
				return
			case <-breakCh:
				logger.Warn("break detected from client")
				cmd.Process.Signal(os.Interrupt)
			case sig := <-signalsCh:
				var signal os.Signal

				switch sig {
				case ssh.SIGABRT:
					signal = syscall.SIGABRT
				case ssh.SIGINT:
					signal = syscall.SIGINT
				case ssh.SIGKILL:
					signal = syscall.SIGKILL
				case ssh.SIGQUIT:
					signal = syscall.SIGQUIT
				case ssh.SIGUSR1:
					signal = syscall.SIGUSR1
				case ssh.SIGUSR2:
					signal = syscall.SIGUSR2
				}

				if signal != nil {
					cmd.Process.Signal(signal)
				}
			}
		}
	}
}
