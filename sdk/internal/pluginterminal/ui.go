package pluginterminal

import (
	"io"

	"github.com/hashicorp/waypoint/sdk/internal/stdio"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// UI is a wrapper around the basic terminal UI so that we can
// get the proper stdout/stderr streams.
type UI struct {
	*terminal.BasicUI
}

// OutputWriters implements UI to return the direct TTY output from
// the parent process and not our os.Stdout/stderr since those are
// redirected over the gRPC interface.
func (ui *UI) OutputWriters() (io.Writer, io.Writer, error) {
	return stdio.Stdout(), stdio.Stderr(), nil
}

var (
	_ terminal.UI = (*UI)(nil)
)
