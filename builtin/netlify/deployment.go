package netlify

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)
