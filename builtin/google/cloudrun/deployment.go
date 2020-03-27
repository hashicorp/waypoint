package cloudrun

import (
	"github.com/golang/protobuf/proto"

	"github.com/mitchellh/devflow/sdk/component"
)

func (d *Deployment) Proto() proto.Message { return d }

var _ component.Deployment = (*Deployment)(nil)
