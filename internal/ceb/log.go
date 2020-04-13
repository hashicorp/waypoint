package ceb

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func (ceb *CEB) initLogStream(ctx context.Context, cfg *config) error {
	// Open our log stream
	ceb.logger.Debug("connecting to log stream")
	client, err := ceb.client.EntrypointLogStream(ctx)
	if err != nil {
		return status.Errorf(codes.Aborted,
			"failed to open a log stream: %s", err)
	}
	ceb.cleanup(func() { client.CloseAndRecv() })
	ceb.logger.Trace("log stream connected")

	// Create our request structure which always has the
	req := &pb.EntrypointLogBatch{InstanceId: ceb.id}

	// Create our conn
	conn := &grpc_net_conn.Conn{
		Stream:  client,
		Request: req,
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			return &msg.(*pb.EntrypointLogBatch).Data
		}),
	}

	// Set our output for the command
	ceb.childCmd.Stdout = conn
	ceb.childCmd.Stderr = conn

	return nil
}
