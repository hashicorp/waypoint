package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
	sshterm "golang.org/x/crypto/ssh/terminal"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SnapshotBackupCommand struct {
	*baseCommand
}

// initWriter inspects args to figure out where the snapshot will be written to. It
// supports args[0] being '-' to force writing to stdout.
func (c *SnapshotBackupCommand) initWriter(args []string) (io.Writer, io.Closer, error) {
	if len(args) >= 1 {
		if args[0] == "-" {
			return os.Stdout, nil, nil
		}

		f, err := os.Create(args[0])
		if err != nil {
			return nil, nil, err
		}

		return f, f, nil
	}

	f := os.Stdout

	if sshterm.IsTerminal(int(f.Fd())) {
		return nil, nil, fmt.Errorf("stdout is a terminal, refusing to pollute (use '-' to force)")
	}

	return f, nil, nil
}

func (c *SnapshotBackupCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	client := c.project.Client()

	w, closer, err := c.initWriter(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output: %s", err)
		return 1
	}

	if closer != nil {
		defer closer.Close()
	}

	stream, err := client.CreateSnapshot(c.Ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate snapshot: %s", err)
		return 1
	}

	resp, err := stream.Recv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to receive snapshot start message: %s", err)
		return 1
	}

	if _, ok := resp.Event.(*pb.CreateSnapshotResponse_Open_); !ok {
		fmt.Fprintf(os.Stderr, "failed to receive snapshot start message: %s", err)
		return 1
	}

	for {
		ev, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Fprintf(os.Stderr, "error receiving snapshot data: %s", err)
			return 1
		}

		chunk, ok := ev.Event.(*pb.CreateSnapshotResponse_Chunk)
		if ok {
			_, err = w.Write(chunk.Chunk)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing snapshot data: %s", err)
				return 1
			}
		} else {
			fmt.Fprintf(os.Stderr, "unexpected protocol value: %T", ev.Event)
			return 1
		}
	}

	if w != os.Stdout {
		c.ui.Output("Snapshot written to '%s'", args[0])
	}

	return 0
}

func (c *SnapshotBackupCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *SnapshotBackupCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictFiles("")
}

func (c *SnapshotBackupCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SnapshotBackupCommand) Synopsis() string {
	return "Write a backup of the server data."
}

func (c *SnapshotBackupCommand) Help() string {
	return formatHelp(`
Usage: waypoint server snapshot [<filenamp>]

	Generate a snapshot from the current server and write it to a file specified
	by the given name. If no name is specified and standard out is not a terminal,
	the backup will written to standard out. Using a name of '-' will force writing
	to standard out.

` + c.Flags().Help())
}
