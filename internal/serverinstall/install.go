package serverinstall

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ServerPlatformInstaller interface {
	Install(ctx context.Context, ui terminal.UI, log hclog.Logger) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, error)
}

func NewServerPlatformInstaller(c *Config, platform string) (ServerPlatformInstaller, error) {
	var p ServerPlatformInstaller
	switch platform {
	case "docker":
		p = &PlatformDocker{config: c}
	case "kubernetes":
		p = &PlatformKubernetes{Config: c}
	case "nomad":
		p = &PlatformNomad{config: c}
	default:
		return nil, fmt.Errorf("unknown server platform: %s", platform)
	}
	return p, nil
}
