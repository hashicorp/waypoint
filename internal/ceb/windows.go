// +build windows

package ceb

import (
	"context"
	"io"
	"os/exec"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) startExecGroup(es []*pb.EntrypointConfig_Exec, env []string) {
}

func (ceb *CEB) initChildCmd(ctx context.Context, cfg *config) error {
	return io.EOF
}

func (ceb *CEB) markChildCmdReady() {}
func (ceb *CEB) copyCmd(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
