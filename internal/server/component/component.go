// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package component has component implementations for the various
// resulting types.
package component

import (
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func Deployment(v *pb.Deployment) component.Deployment {
	return &deployment{Value: v}
}

type deployment struct {
	Value *pb.Deployment
}

func (d *deployment) Proto() proto.Message { return d.Value.Deployment }

var (
	_ component.Deployment     = (*deployment)(nil)
	_ component.ProtoMarshaler = (*deployment)(nil)
)
