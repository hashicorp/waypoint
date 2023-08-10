// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package inlinekeepalive

import (
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestKeepaliveClientStream_RecvMsg(t *testing.T) {
	require := require.New(t)
	log := hclog.Default()
	log.SetLevel(hclog.Trace)

	// messages to play to the interceptor. Will be sent in top-down order.
	messages := []proto.Message{
		&pb.InlineKeepalive{Signature: KeepaliveProtoSignature},
		&pb.LogBatch{
			DeploymentId: "first",
		},

		&pb.InlineKeepalive{Signature: KeepaliveProtoSignature},
		&pb.InlineKeepalive{Signature: KeepaliveProtoSignature},
		&pb.InlineKeepalive{Signature: KeepaliveProtoSignature},
		&pb.LogBatch{
			DeploymentId: "second",
		},
	}

	k := &KeepaliveClientStream{
		log:     log,
		handler: TestClientStream(t, messages),
		sendMx:  &sync.Mutex{},
	}

	msg := &pb.LogBatch{}

	// Can intercept one keepalive

	require.NoError(k.RecvMsg(msg))
	require.Equal(msg.DeploymentId, "first")

	// Can intercept many keepalives

	require.NoError(k.RecvMsg(msg))
	require.Equal(msg.DeploymentId, "second")
}
