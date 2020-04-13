package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/sys/unix"

	"github.com/mitchellh/devflow/internal/ceb"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return 1
	}

	// Init our core logic
	core, err := ceb.New(context.Background(),
		ceb.WithEnvDefaults(),
	)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Error initializing Devflow entrypoint: %s\n", err)
		return 1
	}
	defer core.Close()

	// Exec requires a full path to a binary. If we weren't given an absolute
	// path then we need to look it up via the PATH.
	if !filepath.IsAbs(args[0]) {
		path, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintf(flag.CommandLine.Output(),
				"Error execing process: %s\n", err)
			usage()
			return 1
		}

		args[0] = path
	}

	err = unix.Exec(args[0], args, os.Environ())
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Error execing process: %s\n", err)
		usage()
		return 1
	}

	return 0
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		strings.TrimLeftFunc(usageText, unicode.IsSpace),
		os.Args[0])
	flag.PrintDefaults()
}

const usageText = `
Usage: %[1]s [cmd] [args...]

    This the custom entrypoint to support Devflow. It will re-execute any
    command given after configuring the environment for usage with Devflow.

`
