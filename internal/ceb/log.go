package ceb

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func (ceb *CEB) initLogStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("log")

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// Open our log stream
	log.Debug("connecting to log stream")
	client, err := ceb.client.EntrypointLogStream(ctx)
	if err != nil {
		return status.Errorf(codes.Aborted,
			"failed to open a log stream: %s", err)
	}
	ceb.cleanup(func() { client.CloseAndRecv() })
	log.Trace("log stream connected")

	// Set our output for the command. We use a multiwriter so that we
	// can always send the out/err back to the normal channels so that
	// users can see it.
	ceb.childCmd.Stdout = io.MultiWriter(w, ceb.childCmd.Stdout)
	ceb.childCmd.Stderr = io.MultiWriter(w, ceb.childCmd.Stderr)

	// Start our goroutine that'll read from the pipe.
	// TODO(mitchellh): lots of error handling needs to be added here
	//   - server errors, closing
	//   - reader errors, EOF
	//   - on server error we should keep reading from the pipe otherwise
	//     the child may get a broken pipe error
	go func() {
		defer r.Close()
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadString('\n')
			if err != nil {
				log.Warn("error reading logs", "error", err)
				return
			}

			log.Trace("sending line", "line", line)
			err = client.Send(&pb.EntrypointLogBatch{
				InstanceId: ceb.id,
				Lines: []*pb.LogBatch_Entry{
					&pb.LogBatch_Entry{
						Timestamp: ptypes.TimestampNow(),
						Line:      line,
					},
				},
			})
			if err != nil {
				log.Warn("error sending logs", "error", err)
				return
			}
		}
	}()

	return nil
}
