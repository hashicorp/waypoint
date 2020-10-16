// Package component has component implementations for the various
// resulting types.
package component

import (
	"github.com/golang/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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
