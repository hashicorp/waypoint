package plugincomponent

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/hashicorp/waypoint/sdk/component"
)

// Deployment implements component.Deployment.
type Deployment struct {
	Any *any.Any
}

func (c *Deployment) Proto() proto.Message { return c.Any }
func (c *Deployment) String() string       { return "" }

var _ component.Deployment = (*Deployment)(nil)
