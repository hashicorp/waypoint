package plugin

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

type mix_Registry_Authenticator struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Registry
}
