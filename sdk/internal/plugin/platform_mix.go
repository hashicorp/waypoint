package plugin

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

type mix_Platform_Destroy struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.Destroyer
}

type mix_Platform_Log struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.LogPlatform
}

type mix_Platform_Log_Destroy struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.LogPlatform
	component.Destroyer
}
