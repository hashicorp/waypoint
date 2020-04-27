package execclient

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"
	"golang.org/x/crypto/ssh/terminal"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	internalptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

type Client struct {
	Context      context.Context
	Client       pb.WaypointClient
	DeploymentId string
	Args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

func (c *Client) Run() (int, error) {
	// Determine if we should allocate a pty. If we should, we need to send
	// along a TERM value to the remote end that matches our own.
	var ptyReq *pb.ExecStreamRequest_PTY
	var ptyF *os.File
	if f, ok := c.Stdout.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
		ptyF = f
		ws, err := pty.GetsizeFull(ptyF)
		if err != nil {
			return 0, err
		}

		ptyReq = &pb.ExecStreamRequest_PTY{
			Enable:     true,
			Term:       os.Getenv("TERM"),
			WindowSize: internalptypes.WinsizeProto(ws),
		}

		// We need to go into raw mode with stdin
		if f, ok := c.Stdin.(*os.File); ok {
			oldState, err := terminal.MakeRaw(int(f.Fd()))
			if err != nil {
				return 0, err
			}
			defer func() { _ = terminal.Restore(int(f.Fd()), oldState) }()
		}
	}

	// Start our exec stream
	client, err := c.Client.StartExecStream(c.Context)
	if err != nil {
		return 0, err
	}

	// Send the start event
	if err := client.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				DeploymentId: c.DeploymentId,
				Args:         c.Args,
				Pty:          ptyReq,
			},
		},
	}); err != nil {
		return 0, err
	}

	// Create the context that we'll listen to that lets us cancel our
	// extra goroutines here.
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	// Build our connection. We only build the stdin sending side because
	// we can receive other message types from our recv.
	go io.Copy(&grpc_net_conn.Conn{
		Stream:  client,
		Request: &pb.ExecStreamRequest{},
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			req := msg.(*pb.ExecStreamRequest)
			if req.Event == nil {
				req.Event = &pb.ExecStreamRequest_Input_{
					Input: &pb.ExecStreamRequest_Input{},
				}
			}

			return &req.Event.(*pb.ExecStreamRequest_Input_).Input.Data
		}),
	}, c.Stdin)

	// Add our recv blocker that sends data
	recvCh := make(chan *pb.ExecStreamResponse)
	go func() {
		defer cancel()
		for {
			resp, err := client.Recv()
			if err != nil {
				// TODO: log
				return
			}

			recvCh <- resp
		}
	}()

	// Listen for window change events
	winchCh := make(chan os.Signal, 1)
	signal.Notify(winchCh, syscall.SIGWINCH)
	defer signal.Stop(winchCh)

	// Loop for data
	for {
		select {
		case resp := <-recvCh:
			switch event := resp.Event.(type) {
			case *pb.ExecStreamResponse_Output_:
				// TODO: stderr
				out := c.Stdout
				io.Copy(out, bytes.NewReader(event.Output.Data))

			case *pb.ExecStreamResponse_Exit_:
				return int(event.Exit.Code), nil
			}

			/*
				TODO: send this once the server side handles this
					case <-winchCh:
						// Window change, send new size
						ws, err := pty.GetsizeFull(ptyF)
						if err != nil {
							// Ignore errors
							continue
						}

						// Send the new window size
						if err := client.Send(&pb.ExecStreamRequest{
							Event: &pb.ExecStreamRequest_Winch{
								Winch: internalptypes.WinsizeProto(ws),
							},
						}); err != nil {
							// Ignore this error
							continue
						}
			*/

		case <-ctx.Done():
			return 1, nil
		}
	}
}
