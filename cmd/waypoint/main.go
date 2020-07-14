package main

import (
	// This is unfortunately the hackiest thing I've ever had to do with Go.
	// protobuf-go disallows multiple files of the same filename and issues
	// a warning when they exist. This is extremely annoying because the
	// proto files themselves have declared packages that are unique. This
	// doesn't appear to be configurable in any way.
	//
	// We therefore have to import this first to ensure the init() is called
	// in this package first to silence any log output. This is guaranteed by
	// the Go spec: https://golang.org/ref/spec#Program_initialization_and_execution
	_ "github.com/hashicorp/waypoint/internal/pkg/logsilence"

	"os"
	"path/filepath"

	"github.com/hashicorp/waypoint/internal/cli"
)

func main() {
	// Make args[0] just the name of the executable since it is used in logs.
	os.Args[0] = filepath.Base(os.Args[0])

	os.Exit(cli.Main(os.Args))
}
