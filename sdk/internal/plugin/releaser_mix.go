package plugin

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

type mix_ReleaseManager_Authenticator struct {
	component.Authenticator
	component.ConfigurableNotify
	component.ReleaseManager
	component.Destroyer
	component.WorkspaceDestroyer
	component.Documented
}
