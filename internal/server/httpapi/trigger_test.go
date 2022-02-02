package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/gen/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestHandleTrigger(t *testing.T) {
	require := require.New(t)

	// Get our gRPC server
	impl := &triggerImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	// Mock a request
	resp, err := http.Get(httpServer.URL + "/v1/trigger/123" + "?token=foo-bar-baz&stream=true")
	if err != nil {
		t.Errorf("failed to make http request: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code: %d, expected %d", resp.StatusCode, 200)
	}

	decoder := json.NewDecoder(resp.Body)
	var msg Message
	require.NoError(decoder.Decode(&msg))

	require.Equal(msg.ValueType, "TerminalEventLine")
	require.Equal(msg.Value, "testing")

	msg = Message{}
	require.NoError(decoder.Decode(&msg))

	require.Equal(msg.ValueType, "TerminalEventLine")
	require.Equal(msg.Value, "another one")

	msg = Message{}
	require.NoError(decoder.Decode(&msg))

	require.Equal(msg.ValueType, "TerminalEventStatus")
	require.Equal(msg.Value, "a status msg")
}

func TestHandleTrigger_BadFailures(t *testing.T) {
	// TODO: write a new impl for failures
}

type triggerImpl struct {
	sync.Mutex
	mocks.WaypointServer
}

type streamClientImpl struct {
	Recv []*pb.GetJobStreamResponse
	grpc.ClientStream
}

// RunTrigger mocks out a "good" RunTrigger execution and returns a slice
// of fake job ids that were "queued"
func (v *triggerImpl) RunTrigger(ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	return &pb.RunTriggerResponse{
		JobIds: []string{"000", "123", "999"},
	}, nil
}

func (v *triggerImpl) GetJobStream(req *pb.GetJobStreamRequest, srv pb.Waypoint_GetJobStreamServer) error {
	// send the beginning jobstream messages

	// We expect a start msg before reading the rest of the stream
	startMsg := &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Open_{
			Open: &pb.GetJobStreamResponse_Open{},
		},
	}
	srv.Send(startMsg)

	// Now we send some random stream data to be decoded by the tests

	lineMsg := &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Terminal_{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "testing",
							},
						},
					},
				},
			},
		},
	}
	srv.Send(lineMsg)

	lineMsg = &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Terminal_{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Line_{
							Line: &pb.GetJobStreamResponse_Terminal_Event_Line{
								Msg: "another one",
							},
						},
					},
				},
			},
		},
	}
	srv.Send(lineMsg)

	lineMsg = &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Terminal_{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Status_{
							Status: &pb.GetJobStreamResponse_Terminal_Event_Status{
								Status: "OK",
								Msg:    "a status msg",
							},
						},
					},
				},
			},
		},
	}
	srv.Send(lineMsg)

	// TODO add the rest of the events

	return nil
}
