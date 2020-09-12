package ceb

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) initLogStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("log")

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// Set our output for the command. We use a multiwriter so that we
	// can always send the out/err back to the normal channels so that
	// users can see it.
	ceb.childCmd.Stdout = io.MultiWriter(w, ceb.childCmd.Stdout)
	ceb.childCmd.Stderr = io.MultiWriter(w, ceb.childCmd.Stderr)

	// We need to start a goroutine to read from our pipe. If we don't
	// read from the pipe the child command will get a SIGPIPE and could
	// exit/crash if it doesn't handle it. So even if we don't have a
	// connection to the server, we need to be draining the pipe.
	entryCh := make(chan *pb.LogBatch_Entry, 30)
	go func() {
		defer r.Close()
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadString('\n')
			if err != nil {
				log.Error("error reading logs", "error", err)
				return
			}

			if log.IsTrace() {
				log.Trace("sending line", "line", line[:len(line)-1])
			}

			entry := &pb.LogBatch_Entry{
				Timestamp: ptypes.TimestampNow(),
				Line:      line,
			}

			// Send the entry. We never block here because blocking the
			// pipe is worse. The channel is buffered to help with this.
			select {
			case entryCh <- entry:
			default:
			}
		}
	}()

	// Start up our server stream
	if err := ceb.initLogStreamSender(log, ctx, entryCh); err != nil {
		return err
	}

	return nil
}

func (ceb *CEB) initLogStreamSender(
	log hclog.Logger,
	ctx context.Context,
	entryCh <-chan *pb.LogBatch_Entry,
) error {
	// Open our log stream
	log.Debug("connecting to log stream")
	client, err := ceb.client.EntrypointLogStream(ctx, grpc.WaitForReady(true))
	if err != nil {
		return status.Errorf(codes.Aborted,
			"failed to open a log stream: %s", err)
	}
	ceb.cleanup(func() { client.CloseAndRecv() })
	log.Trace("log stream connected")

	// NOTE(mitchellh): Lots of improvements we can make here one day:
	//   - we can coalesce entryCh receives to send less log updates
	//   - during reconnect we can buffer entryCh receives
	go func() {
		for {
			var entry *pb.LogBatch_Entry
			select {
			case <-ctx.Done():
				return

			case entry = <-entryCh:
			}

			err := client.Send(&pb.EntrypointLogBatch{
				InstanceId: ceb.id,
				Lines:      []*pb.LogBatch_Entry{entry},
			})
			if err == io.EOF || status.Code(err) == codes.Unavailable {
				log.Error("log stream disconnected from server, attempting reconnect")
				err = ceb.initLogStreamSender(log, ctx, entryCh)
				if err == nil {
					return
				}
			}
			if err != nil {
				log.Warn("error sending logs", "error", err)
				return
			}
		}
	}()

	return nil
}
