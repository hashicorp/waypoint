package plugincomponent

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/hashicorp/waypoint/sdk/component"
)

// Artifact implements component.Artifact.
type Artifact struct {
	Any       *any.Any
	LabelsVal map[string]string
}

func (c *Artifact) Proto() proto.Message { return c.Any }

func (c *Artifact) Labels() map[string]string { return c.LabelsVal }

var (
	_ component.Artifact = (*Artifact)(nil)
)
