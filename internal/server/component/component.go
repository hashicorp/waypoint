// Package component has component implementations for the various
// resulting types.
package component

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
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
