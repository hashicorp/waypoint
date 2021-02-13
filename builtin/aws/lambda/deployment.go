package lambda

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

var _ component.Deployment = (*Deployment)(nil)
