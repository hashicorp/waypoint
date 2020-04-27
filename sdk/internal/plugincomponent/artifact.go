package plugincomponent

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/hashicorp/waypoint/sdk/component"
)

// Artifact implements component.Artifact.
type Artifact struct {
	Any *any.Any
}

func (c *Artifact) Proto() proto.Message { return c.Any }

var _ component.Artifact = (*Artifact)(nil)
