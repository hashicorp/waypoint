// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/waypoint/pkg/server/gen/mocks"
	"github.com/stretchr/testify/require"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// For a note about the magic value 445DHu, see exec_test.go

func TestHandleTrigger(t *testing.T) {
	require := require.New(t)

	// Get our gRPC server
	impl := &triggerImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	// Mock a request
	resp, err := http.Get(httpServer.URL + "/v1/trigger/123" + "?token=445DHu&stream=true")
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

	msg = Message{}
	require.NoError(decoder.Decode(&msg))

	require.Equal(msg.ValueType, "TerminalEventStep")
	require.Equal(msg.Value, "step msg")

	msg = Message{}
	require.NoError(decoder.Decode(&msg))

	require.Equal(msg.ValueType, "Complete")
}

func TestHandleTrigger_BadFailures(t *testing.T) {
	require := require.New(t)

	// Get our gRPC server
	impl := &triggerBadImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	// Mock a request
	resp, err := http.Get(httpServer.URL + "/v1/trigger/123" + "?token=445DHu&stream=true")
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

	require.Equal(msg.ValueType, "Error")
	require.Equal(msg.ExitCode, "1")
	require.NotNil(msg.Error)
}

func TestHandleTrigger_CancelStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	require := require.New(t)

	// Get our gRPC server
	impl := &triggerImpl{} // use the good implementation which sends lots of events
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	// Mock a request
	req, err := http.NewRequest("GET", httpServer.URL+"/v1/trigger/123"+"?token=445DHu&stream=true", nil)
	if err != nil {
		t.Errorf("failed to make http request: %s", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("failed to make http request: %s", err)
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

	// cancel the stream
	cancel()
	require.Error(decoder.Decode(&msg))
}

type triggerImpl struct {
	sync.Mutex
	mocks.WaypointServer
	pb.UnsafeWaypointServer
}

type triggerBadImpl struct {
	sync.Mutex
	mocks.WaypointServer
	pb.UnsafeWaypointServer
}

// RunTrigger mocks out a "good" RunTrigger execution and returns a slice
// of fake job ids that were "queued"
func (v *triggerImpl) RunTrigger(ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	return &pb.RunTriggerResponse{
		JobIds: []string{"123"},
	}, nil
}

// RunTrigger mocks out a "good" RunTrigger execution and returns a slice
// of fake job ids that were "queued"
func (v *triggerBadImpl) RunTrigger(ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	return &pb.RunTriggerResponse{
		JobIds: []string{"123"},
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

	lineMsg = &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Terminal_{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Events: []*pb.GetJobStreamResponse_Terminal_Event{
					{
						Event: &pb.GetJobStreamResponse_Terminal_Event_Step_{
							Step: &pb.GetJobStreamResponse_Terminal_Event_Step{
								Id:  1,
								Msg: "step msg",
							},
						},
					},
				},
			},
		},
	}
	srv.Send(lineMsg)

	lineMsg = &pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Complete_{
			Complete: &pb.GetJobStreamResponse_Complete{},
		},
	}
	srv.Send(lineMsg)

	return nil
}

func (v *triggerBadImpl) GetJobStream(req *pb.GetJobStreamRequest, srv pb.Waypoint_GetJobStreamServer) error {
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
		Event: &pb.GetJobStreamResponse_Error_{
			Error: &pb.GetJobStreamResponse_Error{
				Error: &pbstatus.Status{},
			},
		},
	}
	srv.Send(lineMsg)

	return nil
}
