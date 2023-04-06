// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ceb

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint/internal/pkg/gatedwriter"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// initSystemLogger initializes ceb.logger and sets up all the fields
// for streaming system logs to the Waypoint server.
func (ceb *CEB) initSystemLogger() {
	// Create an intercept logger with our default options. This will
	// behave just like hclog.L() (which we use at the time of writing)
	// and let us register additional sinks for streaming and so on.
	opts := *hclog.DefaultOptions
	opts.Name = "entrypoint"
	opts.Level = hclog.Debug
	intercept := hclog.NewInterceptLogger(&opts)
	nonintercept := hclog.New(&opts)
	ceb.logger = intercept

	// Set our initial log level
	if v := os.Getenv(envLogLevel); v != "" {
		level := hclog.LevelFromString(v)
		if level == hclog.NoLevel {
			ceb.logger.Warn("log level provided in env var is invalid", v)
		} else {
			ceb.logger.SetLevel(level)
		}
	}

	// Create our reader/writer that will send to the server log stream.
	r, w := io.Pipe()

	// We set the writer as a gated writer so that we can buffer all our
	// log messages prior to attempting a log stream connection. Once we
	// attempt a log stream connection we flush.
	ceb.logGatedWriter = gatedwriter.NewWriter(w)

	// Create our channel where we can send logs to. We allow some buffering.
	// We then start a goroutine that'll read the logs from this pipe and
	// send them to our channel that will eventually get flushed to the server.
	entryCh := make(chan *pb.LogBatch_Entry, 30)
	ceb.logCh = entryCh
	go ceb.logReader(
		nonintercept.Named("system_log_streamer"),
		r,
		pb.LogBatch_Entry_ENTRYPOINT,
	)

	// Register a sink that will go to the log stream.
	intercept.RegisterSink(hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Name:        "entrypoint",
		Level:       hclog.Info,
		Output:      ceb.logGatedWriter,
		Color:       hclog.ColorOff,
		DisableTime: true, // because we calculate it ourselves for streaming
		Exclude:     ceb.logStreamExclude,
	}))
}

func (ceb *CEB) logStreamExclude(level hclog.Level, msg string, args ...interface{}) bool {
	if level == hclog.Info {
		// We want to exclude some Horizon logs. We don't set the level lower
		// because we want the root logs to show up fine we just don't want
		// to stream them.
		return strings.Contains(msg, "request started") ||
			strings.Contains(msg, "request ended")
	}

	return false
}

func (ceb *CEB) initLogStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("log")

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// Set our output for the command. We use a multiwriter so that we
	// can always send the out/err back to the normal channels so that
	// users can see it.
	ceb.childCmdBase.Stdout = io.MultiWriter(w, ceb.childCmdBase.Stdout)
	ceb.childCmdBase.Stderr = io.MultiWriter(w, ceb.childCmdBase.Stderr)

	// We need to start a goroutine to read from our pipe. If we don't
	// read from the pipe the child command will get a SIGPIPE and could
	// exit/crash if it doesn't handle it. So even if we don't have a
	// connection to the server, we need to be draining the pipe.
	go ceb.logReader(log, r, pb.LogBatch_Entry_APP)

	// Start up our server stream. We do this in a goroutine cause we don't
	// want to block the child command startup on it.
	go ceb.initLogStreamSender(log, ctx)

	return nil
}

func (ceb *CEB) initLogStreamSender(
	log hclog.Logger,
	ctx context.Context,
) error {
	// wait for initial server connection
	serverClient := ceb.waitClient()
	if serverClient == nil {
		return ctx.Err()
	}

	// Open our log stream
	log.Debug("connecting to log stream")
	client, err := serverClient.EntrypointLogStream(ctx, grpc.WaitForReady(true))
	if err != nil {
		return status.Errorf(codes.Aborted,
			"failed to open a log stream: %s", err)
	}
	ceb.cleanup(func() { client.CloseAndRecv() })
	log.Trace("log stream connected")

	// NOTE(mitchellh): Lots of improvements we can make here one day:
	//   - we can coalesce channel receives to send less log updates
	//   - during reconnect we can buffer channel receives
	go func() {
		// Wait for the state that our config stream is connected. Logs are
		// not allowed (and dropped by the server) until we're connected so
		// this lets us get all our startup logs in safely.
		if ceb.waitState(&ceb.stateConfig, true) {
			// Early exit request
			return
		}

		// dropCh is non-nil if the server noted it doesn't support
		// log streaming. We periodically retry to send logs because
		// it is possible the server is reconfigured and restarted to
		// support logs.
		var dropCh <-chan time.Time

		for {
			var entry *pb.LogBatch_Entry
			select {
			case <-ctx.Done():
				return

			case entry = <-ceb.logCh:

			case <-dropCh:
				// Reset dropCh and try to send logs again.
				dropCh = nil
			}

			// If we're dropping logs because the server is rejecting our
			// log stream, then do nothing.
			if dropCh != nil {
				continue
			}

			// Nothing to send? Do nothing. This is possible in the dropCh
			// case in the select.
			if entry == nil {
				continue
			}

			err := client.Send(&pb.EntrypointLogBatch{
				InstanceId: ceb.id,
				Lines:      []*pb.LogBatch_Entry{entry},
			})

			// Server returns this error when it doesn't support application
			// log streaming. This is possible depending on the configuration
			// or state backend.
			if status.Code(err) == codes.Unimplemented {
				log.Warn("log stream unimplemented on server, dropping logs")
				err = nil
				dropCh = time.After(5 * time.Minute)
			}

			if err == io.EOF || status.Code(err) == codes.Unavailable {
				log.Debug("log stream disconnected from server, attempting reconnect",
					"err", err)
				err = ceb.initLogStreamSender(log, ctx)
				if err == nil {
					return
				}

				log.Error("log stream disconnected from server, reconnect failed",
					"err", err)
			}
			if err != nil {
				log.Warn("error sending logs", "error", err)
				return
			}
		}
	}()

	// Open the gated writer since we should now start consuming logs.
	ceb.logGatedWriter.Flush()

	return nil
}

// logReader reads lines from r and sends them to ceb.logCh with the
// proper envelope (pb.LogBatch_Entry). This should be started in a goroutine.
func (ceb *CEB) logReader(
	log hclog.Logger,
	r io.ReadCloser,
	src pb.LogBatch_Entry_Source,
) {
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
			Source:    src,
			Timestamp: timestamppb.Now(),
			Line:      line,
		}

		// Send the entry. We never block here because blocking the
		// pipe is worse. The channel is buffered to help with this.
		select {
		case ceb.logCh <- entry:
		default:
		}
	}
}
