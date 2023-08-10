// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/ceb/execwriter"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (ceb *CEB) startExecGroup(es []*pb.EntrypointConfig_Exec, env []string) {
	// If exec is disabled, log. This should never happen because we advertise
	// disabled exec to the server, and the server should not assign us any
	// exec sessions. However, we don't want to explicitly trust the server
	// so we also safeguard here that we do not exec if we've disabled it.
	if ceb.execDisable {
		ceb.logger.Warn("startExecGroup called but disableExec is true. This should not happen.")
		return
	}

	idx := ceb.execIdx
	for _, exec := range es {
		// Ignore exec sessions we already have
		if exec.Index <= ceb.execIdx {
			continue
		}

		// Track the highest index
		if exec.Index > idx {
			idx = exec.Index
		}

		// Start our session
		go ceb.startExec(exec, env)
	}

	// Store our exec index
	ceb.execIdx = idx
}

func (ceb *CEB) startExec(execConfig *pb.EntrypointConfig_Exec, env []string) {
	log := ceb.logger.Named("exec").With("index", execConfig.Index)

	// wait for initial server connection
	serverClient := ceb.waitClient()
	if serverClient == nil {
		log.Warn("nil client, can't execute")
	}

	// Open the stream
	log.Info("starting exec stream", "args", execConfig.Args)
	client, err := serverClient.EntrypointExecStream(ceb.context)
	if err != nil {
		log.Warn("error opening exec stream", "err", err)
		return
	}
	defer client.CloseSend()

	// Send our open message
	log.Trace("sending open message")
	if err := client.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: ceb.id,
				Index:      execConfig.Index,
			},
		},
	}); err != nil {
		log.Warn("error opening exec stream", "err", err)
		return
	}

	// Build our command
	cmd, err := ceb.buildCmd(ceb.context, execConfig.Args)
	if err != nil {
		log.Warn("error building exec command", "err", err)
		st, ok := status.FromError(err)
		if !ok {
			st = status.New(codes.Unknown, err.Error())
		}

		if err := client.Send(&pb.EntrypointExecRequest{
			Event: &pb.EntrypointExecRequest_Error_{
				Error: &pb.EntrypointExecRequest_Error{
					Error: st.Proto(),
				},
			},
		}); err != nil {
			log.Warn("error sending error message", "err", err)
			return
		}

		return
	}

	// Set our environment variables from our `waypoint config` settings.
	cmd.Env = append(cmd.Env, env...)

	// Create our pipe for stdin so that we can send data
	stdinR, stdinW := io.Pipe()
	defer stdinW.Close()

	// Start our receive data loop
	respCh := make(chan *pb.EntrypointExecResponse)
	go func() {
		defer close(respCh)

		for {
			resp, err := client.Recv()
			if err != nil {
				if err == io.EOF {
					log.Info("exec stream ended by client")
				} else {
					log.Warn("error receiving from server stream", "err", err)
				}
				return
			}

			respCh <- resp
		}
	}()

	// We need to modify our command so the input/output is all over gRPC
	cmd.Stdin = stdinR
	cmd.Stdout = execwriter.Writer(client, pb.EntrypointExecRequest_Output_STDOUT)
	cmd.Stderr = execwriter.Writer(client, pb.EntrypointExecRequest_Output_STDERR)

	// PTY
	var ptyFile *os.File
	if ptyReq := execConfig.Pty; ptyReq != nil && ptyReq.Enable {
		log.Info("pty requested, allocating a pty")

		// If we're setting a pty we'll be overriding our stdin/out/err
		// so we need to get access to the original gRPC writers so we can
		// copy later.
		stdin := cmd.Stdin
		stdout := cmd.Stdout

		// We need to nil these so they get set to the pty below
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil

		// Set our TERM value
		if ptyReq.Term != "" {
			cmd.Env = append(cmd.Env, "TERM="+ptyReq.Term)
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
			Rows: uint16(ptyReq.WindowSize.Rows),
			Cols: uint16(ptyReq.WindowSize.Cols),
			X:    uint16(ptyReq.WindowSize.Width),
			Y:    uint16(ptyReq.WindowSize.Height),
		})
		if err != nil {
			log.Warn("error starting pty", "err", err)
			st, ok := status.FromError(err)
			if !ok {
				st = status.New(codes.Unknown, err.Error())
			}

			if err := client.Send(&pb.EntrypointExecRequest{
				Event: &pb.EntrypointExecRequest_Error_{
					Error: &pb.EntrypointExecRequest_Error{
						Error: st.Proto(),
					},
				},
			}); err != nil {
				log.Warn("error sending error message", "err", err)
				return
			}
		}
		defer ptyFile.Close()

		// Copy stdin to the pty
		go io.Copy(ptyFile, stdin)
		go io.Copy(stdout, ptyFile)
	} else {
		if err := cmd.Start(); err != nil {
			log.Warn("error starting exec command", "err", err)
			st, ok := status.FromError(err)
			if !ok {
				st = status.New(codes.Unknown, err.Error())
			}

			if err := client.Send(&pb.EntrypointExecRequest{
				Event: &pb.EntrypointExecRequest_Error_{
					Error: &pb.EntrypointExecRequest_Error{
						Error: st.Proto(),
					},
				},
			}); err != nil {
				log.Warn("error sending error message", "err", err)
				return
			}
		}
	}

	// Wait for the command to exit in a goroutine so we can handle
	// concurrent events happening below.
	cmdExitCh := make(chan error, 1)
	go func() {
		cmdExitCh <- cmd.Wait()
	}()

	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				// channel is closed, we should terminate our child process.
				log.Info("exec recv stream closed, killing child")

				// terminate the child process first, the last "true" argument
				// ensures we get the ExitError result back rather than nil.
				// We need that to set the proper exec session state.
				err := ceb.termChildCmd(log, cmd, cmdExitCh, true, true)

				// send final exec session updates. This will log if any errors occur.
				ceb.handleExecProcessExit(log, client, err)
				return
			}

			switch event := resp.Event.(type) {
			case *pb.EntrypointExecResponse_Input:
				// Copy the input to stdin
				log.Trace("input received", "data", event.Input)
				io.Copy(stdinW, bytes.NewReader(event.Input))

			case *pb.EntrypointExecResponse_InputEof:
				log.Trace("input EOF, closing stdin")
				stdinW.Close()

			case *pb.EntrypointExecResponse_Winch:
				log.Debug("window size change event, changing")

				sz := pty.Winsize{
					Rows: uint16(event.Winch.Rows),
					Cols: uint16(event.Winch.Cols),
					X:    uint16(event.Winch.Width),
					Y:    uint16(event.Winch.Height),
				}
				if err := pty.Setsize(ptyFile, &sz); err != nil {
					log.Warn("error changing window size, this doesn't quit the stream",
						"err", err)
				}
			}

		case err := <-cmdExitCh:
			ceb.handleExecProcessExit(log, client, err)

			// Exit!
			return
		}
	}

}

// handleExecProcessExit sends the final update to the server with the
// exec process exit information (such as exit code). This will log any
// errors rather than return since there is no reasonable way to handle them.
func (ceb *CEB) handleExecProcessExit(
	log hclog.Logger,
	client pb.Waypoint_EntrypointExecStreamClient,
	err error,
) {
	var exitCode int
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			log.Warn("non ExitError from Wait, process may be dangling",
				"err", err)

			// For the client, treat it as an erroneous exit.
			exitCode = 1
		}
	}

	// Send our exit code
	log.Info("exec stream exited", "code", exitCode)
	if err := client.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Exit_{
			Exit: &pb.EntrypointExecRequest_Exit{
				Code: int32(exitCode),
			},
		},
	}); err != nil {
		log.Warn("error sending exit message", "err", err)
	}
}
