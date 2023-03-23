// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inlinekeepalive

import (
	"context"

	"github.com/mitchellh/go-testing-interface"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// TestClientStream returns a grpc.ClientStream that plays the given messages when Recv is called
func TestClientStream(t testing.T, recvMessages []proto.Message) grpc.ClientStream {
	return &testStream{
		recvMessages: recvMessages,
		t:            t,
	}
}

// TestServerStream returns a grpc.ServerStream that plays the given messages when Recv is called
func TestServerStream(t testing.T, recvMessages []proto.Message) grpc.ServerStream {
	return &testStream{
		recvMessages: recvMessages,
		t:            t,
	}
}

type testStream struct {
	recvMessages []proto.Message
	t            testing.T
}

// RecvMsg intercepts keepalive messages and does not pass them
// along to the handler.
func (t *testStream) RecvMsg(m interface{}) error {
	msg := t.recvMessages[0]
	t.recvMessages = t.recvMessages[1:]

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return proto.Unmarshal(bytes, m.(proto.Message))
}

// Required to implement grpc.ServerStream and grpc.ClientStream, but unused in this test

func (t *testStream) SetHeader(metadata.MD) error {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}

func (t *testStream) SendHeader(metadata.MD) error {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}

func (t *testStream) SetTrailer(metadata.MD) {
	t.t.Fatal("unexpected usage of TestStream")
}

func (t *testStream) Header() (metadata.MD, error) {
	t.t.Fatal("unexpected usage of TestStream")
	return nil, nil
}

func (t *testStream) Trailer() metadata.MD {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}

func (t *testStream) CloseSend() error {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}

func (t *testStream) Context() context.Context {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}

func (t *testStream) SendMsg(m interface{}) error {
	t.t.Fatal("unexpected usage of TestStream")
	return nil
}
