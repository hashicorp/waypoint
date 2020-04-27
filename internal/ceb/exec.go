package ceb

import (
	"io"
	"os/exec"
	"syscall"

	//"github.com/creack/pty"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	// We need to modify our command so the input/output is all over gRPC
	cmd.Stdin = ceb.execInputReader(client)
	cmd.Stdout = ceb.execOutputWriter(client, pb.EntrypointExecRequest_Output_STDOUT)
	cmd.Stderr = ceb.execOutputWriter(client, pb.EntrypointExecRequest_Output_STDERR)

	// Run the command and wait for it to end
	exitCode := 0
	err = cmd.Run()
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
		return
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

func (ceb *CEB) execInputReader(
	client grpc.ClientStream,
) io.Reader {
	return &grpc_net_conn.Conn{
		Stream:   client,
		Response: &pb.EntrypointExecResponse{},
		Decode: grpc_net_conn.SimpleDecoder(func(msg proto.Message) *[]byte {
			return &msg.(*pb.EntrypointExecResponse).Data
		}),
	}
}
