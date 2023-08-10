// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// inlinekeepalive is a package that sends "keepalive" messages over
// existing grpc streams.

// This is designed to work around a very specific problem: AWS ALBs (and perhaps other load balancers)
// do not respect http2 pings (https://stackoverflow.com/questions/66818645/http2-ping-frames-over-aws-alb-grpc-keepalive-ping)
// As a result, if a stream sits "idle" for more than the LB timeout (usually 60 seconds), the
// load balancer times out that specific http2 sub-connection.
// This is bad for waypoint, as we have streams that we expect to sit open and idle
// for long-ish periods of time (i.e. the GetLogStream stream that backs `waypoint logs`),
// and those streams don't yet have transparent reconnect/resume functionality. It's
// a bad experience for users for a log stream with no messages to time out after 60 seconds.

// This package provides GRPC server and client interceptors. Every time a new streaming client is created,
// one interceptor will create a new goroutine to send "keepalive" protobuf messages that use a high-order field.
// The other interceptor will detect and intercept those messages without exposing them to the underlying handlers.
// In this method, we keep some traffic on all of our streams without needing to be aware
// of it in every RPC.

// It's important that neither party (client or server) receives these messages unexpectedly -
// if they don't have the interceptor to handle them, they will get passed onto the underlying
// handlers and cause unexpected behavior (likely panics). Clients are expected to add
// metadata to their outbound contexts indicating they can receive inline keepalives, and
// servers are expected to advertise via GetVersionInfo features to clients that they
// can receive them.

// Long-term, we expect load balancers will more universally support http2 pings. We also
// have improvements planned to make streaming methods able to be resumed transparently
// to users if they are unexpectedly disconnected. When either of those conditions are true,
// we can stop using inline keepalives.

package inlinekeepalive
