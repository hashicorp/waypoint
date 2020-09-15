package plugin

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

type mix_Builder_Authenticator struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Builder
	component.Documented
}
