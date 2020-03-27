package cloudrun

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/mitchellh/devflow/builtin/google/proto"
	"github.com/mitchellh/devflow/sdk/component"
)

type Deployment struct {
	Value *pb.CloudRunDeployment
}

func (d *Deployment) Proto() proto.Message { return d.Value }

var _ component.Deployment = (*Deployment)(nil)
