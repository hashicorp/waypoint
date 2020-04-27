package ceb

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	internalptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func (ceb *CEB) startExecGroup(es []*pb.EntrypointConfig_Exec) {
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
		go ceb.startExec(exec)
	}

	// Store our exec index
	ceb.execIdx = idx
}

func (ceb *CEB) startExec(execConfig *pb.EntrypointConfig_Exec) {
	log := ceb.logger.Named("exec").With("index", execConfig.Index)

	// Open the stream
	log.Info("starting exec stream", "args", execConfig.Args)
	client, err := ceb.client.EntrypointExecStream(ceb.context)
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
	}

	// Create our pipe for stdin so that we can send data
	stdinR, stdinW := io.Pipe()
	defer stdinW.Close()

	// Start our receive data loop
	respCh := make(chan *pb.EntrypointExecResponse)
	go func() {
		for {
			resp, err := client.Recv()
			if err != nil {
				// TODO
				return
			}

			respCh <- resp
		}
	}()

	// We need to modify our command so the input/output is all over gRPC
	cmd.Stdin = stdinR
	cmd.Stdout = ceb.execOutputWriter(client, pb.EntrypointExecRequest_Output_STDOUT)
	cmd.Stderr = ceb.execOutputWriter(client, pb.EntrypointExecRequest_Output_STDERR)

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

		// Start with a pty
		ptyFile, err = pty.StartWithSize(cmd, internalptypes.Winsize(ptyReq.WindowSize))
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
		}
		defer ptyFile.Close()

		// Copy stdin to the pty
		go io.Copy(ptyFile, stdin)
		go io.Copy(stdout, ptyFile)
	} else {
		if err := cmd.Start(); err != nil {
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
		case resp := <-respCh:
			switch event := resp.Event.(type) {
			case *pb.EntrypointExecResponse_Input:
				// Copy the input to stdin
				log.Trace("input received", "data", event.Input)
				io.Copy(stdinW, bytes.NewReader(event.Input))

			case *pb.EntrypointExecResponse_Winch:
				log.Debug("window size change event, changing")
				if err := pty.Setsize(ptyFile, internalptypes.Winsize(event.Winch)); err != nil {
					log.Warn("error changing window size, this doesn't quit the stream",
						"err", err)
				}
			}

		case err := <-cmdExitCh:
			var exitCode int
			if err != nil {
				if exiterr, ok := err.(*exec.ExitError); ok {
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						exitCode = status.ExitStatus()
					}
				} else {
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

			// Exit!
			return
		}
	}

}

func (ceb *CEB) execOutputWriter(
	client grpc.ClientStream,
	channel pb.EntrypointExecRequest_Output_Channel,
) io.Writer {
	return &grpc_net_conn.Conn{
		Stream:  client,
		Request: &pb.EntrypointExecRequest{},
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			req := msg.(*pb.EntrypointExecRequest)
			if req.Event == nil {
				req.Event = &pb.EntrypointExecRequest_Output_{
					Output: &pb.EntrypointExecRequest_Output{
						Channel: channel,
					},
				}
			}

			return &req.Event.(*pb.EntrypointExecRequest_Output_).Output.Data
		}),
	}
}
