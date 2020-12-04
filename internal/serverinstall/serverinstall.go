package serverinstall

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type Installer interface {
	Install(ctx context.Context, ui terminal.UI, log hclog.Logger) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, error)
	InstallFlags(*flag.Set)
}

var Platforms = map[string]Installer{
	"kubernetes": &K8sInstaller{},
	"nomad":      &NomadInstaller{},
	"docker":     &DockerInstaller{},
}
