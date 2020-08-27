package plugin

import (
	"github.com/hashicorp/waypoint/sdk/component"
)

type mix_Platform_Authenticator struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.PlatformReleaser
	component.WorkspaceDestroyer
}

type mix_Platform_Destroy struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.PlatformReleaser
	component.Destroyer
	component.WorkspaceDestroyer
}

type mix_Platform_Log struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.PlatformReleaser
	component.LogPlatform
	component.WorkspaceDestroyer
}

type mix_Platform_Log_Destroy struct {
	component.Authenticator
	component.ConfigurableNotify
	component.Platform
	component.PlatformReleaser
	component.LogPlatform
	component.Destroyer
	component.WorkspaceDestroyer
}
