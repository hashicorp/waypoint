package plugincomponent

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/hashicorp/waypoint/sdk/component"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// Release implements component.Release.
type Release struct {
	Any     *any.Any
	Release *pb.Release
}

func (c *Release) Proto() proto.Message { return c.Any }
func (c *Release) URL() string          { return c.Release.Url }

var _ component.Release = (*Release)(nil)
