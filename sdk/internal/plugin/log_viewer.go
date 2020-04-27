package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// logViewerClient is an implementation of component.LogViewer over gRPC.
type logViewerClient struct {
	Client pb.LogViewerClient
	Logger hclog.Logger
}

func (c *logViewerClient) NextLogBatch(ctx context.Context) ([]component.LogEvent, error) {
	resp, err := c.Client.NextLogBatch(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var events []component.LogEvent
	for _, event := range resp.Events {
		ts, err := ptypes.Timestamp(event.Timestamp)
		if err != nil {
			// This is only possible with a poorly behaved plugin because
			// the plugin must've converted a timestamp propertly to this format.
			return nil, err
		}

		events = append(events, component.LogEvent{
			Partition: event.Partition,
			Timestamp: ts,
			Message:   event.Contents,
		})
	}

	return events, nil
}

// logViewerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type logViewerServer struct {
	Impl   component.LogViewer
	Logger hclog.Logger
}

func (s *logViewerServer) NextLogBatch(
	ctx context.Context,
	args *empty.Empty,
) (*pb.Logs_NextBatchResp, error) {
	events, err := s.Impl.NextLogBatch(ctx)
	if err != nil {
		return nil, err
	}

	var result pb.Logs_NextBatchResp
	for _, event := range events {
		ts, err := ptypes.TimestampProto(event.Timestamp)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"log entry has invalid timestamp: %s",
				err)
		}

		result.Events = append(result.Events, &pb.Logs_Event{
			Partition: event.Partition,
			Timestamp: ts,
			Contents:  event.Message,
		})
	}

	return &result, nil
}

var (
	_ pb.LogViewerServer  = (*logViewerServer)(nil)
	_ component.LogViewer = (*logViewerClient)(nil)
)
