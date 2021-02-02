package execclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"

	"github.com/containerd/console"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	grpc_net_conn "github.com/mitchellh/go-grpc-net-conn"
	sshterm "golang.org/x/crypto/ssh/terminal"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type Client struct {
	Logger  hclog.Logger
	UI      terminal.UI
	Context context.Context
	Client  pb.WaypointClient
	Args    []string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer

	// Either DeploymentId or InstanceId have to be set. If both are set, then
	// InstanceId takes priority.
	//
	// These identify a deployment that is used to search for an instance on
	// server side. We target a specific deployment so that the exec session
	// enters the correct application code.
	DeploymentId  string
	DeploymentSeq uint64

	// If set, will cause the server to connect to this specific
	// instance.
	InstanceId string
}

func (c *Client) Run() (int, error) {
	// Determine if we should allocate a pty. If we should, we need to send
	// along a TERM value to the remote end that matches our own.
	var ptyReq *pb.ExecStreamRequest_PTY
	var ptyF *os.File
	var status terminal.Status

	if f, ok := c.Stdout.(*os.File); ok && sshterm.IsTerminal(int(f.Fd())) {
		status = c.UI.Status()
		defer status.Close()
		status.Update(fmt.Sprintf("Connecting to deployment v%d...", c.DeploymentSeq))

		ptyF = f
		c, err := console.ConsoleFromFile(ptyF)
		if err != nil {
			return 0, err
		}

		sz, err := c.Size()
		c = nil
		if err != nil {
			return 0, err
		}

		ptyReq = &pb.ExecStreamRequest_PTY{
			Enable: true,
			Term:   os.Getenv("TERM"),
			WindowSize: &pb.ExecStreamRequest_WindowSize{
				Rows:   int32(sz.Height),
				Cols:   int32(sz.Width),
				Height: int32(sz.Height),
				Width:  int32(sz.Width),
			},
		}
	}

	// Start our exec stream
	client, err := c.Client.StartExecStream(c.Context)
	if err != nil {
		return 0, err
	}

	defer client.CloseSend()

	if status != nil {
		status.Update("Initializing session...")
	}

	start := &pb.ExecStreamRequest_Start{
		Args: c.Args,
		Pty:  ptyReq,
	}

	if c.InstanceId != "" {
		start.Target = &pb.ExecStreamRequest_Start_InstanceId{
			InstanceId: c.InstanceId,
		}
	} else {
		start.Target = &pb.ExecStreamRequest_Start_DeploymentId{
			DeploymentId: c.DeploymentId,
		}
	}

	// Send the start event
	if err := client.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: start,
		},
	}); err != nil {
		return 0, err
	}

	if status != nil {
		status.Update("Waiting for instance assignment...")
	}

	// Receive our open message. If this fails then we weren't assigned.
	resp, err := client.Recv()
	if err != nil {
		return 1, err
	}
	if _, ok := resp.Event.(*pb.ExecStreamResponse_Open_); !ok {
		return 1, fmt.Errorf("internal protocol error: unexpected opening message")
	}

	if ptyF != nil {
		status.Close()
		c.UI.Output("Connected to deployment v%d", c.DeploymentSeq, terminal.WithSuccessStyle())
	}

	// Close our UI if we can
	if closer, ok := c.UI.(io.Closer); ok {
		closer.Close()
	}

	if ptyF != nil {
		// We need to go into raw mode with stdin
		if f, ok := c.Stdout.(*os.File); ok {
			oldState, err := sshterm.MakeRaw(int(f.Fd()))
			if err != nil {
				return 0, err
			}
			defer sshterm.Restore(int(f.Fd()), oldState)
		}

		fmt.Fprintf(c.Stdout, "\r")
	}

	// Create the context that we'll listen to that lets us cancel our
	// extra goroutines here.
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	// If we have interactive stdin WITHOUT a tty output, we treat stdin
	// as if its closed. This allows commands such as "ls" to work without
	// a terminal otherwise they'll hang open forever on stdin. Conversely,
	// we allow interactive stdin with a TTY, and we also allow non-interactive
	// stdin always since that'll end in an EOF at some point.
	copyStdin := true
	if stdinF, ok := c.Stdin.(*os.File); ok {
		fi, err := stdinF.Stat()
		if err != nil {
			return 0, err
		}

		if fi.Mode()&os.ModeCharDevice != 0 && ptyF == nil {
			// Stdin is from a terminal but we don't have a pty output.
			// In this case, we treat stdin as if its closed.
			c.Logger.Info("terminal stdin without a pty, not using input for command")
			copyStdin = false
		}
	}

	// We need to lock access to our stream writer since it is unsafe in
	// gRPC to concurrently send data.
	var streamLock sync.Mutex

	// Build our connection. We only build the stdin sending side because
	// we can receive other message types from our recv.
	go func() {
		input := &EscapeWatcher{Cancel: cancel, Input: c.Stdin}

		// If we're copying stdin then start that copy process.
		if copyStdin {
			io.Copy(&grpc_net_conn.Conn{
				Stream:       client,
				Request:      &pb.ExecStreamRequest{},
				ResponseLock: &streamLock,
				Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
					req := msg.(*pb.ExecStreamRequest)
					if req.Event == nil {
						req.Event = &pb.ExecStreamRequest_Input_{
							Input: &pb.ExecStreamRequest_Input{},
						}
					}

					return &req.Event.(*pb.ExecStreamRequest_Input_).Input.Data
				}),
			}, input)
		} else {
			// If we're NOT copying, we still start a copy to discard
			// in the background so that we still handle escape sequences.
			go io.Copy(ioutil.Discard, input)
		}

		// After the copy ends, no matter what we send an EOF to the
		// remote end because there will be no more input.
		c.Logger.Debug("stdin closed, sending input EOF event")
		streamLock.Lock()
		defer streamLock.Unlock()
		if err := client.Send(&pb.ExecStreamRequest{
			Event: &pb.ExecStreamRequest_InputEof{
				InputEof: &empty.Empty{},
			},
		}); err != nil {
			c.Logger.Warn("error sending InputEOF event", "err", err)
		}
	}()

	// Add our recv blocker that sends data
	recvCh := make(chan *pb.ExecStreamResponse)
	go func() {
		defer cancel()
		for {
			resp, err := client.Recv()
			if err != nil {
				if err != io.EOF {
					c.Logger.Error("receive error", "err", err)
				}
				return
			}

			recvCh <- resp
		}
	}()

	// Listen for window change events
	winchCh := make(chan os.Signal, 1)
	registerSigwinch(winchCh)
	defer signal.Stop(winchCh)

	// Loop for data
	for {
		select {
		case resp := <-recvCh:
			switch event := resp.Event.(type) {
			case *pb.ExecStreamResponse_Output_:
				switch event.Output.Channel {
				case pb.ExecStreamResponse_Output_STDOUT:
					io.Copy(c.Stdout, bytes.NewReader(event.Output.Data))
				case pb.ExecStreamResponse_Output_STDERR:
					io.Copy(c.Stderr, bytes.NewReader(event.Output.Data))
				}
			case *pb.ExecStreamResponse_Exit_:
				return int(event.Exit.Code), nil

			default:
				c.Logger.Warn("unknown event type",
					"type", fmt.Sprintf("%T", resp.Event))
			}

		case <-winchCh:
			// Window change, send new size
			c, err := console.ConsoleFromFile(ptyF)
			if err != nil {
				continue
			}

			sz, err := c.Size()
			if err != nil {
				continue
			}

			// Send the new window size
			streamLock.Lock()
			err = client.Send(&pb.ExecStreamRequest{
				Event: &pb.ExecStreamRequest_Winch{
					Winch: &pb.ExecStreamRequest_WindowSize{
						Rows:   int32(sz.Height),
						Cols:   int32(sz.Width),
						Height: int32(sz.Height),
						Width:  int32(sz.Width),
					},
				},
			})
			streamLock.Unlock()
			if err != nil {
				// Ignore this error
				continue
			}

		case <-ctx.Done():
			return 1, nil
		}
	}
}
