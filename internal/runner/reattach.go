// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// reattachClient wraps another RunnerJobStream gRPC stream client and
// performs automatic reattach on server unavailable errors.
//
// This should be only used to wrap a client around a job that is
// already acked, not immediately after job stream opening.
//
// This maintains the same concurrency requirements as a normal gRPC
// stream client. See the gRPC docs for specific details.
//
// This client is tested via the accept tests testing server down behavior.
type reattachClient struct {
	// Set these
	ctx    context.Context
	client pb.Waypoint_RunnerJobStreamClient
	log    hclog.Logger
	runner *Runner
	jobId  string

	// Internals, do not set
	clientLock sync.Mutex
	clientGen  uint64
}

func (c *reattachClient) Send(req *pb.RunnerJobStreamRequest) error {
	return c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		return client.Send(req)
	})
}

func (c *reattachClient) Recv() (*pb.RunnerJobStreamResponse, error) {
	var resp *pb.RunnerJobStreamResponse
	err := c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		var err error
		resp, err = client.Recv()
		return err
	})

	return resp, err
}

func (c *reattachClient) Header() (metadata.MD, error) {
	var resp metadata.MD
	err := c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		var err error
		resp, err = client.Header()
		return err
	})

	return resp, err
}

func (c *reattachClient) Trailer() metadata.MD {
	var resp metadata.MD
	err := c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		resp = client.Trailer()
		return nil
	})
	if err != nil {
		// Should never happen.
		panic(err)
	}

	return resp
}

func (c *reattachClient) Context() context.Context {
	var resp context.Context
	err := c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		resp = client.Context()
		return nil
	})
	if err != nil {
		// Should never happen.
		panic(err)
	}

	return resp
}

func (c *reattachClient) CloseSend() error {
	return c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		return client.CloseSend()
	})
}

func (c *reattachClient) SendMsg(m interface{}) error {
	return c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		return client.SendMsg(m)
	})
}

func (c *reattachClient) RecvMsg(m interface{}) error {
	return c.do(func(client pb.Waypoint_RunnerJobStreamClient) error {
		return client.RecvMsg(m)
	})
}

func (c *reattachClient) do(f func(client pb.Waypoint_RunnerJobStreamClient) error) error {
	// Shorthand for common attribute access
	log := c.log
	r := c.runner

	// Get our configuration state value. We use this so that we can detect
	// when we've reconnected during failures.
	stateGen := r.readState(&r.stateConfig)

	// We use this as a way to avoid having unlocks everywhere since we're
	// using a loop.
	locked := false
	defer func() {
		if locked {
			c.clientLock.Unlock()
			locked = false
		}
	}()

	for {
		// If we don't have a lock (our first time through), then we grab it
		// so we can get our client. If we have a lock (a connection retry)
		// then we avoid this. We always unlock after this.
		if !locked {
			c.clientLock.Lock()
		}

		// Get our current client
		client := c.client
		clientGen := c.clientGen
		c.clientLock.Unlock()
		locked = false

		// Do operation
		err := f(client)
		if err == nil {
			// No error means we're done.
			return nil
		}

		// We attempt a reconnect on Unavailable and NotFound but otherwise
		// we return the error as-is.
		if status.Code(err) != codes.Unavailable &&
			status.Code(err) != codes.NotFound {
			return err
		}
		log.Warn("server down during job API, will reconnect")

		// Lock so we only ever reconnect once at a time.
		c.clientLock.Lock()
		locked = true
		if c.clientGen != clientGen {
			// If the client gen changed, we already reconnected.
			continue
		}

		// This is the label that should be jumped to if reconnection
		// fails below. THE LOCK MUST STILL BE HELD when this is jumped to.
		// This lets us retry the reconnect without having to retry
		// the argument `f()` logic.
	RETRY_RECONNECT:

		// Since this is a disconnect, we have to wait for our
		// RunnerConfig stream to re-establish. We wait for the config
		// generation to increment.
		log.Info("waiting for runner to re-register before reconnecting")
		if r.waitStateGreater(&r.stateConfig, stateGen) {
			return status.Error(codes.Internal, "early exit while waiting for reconnect")
		}

		// Get the new state generation in case we get another disconnect
		// we can continue retrying and detecting new connections.
		stateGen = r.readState(&r.stateConfig)

		log.Info("opening job stream")
		client, err = r.client.RunnerJobStream(c.ctx, grpc.WaitForReady(true))
		if err != nil {
			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.NotFound {
				goto RETRY_RECONNECT
			}

			return err
		}

		// Send our request
		log.Info("sending job request with reattach")
		if err := client.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Request_{
				Request: &pb.RunnerJobStreamRequest_Request{
					RunnerId:      r.id,
					ReattachJobId: c.jobId,
				},
			},
		}); err != nil {
			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.NotFound {
				goto RETRY_RECONNECT
			}

			return err
		}

		// Wait for an assignment
		log.Info("waiting for job assignment")
		resp, err := client.Recv()
		if err != nil {
			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.NotFound {
				goto RETRY_RECONNECT
			}

			return err
		}

		// We received an assignment!
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		if !ok {
			return status.Errorf(codes.Aborted,
				"expected job assignment, server sent %T",
				resp.Event)
		}
		if assignment.Assignment.Job.Id != c.jobId {
			// This should never happen (means the server is buggy).
			// Therefore, we don't attempt any graceful cleanup we just
			// exit and let the server timeout the job.
			return status.Errorf(codes.Aborted,
				"received invalid job assignment, aborting")
		}
		log.Info("job assignment received matching reattach job ID")

		// Ack the assignment
		log.Info("acking job assignment")
		if err := client.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}); err != nil {
			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.NotFound {
				goto RETRY_RECONNECT
			}

			return err
		}

		// Increment our generation and set our new client so everyone uses it
		c.clientGen++
		c.client = client

		// We are now re-established with the server and can continue
		// our job execution.
	}
}

var _ pb.Waypoint_RunnerJobStreamClient = (*reattachClient)(nil)
