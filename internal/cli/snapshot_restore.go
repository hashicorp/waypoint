// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/waypoint/internal/clisnapshot"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/posener/complete"
	sshterm "golang.org/x/crypto/ssh/terminal"
)

type SnapshotRestoreCommand struct {
	*baseCommand

	// set via -exit, indicates we should tell the server to exit now to be restarted.
	flagExit bool
}

// initWriter inspects args to figure out where the snapshot will be read from. It
// supports args[0] being '-' to force reading from stdin.
func (c *SnapshotRestoreCommand) initReader(args []string) (io.Reader, io.Closer, error) {
	if len(args) >= 1 {
		if args[0] == "-" {
			return os.Stdin, nil, nil
		}

		f, err := os.Open(args[0])
		if err != nil {
			return nil, nil, err
		}

		return f, f, nil
	}

	f := os.Stdin

	if sshterm.IsTerminal(int(f.Fd())) {
		return nil, nil, fmt.Errorf("stdin is a terminal, refusing to use (use '-' to force)")
	}

	return f, nil, nil
}

func (c *SnapshotRestoreCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	r, closer, err := c.initReader(c.args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output: %s", err)
		return 1
	}

	if closer != nil {
		defer closer.Close()
	}

	if err := clisnapshot.ReadSnapshot(c.Ctx, c.project.Client(), r, c.flagExit); err != nil {
		fmt.Fprintf(os.Stderr, "Error restoring Snapshot: %s", err)
		return 1
	}

	if r == os.Stdin {
		c.ui.Output("Server data restored.")
	} else {
		c.ui.Output("Server data restored from '%s'.", c.args[0])
	}

	return 0
}

func (c *SnapshotRestoreCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "exit",
			Target:  &c.flagExit,
			Usage:   "After restoring, the server should exit so it can be restarted.",
			Default: false,
		})
	})
}

func (c *SnapshotRestoreCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictFiles("")
}

func (c *SnapshotRestoreCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SnapshotRestoreCommand) Synopsis() string {
	return "Stage a snapshot on the current server for data restoration"
}

func (c *SnapshotRestoreCommand) Help() string {
	return formatHelp(`
Usage: waypoint server restore [-exit] [<filename>]

Stage a backup snapshot within the current server. The data in the snapshot is not restored
immediately, but rather staged such that on the next server start, it will be restored.

If -exit is passed, the server process will exit after staging the data. This allows a process
monitor to restart the server, where it will see the staged snapshot and restore the data.

If -exit is not passed, an operator must restart the server manually to finish the restoration
process.

The argument should be to a file written previously by 'waypoint server snapshot'.
If no name is specified and standard input is not a terminal, the backup will read from
standard input. Using a name of '-' will force reading from standard input.

` + c.Flags().Help())
}
