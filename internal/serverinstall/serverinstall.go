package serverinstall

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverinstall/config"
	"github.com/hashicorp/waypoint/internal/serverinstall/docker"
	"github.com/hashicorp/waypoint/internal/serverinstall/k8s"
	"github.com/hashicorp/waypoint/internal/serverinstall/nomad"
)

type ServerPlatformInstaller interface {
	Install(ctx context.Context, ui terminal.UI, log hclog.Logger) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, error)
}

// func (c *BaseConfig) NewServerPlatformInstaller
func NewServerPlatformInstaller(c *config.BaseConfig, platform string) (ServerPlatformInstaller, error) {
	var p ServerPlatformInstaller
	switch platform {
	case "docker":
		p = &docker.Platform{Config: c}
	case "kubernetes":
		p = &k8s.Platform{Config: c}
	case "nomad":
		p = &nomad.Platform{Config: c}
	default:
		return nil, fmt.Errorf("unknown server platform: %s", platform)
	}
	return p, nil
}
